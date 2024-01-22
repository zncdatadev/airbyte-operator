package controller

import (
	stackv1alpha1 "github.com/zncdata-labs/airbyte-operator/api/v1alpha1"
	"github.com/zncdata-labs/operator-go/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"strconv"
)

func (r *AirbyteReconciler) makeConnectorBuilderServerService(instance *stackv1alpha1.Airbyte) ([]*corev1.Service, error) {
	var services []*corev1.Service

	if instance.Spec.ConnectorBuilderServer.RoleGroups != nil {
		for roleGroupName, roleGroup := range instance.Spec.ConnectorBuilderServer.RoleGroups {
			svc, err := r.makeConnectorBuilderServerServiceForRoleGroup(instance, roleGroupName, roleGroup, r.Scheme)
			if err != nil {
				return nil, err
			}
			services = append(services, svc)
		}
	}

	return services, nil
}

func (r *AirbyteReconciler) makeConnectorBuilderServerServiceForRoleGroup(instance *stackv1alpha1.Airbyte, roleGroupName string, roleGroup *stackv1alpha1.RoleGroupConnectorBuilderServerSpec, schema *runtime.Scheme) (*corev1.Service, error) {
	labels := instance.GetLabels()

	additionalLabels := make(map[string]string)

	if roleGroup.Config != nil && roleGroup.Config.MatchLabels != nil {
		for k, v := range roleGroup.Config.MatchLabels {
			additionalLabels[k] = v
		}
	}

	mergedLabels := make(map[string]string)
	for key, value := range labels {
		mergedLabels[key] = value
	}
	for key, value := range additionalLabels {
		mergedLabels[key] = value
	}

	var port int32
	var serviceType corev1.ServiceType
	var annotations map[string]string

	if roleGroup != nil && roleGroup.Config != nil && roleGroup.Config.Service != nil {
		port = roleGroup.Config.Service.Port
		serviceType = roleGroup.Config.Service.Type
		annotations = roleGroup.Config.Service.Annotations
	} else {
		port = instance.Spec.ConnectorBuilderServer.Service.Port
		serviceType = instance.Spec.ConnectorBuilderServer.Service.Type
		annotations = instance.Spec.ConnectorBuilderServer.Service.Annotations
	}

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        instance.GetNameWithSuffix("-airbyte-connector-builder-server-svc" + "-" + roleGroupName),
			Namespace:   instance.Namespace,
			Labels:      mergedLabels,
			Annotations: annotations,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port:     port,
					Name:     "http",
					Protocol: "TCP",
				},
			},
			Selector: mergedLabels,
			Type:     serviceType,
		},
	}
	err := ctrl.SetControllerReference(instance, svc, schema)
	if err != nil {
		r.Log.Error(err, "Failed to set controller reference for service")
		return nil, errors.Wrap(err, "Failed to set controller reference for service")
	}
	return svc, nil
}

func (r *AirbyteReconciler) makeConnectorBuilderServerDeploymentForRoleGroup(instance *stackv1alpha1.Airbyte, roleGroupName string, roleGroup *stackv1alpha1.RoleGroupConnectorBuilderServerSpec, schema *runtime.Scheme) *appsv1.Deployment {
	labels := instance.GetLabels()

	additionalLabels := make(map[string]string)

	if roleGroup != nil && roleGroup.Config.MatchLabels != nil {
		for k, v := range roleGroup.Config.MatchLabels {
			additionalLabels[k] = v
		}
	}

	mergedLabels := make(map[string]string)
	for key, value := range labels {
		mergedLabels[key] = value
	}
	for key, value := range additionalLabels {
		mergedLabels[key] = value
	}

	var image stackv1alpha1.ImageSpec
	var securityContext *corev1.PodSecurityContext
	var containerPorts []corev1.ContainerPort

	if roleGroup != nil && roleGroup.Config != nil && roleGroup.Config.Image != nil {
		image = *roleGroup.Config.Image
	} else {
		image = *instance.Spec.ConnectorBuilderServer.Image
	}

	if roleGroup != nil && roleGroup.Config != nil && roleGroup.Config.SecurityContext != nil {
		securityContext = roleGroup.Config.SecurityContext
	} else {
		securityContext = instance.Spec.ConnectorBuilderServer.SecurityContext
	}

	if roleGroup != nil && roleGroup.Config != nil && roleGroup.Config.Service != nil && roleGroup.Config.Service.Port != 0 {
		containerPorts = append(containerPorts, corev1.ContainerPort{
			ContainerPort: roleGroup.Config.Service.Port,
			Name:          "http",
			Protocol:      corev1.ProtocolTCP,
		})

	} else {
		containerPorts = append(containerPorts, corev1.ContainerPort{
			ContainerPort: instance.Spec.ConnectorBuilderServer.Service.Port,
			Name:          "http",
			Protocol:      corev1.ProtocolTCP,
		})
	}

	var envVars []corev1.EnvVar

	if roleGroup.Config != nil && roleGroup.Config.Debug != nil {
		if roleGroup.Config.Debug.Enabled {
			envVars = append(envVars, corev1.EnvVar{
				Name:  "JAVA_TOOL_OPTIONS",
				Value: "-Xdebug -agentlib:jdwp=transport=dt_socket,address=0.0.0.0:" + strconv.FormatInt(int64(roleGroup.Config.Debug.RemoteDebugPort), 10) + ",server=y,suspend=n",
			})
		} else if instance != nil && instance.Spec.ConnectorBuilderServer != nil && instance.Spec.ConnectorBuilderServer.RoleConfig.Debug != nil {
			if instance.Spec.ConnectorBuilderServer.RoleConfig.Debug.Enabled {
				envVars = append(envVars, corev1.EnvVar{
					Name:  "JAVA_TOOL_OPTIONS",
					Value: "-Xdebug -agentlib:jdwp=transport=dt_socket,address=0.0.0.0:" + strconv.FormatInt(int64(instance.Spec.ConnectorBuilderServer.RoleConfig.Debug.RemoteDebugPort), 10) + ",server=y,suspend=n",
				})
			}
		}
	}

	if instance.Spec.ClusterConfig != nil {
		if instance.Spec.ClusterConfig.DeploymentMode == "oss" {
			envVars = append(envVars, corev1.EnvVar{
				Name: "AIRBYTE_VERSION",
				ValueFrom: &corev1.EnvVarSource{
					ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
						Key: "AIRBYTE_VERSION",
						LocalObjectReference: corev1.LocalObjectReference{
							Name: instance.GetNameWithSuffix("-airbyte-env"),
						},
					},
				},
			})
		}
	}

	if roleGroup.Config != nil && roleGroup.Config.Secret != nil {
		envVars = appendEnvVarsFromSecret(envVars, roleGroup.Config.Secret, "connector-builder-server-secrets", roleGroupName)
	} else if instance.Spec.ConnectorBuilderServer != nil && instance.Spec.ConnectorBuilderServer.Secret != nil {
		envVars = appendEnvVarsFromSecret(envVars, instance.Spec.ConnectorBuilderServer.Secret, "connector-builder-server-secrets", roleGroupName)
	}

	if instance != nil && (instance.Spec.ConnectorBuilderServer != nil || instance.Spec.ClusterConfig != nil) {
		envVarsMap := make(map[string]string)

		if roleGroup.Config != nil && roleGroup.Config.EnvVars != nil {
			for key, value := range roleGroup.Config.EnvVars {
				envVarsMap[key] = value
			}
		} else if instance.Spec.ConnectorBuilderServer != nil && instance.Spec.ConnectorBuilderServer.RoleConfig.EnvVars != nil {
			for key, value := range instance.Spec.ConnectorBuilderServer.RoleConfig.EnvVars {
				envVarsMap[key] = value
			}
		}

		if instance.Spec.ClusterConfig != nil && instance.Spec.ClusterConfig.EnvVars != nil {
			for key, value := range instance.Spec.ClusterConfig.EnvVars {
				envVarsMap[key] = value
			}
		}

		for key, value := range envVarsMap {
			envVars = append(envVars, corev1.EnvVar{
				Name:  key,
				Value: value,
			})
		}
	}

	if roleGroup.Config != nil && roleGroup.Config.ExtraEnv != (corev1.EnvVar{}) {
		envVars = append(envVars, roleGroup.Config.ExtraEnv)
	} else if instance.Spec.ConnectorBuilderServer != nil && instance.Spec.ConnectorBuilderServer.RoleConfig.ExtraEnv != (corev1.EnvVar{}) {
		envVars = append(envVars, instance.Spec.ConnectorBuilderServer.RoleConfig.ExtraEnv)
	}

	if roleGroup.Config != nil && roleGroup.Config.Debug != nil && roleGroup.Config.Debug.Enabled {
		containerPorts = append(containerPorts, corev1.ContainerPort{
			ContainerPort: roleGroup.Config.Debug.RemoteDebugPort,
			Name:          "debug",
			Protocol:      corev1.ProtocolTCP,
		})
	} else if instance.Spec.ConnectorBuilderServer != nil && instance.Spec.ConnectorBuilderServer.RoleConfig.Debug != nil && instance.Spec.ConnectorBuilderServer.RoleConfig.Debug.Enabled {
		containerPorts = append(containerPorts, corev1.ContainerPort{
			ContainerPort: instance.Spec.ConnectorBuilderServer.RoleConfig.Debug.RemoteDebugPort,
			Name:          "debug",
			Protocol:      corev1.ProtocolTCP,
		})
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.GetNameWithSuffix("-connector-builder-server"),
			Namespace: instance.Namespace,
			Labels:    mergedLabels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &roleGroup.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: mergedLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: mergedLabels,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: instance.GetNameWithSuffix("-admin"),
					SecurityContext:    securityContext,
					Containers: []corev1.Container{
						{
							Name:            instance.Name,
							Image:           image.Repository + ":" + image.Tag,
							ImagePullPolicy: image.PullPolicy,
							Resources:       *roleGroup.Config.Resources,
							Env:             envVars,
							Ports:           containerPorts,
						},
					},
				},
			},
		},
	}
	ConnectorBuilderServerScheduler(instance, dep)

	err := ctrl.SetControllerReference(instance, dep, schema)
	if err != nil {
		r.Log.Error(err, "Failed to set controller reference for deployment")
		return nil
	}
	return dep
}

func (r *AirbyteReconciler) makeConnectorBuilderServerSecret(instance *stackv1alpha1.Airbyte) ([]*corev1.Secret, error) {
	var secrets []*corev1.Secret

	if instance.Spec.ConnectorBuilderServer.RoleGroups != nil {
		for roleGroupName, roleGroup := range instance.Spec.ConnectorBuilderServer.RoleGroups {
			secret, err := r.makeConnectorBuilderServerSecretForRoleGroup(instance, roleGroupName, roleGroup, r.Scheme)
			if err != nil {
				return nil, err
			}
			secrets = append(secrets, secret)
		}
	}

	return secrets, nil
}

func (r *AirbyteReconciler) makeConnectorBuilderServerSecretForRoleGroup(instance *stackv1alpha1.Airbyte, roleGroupName string, roleGroup *stackv1alpha1.RoleGroupConnectorBuilderServerSpec, schema *runtime.Scheme) (*corev1.Secret, error) {
	labels := instance.GetLabels()

	additionalLabels := make(map[string]string)

	if roleGroup != nil && roleGroup.Config.MatchLabels != nil {
		for k, v := range roleGroup.Config.MatchLabels {
			additionalLabels[k] = v
		}
	}

	mergedLabels := make(map[string]string)
	for key, value := range labels {
		mergedLabels[key] = value
	}
	for key, value := range additionalLabels {
		mergedLabels[key] = value
	}

	if instance.Spec.ConnectorBuilderServer.Secret != nil || roleGroup.Config.Secret != nil {
		var data map[string][]byte
		if roleGroup.Config != nil && roleGroup.Config.Secret != nil {
			data = createSecretData(roleGroup.Config.Secret)
		} else if instance.Spec.ConnectorBuilderServer.Secret != nil {
			data = createSecretData(instance.Spec.ConnectorBuilderServer.Secret)
		}

		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "connector-builder-server-secrets" + "-" + roleGroupName,
				Namespace: instance.Namespace,
				Labels:    mergedLabels,
			},
			Type: corev1.SecretTypeOpaque,
			Data: data,
		}

		if err := ctrl.SetControllerReference(instance, secret, schema); err != nil {
			return nil, errors.Wrap(err, "Failed to set controller reference for secret")
		}
		return secret, nil
	}

	return nil, nil
}

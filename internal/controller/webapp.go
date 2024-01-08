package controller

import (
	stackv1alpha1 "github.com/zncdata-labs/airbyte-operator/api/v1alpha1"
	"github.com/zncdata-labs/operator-go/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *AirbyteReconciler) makeWebAppService(instance *stackv1alpha1.Airbyte) ([]*corev1.Service, error) {
	var services []*corev1.Service

	if instance.Spec.WebApp.RoleGroups != nil {
		for roleGroupName, roleGroup := range instance.Spec.WebApp.RoleGroups {
			svc, err := r.makeWebAppServiceForRoleGroup(instance, roleGroupName, roleGroup, r.Scheme)
			if err != nil {
				return nil, err
			}
			services = append(services, svc)
		}
	}

	return services, nil
}

func (r *AirbyteReconciler) makeWebAppServiceForRoleGroup(instance *stackv1alpha1.Airbyte, roleGroupName string, roleGroup *stackv1alpha1.RoleGroupWebAppSpec, schema *runtime.Scheme) (*corev1.Service, error) {
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
		port = instance.Spec.WebApp.Service.Port
		serviceType = instance.Spec.WebApp.Service.Type
		annotations = instance.Spec.WebApp.Service.Annotations
	}

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        instance.GetNameWithSuffix("-airbyte-webapp-svc" + "-" + roleGroupName),
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

func (r *AirbyteReconciler) makeWebAppDeployment(instance *stackv1alpha1.Airbyte) ([]*appsv1.Deployment, error) {
	var deployments []*appsv1.Deployment

	if instance.Spec.WebApp.RoleGroups != nil {
		for roleGroupName, roleGroup := range instance.Spec.WebApp.RoleGroups {
			dep := r.makeWebAppDeploymentForRoleGroup(instance, roleGroupName, roleGroup, r.Scheme)
			if dep != nil {
				deployments = append(deployments, dep)
			}
		}
	}

	return deployments, nil
}

func (r *AirbyteReconciler) makeWebAppDeploymentForRoleGroup(instance *stackv1alpha1.Airbyte, roleGroupName string, roleGroup *stackv1alpha1.RoleGroupWebAppSpec, schema *runtime.Scheme) *appsv1.Deployment {
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
		image = *instance.Spec.WebApp.Image
	}

	if roleGroup != nil && roleGroup.Config != nil && roleGroup.Config.SecurityContext != nil {
		securityContext = roleGroup.Config.SecurityContext
	} else {
		securityContext = instance.Spec.WebApp.SecurityContext
	}

	if roleGroup != nil && roleGroup.Config != nil && roleGroup.Config.Service != nil && roleGroup.Config.Service.Port != 0 {
		containerPorts = append(containerPorts, corev1.ContainerPort{
			ContainerPort: roleGroup.Config.Service.Port,
			Name:          "http",
			Protocol:      corev1.ProtocolTCP,
		})

	} else {
		containerPorts = append(containerPorts, corev1.ContainerPort{
			ContainerPort: instance.Spec.WebApp.Service.Port,
			Name:          "http",
			Protocol:      corev1.ProtocolTCP,
		})
	}

	envVarNames := []string{"TRACKING_STRATEGY", "INTERNAL_API_HOST", "KEYCLOAK_INTERNAL_HOST", "CONNECTOR_BUILDER_API_HOST",
		"AIRBYTE_VERSION", "API_URL", "CONNECTOR_BUILDER_API_URL"}
	var envVars []corev1.EnvVar

	if instance != nil && instance.Spec.ClusterConfig != nil {
		if instance.Spec.ClusterConfig.DeploymentMode == "oss" {
			for _, envVarName := range envVarNames {
				envVar := corev1.EnvVar{
					Name: envVarName,
					ValueFrom: &corev1.EnvVarSource{
						ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
							Key: envVarName,
							LocalObjectReference: corev1.LocalObjectReference{
								Name: instance.GetNameWithSuffix("-airbyte-env"),
							},
						},
					},
				}
				envVars = append(envVars, envVar)
			}
		}
	}

	if roleGroup.Config != nil && roleGroup.Config.Secret != nil {
		envVars = appendEnvVarsFromSecret(envVars, roleGroup.Config.Secret, "webapp-secrets", roleGroupName)
	} else if instance.Spec.WebApp != nil && instance.Spec.WebApp.Secret != nil {
		envVars = appendEnvVarsFromSecret(envVars, instance.Spec.WebApp.Secret, "webapp-secrets", roleGroupName)
	}

	if instance != nil && (instance.Spec.WebApp != nil || instance.Spec.ClusterConfig != nil) {
		envVarsMap := make(map[string]string)

		if roleGroup.Config != nil && roleGroup.Config.EnvVars != nil {
			for key, value := range roleGroup.Config.EnvVars {
				envVarsMap[key] = value
			}
		} else if instance.Spec.WebApp != nil && instance.Spec.WebApp.RoleConfig.EnvVars != nil {
			for key, value := range instance.Spec.WebApp.RoleConfig.EnvVars {
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
	} else if instance != nil && instance.Spec.WebApp != nil && instance.Spec.WebApp.RoleConfig.ExtraEnv != (corev1.EnvVar{}) {
		envVars = append(envVars, instance.Spec.WebApp.RoleConfig.ExtraEnv)
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.GetNameWithSuffix("-webapp" + "-" + roleGroupName),
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
					ServiceAccountName: "airbyte-admin",
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

func (r *AirbyteReconciler) makeWebAppIngress(instance *stackv1alpha1.Airbyte) ([]*v1.Ingress, error) {
	var ing []*v1.Ingress

	if instance.Spec.WebApp.RoleGroups != nil {
		for roleGroupName, roleGroup := range instance.Spec.WebApp.RoleGroups {
			i, err := r.makeWebAppIngressForRoleGroup(instance, roleGroupName, roleGroup, r.Scheme)
			if err != nil {
				return nil, err
			}
			ing = append(ing, i)
		}
	}
	return ing, nil
}

func (r *AirbyteReconciler) makeWebAppIngressForRoleGroup(instance *stackv1alpha1.Airbyte, roleGroupName string, roleGroup *stackv1alpha1.RoleGroupWebAppSpec, schema *runtime.Scheme) (*v1.Ingress, error) {
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

	pt := v1.PathTypeImplementationSpecific

	var host string
	var port int32

	if roleGroup != nil && roleGroup.Config != nil && roleGroup.Config.Ingress != nil {
		if !roleGroup.Config.Ingress.Enabled {
			return nil, nil
		}
		host = roleGroup.Config.Ingress.Host
		if roleGroup.Config.Service != nil {
			port = roleGroup.Config.Service.Port
		} else if instance.Spec.WebApp.Service != nil {
			port = instance.Spec.WebApp.Service.Port
		}
	} else {
		if instance.Spec.WebApp.Ingress != nil && !instance.Spec.WebApp.Ingress.Enabled {
			return nil, nil
		}
		host = instance.Spec.WebApp.Ingress.Host
		if instance.Spec.WebApp.Service != nil {
			port = instance.Spec.WebApp.Service.Port
		}
	}

	ing := &v1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.GetNameWithSuffix(roleGroupName),
			Namespace: instance.Namespace,
			Labels:    mergedLabels,
		},
		Spec: v1.IngressSpec{
			Rules: []v1.IngressRule{
				{
					Host: host,
					IngressRuleValue: v1.IngressRuleValue{
						HTTP: &v1.HTTPIngressRuleValue{
							Paths: []v1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pt,
									Backend: v1.IngressBackend{
										Service: &v1.IngressServiceBackend{
											Name: instance.GetNameWithSuffix(roleGroupName),
											Port: v1.ServiceBackendPort{
												Number: port,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	err := ctrl.SetControllerReference(instance, ing, schema)
	if err != nil {
		r.Log.Error(err, "Failed to set controller reference for ingress")
		return nil, errors.Wrap(err, "Failed to set controller reference for ingress")
	}
	return ing, nil
}

func (r *AirbyteReconciler) makeWebAppSecret(instance *stackv1alpha1.Airbyte) ([]*corev1.Secret, error) {
	var secrets []*corev1.Secret

	if instance.Spec.WebApp.RoleGroups != nil {
		for roleGroupName, roleGroup := range instance.Spec.WebApp.RoleGroups {
			secret, err := r.makeWebAppSecretForRoleGroup(instance, roleGroupName, roleGroup, r.Scheme)
			if err != nil {
				return nil, err
			}
			secrets = append(secrets, secret)
		}
	}

	return secrets, nil
}

func (r *AirbyteReconciler) makeWebAppSecretForRoleGroup(instance *stackv1alpha1.Airbyte, roleGroupName string, roleGroup *stackv1alpha1.RoleGroupWebAppSpec, schema *runtime.Scheme) (*corev1.Secret, error) {
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

	if instance.Spec.WebApp.Secret != nil || roleGroup.Config.Secret != nil {
		var data map[string][]byte
		if roleGroup.Config != nil && roleGroup.Config.Secret != nil {
			data = createSecretData(roleGroup.Config.Secret)
		} else if instance.Spec.WebApp.Secret != nil {
			data = createSecretData(instance.Spec.WebApp.Secret)
		}

		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "webapp-secrets" + "-" + roleGroupName,
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

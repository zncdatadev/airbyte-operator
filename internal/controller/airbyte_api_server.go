package controller

import (
	"context"
	stackv1alpha1 "github.com/zncdata-labs/airbyte-operator/api/v1alpha1"
	"github.com/zncdata-labs/operator-go/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
)

func (r *AirbyteReconciler) reconcileApiServerService(ctx context.Context,
	instance *stackv1alpha1.Airbyte) error {
	roleGroups := convertRoleGroupToRoleConfigObject(instance, ApiServer)
	reconcileParams := r.createReconcileParams(ctx, roleGroups, instance, r.extractApiServerServiceForRoleGroup)
	if err := reconcileParams.createOrUpdateResource(); err != nil {
		return err
	}
	return nil
}

// extract Api server service for role group
func (r *AirbyteReconciler) extractApiServerServiceForRoleGroup(params ExtractorParams) (client.Object, error) {
	instance := params.instance.(*stackv1alpha1.Airbyte)
	server := instance.Spec.ApiServer
	groupCfg := params.roleGroup
	roleCfg := server.RoleConfig
	clusterCfg := params.cluster
	roleGroupName := params.roleGroupName
	mergedLabels := r.mergeLabels(groupCfg, instance.GetLabels(), clusterCfg, ApiServer)
	schema := params.scheme

	port, serviceType, annotations := getServiceInfo(groupCfg, roleCfg, clusterCfg)
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        createSvcNameForRoleGroup(instance, ApiServer, roleGroupName),
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

// reconcile Airbyte api server deployment
func (r *AirbyteReconciler) reconcileAirbyteApiServerDeployment(ctx context.Context,
	instance *stackv1alpha1.Airbyte) error {
	roleGroups := convertRoleGroupToRoleConfigObject(instance, ApiServer)
	reconcileParams := r.createReconcileParams(ctx, roleGroups, instance, r.extractApiServerDeploymentForRoleGroup)
	if err := reconcileParams.createOrUpdateResource(); err != nil {
		return err
	}
	return nil
}

func mergeEnvVarsForApiServerDeployment(instance *stackv1alpha1.Airbyte, roleGroup *stackv1alpha1.ApiServerRoleConfigSpec,
	roleGroupName string, containerPorts *[]corev1.ContainerPort) []corev1.EnvVar {
	var envVars []corev1.EnvVar

	//INTERNAL_API_HOST
	envVars = append(envVars, corev1.EnvVar{
		Name:  "INTERNAL_API_HOST",
		Value: getSvcHost(instance, Server, roleGroupName),
	})
	envVars = append(envVars, corev1.EnvVar{
		Name:  "AIRBYTE_API_HOST",
		Value: getSvcHost(instance, ApiServer, roleGroupName),
	})

	if roleGroup.Config != nil && roleGroup.Debug != nil {
		if roleGroup.Debug.Enabled {
			envVars = append(envVars, corev1.EnvVar{
				Name:  "JAVA_TOOL_OPTIONS",
				Value: "-Xdebug -agentlib:jdwp=transport=dt_socket,address=0.0.0.0:" + strconv.FormatInt(int64(roleGroup.Debug.RemoteDebugPort), 10) + ",server=y,suspend=n",
			})
		} else if instance != nil && instance.Spec.ApiServer != nil && instance.Spec.ApiServer.RoleConfig.Debug != nil {
			if instance.Spec.ApiServer.RoleConfig.Debug.Enabled {
				envVars = append(envVars, corev1.EnvVar{
					Name:  "JAVA_TOOL_OPTIONS",
					Value: "-Xdebug -agentlib:jdwp=transport=dt_socket,address=0.0.0.0:" + strconv.FormatInt(int64(instance.Spec.ApiServer.RoleConfig.Debug.RemoteDebugPort), 10) + ",server=y,suspend=n",
				})
			}
		}
	}

	if instance != nil && instance.Spec.ClusterConfig != nil {
		if instance.Spec.ClusterConfig.DeploymentMode == "oss" {
			envVars = append(envVars, corev1.EnvVar{
				Name: "AIRBYTE_VERSION",
				ValueFrom: &corev1.EnvVarSource{
					ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
						Key: "AIRBYTE_VERSION",
						LocalObjectReference: corev1.LocalObjectReference{
							Name: createNameForRoleGroup(instance, "env", ""),
						},
					},
				},
			})
		}
	}

	if roleGroup.Config != nil && roleGroup.Secret != nil {
		envVars = appendEnvVarsFromSecret(envVars, roleGroup.Secret, "server-secrets", roleGroupName)
	} else if instance.Spec.ApiServer != nil && instance.Spec.ApiServer.RoleConfig.Secret != nil {
		envVars = appendEnvVarsFromSecret(envVars, instance.Spec.ApiServer.RoleConfig.Secret, "server-secrets", roleGroupName)
	}

	if instance != nil && (instance.Spec.ApiServer != nil || instance.Spec.ClusterConfig != nil) {
		envVarsMap := make(map[string]string)

		if roleGroup.Config != nil && roleGroup.EnvVars != nil {
			for key, value := range roleGroup.EnvVars {
				envVarsMap[key] = value
			}
		} else if instance.Spec.ApiServer != nil && instance.Spec.ApiServer.RoleConfig.EnvVars != nil {
			for key, value := range instance.Spec.ApiServer.RoleConfig.EnvVars {
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

	if roleGroup.Config != nil && roleGroup.ExtraEnv != nil {
		envVars = append(envVars, *roleGroup.ExtraEnv)
	} else if instance.Spec.ApiServer != nil && instance.Spec.ApiServer.RoleConfig.ExtraEnv != nil {
		envVars = append(envVars, *instance.Spec.ApiServer.RoleConfig.ExtraEnv)
	}

	if roleGroup.Config != nil && roleGroup.Debug != nil && roleGroup.Debug.Enabled {
		*containerPorts = append(*containerPorts, corev1.ContainerPort{
			ContainerPort: roleGroup.Debug.RemoteDebugPort,
			Name:          "debug",
			Protocol:      corev1.ProtocolTCP,
		})
	} else if instance.Spec.ApiServer != nil && instance.Spec.ApiServer.RoleConfig.Debug != nil && instance.Spec.ApiServer.RoleConfig.Debug.Enabled {
		*containerPorts = append(*containerPorts, corev1.ContainerPort{
			ContainerPort: instance.Spec.ApiServer.RoleConfig.Debug.RemoteDebugPort,
			Name:          "debug",
			Protocol:      corev1.ProtocolTCP,
		})
	}
	return envVars
}

// extract Api server deployment for role group
func (r *AirbyteReconciler) extractApiServerDeploymentForRoleGroup(params ExtractorParams) (client.Object, error) {
	instance := params.instance.(*stackv1alpha1.Airbyte)
	server := instance.Spec.ApiServer
	groupCfg := params.roleGroup
	roleCfg := server.RoleConfig
	clusterCfg := params.cluster
	roleGroupName := params.roleGroupName
	mergedLabels := r.mergeLabels(groupCfg, instance.GetLabels(), clusterCfg, ApiServer)
	schema := params.scheme

	realGroupCfg := groupCfg.(*stackv1alpha1.ApiServerRoleConfigSpec)
	image, securityContext, replicas, resources, containerPorts := getDeploymentInfo(groupCfg, roleCfg, clusterCfg)
	envVars := mergeEnvVarsForApiServerDeployment(instance, realGroupCfg, roleGroupName, &containerPorts)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      createNameForRoleGroup(instance, "api-server", roleGroupName),
			Namespace: instance.Namespace,
			Labels:    mergedLabels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: mergedLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: mergedLabels,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: createNameForRoleGroup(instance, "admin", ""),
					SecurityContext:    securityContext,
					Containers: []corev1.Container{
						{
							Name:            instance.Name,
							Image:           image.Repository + ":" + image.Tag,
							ImagePullPolicy: image.PullPolicy,
							Resources:       *resources,
							Env:             envVars,
							Ports:           containerPorts,
						},
					},
				},
			},
		},
	}

	ApiServerScheduler(instance, dep, realGroupCfg)

	err := ctrl.SetControllerReference(instance, dep, schema)
	if err != nil {
		r.Log.Error(err, "Failed to set controller reference for deployment")
		return nil, err
	}
	return dep, nil
}

// reconcile Airbyte api server secret
func (r *AirbyteReconciler) reconcileAirbyteApiServerSecret(ctx context.Context,
	instance *stackv1alpha1.Airbyte) error {
	roleGroups := convertRoleGroupToRoleConfigObject(instance, ApiServer)
	reconcileParams := r.createReconcileParams(ctx, roleGroups, instance, r.extractAirbyteApiServerSecretForRoleGroup)
	if err := reconcileParams.createOrUpdateResource(); err != nil {
		return err
	}
	return nil
}

// extract Airbyte api server secret for role group
func (r *AirbyteReconciler) extractAirbyteApiServerSecretForRoleGroup(params ExtractorParams) (client.Object, error) {
	instance := params.instance.(*stackv1alpha1.Airbyte)
	server := instance.Spec.ApiServer
	groupCfg := params.roleGroup
	roleCfg := server.RoleConfig
	clusterCfg := params.cluster
	roleGroupName := params.roleGroupName
	mergedLabels := r.mergeLabels(groupCfg, instance.GetLabels(), clusterCfg, ApiServer)
	schema := params.scheme

	realGroupCfg := groupCfg.(*stackv1alpha1.ApiServerRoleConfigSpec)
	secretField := roleCfg.Secret
	if realGroupCfg.Secret != nil {
		secretField = realGroupCfg.Secret
	}
	data := createSecretData(secretField)

	if data != nil {
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "server-secrets" + "-" + roleGroupName,
				Namespace: instance.Namespace,
				Labels:    mergedLabels,
			},
			Type: corev1.SecretTypeOpaque,
			Data: data,
		}
		if err := ctrl.SetControllerReference(instance, secret, schema); err != nil {
			return nil, errors.Wrap(err, "Failed to set controller reference for service")
		}
		return secret, nil
	}
	return nil, nil
}

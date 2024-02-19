package controller

import (
	"context"
	stackv1alpha1 "github.com/zncdata-labs/airbyte-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// reconcile Cron deployment
func (r *AirbyteReconciler) reconcileCronDeployment(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	roleGroups := convertRoleGroupToRoleConfigObject(instance, Cron)
	reconcileParams := r.createReconcileParams(ctx, roleGroups, instance, r.extractCronDeployment)
	if err := reconcileParams.createOrUpdateResource(); err != nil {
		return err
	}
	return nil
}

// extract Cron deployment for role group
func (r *AirbyteReconciler) extractCronDeployment(params ExtractorParams) (client.Object, error) {
	instance := params.instance.(*stackv1alpha1.Airbyte)
	server := instance.Spec.ApiServer
	groupCfg := params.roleGroup
	roleCfg := server.RoleConfig
	clusterCfg := params.cluster
	mergedLabels := r.mergeLabels(groupCfg, instance.GetLabels(), clusterCfg, Cron)
	schema := params.scheme

	image, securityContext, replicas, resources, _ := getDeploymentInfo(groupCfg, roleCfg, clusterCfg)
	envVars := r.mergeCronEnvVars(instance, params.roleGroupName)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      createNameForRoleGroup(instance, "cron", params.roleGroupName),
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
							Name:            "airbyte-cron",
							Image:           image.Repository + ":" + image.Tag,
							ImagePullPolicy: image.PullPolicy,
							Resources:       *resources,
							Env:             envVars,
						},
					},
				},
			},
		},
	}

	CronScheduler(instance, dep, groupCfg)

	err := ctrl.SetControllerReference(instance, dep, schema)
	if err != nil {
		r.Log.Error(err, "Failed to set controller reference for deployment")
		return nil, err
	}
	return dep, nil
}

// merge env var for cron deployment
func (r *AirbyteReconciler) mergeCronEnvVars(instance *stackv1alpha1.Airbyte, groupName string) []corev1.EnvVar {
	envVarNames := []string{"AIRBYTE_VERSION", "CONFIGS_DATABASE_MINIMUM_FLYWAY_MIGRATION_VERSION", "DATABASE_URL",
		"TRACKING_STRATEGY", "WORKFLOW_FAILURE_RESTART_DELAY_SECONDS", "WORKSPACE_DOCKER_MOUNT", "WORKSPACE_ROOT"}
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
								Name: createNameForRoleGroup(instance, "env", ""),
							},
						},
					},
				}
				envVars = append(envVars, envVar)
			}

			envVars = append(envVars, corev1.EnvVar{
				Name:  "TEMPORAL_HOST",
				Value: getSvcHost(instance, Temporal, groupName),
			})

			envVars = append(envVars, corev1.EnvVar{
				Name: "MICRONAUT_ENVIRONMENTS",
				ValueFrom: &corev1.EnvVarSource{
					ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
						Key: "CRON_MICRONAUT_ENVIRONMENTS",
						LocalObjectReference: corev1.LocalObjectReference{
							Name: createNameForRoleGroup(instance, "env", ""),
						},
					},
				},
			})

			envVars = append(envVars, corev1.EnvVar{
				Name: "DATABASE_USER",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						Key: "DATABASE_USER",
						LocalObjectReference: corev1.LocalObjectReference{
							Name: createNameForRoleGroup(instance, "secrets", ""),
						},
					},
				},
			})

			secretName := createNameForRoleGroup(instance, "secrets", "")
			secretKey := "DATABASE_PASSWORD"

			envVars = append(envVars, corev1.EnvVar{
				Name: "DATABASE_PASSWORD",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: secretName,
						},
						Key: secretKey,
					},
				},
			})
		}
	}

	if instance != nil && instance.Spec.Cron != nil && instance.Spec.Cron.RoleConfig.Secret != nil {
		for key := range instance.Spec.Cron.RoleConfig.Secret {
			envVars = append(envVars, corev1.EnvVar{
				Name: key,
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "cron-secrets",
						},
						Key: key,
					},
				},
			})
		}
	}

	if instance != nil && (instance.Spec.Cron != nil || instance.Spec.ClusterConfig != nil) {
		envVarsMap := make(map[string]string)

		if instance.Spec.Cron != nil && instance.Spec.Cron.RoleConfig.EnvVars != nil {
			for key, value := range instance.Spec.Cron.RoleConfig.EnvVars {
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

	if instance != nil && instance.Spec.Cron != nil && instance.Spec.Cron.RoleConfig.ExtraEnv != (corev1.EnvVar{}) {
		envVars = append(envVars, instance.Spec.Cron.RoleConfig.ExtraEnv)
	}
	return envVars
}

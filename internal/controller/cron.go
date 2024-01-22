package controller

import (
	stackv1alpha1 "github.com/zncdata-labs/airbyte-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *AirbyteReconciler) makeCronDeployment(instance *stackv1alpha1.Airbyte, schema *runtime.Scheme) *appsv1.Deployment {
	labels := instance.GetLabels()

	envVarNames := []string{"AIRBYTE_VERSION", "CONFIGS_DATABASE_MINIMUM_FLYWAY_MIGRATION_VERSION", "DATABASE_URL",
		"TEMPORAL_HOST", "TRACKING_STRATEGY", "WORKFLOW_FAILURE_RESTART_DELAY_SECONDS", "WORKSPACE_DOCKER_MOUNT", "WORKSPACE_ROOT"}
	var envVars []corev1.EnvVar

	if instance != nil && instance.Spec.Global != nil {
		if instance.Spec.Global.DeploymentMode == "oss" {
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

			envVars = append(envVars, corev1.EnvVar{
				Name: "MICRONAUT_ENVIRONMENTS",
				ValueFrom: &corev1.EnvVarSource{
					ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
						Key: "CRON_MICRONAUT_ENVIRONMENTS",
						LocalObjectReference: corev1.LocalObjectReference{
							Name: instance.GetNameWithSuffix("-airbyte-env"),
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
							Name: instance.GetNameWithSuffix("-airbyte-secrets"),
						},
					},
				},
			})

			secretName := instance.GetNameWithSuffix("-airbyte-secrets")
			secretKey := "DATABASE_PASSWORD"

			if instance != nil && instance.Spec.Global != nil && instance.Spec.Global.Database != nil {
				if instance.Spec.Global.Database.SecretName != "" {
					secretName = instance.Spec.Global.Database.SecretName
				}
				if instance.Spec.Global.Database.SecretValue != "" {
					secretKey = instance.Spec.Global.Database.SecretValue
				}
			}

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

	if instance != nil && instance.Spec.Cron != nil && instance.Spec.Cron.Secret != nil {
		for key := range instance.Spec.Cron.Secret {
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

	if instance != nil && (instance.Spec.Cron != nil || instance.Spec.Global != nil) {
		envVarsMap := make(map[string]string)

		if instance.Spec.Cron != nil && instance.Spec.Cron.EnvVars != nil {
			for key, value := range instance.Spec.Cron.EnvVars {
				envVarsMap[key] = value
			}
		}

		if instance.Spec.Global != nil && instance.Spec.Global.EnvVars != nil {
			for key, value := range instance.Spec.Global.EnvVars {
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

	if instance != nil && instance.Spec.Cron != nil && instance.Spec.Cron.ExtraEnv != (corev1.EnvVar{}) {
		envVars = append(envVars, instance.Spec.Cron.ExtraEnv)
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.GetNameWithSuffix("-airbyte-cron"),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &instance.Spec.Cron.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: instance.GetNameWithSuffix("-admin"),
					SecurityContext:    instance.Spec.SecurityContext,
					Containers: []corev1.Container{
						{
							Name:            "airbyte-cron",
							Image:           instance.Spec.Cron.Image.Repository + ":" + instance.Spec.Cron.Image.Tag,
							ImagePullPolicy: instance.Spec.Cron.Image.PullPolicy,
							Resources:       *instance.Spec.Cron.Resources,
							Env:             envVars,
						},
					},
				},
			},
		},
	}
	CronScheduler(instance, dep)

	err := ctrl.SetControllerReference(instance, dep, schema)
	if err != nil {
		r.Log.Error(err, "Failed to set controller reference for deployment")
		return nil
	}
	return dep
}

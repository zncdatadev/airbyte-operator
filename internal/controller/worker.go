package controller

import (
	stackv1alpha1 "github.com/zncdata-labs/airbyte-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *AirbyteReconciler) makeWorkerDeployments(instance *stackv1alpha1.Airbyte) []*appsv1.Deployment {
	var deployments []*appsv1.Deployment

	if instance.Spec.Worker.RoleGroups != nil {
		for roleGroupName, roleGroup := range instance.Spec.Worker.RoleGroups {
			dep := r.makeWorkerDeploymentForRoleGroup(instance, roleGroupName, roleGroup, r.Scheme)
			if dep != nil {
				deployments = append(deployments, dep)
			}
		}
	}

	return deployments
}

func (r *AirbyteReconciler) makeWorkerDeploymentForRoleGroup(instance *stackv1alpha1.Airbyte, roleGroupName string, roleGroup *stackv1alpha1.RoleGroupWorkerSpec, schema *runtime.Scheme) *appsv1.Deployment {
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

	if roleGroup != nil && roleGroup.Config != nil && roleGroup.Config.Image != nil {
		image = *roleGroup.Config.Image
	} else {
		image = *instance.Spec.Worker.Image
	}

	if roleGroup != nil && roleGroup.Config != nil && roleGroup.Config.SecurityContext != nil {
		securityContext = roleGroup.Config.SecurityContext
	} else {
		securityContext = instance.Spec.SecurityContext
	}

	portVarNames := []int32{9000, 9001, 9002, 9003, 9004, 9005, 9006, 9007, 9008, 9009, 9010,
		9011, 9012, 9013, 9014, 9015, 9016, 9017, 9018, 9019, 9020,
		9021, 9022, 9023, 9024, 9025, 9026, 9027, 9028, 9029, 9030}

	var containerPorts []corev1.ContainerPort

	// 添加带有名称的第一个端口
	containerPorts = append(containerPorts, corev1.ContainerPort{
		Name:          "heartbeat",
		ContainerPort: portVarNames[0],
	})

	// 添加其他端口，不带名称
	for _, portVarName := range portVarNames[1:] {
		containerPorts = append(containerPorts, corev1.ContainerPort{
			ContainerPort: portVarName,
		})
	}

	envVarNames := []string{"AIRBYTE_VERSION", "CONFIG_ROOT", "DATABASE_HOST", "DATABASE_PORT", "DATABASE_URL",
		"TRACKING_STRATEGY", "WORKSPACE_ROOT", "LOCAL_ROOT", "WEBAPP_URL", "TEMPORAL_HOST",
		"TEMPORAL_WORKER_PORTS", "S3_PATH_STYLE_ACCESS", "JOB_MAIN_CONTAINER_CPU_REQUEST", "JOB_MAIN_CONTAINER_CPU_LIMIT",
		"JOB_MAIN_CONTAINER_MEMORY_REQUEST", "JOB_MAIN_CONTAINER_MEMORY_LIMIT", "S3_LOG_BUCKET", "S3_LOG_BUCKET_REGION",
		"S3_MINIO_ENDPOINT", "GOOGLE_APPLICATION_CREDENTIALS", "GCS_LOG_BUCKET", "STATE_STORAGE_MINIO_BUCKET_NAME",
		"STATE_STORAGE_MINIO_ENDPOINT", "INTERNAL_API_HOST",
		"CONFIGS_DATABASE_MINIMUM_FLYWAY_MIGRATION_VERSION", "JOBS_DATABASE_MINIMUM_FLYWAY_MIGRATION_VERSION",
		"METRIC_CLIENT", "OTEL_COLLECTOR_ENDPOINT", "ACTIVITY_MAX_ATTEMPT", "ACTIVITY_INITIAL_DELAY_BETWEEN_ATTEMPTS_SECONDS",
		"ACTIVITY_MAX_DELAY_BETWEEN_ATTEMPTS_SECONDS", "WORKFLOW_FAILURE_RESTART_DELAY_SECONDS", "AUTO_DETECT_SCHEMA",
		"USE_STREAM_CAPABLE_STATE", "SHOULD_RUN_NOTIFY_WORKFLOWS", "WORKER_LOGS_STORAGE_TYPE", "WORKER_STATE_STORAGE_TYPE"}
	secretVarNames := []string{"DATABASE_PASSWORD", "DATABASE_USER", "STATE_STORAGE_MINIO_ACCESS_KEY", "STATE_STORAGE_MINIO_SECRET_ACCESS_KEY"}
	var envVars []corev1.EnvVar

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
				Key: "WORKERS_MICRONAUT_ENVIRONMENTS",
				LocalObjectReference: corev1.LocalObjectReference{
					Name: instance.GetNameWithSuffix("-airbyte-env"),
				},
			},
		},
	})

	envVars = append(envVars, corev1.EnvVar{
		Name: "JOB_KUBE_NAMESPACE",
		ValueFrom: &corev1.EnvVarSource{
			FieldRef: &corev1.ObjectFieldSelector{
				FieldPath: "metadata.namespace",
			},
		},
	})

	envVars = append(envVars, corev1.EnvVar{
		Name:  "WORKSPACE_DOCKER_MOUNT",
		Value: "workspace",
	}, corev1.EnvVar{
		Name:  "LOG_LEVEL",
		Value: "INFO",
	}, corev1.EnvVar{
		Name:  "JOB_KUBE_SERVICEACCOUNT",
		Value: instance.Spec.ClusterConfig.ServiceAccountName,
	})

	envVars = append(envVars, corev1.EnvVar{
		Name:  "LOG_LEVEL",
		Value: "INFO",
	})

	envVars = append(envVars, corev1.EnvVar{
		Name: "AWS_ACCESS_KEY_ID",
		ValueFrom: &corev1.EnvVarSource{
			SecretKeyRef: &corev1.SecretKeySelector{
				Key: "STATE_STORAGE_MINIO_ACCESS_KEY",
				LocalObjectReference: corev1.LocalObjectReference{
					Name: instance.GetNameWithSuffix("-airbyte-secrets"),
				},
			},
		},
	})

	envVars = append(envVars, corev1.EnvVar{
		Name: "AWS_SECRET_ACCESS_KEY",
		ValueFrom: &corev1.EnvVarSource{
			SecretKeyRef: &corev1.SecretKeySelector{
				Key: "STATE_STORAGE_MINIO_SECRET_ACCESS_KEY",
				LocalObjectReference: corev1.LocalObjectReference{
					Name: instance.GetNameWithSuffix("-airbyte-secrets"),
				},
			},
		},
	})

	for _, secretVarName := range secretVarNames {
		envVar := corev1.EnvVar{
			Name: secretVarName,
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					Key: secretVarName,
					LocalObjectReference: corev1.LocalObjectReference{
						Name: instance.GetNameWithSuffix("-airbyte-secrets"),
					},
				},
			},
		}
		envVars = append(envVars, envVar)
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.GetNameWithSuffix("-worker" + "-" + roleGroupName),
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
							Name:            "airbyte-worker-container",
							Image:           image.Repository + ":" + image.Tag,
							ImagePullPolicy: image.PullPolicy,
							Resources:       *roleGroup.Config.Resources,
							Env:             envVars,
							Ports:           containerPorts,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "gcs-log-creds-volume",
									MountPath: "/secrets/gcs-log-creds",
									ReadOnly:  true,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "gcs-log-creds-volume",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName:  instance.GetNameWithSuffix("-gcs-log-creds"),
									DefaultMode: func() *int32 { mode := int32(420); return &mode }(),
								},
							},
						},
					},
				},
			},
		},
	}
	WorkerScheduler(instance, dep)

	err := ctrl.SetControllerReference(instance, dep, schema)
	if err != nil {
		r.Log.Error(err, "Failed to set controller reference for deployment")
		return nil
	}
	return dep
}

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

// reconcile worker deployment
func (r *AirbyteReconciler) reconcileWorkerDeployment(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	roleGroups := convertRoleGroupToRoleConfigObject(instance, Worker)
	reconcileParams := r.createReconcileParams(ctx, roleGroups, instance, r.extractWorkerDeploymentForRoleGroup)
	if err := reconcileParams.createOrUpdateResource(); err != nil {
		return err
	}
	return nil
}

// extract worker deployment for role group
func (r *AirbyteReconciler) extractWorkerDeploymentForRoleGroup(params ExtractorParams) (client.Object, error) {
	instance := params.instance.(*stackv1alpha1.Airbyte)
	server := instance.Spec.ApiServer
	groupCfg := params.roleGroup
	roleCfg := server.RoleConfig
	clusterCfg := params.cluster
	roleGroupName := params.roleGroupName
	mergedLabels := r.mergeLabels(groupCfg, instance.GetLabels(), clusterCfg)
	schema := params.scheme

	image, securityContext, replicas, resources, containerPorts := getDeploymentInfo(groupCfg, roleCfg, clusterCfg)
	envVars := mergeEnvVarsForWorkerWebAppDeployment(instance, roleGroupName, &containerPorts)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      createNameForRoleGroup(instance, "worker", roleGroupName),
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
							Name:            "airbyte-worker-container",
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
	WorkerScheduler(instance, dep, groupCfg)

	err := ctrl.SetControllerReference(instance, dep, schema)
	if err != nil {
		r.Log.Error(err, "Failed to set controller reference for deployment")
		return nil, err
	}
	return dep, nil
}

// merge env var for worker webapp deployment
func mergeEnvVarsForWorkerWebAppDeployment(instance *stackv1alpha1.Airbyte, groupName string, containerPorts *[]corev1.ContainerPort) []corev1.EnvVar {
	portVarNames := []int32{9000, 9001, 9002, 9003, 9004, 9005, 9006, 9007, 9008, 9009, 9010,
		9011, 9012, 9013, 9014, 9015, 9016, 9017, 9018, 9019, 9020,
		9021, 9022, 9023, 9024, 9025, 9026, 9027, 9028, 9029, 9030}

	// Append a new ContainerPort to the containerPorts slice.
	// The new ContainerPort is configured with the name "heartbeat" and its port number is the first element from the portVarNames slice.
	*containerPorts = append(*containerPorts, corev1.ContainerPort{
		Name:          "heartbeat",
		ContainerPort: portVarNames[0],
	})

	// Iterate over the remaining elements in the portVarNames slice, starting from the second element.
	// For each portVarName, a new ContainerPort is appended to the containerPorts slice.
	// The new ContainerPort is configured with its port number set to the current portVarName.
	// Note: These additional ContainerPorts do not have a name.
	for _, portVarName := range portVarNames[1:] {
		*containerPorts = append(*containerPorts, corev1.ContainerPort{
			ContainerPort: portVarName,
		})
	}

	envVarNames := []string{"AIRBYTE_VERSION", "CONFIG_ROOT", "DATABASE_HOST", "DATABASE_PORT", "DATABASE_URL",
		"TRACKING_STRATEGY", "WORKSPACE_ROOT", "LOCAL_ROOT", "WEBAPP_URL", "TEMPORAL_HOST",
		"TEMPORAL_WORKER_PORTS", "S3_PATH_STYLE_ACCESS", "JOB_MAIN_CONTAINER_CPU_REQUEST", "JOB_MAIN_CONTAINER_CPU_LIMIT",
		"JOB_MAIN_CONTAINER_MEMORY_REQUEST", "JOB_MAIN_CONTAINER_MEMORY_LIMIT", "S3_LOG_BUCKET", "S3_LOG_BUCKET_REGION",
		"S3_ENDPOINT", "INTERNAL_API_HOST",
		"CONFIGS_DATABASE_MINIMUM_FLYWAY_MIGRATION_VERSION", "JOBS_DATABASE_MINIMUM_FLYWAY_MIGRATION_VERSION",
		"METRIC_CLIENT", "OTEL_COLLECTOR_ENDPOINT", "ACTIVITY_MAX_ATTEMPT", "ACTIVITY_INITIAL_DELAY_BETWEEN_ATTEMPTS_SECONDS",
		"ACTIVITY_MAX_DELAY_BETWEEN_ATTEMPTS_SECONDS", "WORKFLOW_FAILURE_RESTART_DELAY_SECONDS", "AUTO_DETECT_SCHEMA",
		"USE_STREAM_CAPABLE_STATE", "SHOULD_RUN_NOTIFY_WORKFLOWS", "WORKER_LOGS_STORAGE_TYPE", "WORKER_STATE_STORAGE_TYPE"}
	secretVarNames := []string{"DATABASE_PASSWORD", "DATABASE_USER", "AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY"}
	var envVars []corev1.EnvVar

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
		Name: "MICRONAUT_ENVIRONMENTS",
		ValueFrom: &corev1.EnvVarSource{
			ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
				Key: "WORKERS_MICRONAUT_ENVIRONMENTS",
				LocalObjectReference: corev1.LocalObjectReference{
					Name: createNameForRoleGroup(instance, "env", ""),
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

	for _, secretVarName := range secretVarNames {
		envVar := corev1.EnvVar{
			Name: secretVarName,
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					Key: secretVarName,
					LocalObjectReference: corev1.LocalObjectReference{
						Name: createNameForRoleGroup(instance, "secrets", ""),
					},
				},
			},
		}
		envVars = append(envVars, envVar)
	}
	return envVars
}

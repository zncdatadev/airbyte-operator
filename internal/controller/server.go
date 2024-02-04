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
)

// reconcile server service
func (r *AirbyteReconciler) reconcileServerService(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	roleGroups := convertRoleGroupToRoleConfigObject(instance, Server)
	reconcileParams := r.createReconcileParams(ctx, roleGroups, instance, r.extractServerServiceForRoleGroup)
	if err := reconcileParams.createOrUpdateResource(); err != nil {
		return err
	}
	return nil
}

// extract airbyte server service for role group
func (r *AirbyteReconciler) extractServerServiceForRoleGroup(params ExtractorParams) (client.Object, error) {
	instance := params.instance.(*stackv1alpha1.Airbyte)
	server := instance.Spec.Server
	groupCfg := params.roleGroup
	roleCfg := server.RoleConfig
	clusterCfg := params.cluster
	roleGroupName := params.roleGroupName
	mergedLabels := r.mergeLabels(groupCfg, instance.GetLabels(), clusterCfg)
	schema := params.scheme

	port, serviceType, annotations := getServiceInfo(groupCfg, roleCfg, clusterCfg)
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        createNameForRoleGroup(instance, "server", roleGroupName),
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

// reconcile server deployment
func (r *AirbyteReconciler) reconcileServerDeployment(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	roleGroups := convertRoleGroupToRoleConfigObject(instance, Server)
	reconcileParams := r.createReconcileParams(ctx, roleGroups, instance, r.extractServerDeploymentForRoleGroup)
	if err := reconcileParams.createOrUpdateResource(); err != nil {
		return err
	}
	return nil
}

// extract airbyte server deployment for role group
func (r *AirbyteReconciler) extractServerDeploymentForRoleGroup(params ExtractorParams) (client.Object, error) {
	instance := params.instance.(*stackv1alpha1.Airbyte)
	server := instance.Spec.Server
	groupCfg := params.roleGroup
	roleCfg := server.RoleConfig
	clusterCfg := params.cluster
	roleGroupName := params.roleGroupName
	mergedLabels := r.mergeLabels(groupCfg, instance.GetLabels(), clusterCfg)
	schema := params.scheme

	image, securityContext, replicas, resources, containerPorts := getDeploymentInfo(groupCfg, roleCfg, clusterCfg)
	envVars := mergeEnvVarsForServerDeployment(instance, roleGroupName)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      createNameForRoleGroup(instance, "server", roleGroupName),
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
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      createNameForRoleGroup(instance, "yml-volume", roleGroupName),
									MountPath: "/app/configs/airbyte.yml",
									SubPath:   "fileContents",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: createNameForRoleGroup(instance, "yml-volume", roleGroupName), // 第二个Volume的名称
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									DefaultMode: func() *int32 { mode := int32(420); return &mode }(),
									LocalObjectReference: corev1.LocalObjectReference{
										Name: createNameForRoleGroup(instance, "yml", ""),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	ServerScheduler(instance, dep, groupCfg)

	err := ctrl.SetControllerReference(instance, dep, schema)
	if err != nil {
		r.Log.Error(err, "Failed to set controller reference for deployment")
		return nil, err
	}
	return dep, nil
}

// merge env vars for airbyte server deployment
func mergeEnvVarsForServerDeployment(instance *stackv1alpha1.Airbyte, groupName string) []corev1.EnvVar {
	envVarNames := []string{"AIRBYTE_VERSION", "AIRBYTE_EDITION", "AUTO_DETECT_SCHEMA", "CONFIG_ROOT", "DATABASE_URL",
		"TRACKING_STRATEGY", "WORKER_ENVIRONMENT", "WORKSPACE_ROOT", "WEBAPP_URL", "TEMPORAL_HOST", "JOB_MAIN_CONTAINER_CPU_REQUEST",
		"JOB_MAIN_CONTAINER_CPU_LIMIT", "JOB_MAIN_CONTAINER_MEMORY_REQUEST", "JOB_MAIN_CONTAINER_MEMORY_LIMIT", "S3_LOG_BUCKET", "S3_LOG_BUCKET_REGION",
		"S3_ENDPOINT", "STATE_STORAGE_MINIO_BUCKET_NAME", "S3_PATH_STYLE_ACCESS",
		"CONFIGS_DATABASE_MINIMUM_FLYWAY_MIGRATION_VERSION",
		"JOBS_DATABASE_MINIMUM_FLYWAY_MIGRATION_VERSION", "WORKER_LOGS_STORAGE_TYPE", "WORKER_STATE_STORAGE_TYPE", "KEYCLOAK_INTERNAL_HOST"}
	secretVarNames := []string{"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "DATABASE_PASSWORD", "DATABASE_USER"}

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
		Name:  "LOG_LEVEL",
		Value: "INFO",
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

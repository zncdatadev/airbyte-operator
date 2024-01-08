package controller

import (
	stackv1alpha1 "github.com/zncdata-labs/airbyte-operator/api/v1alpha1"
	"github.com/zncdata-labs/operator-go/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *AirbyteReconciler) makeServerService(instance *stackv1alpha1.Airbyte) ([]*corev1.Service, error) {
	var services []*corev1.Service

	if instance.Spec.Server.RoleGroups != nil {
		for roleGroupName, roleGroup := range instance.Spec.Server.RoleGroups {
			svc, err := r.makeServerServiceForRoleGroup(instance, roleGroupName, roleGroup, r.Scheme)
			if err != nil {
				return nil, err
			}
			services = append(services, svc)
		}
	}

	return services, nil
}

func (r *AirbyteReconciler) makeServerServiceForRoleGroup(instance *stackv1alpha1.Airbyte, roleGroupName string, roleGroup *stackv1alpha1.RoleGroupServerSpec, schema *runtime.Scheme) (*corev1.Service, error) {
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
		port = instance.Spec.Server.Service.Port
		serviceType = instance.Spec.Server.Service.Type
		annotations = instance.Spec.Server.Service.Annotations
	}

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        instance.GetNameWithSuffix("-airbyte-server-svc" + "-" + roleGroupName),
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

func (r *AirbyteReconciler) makeCoordinatorDeployments(instance *stackv1alpha1.Airbyte) []*appsv1.Deployment {
	var deployments []*appsv1.Deployment

	if instance.Spec.Server.RoleGroups != nil {
		for roleGroupName, roleGroup := range instance.Spec.Server.RoleGroups {
			dep := r.makeServerDeploymentForRoleGroup(instance, roleGroupName, roleGroup, r.Scheme)
			if dep != nil {
				deployments = append(deployments, dep)
			}
		}
	}

	return deployments
}

func (r *AirbyteReconciler) makeServerDeploymentForRoleGroup(instance *stackv1alpha1.Airbyte, roleGroupName string, roleGroup *stackv1alpha1.RoleGroupServerSpec, schema *runtime.Scheme) *appsv1.Deployment {
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
	var containerPort int32

	if roleGroup != nil && roleGroup.Config != nil && roleGroup.Config.Image != nil {
		image = *roleGroup.Config.Image
	} else {
		image = *instance.Spec.Server.Image
	}

	if roleGroup != nil && roleGroup.Config != nil && roleGroup.Config.SecurityContext != nil {
		securityContext = roleGroup.Config.SecurityContext
	} else {
		securityContext = instance.Spec.SecurityContext
	}

	if roleGroup != nil && roleGroup.Config != nil && roleGroup.Config.Service != nil && roleGroup.Config.Service.Port != 0 {
		containerPort = roleGroup.Config.Service.Port
	} else {
		containerPort = instance.Spec.Server.Service.Port
	}

	envVarNames := []string{"AIRBYTE_VERSION", "AIRBYTE_EDITION", "AUTO_DETECT_SCHEMA", "CONFIG_ROOT", "DATABASE_URL",
		"TRACKING_STRATEGY", "WORKER_ENVIRONMENT", "WORKSPACE_ROOT", "WEBAPP_URL", "TEMPORAL_HOST", "JOB_MAIN_CONTAINER_CPU_REQUEST",
		"JOB_MAIN_CONTAINER_CPU_LIMIT", "JOB_MAIN_CONTAINER_MEMORY_REQUEST", "JOB_MAIN_CONTAINER_MEMORY_LIMIT", "S3_LOG_BUCKET", "S3_LOG_BUCKET_REGION",
		"S3_MINIO_ENDPOINT", "STATE_STORAGE_MINIO_BUCKET_NAME", "STATE_STORAGE_MINIO_ENDPOINT", "S3_PATH_STYLE_ACCESS", "GOOGLE_APPLICATION_CREDENTIALS",
		"GCS_LOG_BUCKET", "CONFIGS_DATABASE_MINIMUM_FLYWAY_MIGRATION_VERSION",
		"JOBS_DATABASE_MINIMUM_FLYWAY_MIGRATION_VERSION", "WORKER_LOGS_STORAGE_TYPE", "WORKER_STATE_STORAGE_TYPE", "KEYCLOAK_INTERNAL_HOST"}
	secretVarNames := []string{"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "DATABASE_PASSWORD", "DATABASE_USER", "STATE_STORAGE_MINIO_ACCESS_KEY"}
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
						Name: instance.GetNameWithSuffix("-airbyte-secrets"),
					},
				},
			},
		}
		envVars = append(envVars, envVar)
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.GetNameWithSuffix("-server" + "-" + roleGroupName),
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
					SecurityContext: securityContext,
					Containers: []corev1.Container{
						{
							Name:            instance.Name,
							Image:           image.Repository + ":" + image.Tag,
							ImagePullPolicy: image.PullPolicy,
							Resources:       *roleGroup.Config.Resources,
							Env:             envVars,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: containerPort,
									Name:          "http",
									Protocol:      "TCP",
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "gcs-log-creds-volume",
									MountPath: "/secrets/gcs-log-creds",
									ReadOnly:  true,
								},
								{
									Name:      instance.GetNameWithSuffix("-yml-volume"),
									MountPath: "/app/configs/airbyte.yml",
									SubPath:   "fileContents",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "gcs-log-creds-volume", // 第一个Volume的名称
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName:  instance.GetNameWithSuffix("-gcs-log-creds"),
									DefaultMode: func() *int32 { mode := int32(420); return &mode }(),
								},
							},
						},
						{
							Name: instance.GetNameWithSuffix("-yml-volume"), // 第二个Volume的名称
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									DefaultMode: func() *int32 { mode := int32(420); return &mode }(),
									LocalObjectReference: corev1.LocalObjectReference{
										Name: instance.GetNameWithSuffix("-airbyte-yml"),
									},
								},
							},
						},
					},
				},
			},
		},
	}
	ServerScheduler(instance, dep)

	err := ctrl.SetControllerReference(instance, dep, schema)
	if err != nil {
		r.Log.Error(err, "Failed to set controller reference for deployment")
		return nil
	}
	return dep
}

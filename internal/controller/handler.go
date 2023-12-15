package controller

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"

	stackv1alpha1 "github.com/zncdata-labs/airbyte-operator/api/v1alpha1"
)

// make service
func (r *AirbyteReconciler) makeServerService(instance *stackv1alpha1.Airbyte, schema *runtime.Scheme) *corev1.Service {
	labels := instance.GetLabels()
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        instance.Name,
			Namespace:   instance.Namespace,
			Labels:      labels,
			Annotations: instance.Spec.Server.Service.Annotations,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port:     instance.Spec.Server.Service.Port,
					Name:     "http",
					Protocol: "TCP",
				},
			},
			Selector: labels,
			Type:     instance.Spec.Server.Service.Type,
		},
	}
	err := ctrl.SetControllerReference(instance, svc, schema)
	if err != nil {
		r.Log.Error(err, "Failed to set controller reference for service")
		return nil
	}
	return svc
}

func (r *AirbyteReconciler) makeAirbyteApiServerService(instance *stackv1alpha1.Airbyte, schema *runtime.Scheme) *corev1.Service {
	labels := instance.GetLabels()
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        instance.GetNameWithSuffix("-airbyte-api-server-svc"),
			Namespace:   instance.Namespace,
			Labels:      labels,
			Annotations: instance.Spec.AirbyteApiServer.Service.Annotations,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port:     instance.Spec.AirbyteApiServer.Service.Port,
					Name:     "http",
					Protocol: "TCP",
				},
			},
			Selector: labels,
			Type:     instance.Spec.AirbyteApiServer.Service.Type,
		},
	}
	err := ctrl.SetControllerReference(instance, svc, schema)
	if err != nil {
		r.Log.Error(err, "Failed to set controller reference for service")
		return nil
	}
	return svc
}

func (r *AirbyteReconciler) makeConnectorBuilderServerService(instance *stackv1alpha1.Airbyte, schema *runtime.Scheme) *corev1.Service {
	labels := instance.GetLabels()
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        instance.GetNameWithSuffix("-airbyte-connector-builder-server-svc"),
			Namespace:   instance.Namespace,
			Labels:      labels,
			Annotations: instance.Spec.ConnectorBuilderServer.Service.Annotations,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port:     instance.Spec.ConnectorBuilderServer.Service.Port,
					Name:     "http",
					Protocol: "TCP",
				},
			},
			Selector: labels,
			Type:     instance.Spec.ConnectorBuilderServer.Service.Type,
		},
	}
	err := ctrl.SetControllerReference(instance, svc, schema)
	if err != nil {
		r.Log.Error(err, "Failed to set controller reference for service")
		return nil
	}
	return svc
}

func (r *AirbyteReconciler) makeTemporalService(instance *stackv1alpha1.Airbyte, schema *runtime.Scheme) *corev1.Service {
	labels := instance.GetLabels()
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        instance.GetNameWithSuffix("-temporal"),
			Namespace:   instance.Namespace,
			Labels:      labels,
			Annotations: instance.Spec.Temporal.Service.Annotations,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port:     instance.Spec.Temporal.Service.Port,
					Name:     "http",
					Protocol: "TCP",
				},
			},
			Selector: labels,
			Type:     instance.Spec.Temporal.Service.Type,
		},
	}
	err := ctrl.SetControllerReference(instance, svc, schema)
	if err != nil {
		r.Log.Error(err, "Failed to set controller reference for service")
		return nil
	}
	return svc
}

func (r *AirbyteReconciler) makeWebAppService(instance *stackv1alpha1.Airbyte, schema *runtime.Scheme) *corev1.Service {
	labels := instance.GetLabels()
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        instance.GetNameWithSuffix("-airbyte-webapp-svc"),
			Namespace:   instance.Namespace,
			Labels:      labels,
			Annotations: instance.Spec.WebApp.Service.Annotations,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port:     instance.Spec.WebApp.Service.Port,
					Name:     "http",
					Protocol: "TCP",
				},
			},
			Selector: labels,
			Type:     instance.Spec.WebApp.Service.Type,
		},
	}
	err := ctrl.SetControllerReference(instance, svc, schema)
	if err != nil {
		r.Log.Error(err, "Failed to set controller reference for service")
		return nil
	}
	return svc
}

func (r *AirbyteReconciler) reconcileService(ctx context.Context, instance *stackv1alpha1.Airbyte) error {

	ServerService := r.makeServerService(instance, r.Scheme)
	if ServerService == nil {
		return nil
	}

	AirbyteApiServerService := r.makeAirbyteApiServerService(instance, r.Scheme)
	if AirbyteApiServerService == nil {
		return nil
	}

	ConnectorBuilderServerService := r.makeConnectorBuilderServerService(instance, r.Scheme)
	if ConnectorBuilderServerService == nil {
		return nil
	}

	TemporalService := r.makeTemporalService(instance, r.Scheme)
	if TemporalService == nil {
		return nil
	}

	WebAppService := r.makeWebAppService(instance, r.Scheme)
	if WebAppService == nil {
		return nil
	}

	serverService := r.makeServerService(instance, r.Scheme)
	if serverService != nil {
		if err := CreateOrUpdate(ctx, r.Client, serverService); err != nil {
			r.Log.Error(err, "Failed to create or update service")
			return err
		}
	}

	airbyteApiServerService := r.makeAirbyteApiServerService(instance, r.Scheme)
	if airbyteApiServerService != nil {
		if err := CreateOrUpdate(ctx, r.Client, airbyteApiServerService); err != nil {
			r.Log.Error(err, "Failed to create or update airbyteApiServer service")
			return err
		}
	}

	connectorBuilderServerService := r.makeConnectorBuilderServerService(instance, r.Scheme)
	if connectorBuilderServerService != nil {
		if err := CreateOrUpdate(ctx, r.Client, connectorBuilderServerService); err != nil {
			r.Log.Error(err, "Failed to create or update connectorBuilderServerService service")
			return err
		}
	}

	temporalService := r.makeTemporalService(instance, r.Scheme)
	if temporalService != nil {
		if err := CreateOrUpdate(ctx, r.Client, temporalService); err != nil {
			r.Log.Error(err, "Failed to create or update temporalService service")
			return err
		}
	}

	webAppService := r.makeWebAppService(instance, r.Scheme)
	if webAppService != nil {
		if err := CreateOrUpdate(ctx, r.Client, webAppService); err != nil {
			r.Log.Error(err, "Failed to create or update webAppService service")
			return err
		}
	}
	return nil
}

func (r *AirbyteReconciler) makeServerDeployment(instance *stackv1alpha1.Airbyte, schema *runtime.Scheme) *appsv1.Deployment {
	labels := instance.GetLabels()
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
			Name:      instance.GetNameWithSuffix("-server"),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &instance.Spec.Server.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					SecurityContext: instance.Spec.SecurityContext,
					Containers: []corev1.Container{
						{
							Name:            instance.Name,
							Image:           instance.Spec.Server.Image.Repository + ":" + instance.Spec.Server.Image.Tag,
							ImagePullPolicy: instance.Spec.Server.Image.PullPolicy,
							Resources:       *instance.Spec.Server.Resources,
							Env:             envVars,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: instance.Spec.Server.Service.Port,
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

func (r *AirbyteReconciler) makeWorkerDeployment(instance *stackv1alpha1.Airbyte, schema *runtime.Scheme) *appsv1.Deployment {
	labels := instance.GetLabels()

	portVarNames := []int32{9000, 9001, 9002, 9003, 9004, 9005, 9006,
		9007, 9008, 9010, 9011, 9012, 9013, 9014, 9015, 9016, 9020,
		9021, 9022, 9023, 9024, 9025, 9026, 9027, 9028, 9030}

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
				Key: "CRON_MICRONAUT_ENVIRONMENTS",
				LocalObjectReference: corev1.LocalObjectReference{
					Name: instance.GetNameWithSuffix("-airbyte-env"),
				},
			},
		},
	})

	envVars = append(envVars, corev1.EnvVar{
		Name:  "WORKSPACE_DOCKER_MOUNT",
		Value: "workspace",
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
			Name:      instance.GetNameWithSuffix("-worker"),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &instance.Spec.Worker.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					//ServiceAccountName: instance.GetNameWithSuffix("-admin"),
					SecurityContext: instance.Spec.SecurityContext,
					Containers: []corev1.Container{
						{
							Name:            instance.Name,
							Image:           instance.Spec.Worker.Image.Repository + ":" + instance.Spec.Worker.Image.Tag,
							ImagePullPolicy: instance.Spec.Worker.Image.PullPolicy,
							Resources:       *instance.Spec.Worker.Resources,
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

func (r *AirbyteReconciler) makeAirbyteApiServerDeployment(instance *stackv1alpha1.Airbyte, schema *runtime.Scheme) *appsv1.Deployment {
	labels := instance.GetLabels()

	envVarNames := []string{"AIRBYTE_API_HOST", "INTERNAL_API_HOST"}
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

	if instance != nil && instance.Spec.AirbyteApiServer != nil && instance.Spec.AirbyteApiServer.Debug != nil {
		if instance.Spec.AirbyteApiServer.Debug.Enabled {
			envVars = append(envVars, corev1.EnvVar{
				Name:  "JAVA_TOOL_OPTIONS",
				Value: "-Xdebug -agentlib:jdwp=transport=dt_socket,address=0.0.0.0:" + strconv.FormatInt(int64(instance.Spec.AirbyteApiServer.Debug.RemoteDebugPort), 10) + ",server=y,suspend=n",
			})
		}
	}

	if instance != nil && instance.Spec.Global != nil {
		if instance.Spec.Global.DeploymentMode == "oss" {
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

	if instance != nil && instance.Spec.AirbyteApiServer != nil && instance.Spec.AirbyteApiServer.Secret != nil {
		for key := range instance.Spec.AirbyteApiServer.Secret {
			envVars = append(envVars, corev1.EnvVar{
				Name: key,
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "server-secrets",
						},
						Key: key,
					},
				},
			})
		}
	}

	if instance != nil && (instance.Spec.AirbyteApiServer != nil || instance.Spec.Global != nil) {
		envVarsMap := make(map[string]string)

		if instance.Spec.AirbyteApiServer != nil && instance.Spec.AirbyteApiServer.EnvVars != nil {
			for key, value := range instance.Spec.AirbyteApiServer.EnvVars {
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

	if instance != nil && instance.Spec.AirbyteApiServer != nil && instance.Spec.AirbyteApiServer.ExtraEnv != (corev1.EnvVar{}) {
		envVars = append(envVars, instance.Spec.AirbyteApiServer.ExtraEnv)
	}

	containerPorts := []corev1.ContainerPort{
		{
			ContainerPort: 8006,
			Name:          "http",
			Protocol:      corev1.ProtocolTCP,
		},
	}

	if instance != nil && instance.Spec.AirbyteApiServer != nil && instance.Spec.AirbyteApiServer.Debug != nil && instance.Spec.AirbyteApiServer.Debug.Enabled {
		containerPorts = append(containerPorts, corev1.ContainerPort{
			ContainerPort: instance.Spec.AirbyteApiServer.Debug.RemoteDebugPort,
			Name:          "debug",
			Protocol:      corev1.ProtocolTCP,
		})
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.GetNameWithSuffix("-airbyte-api-server"),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &instance.Spec.AirbyteApiServer.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					//ServiceAccountName: instance.GetNameWithSuffix("-admin"),
					SecurityContext: instance.Spec.SecurityContext,
					Containers: []corev1.Container{
						{
							Name:            instance.Name,
							Image:           instance.Spec.AirbyteApiServer.Image.Repository + ":" + instance.Spec.AirbyteApiServer.Image.Tag,
							ImagePullPolicy: instance.Spec.AirbyteApiServer.Image.PullPolicy,
							Resources:       *instance.Spec.AirbyteApiServer.Resources,
							Env:             envVars,
							Ports:           containerPorts,
						},
					},
				},
			},
		},
	}
	AirbyteApiServerScheduler(instance, dep)

	err := ctrl.SetControllerReference(instance, dep, schema)
	if err != nil {
		r.Log.Error(err, "Failed to set controller reference for deployment")
		return nil
	}
	return dep
}

func (r *AirbyteReconciler) makeConnectorBuilderServerDeployment(instance *stackv1alpha1.Airbyte, schema *runtime.Scheme) *appsv1.Deployment {
	labels := instance.GetLabels()

	var envVars []corev1.EnvVar

	if instance != nil && instance.Spec.ConnectorBuilderServer != nil && instance.Spec.ConnectorBuilderServer.Debug != nil {
		if instance.Spec.ConnectorBuilderServer.Debug.Enabled {
			envVars = append(envVars, corev1.EnvVar{
				Name:  "JAVA_TOOL_OPTIONS",
				Value: "-Xdebug -agentlib:jdwp=transport=dt_socket,address=0.0.0.0:" + strconv.FormatInt(int64(instance.Spec.ConnectorBuilderServer.Debug.RemoteDebugPort), 10) + ",server=y,suspend=n",
			})
		}
	}

	if instance != nil && instance.Spec.Global != nil {
		if instance.Spec.Global.DeploymentMode == "oss" {
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

	if instance != nil && instance.Spec.ConnectorBuilderServer != nil && instance.Spec.ConnectorBuilderServer.Secret != nil {
		for key := range instance.Spec.ConnectorBuilderServer.Secret {
			envVars = append(envVars, corev1.EnvVar{
				Name: key,
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "server-secrets",
						},
						Key: key,
					},
				},
			})
		}
	}

	if instance != nil && (instance.Spec.ConnectorBuilderServer != nil || instance.Spec.Global != nil) {
		envVarsMap := make(map[string]string)

		if instance.Spec.ConnectorBuilderServer != nil && instance.Spec.ConnectorBuilderServer.EnvVars != nil {
			for key, value := range instance.Spec.ConnectorBuilderServer.EnvVars {
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

	if instance != nil && instance.Spec.ConnectorBuilderServer != nil && instance.Spec.ConnectorBuilderServer.ExtraEnv != (corev1.EnvVar{}) {
		envVars = append(envVars, instance.Spec.ConnectorBuilderServer.ExtraEnv)
	}

	containerPorts := []corev1.ContainerPort{
		{
			ContainerPort: 8006,
			Name:          "http",
			Protocol:      corev1.ProtocolTCP,
		},
	}

	if instance != nil && instance.Spec.ConnectorBuilderServer != nil && instance.Spec.ConnectorBuilderServer.Debug != nil && instance.Spec.ConnectorBuilderServer.Debug.Enabled {
		containerPorts = append(containerPorts, corev1.ContainerPort{
			ContainerPort: instance.Spec.ConnectorBuilderServer.Debug.RemoteDebugPort,
			Name:          "debug",
			Protocol:      corev1.ProtocolTCP,
		})
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.GetNameWithSuffix("-connector-builder-server"),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &instance.Spec.ConnectorBuilderServer.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					//ServiceAccountName: instance.GetNameWithSuffix("-admin"),
					SecurityContext: instance.Spec.SecurityContext,
					Containers: []corev1.Container{
						{
							Name:            instance.Name,
							Image:           instance.Spec.ConnectorBuilderServer.Image.Repository + ":" + instance.Spec.ConnectorBuilderServer.Image.Tag,
							ImagePullPolicy: instance.Spec.ConnectorBuilderServer.Image.PullPolicy,
							Resources:       *instance.Spec.ConnectorBuilderServer.Resources,
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
					//ServiceAccountName: instance.GetNameWithSuffix("-admin"),
					SecurityContext: instance.Spec.SecurityContext,
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

func (r *AirbyteReconciler) makePodSweeperDeployment(instance *stackv1alpha1.Airbyte, schema *runtime.Scheme) *appsv1.Deployment {
	labels := instance.GetLabels()

	var envVars []corev1.EnvVar

	envVars = append(envVars, corev1.EnvVar{
		Name: "KUBE_NAMESPACE",
		ValueFrom: &corev1.EnvVarSource{
			FieldRef: &corev1.ObjectFieldSelector{
				FieldPath: "metadata.namespace",
			},
		},
	})

	if instance != nil && instance.Spec.PodSweeper != nil && instance.Spec.PodSweeper.TimeToDeletePods != nil {
		envVars = append(envVars, corev1.EnvVar{
			Name:  "RUNNING_TTL_MINUTES",
			Value: instance.Spec.PodSweeper.TimeToDeletePods.Running,
		})
		envVars = append(envVars, corev1.EnvVar{
			Name:  "SUCCEEDED_TTL_MINUTES",
			Value: strconv.Itoa(int(instance.Spec.PodSweeper.TimeToDeletePods.Succeeded)),
		})
		envVars = append(envVars, corev1.EnvVar{
			Name:  "UNSUCCESSFUL_TTL_MINUTES",
			Value: strconv.Itoa(int(instance.Spec.PodSweeper.TimeToDeletePods.Unsuccessful)),
		})
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.GetNameWithSuffix("-pod-sweeper"),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &instance.Spec.PodSweeper.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					//ServiceAccountName: instance.GetNameWithSuffix("-admin"),
					SecurityContext: instance.Spec.SecurityContext,
					Containers: []corev1.Container{
						{
							Name:            "airbyte-pod-sweeper",
							Image:           instance.Spec.PodSweeper.Image.Repository + ":" + instance.Spec.PodSweeper.Image.Tag,
							ImagePullPolicy: instance.Spec.PodSweeper.Image.PullPolicy,
							Resources:       *instance.Spec.PodSweeper.Resources,
							Env:             envVars,
							Command:         []string{"/bin/bash", "-c", "/script/sweep-pod.sh"},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "sweep-pod-script",
									MountPath: "/script/sweep-pod.sh",
									SubPath:   "sweep-pod.sh",
								},
								{
									Name:      "kube-config",
									MountPath: "/.kube",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "kube-config",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
						{
							Name: "sweep-pod-script",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: instance.GetNameWithSuffix("-sweep-pod-script"),
									},
									DefaultMode: func() *int32 { mode := int32(0755); return &mode }(),
								},
							},
						},
					},
				},
			},
		},
	}

	CronScheduler(instance, dep)

	if instance.Spec.PodSweeper.LivenessProbe != nil && instance.Spec.PodSweeper.LivenessProbe.Enabled {
		dep.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				Exec: &corev1.ExecAction{
					Command: []string{
						"/bin/sh",
						"-c",
						"grep -q sweep-pod.sh /proc/1/cmdline",
					},
				},
			},
			InitialDelaySeconds: instance.Spec.PodSweeper.LivenessProbe.InitialDelaySeconds,
			PeriodSeconds:       instance.Spec.PodSweeper.LivenessProbe.PeriodSeconds,
			TimeoutSeconds:      instance.Spec.PodSweeper.LivenessProbe.TimeoutSeconds,
			SuccessThreshold:    instance.Spec.PodSweeper.LivenessProbe.SuccessThreshold,
			FailureThreshold:    instance.Spec.PodSweeper.LivenessProbe.FailureThreshold,
		}
	}

	if instance.Spec.PodSweeper.ReadinessProbe != nil && instance.Spec.PodSweeper.ReadinessProbe.Enabled {
		dep.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				Exec: &corev1.ExecAction{
					Command: []string{
						"/bin/sh",
						"-c",
						"grep -q sweep-pod.sh /proc/1/cmdline",
					},
				},
			},
			InitialDelaySeconds: instance.Spec.PodSweeper.ReadinessProbe.InitialDelaySeconds,
			PeriodSeconds:       instance.Spec.PodSweeper.ReadinessProbe.PeriodSeconds,
			TimeoutSeconds:      instance.Spec.PodSweeper.ReadinessProbe.TimeoutSeconds,
			SuccessThreshold:    instance.Spec.PodSweeper.ReadinessProbe.SuccessThreshold,
			FailureThreshold:    instance.Spec.PodSweeper.ReadinessProbe.FailureThreshold,
		}
	}

	err := ctrl.SetControllerReference(instance, dep, schema)

	if err != nil {
		r.Log.Error(err, "Failed to set controller reference for deployment")
		return nil
	}
	return dep
}

func (r *AirbyteReconciler) makeTemporalDeployment(instance *stackv1alpha1.Airbyte, schema *runtime.Scheme) *appsv1.Deployment {
	labels := instance.GetLabels()

	var envVars []corev1.EnvVar

	if instance != nil && instance.Spec.Global != nil {
		if instance.Spec.Global.DeploymentMode == "oss" {

			envVars = append(envVars, corev1.EnvVar{
				Name:  "AUTO_SETUP",
				Value: "true",
			}, corev1.EnvVar{
				Name:  "DB",
				Value: "postgresql",
			})

			if instance.Spec.Global.Database != nil {
				if instance.Spec.Global.Database.Port != 0 {
					envVars = append(envVars, corev1.EnvVar{
						Name:  "DB_PORT",
						Value: strconv.Itoa(int(instance.Spec.Global.Database.Port)),
					})
				} else if instance.Spec.Global.ConfigMapName != "" {
					envVars = append(envVars, corev1.EnvVar{
						Name: "DB_PORT",
						ValueFrom: &corev1.EnvVarSource{
							ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: instance.Spec.Global.ConfigMapName,
								},
								Key: "DATABASE_PORT",
							},
						},
					})
				} else {
					envVars = append(envVars, corev1.EnvVar{
						Name: "DB_PORT",
						ValueFrom: &corev1.EnvVarSource{
							ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: instance.GetNameWithSuffix("-airbyte-env"),
								},
								Key: "DATABASE_PORT",
							},
						},
					})
				}
			}

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

			if instance.Spec.Global.Database != nil {
				secretName := instance.GetNameWithSuffix("-airbyte-secrets")
				secretKey := "DATABASE_PASSWORD"

				if instance.Spec.Global.Database.SecretName != "" {
					secretName = instance.Spec.Global.Database.SecretName
				}
				if instance.Spec.Global.Database.SecretValue != "" {
					secretKey = instance.Spec.Global.Database.SecretValue
				}

				envVars = append(envVars, corev1.EnvVar{
					Name: "POSTGRES_PWD",
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							Key: secretKey,
							LocalObjectReference: corev1.LocalObjectReference{
								Name: secretName,
							},
						},
					},
				})
			}

			configMapName := instance.GetNameWithSuffix("-airbyte-env")
			if instance.Spec.Global.ConfigMapName != "" {
				configMapName = instance.Spec.Global.ConfigMapName
			}

			envVars = append(envVars, corev1.EnvVar{
				Name: "POSTGRES_SEEDS",
				ValueFrom: &corev1.EnvVarSource{
					ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
						Key: "DATABASE_HOST",
						LocalObjectReference: corev1.LocalObjectReference{
							Name: configMapName,
						},
					},
				},
			})

			envVars = append(envVars, corev1.EnvVar{
				Name:  "DYNAMIC_CONFIG_FILE_PATH",
				Value: "config/dynamicconfig/development.yaml",
			})
		}
	}

	if instance != nil && instance.Spec.Temporal != nil && instance.Spec.Temporal.Secret != nil {
		for key := range instance.Spec.Temporal.Secret {
			envVars = append(envVars, corev1.EnvVar{
				Name: key,
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "temporal-secrets",
						},
						Key: key,
					},
				},
			})
		}
	}

	if instance != nil && (instance.Spec.Temporal != nil || instance.Spec.Global != nil) {
		envVarsMap := make(map[string]string)

		if instance.Spec.Temporal != nil && instance.Spec.Temporal.EnvVars != nil {
			for key, value := range instance.Spec.Temporal.EnvVars {
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

	if instance != nil && instance.Spec.Temporal != nil && instance.Spec.Temporal.ExtraEnv != (corev1.EnvVar{}) {
		envVars = append(envVars, instance.Spec.Temporal.ExtraEnv)
	}

	containerPorts := []corev1.ContainerPort{
		{
			ContainerPort: 7233,
		},
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.GetNameWithSuffix("-temporal"),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &instance.Spec.Temporal.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					//ServiceAccountName: instance.GetNameWithSuffix("-admin"),
					SecurityContext: instance.Spec.SecurityContext,
					Containers: []corev1.Container{
						{
							Name:            instance.Name,
							Image:           instance.Spec.Temporal.Image.Repository + ":" + instance.Spec.Temporal.Image.Tag,
							ImagePullPolicy: instance.Spec.Temporal.Image.PullPolicy,
							Resources:       *instance.Spec.Temporal.Resources,
							Env:             envVars,
							Ports:           containerPorts,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "airbyte-temporal-dynamicconfig",
									MountPath: "/config/dynamicconfig",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "airbyte-temporal-dynamicconfig",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: instance.GetNameWithSuffix("-tempora-dynamicconfig"),
									},
									Items: []corev1.KeyToPath{
										{
											Key:  "development.yaml",
											Path: "development.yaml",
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
	ConnectorBuilderServerScheduler(instance, dep)

	err := ctrl.SetControllerReference(instance, dep, schema)
	if err != nil {
		r.Log.Error(err, "Failed to set controller reference for deployment")
		return nil
	}
	return dep
}

func (r *AirbyteReconciler) makeWebAppDeployment(instance *stackv1alpha1.Airbyte, schema *runtime.Scheme) *appsv1.Deployment {
	labels := instance.GetLabels()

	envVarNames := []string{"AIRBYTE_API_HOST", "INTERNAL_API_HOST"}
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
		}
	}

	if instance != nil && instance.Spec.WebApp != nil && instance.Spec.WebApp.Secret != nil {
		for key := range instance.Spec.WebApp.Secret {
			envVars = append(envVars, corev1.EnvVar{
				Name: key,
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "webapp-secrets",
						},
						Key: key,
					},
				},
			})
		}
	}

	if instance != nil && (instance.Spec.WebApp != nil || instance.Spec.Global != nil) {
		envVarsMap := make(map[string]string)

		if instance.Spec.WebApp != nil && instance.Spec.WebApp.EnvVars != nil {
			for key, value := range instance.Spec.WebApp.EnvVars {
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

	if instance != nil && instance.Spec.WebApp != nil && instance.Spec.WebApp.ExtraEnv != (corev1.EnvVar{}) {
		envVars = append(envVars, instance.Spec.WebApp.ExtraEnv)
	}

	containerPorts := []corev1.ContainerPort{
		{
			ContainerPort: 80,
			Name:          "http",
			Protocol:      corev1.ProtocolTCP,
		},
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.GetNameWithSuffix("-webapp"),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &instance.Spec.WebApp.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					//ServiceAccountName: instance.GetNameWithSuffix("-admin"),
					SecurityContext: instance.Spec.SecurityContext,
					Containers: []corev1.Container{
						{
							Name:            instance.Name,
							Image:           instance.Spec.WebApp.Image.Repository + ":" + instance.Spec.WebApp.Image.Tag,
							ImagePullPolicy: instance.Spec.WebApp.Image.PullPolicy,
							Resources:       *instance.Spec.WebApp.Resources,
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

func (r *AirbyteReconciler) reconcileDeployment(ctx context.Context, instance *stackv1alpha1.Airbyte) error {

	ServerDeployment := r.makeServerDeployment(instance, r.Scheme)
	if ServerDeployment == nil {
		return nil
	}

	WorkerDeployment := r.makeWorkerDeployment(instance, r.Scheme)
	if WorkerDeployment == nil {
		return nil
	}

	AirbyteApiServerDeployment := r.makeAirbyteApiServerDeployment(instance, r.Scheme)
	if AirbyteApiServerDeployment == nil {
		return nil
	}

	ConnectorBuilderServerDeployment := r.makeConnectorBuilderServerDeployment(instance, r.Scheme)
	if ConnectorBuilderServerDeployment == nil {
		return nil
	}

	CronDeployment := r.makeCronDeployment(instance, r.Scheme)
	if CronDeployment == nil {
		return nil
	}

	PodSweeperDeployment := r.makePodSweeperDeployment(instance, r.Scheme)
	if PodSweeperDeployment == nil {
		return nil
	}

	TemporalDeployment := r.makeTemporalDeployment(instance, r.Scheme)
	if TemporalDeployment == nil {
		return nil
	}

	WebAppDeployment := r.makeWebAppDeployment(instance, r.Scheme)
	if WebAppDeployment == nil {
		return nil
	}

	serverDeployment := r.makeServerDeployment(instance, r.Scheme)
	if serverDeployment != nil {
		if err := CreateOrUpdate(ctx, r.Client, serverDeployment); err != nil {
			r.Log.Error(err, "Failed to create or update Serverdeployment")
			return err
		}
	}

	workerDeployment := r.makeWorkerDeployment(instance, r.Scheme)
	if workerDeployment != nil {
		if err := CreateOrUpdate(ctx, r.Client, workerDeployment); err != nil {
			r.Log.Error(err, "Failed to create or update Workerdeployment")
			return err
		}
	}

	airbyteApiServerDeployment := r.makeAirbyteApiServerDeployment(instance, r.Scheme)
	if airbyteApiServerDeployment != nil {
		if err := CreateOrUpdate(ctx, r.Client, airbyteApiServerDeployment); err != nil {
			r.Log.Error(err, "Failed to create or update AirbyteApiServer")
			return err
		}
	}

	connectorBuilderServerDeployment := r.makeConnectorBuilderServerDeployment(instance, r.Scheme)
	if connectorBuilderServerDeployment != nil {
		if err := CreateOrUpdate(ctx, r.Client, connectorBuilderServerDeployment); err != nil {
			r.Log.Error(err, "Failed to create or update ConnectorBuilderServerDeployment")
			return err
		}
	}

	cronDeployment := r.makeCronDeployment(instance, r.Scheme)
	if cronDeployment != nil {
		if err := CreateOrUpdate(ctx, r.Client, cronDeployment); err != nil {
			r.Log.Error(err, "Failed to create or update cronDeployment")
			return err
		}
	}

	podSweeperDeployment := r.makePodSweeperDeployment(instance, r.Scheme)
	if podSweeperDeployment != nil {
		if err := CreateOrUpdate(ctx, r.Client, podSweeperDeployment); err != nil {
			r.Log.Error(err, "Failed to create or update podSweeperDeployment")
			return err
		}
	}

	temporalDeployment := r.makeTemporalDeployment(instance, r.Scheme)
	if temporalDeployment != nil {
		if err := CreateOrUpdate(ctx, r.Client, temporalDeployment); err != nil {
			r.Log.Error(err, "Failed to create or update temporalDeployment")
			return err
		}
	}

	webAppDeployment := r.makeWebAppDeployment(instance, r.Scheme)
	if webAppDeployment != nil {
		if err := CreateOrUpdate(ctx, r.Client, webAppDeployment); err != nil {
			r.Log.Error(err, "Failed to create or update webAppDeployment")
			return err
		}
	}
	return nil
}

func (r *AirbyteReconciler) makeSecret(instance *stackv1alpha1.Airbyte) []*corev1.Secret {
	var secrets []*corev1.Secret
	labels := instance.GetLabels()
	if instance.Spec.Secret.AirbyteSecrets != nil {
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      instance.GetNameWithSuffix("-gcs-log-creds"),
				Namespace: instance.Namespace,
				Labels:    labels,
			},
			Type: corev1.SecretTypeOpaque,
			Data: map[string][]byte{
				"gcp.json": []byte(instance.Spec.Secret.GcsLogCreds["gcp.json"]),
			},
		}
		secrets = append(secrets, secret)

	}

	if instance.Spec.Secret.AirbyteSecrets != nil {
		data := make(map[string][]byte)
		data["AWS_ACCESS_KEY_ID"] = []byte(instance.Spec.Secret.AirbyteSecrets["AWS_ACCESS_KEY_ID"])
		data["AWS_SECRET_ACCESS_KEY"] = []byte(instance.Spec.Secret.AirbyteSecrets["AWS_SECRET_ACCESS_KEY"])
		data["DATABASE_PASSWORD"] = []byte(instance.Spec.Secret.AirbyteSecrets["DATABASE_PASSWORD"])
		data["DATABASE_USER"] = []byte(instance.Spec.Secret.AirbyteSecrets["DATABASE_USER"])
		data["STATE_STORAGE_MINIO_ACCESS_KEY"] = []byte(instance.Spec.Secret.AirbyteSecrets["STATE_STORAGE_MINIO_ACCESS_KEY"])
		data["STATE_STORAGE_MINIO_SECRET_ACCESS_KEY"] = []byte(instance.Spec.Secret.AirbyteSecrets["STATE_STORAGE_MINIO_SECRET_ACCESS_KEY"])

		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      instance.GetNameWithSuffix("-airbyte-secrets"),
				Namespace: instance.Namespace,
				Labels:    labels,
			},
			Type: corev1.SecretTypeOpaque,
			Data: data,
		}
		secrets = append(secrets, secret)

	}

	if instance.Spec.Temporal.Secret != nil {
		// Create a Data map to hold the secret data
		Data := make(map[string][]byte)

		// Iterate over instance.Spec.Temporal.Secret and create Secret data items
		for k, v := range instance.Spec.Temporal.Secret {
			value := ""
			if v != "" {
				value = base64.StdEncoding.EncodeToString([]byte(v))
			}
			Data[k] = []byte(value)
		}
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "temporal-secrets",
				Namespace: instance.Namespace,
				Labels:    labels,
			},
			Type: corev1.SecretTypeOpaque,
			Data: Data,
		}
		secrets = append(secrets, secret)

	}

	return secrets
}

func (r *AirbyteReconciler) reconcileSecret(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	logger := log.FromContext(ctx)
	objs := r.makeSecret(instance)

	if objs == nil {
		return nil
	}
	for _, obj := range objs {
		if err := CreateOrUpdate(ctx, r.Client, obj); err != nil {
			logger.Error(err, "Failed to create or update secret")
			return err
		}
	}
	// Create or update the secret

	return nil
}

func (r *AirbyteReconciler) makeYmlConfigMap(instance *stackv1alpha1.Airbyte, schema *runtime.Scheme) *corev1.ConfigMap {
	labels := instance.GetLabels()
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.GetNameWithSuffix("-airbyte-yml"),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Data: map[string]string{
			"fileContents": "",
		},
	}
	err := ctrl.SetControllerReference(instance, configMap, schema)
	if err != nil {
		r.Log.Error(err, "Failed to set controller reference for configmap")
		return nil
	}
	return configMap
}

func (r *AirbyteReconciler) makeSweepPodScriptConfigMap(instance *stackv1alpha1.Airbyte, schema *runtime.Scheme) *corev1.ConfigMap {
	labels := instance.GetLabels()
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.GetNameWithSuffix("-sweep-pod-script"),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Data: map[string]string{
			"sweep-pod.sh": `#!/bin/bash
get_job_pods () {
    kubectl -n ${KUBE_NAMESPACE} -L airbyte -l airbyte=job-pod \
      get pods \
      -o=jsonpath='{range .items[*]} {.metadata.name} {.status.phase} {.status.conditions[0].lastTransitionTime} {.status.startTime}{"\n"}{end}'
}
delete_pod() {
    printf "From status '%s' since '%s', " $2 $3
    echo "$1" | grep -v "STATUS" | awk '{print $1}' | xargs --no-run-if-empty kubectl -n ${KUBE_NAMESPACE} delete pod
}
while :
do
    echo "Starting pod sweeper cycle:"

    if [ -n "$${RUNNING_TTL_MINUTES}" ]; then 
      # Time window for running pods
      RUNNING_DATE_STR=` + "`date -d \"now - $${RUNNING_TTL_MINUTES} minutes\" --utc -Ins`" + `
      RUNNING_DATE=` + "`date -d $$RUNNING_DATE_STR +%s`" + `
      echo "Will sweep running pods from before $${RUNNING_DATE_STR}"
    fi

    if [ -n "${SUCCEEDED_TTL_MINUTES}" ]; then
      # Shorter time window for succeeded pods
      SUCCESS_DATE_STR=` + "`date -d \"now - $${SUCCEEDED_TTL_MINUTES} minutes\" --utc -Ins`" + `
      SUCCESS_DATE=` + "`date -d $$SUCCESS_DATE_STR +%s`" + `
      echo "Will sweep succeeded pods from before ${SUCCESS_DATE_STR}"
    fi

    if [ -n "${UNSUCCESSFUL_TTL_MINUTES}" ]; then
      # Longer time window for unsuccessful pods (to debug)
      NON_SUCCESS_DATE_STR=` + "`date -d \"now - $${UNSUCCESSFUL_TTL_MINUTES} minutes\" --utc -Ins`" + `
      NON_SUCCESS_DATE=` + "`date -d $$NON_SUCCESS_DATE_STR +%s`" + `
      echo "Will sweep unsuccessful pods from before ${NON_SUCCESS_DATE_STR}"
    fi
    (
        IFS=$'\n'
        for POD in ` + "`get_job_pods`" + `; do
        IFS=' '
        POD_NAME=` + "`echo $$POD | cut -d \" \" -f 1`" + `
        POD_STATUS=` + "`echo $$POD | cut -d \" \" -f 2`" + `
        POD_DATE_STR=` + "`echo $$POD | cut -d \" \" -f 3`" + `
        POD_START_DATE_STR=` + "`echo $$POD | cut -d \" \" -f 4`" + `
        POD_DATE=` + "`date -d $${POD_DATE_STR:-$$POD_START_DATE_STR} '+%s'`" + `
            if [ -n "${RUNNING_TTL_MINUTES}" ] && [ "$POD_STATUS" = "Running" ]; then
              if [ "$POD_DATE" -lt "$RUNNING_DATE" ]; then
                  delete_pod "$POD_NAME" "$POD_STATUS" "$POD_DATE_STR"
              fi
            elif [ -n "${SUCCEEDED_TTL_MINUTES}" ] && [ "$POD_STATUS" = "Succeeded" ]; then
              if [ "$POD_DATE" -lt "$SUCCESS_DATE" ]; then
                  delete_pod "$POD_NAME" "$POD_STATUS" "$POD_DATE_STR"
              fi
            elif [ -n "${UNSUCCESSFUL_TTL_MINUTES}" ] && [ "$POD_STATUS" != "Running" ] && [ "$POD_STATUS" != "Succeeded" ]; then
              if [ "$POD_DATE" -lt "$NON_SUCCESS_DATE" ]; then
                  delete_pod "$POD_NAME" "$POD_STATUS" "$POD_DATE_STR"
              fi
            fi
        done
    )
    echo "Completed pod sweeper cycle.  Sleeping for 60 seconds..."
    sleep 60
done`,
		},
	}
	err := ctrl.SetControllerReference(instance, configMap, schema)
	if err != nil {
		r.Log.Error(err, "Failed to set controller reference for airbyte-pod-sweeper-sweep-pod-script configmap")
		return nil
	}
	return configMap
}

func (r *AirbyteReconciler) makeTemporalDynamicconfigConfigMap(instance *stackv1alpha1.Airbyte, schema *runtime.Scheme) *corev1.ConfigMap {
	labels := instance.GetLabels()
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.GetNameWithSuffix("-tempora-dynamicconfig"),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Data: map[string]string{
			"development.yaml": `
# when modifying, remember to update the docker-compose version of this file in temporal/dynamicconfig/development.yaml
frontend.enableClientVersionCheck:
  - value: true
    constraints: {}
history.persistenceMaxQPS:
  - value: 3000
    constraints: {}
frontend.persistenceMaxQPS:
  - value: 3000
    constraints: {}
frontend.historyMgrNumConns:
  - value: 30
    constraints: {}
frontend.throttledLogRPS:
  - value: 20
    constraints: {}
history.historyMgrNumConns:
  - value: 50
    constraints: {}
system.advancedVisibilityWritingMode:
  - value: "off"
    constraints: {}
history.defaultActivityRetryPolicy:
  - value:
      InitialIntervalInSeconds: 1
      MaximumIntervalCoefficient: 100.0
      BackoffCoefficient: 2.0
      MaximumAttempts: 0
history.defaultWorkflowRetryPolicy:
  - value:
      InitialIntervalInSeconds: 1
      MaximumIntervalCoefficient: 100.0
      BackoffCoefficient: 2.0
      MaximumAttempts: 0
# Limit for responses. This mostly impacts discovery jobs since they have the largest responses.
limit.blobSize.error:
  - value: 15728640 # 15MB
    constraints: {}
limit.blobSize.warn:
  - value: 10485760 # 10MB
    constraints: {}
`,
		},
	}
	err := ctrl.SetControllerReference(instance, configMap, schema)
	if err != nil {
		r.Log.Error(err, "Failed to set controller reference for airbyte-pod-sweeper-sweep-pod-script configmap")
		return nil
	}
	return configMap
}

func addMapToData(data map[string]string, key string, value map[string]string) {
	if value != nil {
		var strList []string
		for k, v := range value {
			strList = append(strList, k+"="+v)
		}
		data[key] = strings.Join(strList, ",")
	}
}

func (r *AirbyteReconciler) makeEnvConfigMap(instance *stackv1alpha1.Airbyte, schema *runtime.Scheme) *corev1.ConfigMap {
	labels := instance.GetLabels()

	data := map[string]string{
		"AIRBYTE_VERSION": instance.Spec.Server.Image.Tag,
		"AIRBYTE_EDITION": instance.Spec.Global.Edition,
		"CONFIG_ROOT":     "/configs",
		"CONFIGS_DATABASE_MINIMUM_FLYWAY_MIGRATION_VERSION": "0.35.15.001",
		"DATA_DOCKER_MOUNT":          "airbyte_data",
		"DATABASE_DB":                instance.Spec.Postgres.DataBase,
		"DATABASE_HOST":              instance.Spec.Postgres.Host,
		"DATABASE_PORT":              instance.Spec.Postgres.Port,
		"DATABASE_URL":               "jdbc:postgresql://" + instance.Spec.Postgres.Host + ":" + instance.Spec.Postgres.Port + "/" + instance.Spec.Postgres.DataBase,
		"DB_DOCKER_MOUNT":            "airbyte_db",
		"INTERNAL_API_HOST":          instance.Name + "-airbyte-server-svc:" + strconv.FormatInt(int64(instance.Spec.Server.Service.Port), 10),
		"CONNECTOR_BUILDER_API_HOST": instance.Name + "-airbyte-connector-builder-server-svc:" + strconv.FormatInt(int64(instance.Spec.ConnectorBuilderServer.Service.Port), 10),
		"AIRBYTE_API_HOST":           instance.Name + "-airbyte-api-server-svc:" + strconv.FormatInt(int64(instance.Spec.AirbyteApiServer.Service.Port), 10),
		"JOBS_DATABASE_MINIMUM_FLYWAY_MIGRATION_VERSION": "0.29.15.001",
		"LOCAL_ROOT":                             "/tmp/airbyte_local",
		"STATE_STORAGE_MINIO_BUCKET_NAME":        "airbyte-state-storage",
		"TEMPORAL_HOST":                          instance.Name + "-temporal:" + strconv.FormatInt(int64(instance.Spec.Temporal.Service.Port), 10),
		"TEMPORAL_WORKER_PORTS":                  "9000-9040",
		"TRACKING_STRATEGY":                      "segment",
		"WORKER_ENVIRONMENT":                     "kubernetes",
		"WORKSPACE_DOCKER_MOUNT":                 "airbyte_workspace",
		"WORKSPACE_ROOT":                         "/workspace",
		"WORKFLOW_FAILURE_RESTART_DELAY_SECONDS": "",
		"USE_STREAM_CAPABLE_STATE":               "true",
		"AUTO_DETECT_SCHEMA":                     "true",
		"WORKERS_MICRONAUT_ENVIRONMENTS":         "control-plane",
		"CRON_MICRONAUT_ENVIRONMENTS":            "control-plane",
		"SHOULD_RUN_NOTIFY_WORKFLOWS":            "true",
	}

	stateStorageType := ""
	if instance.Spec.Global.StateStorageType != "" {
		stateStorageType = instance.Spec.Global.StateStorageType
	}
	data["WORKER_STATE_STORAGE_TYPE"] = stateStorageType

	containerOrchestratorImage := ""
	if instance != nil && instance.Spec.Worker != nil && instance.Spec.Worker.ContainerOrchestrator != nil {
		if instance.Spec.Worker.ContainerOrchestrator.Image != "" {
			containerOrchestratorImage = instance.Spec.Worker.ContainerOrchestrator.Image
		}
	}
	data["CONTAINER_ORCHESTRATOR_IMAGE"] = containerOrchestratorImage

	containerOrchestratorEnabled := "true"
	if instance != nil && instance.Spec.Worker != nil && instance.Spec.Worker.ContainerOrchestrator != nil {
		if instance.Spec.Worker.ContainerOrchestrator.Enabled {
			containerOrchestratorEnabled = strconv.FormatBool(instance.Spec.Worker.ContainerOrchestrator.Enabled)
		}
	}
	data["CONTAINER_ORCHESTRATOR_ENABLED"] = containerOrchestratorEnabled

	logStorageType := ""
	if instance != nil && instance.Spec.Global != nil && instance.Spec.Global.Log != nil {
		if instance.Spec.Global.Log.StorageType != "" {
			logStorageType = instance.Spec.Global.Log.StorageType
		}
	}
	data["WORKER_LOGS_STORAGE_TYPE"] = logStorageType

	apiUrl := ""
	if instance != nil && instance.Spec.WebApp != nil {
		if instance.Spec.WebApp.ApiUrl != "" {
			apiUrl = instance.Spec.WebApp.ApiUrl
		}
	}
	data["API_URL"] = apiUrl

	ServerUrl := ""
	if instance != nil && instance.Spec.WebApp != nil {
		if instance.Spec.WebApp.ConnectorBuilderServerUrl != "" {
			ServerUrl = instance.Spec.WebApp.ConnectorBuilderServerUrl
		}
	}
	data["CONNECTOR_BUILDER_API_URL"] = ServerUrl

	gcsBucket := ""
	if instance != nil && instance.Spec.Global != nil && instance.Spec.Global.Log != nil && instance.Spec.Global.Log.Gcs != nil {
		if instance.Spec.Global.Log.Gcs.Bucket != "" {
			gcsBucket = instance.Spec.Global.Log.Gcs.Bucket
		}
	}
	data["GCS_LOG_BUCKET"] = gcsBucket

	gcsCredentials := ""
	if instance != nil && instance.Spec.Global != nil && instance.Spec.Global.Log != nil && instance.Spec.Global.Log.Gcs != nil {
		if instance.Spec.Global.Log.Gcs.Credentials != "" {
			gcsCredentials = instance.Spec.Global.Log.Gcs.Credentials
		}
	}
	data["GOOGLE_APPLICATION_CREDENTIALS"] = gcsCredentials

	s3Bucket := ""
	if instance != nil && instance.Spec.Global != nil && instance.Spec.Global.Log != nil && instance.Spec.Global.Log.S3 != nil {
		if instance.Spec.Global.Log.S3.Bucket != "" {
			s3Bucket = instance.Spec.Global.Log.S3.Bucket
		}
	}
	data["S3_LOG_BUCKET"] = s3Bucket

	s3BucketRegion := ""
	if instance != nil && instance.Spec.Global != nil && instance.Spec.Global.Log != nil && instance.Spec.Global.Log.S3 != nil {
		if instance.Spec.Global.Log.S3.BucketRegion != "" {
			s3BucketRegion = instance.Spec.Global.Log.S3.BucketRegion
		}
	}
	data["S3_LOG_BUCKET_REGION"] = s3BucketRegion

	if instance != nil && instance.Spec.Global != nil {
		if instance.Spec.Global.Edition == "pro" {
			if instance.Spec.Keycloak != nil && instance.Spec.Keycloak.Service != nil {
				data["KEYCLOAK_INTERNAL_HOST"] = instance.Name + "-airbyte-keycloak-svc:" + strconv.FormatInt(int64(instance.Spec.Keycloak.Service.Port), 10)
				data["KEYCLOAK_PORT"] = strconv.FormatInt(int64(instance.Spec.Keycloak.Service.Port), 10)
				data["KEYCLOAK_HOSTNAME_URL"] = "webapp-url/auth"
			}
		} else {
			data["KEYCLOAK_INTERNAL_HOST"] = "localhost"
		}
	}

	if instance != nil && instance.Spec.Global != nil && instance.Spec.Global.Jobs != nil && instance.Spec.Global.Jobs.Kube != nil {
		if instance.Spec.Global.Jobs.Kube.Annotations != nil {
			addMapToData(data, "JOB_KUBE_ANNOTATIONS", instance.Spec.Global.Jobs.Kube.Annotations)
		}

		if instance.Spec.Global.Jobs.Kube.Labels != nil {
			addMapToData(data, "JOB_KUBE_LABELS", instance.Spec.Global.Jobs.Kube.Labels)
		}

		if instance.Spec.Global.Jobs.Kube.NodeSelector != nil {
			addMapToData(data, "JOB_KUBE_NODE_SELECTOR", instance.Spec.Global.Jobs.Kube.NodeSelector)
		}

		if instance.Spec.Global.Jobs.Kube.MainContainerImagePullSecret != "" {
			data["JOB_KUBE_MAIN_CONTAINER_IMAGE_PULL_SECRET"] = instance.Spec.Global.Jobs.Kube.MainContainerImagePullSecret
		}

		if instance.Spec.Global.Jobs.Kube.Tolerations != nil {
			tolerations, err := json.Marshal(instance.Spec.Global.Jobs.Kube.Tolerations)
			if err != nil {
				// 处理错误
			}
			data["JOB_KUBE_TOLERATIONS"] = string(tolerations)
		}

		if instance.Spec.Global.Jobs.Kube.Images.Busybox != "" {
			data["JOB_KUBE_BUSYBOX_IMAGE"] = instance.Spec.Global.Jobs.Kube.Images.Busybox
		}

		if instance.Spec.Global.Jobs.Kube.Images.Socat != "" {
			data["JOB_KUBE_SOCAT_IMAGE"] = instance.Spec.Global.Jobs.Kube.Images.Socat
		}

		if instance.Spec.Global.Jobs.Kube.Images.Curl != "" {
			data["JOB_KUBE_CURL_IMAGE"] = instance.Spec.Global.Jobs.Kube.Images.Curl
		}
	}

	cpuLimit := ""
	if instance != nil && instance.Spec.Global != nil && instance.Spec.Global.Jobs != nil && instance.Spec.Global.Jobs.Resources != nil && instance.Spec.Global.Jobs.Resources.Limits != nil {
		if cpu := instance.Spec.Global.Jobs.Resources.Limits.Cpu(); cpu != nil {
			cpuLimit = cpu.String()
		}
	}
	data["JOB_MAIN_CONTAINER_CPU_LIMIT"] = cpuLimit

	cpuRequest := ""
	if instance != nil && instance.Spec.Global != nil && instance.Spec.Global.Jobs != nil && instance.Spec.Global.Jobs.Resources != nil && instance.Spec.Global.Jobs.Resources.Requests != nil {
		if cpu := instance.Spec.Global.Jobs.Resources.Requests.Cpu(); cpu != nil {
			cpuRequest = cpu.String()
		}
	}
	data["JOB_MAIN_CONTAINER_CPU_REQUEST"] = cpuRequest

	memoryLimit := ""
	if instance != nil && instance.Spec.Global != nil && instance.Spec.Global.Jobs != nil && instance.Spec.Global.Jobs.Resources != nil && instance.Spec.Global.Jobs.Resources.Limits != nil {
		if mem := instance.Spec.Global.Jobs.Resources.Limits.Memory(); mem != nil {
			memoryLimit = mem.String()
		}
	}
	data["JOB_MAIN_CONTAINER_MEMORY_LIMIT"] = memoryLimit

	memoryRequest := ""
	if instance != nil && instance.Spec.Global != nil && instance.Spec.Global.Jobs != nil && instance.Spec.Global.Jobs.Resources != nil && instance.Spec.Global.Jobs.Resources.Requests != nil {
		if mem := instance.Spec.Global.Jobs.Resources.Requests.Memory(); mem != nil {
			memoryRequest = mem.String()
		}
	}
	data["JOB_MAIN_CONTAINER_MEMORY_REQUEST"] = memoryRequest

	runDatabaseMigrationOnStartup := "true"
	if instance.Spec.AirbyteBootloader.RunDatabaseMigrationsOnStartup != nil {
		runDatabaseMigrationOnStartup = strconv.FormatBool(*instance.Spec.AirbyteBootloader.RunDatabaseMigrationsOnStartup)
	}
	data["RUN_DATABASE_MIGRATION_ON_STARTUP"] = runDatabaseMigrationOnStartup

	s3MinioEndpoint := ""
	s3PathStyleAccess := false
	if instance != nil && instance.Spec.Global != nil && instance.Spec.Global.Log != nil {
		if instance.Spec.Global.Log.Minio != nil && instance.Spec.Global.Log.Minio.Enabled {
			s3MinioEndpoint = "http://airbyte-minio-svc:9000"
			s3PathStyleAccess = true
		} else if instance.Spec.Global.Log.ExternalMinio != nil {
			if instance.Spec.Global.Log.ExternalMinio.Endpoint != "" {
				s3MinioEndpoint = instance.Spec.Global.Log.ExternalMinio.Endpoint
			} else if instance.Spec.Global.Log.ExternalMinio.Enabled {
				s3MinioEndpoint = "http://" + instance.Spec.Global.Log.ExternalMinio.Host + ":" + strconv.Itoa(int(instance.Spec.Global.Log.ExternalMinio.Port))
				s3PathStyleAccess = true
			}
		}
	}
	data["S3_MINIO_ENDPOINT"] = s3MinioEndpoint
	data["STATE_STORAGE_MINIO_ENDPOINT"] = s3MinioEndpoint
	data["S3_PATH_STYLE_ACCESS"] = strconv.FormatBool(s3PathStyleAccess)

	webAppUrl := ""
	if instance != nil && instance.Spec.WebApp != nil {
		if instance.Spec.WebApp.Url != "" {
			webAppUrl = instance.Spec.WebApp.Url
		} else if instance.Spec.WebApp.Service != nil {
			webAppUrl = fmt.Sprintf("http://%s-airbyte-webapp-svc:%d", instance.Name, instance.Spec.WebApp.Service.Port)
		}
	}
	data["WEBAPP_URL"] = webAppUrl

	metricClient := ""
	if instance != nil && instance.Spec.Global != nil && instance.Spec.Global.Metrics != nil {
		if instance.Spec.Global.Metrics.MetricClient != "" {
			metricClient = instance.Spec.Global.Metrics.MetricClient
		}
	}
	data["METRIC_CLIENT"] = metricClient

	otelCollectorEndpoint := ""
	if instance != nil && instance.Spec.Global != nil && instance.Spec.Global.Metrics != nil {
		if instance.Spec.Global.Metrics.OtelCollectorEndpoint != "" {
			otelCollectorEndpoint = instance.Spec.Global.Metrics.OtelCollectorEndpoint
		}
	}
	data["OTEL_COLLECTOR_ENDPOINT"] = otelCollectorEndpoint

	activityMaxAttempt := ""
	activityInitialDelayBetweenAttemptsSeconds := ""
	activityMaxDelayBetweenAttemptsSeconds := ""
	maxNotifyWorkers := "5"
	if instance != nil && instance.Spec.Worker != nil {
		if instance.Spec.Worker.ActivityMaxAttempt != "" {
			activityMaxAttempt = instance.Spec.Worker.ActivityMaxAttempt
		}

		if instance.Spec.Worker.ActivityInitialDelayBetweenAttemptsSeconds != "" {
			activityInitialDelayBetweenAttemptsSeconds = instance.Spec.Worker.ActivityInitialDelayBetweenAttemptsSeconds
		}

		if instance.Spec.Worker.ActivityMaxDelayBetweenAttemptsSeconds != "" {
			activityMaxDelayBetweenAttemptsSeconds = instance.Spec.Worker.ActivityMaxDelayBetweenAttemptsSeconds
		}

		if instance.Spec.Worker.MaxNotifyWorkers != "" {
			maxNotifyWorkers = instance.Spec.Worker.MaxNotifyWorkers
		}
	}
	data["ACTIVITY_MAX_ATTEMPT"] = activityMaxAttempt
	data["ACTIVITY_INITIAL_DELAY_BETWEEN_ATTEMPTS_SECONDS"] = activityInitialDelayBetweenAttemptsSeconds
	data["ACTIVITY_MAX_DELAY_BETWEEN_ATTEMPTS_SECONDS"] = activityMaxDelayBetweenAttemptsSeconds
	data["MAX_NOTIFY_WORKERS"] = maxNotifyWorkers

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.GetNameWithSuffix("-airbyte-env"),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Data: data,
	}
	err := ctrl.SetControllerReference(instance, configMap, schema)
	if err != nil {
		r.Log.Error(err, "Failed to set controller reference for configmap")
		return nil
	}
	return configMap
}

func (r *AirbyteReconciler) reconcileConfigMap(ctx context.Context, instance *stackv1alpha1.Airbyte) error {

	YmlConfigMap := r.makeYmlConfigMap(instance, r.Scheme)
	if YmlConfigMap == nil {
		return nil
	}

	EnvConfigMap := r.makeEnvConfigMap(instance, r.Scheme)
	if EnvConfigMap == nil {
		return nil
	}

	SweepPodScriptConfigMap := r.makeSweepPodScriptConfigMap(instance, r.Scheme)
	if SweepPodScriptConfigMap == nil {
		return nil
	}

	TemporalDynamicconfigConfigMap := r.makeTemporalDynamicconfigConfigMap(instance, r.Scheme)
	if TemporalDynamicconfigConfigMap == nil {
		return nil
	}

	ymlConfigMap := r.makeYmlConfigMap(instance, r.Scheme)
	if ymlConfigMap != nil {
		if err := CreateOrUpdate(ctx, r.Client, ymlConfigMap); err != nil {
			r.Log.Error(err, "Failed to create or update airbyte-yml configmap")
			return err
		}
	}

	envConfigMap := r.makeEnvConfigMap(instance, r.Scheme)
	if envConfigMap != nil {
		if err := CreateOrUpdate(ctx, r.Client, envConfigMap); err != nil {
			r.Log.Error(err, "Failed to create or update airbyte-env configmap")
			return err
		}
	}

	sweepPodScriptConfigMap := r.makeSweepPodScriptConfigMap(instance, r.Scheme)
	if sweepPodScriptConfigMap != nil {
		if err := CreateOrUpdate(ctx, r.Client, sweepPodScriptConfigMap); err != nil {
			r.Log.Error(err, "Failed to create or update airbyte-sweep-pod-script configmap")
			return err
		}
	}

	temporalDynamicconfigConfigMap := r.makeTemporalDynamicconfigConfigMap(instance, r.Scheme)
	if temporalDynamicconfigConfigMap != nil {
		if err := CreateOrUpdate(ctx, r.Client, temporalDynamicconfigConfigMap); err != nil {
			r.Log.Error(err, "Failed to create or update airbyte-temporal-dynamicconfig configmap")
			return err
		}
	}
	return nil
}

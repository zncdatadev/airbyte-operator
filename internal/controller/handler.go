package controller

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"

	stackv1alpha1 "github.com/zncdata-labs/airbyte-operator/api/v1alpha1"
)

// make pvc
func makePVC(ctx context.Context, instance *stackv1alpha1.Airbyte, schema *runtime.Scheme) *corev1.PersistentVolumeClaim {
	logger := log.FromContext(ctx)
	labels := instance.GetLabels()
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.GetPvcName(),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			StorageClassName: instance.Spec.Persistence.StorageClassName,
			AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.PersistentVolumeAccessMode(instance.Spec.Persistence.AccessMode)},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(instance.Spec.Persistence.StorageSize),
				},
			},
			VolumeMode: instance.Spec.Persistence.VolumeMode,
		},
	}
	err := ctrl.SetControllerReference(instance, pvc, schema)
	if err != nil {
		logger.Error(err, "Failed to set controller reference")
		return nil
	}
	return pvc
}

// reconcilePVC
func (r *AirbyteReconciler) reconcilePVC(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	logger := log.FromContext(ctx)
	pvc := &corev1.PersistentVolumeClaim{}
	err := r.Client.Get(ctx, types.NamespacedName{Namespace: instance.Namespace, Name: instance.GetPvcName()}, pvc)
	if err != nil && errors.IsNotFound(err) {
		pvc := makePVC(ctx, instance, r.Scheme)
		logger.Info("Creating a new PVC", "PVC.Namespace", pvc.Namespace, "PVC.Name", pvc.Name)
		err := r.Client.Create(ctx, pvc)
		if err != nil {
			return err
		}
	} else if err != nil {
		logger.Error(err, "Failed to get PVC")
		return err
	}
	return nil
}

// make service
func makeService(ctx context.Context, instance *stackv1alpha1.Airbyte, schema *runtime.Scheme) *corev1.Service {
	labels := instance.GetLabels()
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        instance.Name,
			Namespace:   instance.Namespace,
			Labels:      labels,
			Annotations: instance.Spec.Service.Annotations,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port:     instance.Spec.Service.Port,
					Name:     "http",
					Protocol: "TCP",
				},
			},
			Selector: labels,
			Type:     instance.Spec.Service.Type,
		},
	}
	err := ctrl.SetControllerReference(instance, svc, schema)
	if err != nil {
		return nil
	}
	return svc
}

func (r *AirbyteReconciler) reconcileService(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	logger := log.FromContext(ctx)
	obj := makeService(ctx, instance, r.Scheme)
	if obj == nil {
		return nil
	}

	if err := CreateOrUpdate(ctx, r.Client, obj); err != nil {
		logger.Error(err, "Failed to create or update service")
		return err
	}
	return nil
}

func makeDeployment(ctx context.Context, instance *stackv1alpha1.Airbyte, schema *runtime.Scheme) *appsv1.Deployment {
	labels := instance.GetLabels()
	envVarNames := []string{"AIRBYTE_VERSION", "AIRBYTE_EDITION", "AUTO_DETECT_SCHEMA", "CONFIG_ROOT", "DATABASE_URL",
		"TRACKING_STRATEGY", "WORKER_ENVIRONMENT", "WORKSPACE_ROOT", "WEBAPP_URL", "TEMPORAL_HOST", "JOB_MAIN_CONTAINER_CPU_REQUEST",
		"JOB_MAIN_CONTAINER_CPU_LIMIT", "JOB_MAIN_CONTAINER_MEMORY_REQUEST", "JOB_MAIN_CONTAINER_MEMORY_LIMIT", "S3_LOG_BUCKET", "S3_LOG_BUCKET_REGION",
		"S3_MINIO_ENDPOINT", "STATE_STORAGE_MINIO_BUCKET_NAME", "STATE_STORAGE_MINIO_ENDPOINT", "S3_PATH_STYLE_ACCESS", "GOOGLE_APPLICATION_CREDENTIALS",
		"GCS_LOG_BUCKET", "CONFIGS_DATABASE_MINIMUM_FLYWAY_MIGRATION_VERSION",
		"JOBS_DATABASE_MINIMUM_FLYWAY_MIGRATION_VERSION", "WORKER_LOGS_STORAGE_TYPE", "WORKER_STATE_STORAGE_TYPE", "KEYCLOAK_INTERNAL_HOST"}
	secretVarNames := []string{"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "DATABASE_PASSWORD", "DATABASE_USER", "STATE_STORAGE_MINIO_ACCESS_KEY"}
	envVars := []corev1.EnvVar{}

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
			Name:      instance.Name,
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &instance.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            instance.Name,
							Image:           instance.GetImageTag(),
							ImagePullPolicy: instance.GetImagePullPolicy(),
							Env:             envVars,
							// Args: []string{
							// 	"/opt/bitnami/scripts/start-scripts/start-master.sh",
							// },
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: instance.Spec.Service.Port,
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
					Tolerations: instance.Spec.Tolerations,
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
			// },
		},
	}
	err := ctrl.SetControllerReference(instance, dep, schema)
	if err != nil {
		return nil
	}
	return dep
}

func (r *AirbyteReconciler) reconcileDeployment(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	logger := log.FromContext(ctx)
	obj := makeDeployment(ctx, instance, r.Scheme)
	if obj == nil {
		return nil
	}
	if err := CreateOrUpdate(ctx, r.Client, obj); err != nil {
		logger.Error(err, "Failed to create or update deployment")
		return err
	}
	return nil
}

func makeSecret(ctx context.Context, instance *stackv1alpha1.Airbyte, schema *runtime.Scheme) []*corev1.Secret {
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

	return secrets
}

func (r *AirbyteReconciler) reconcileSecret(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	logger := log.FromContext(ctx)
	objs := makeSecret(ctx, instance, r.Scheme)

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

func (r *AirbyteReconciler) makeConfigMap(ctx context.Context, instance *stackv1alpha1.Airbyte, schema *runtime.Scheme) *corev1.ConfigMap {
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

func (r *AirbyteReconciler) reconcileConfigMap(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	obj := r.makeConfigMap(ctx, instance, r.Scheme)
	if obj == nil {
		return nil
	}

	if err := CreateOrUpdate(ctx, r.Client, obj); err != nil {
		r.Log.Error(err, "Failed to create or update service")
		return err
	}
	return nil
}

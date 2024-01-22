package controller

import (
	"encoding/base64"
	stackv1alpha1 "github.com/zncdata-labs/airbyte-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"strconv"
)

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
				}
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

			envVars = append(envVars, corev1.EnvVar{
				Name: "POSTGRES_USER",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						Key: "DATABASE_USER",
						LocalObjectReference: corev1.LocalObjectReference{
							Name: instance.GetNameWithSuffix("-airbyte-secrets"),
						},
					},
				},
			})

			envVars = append(envVars, corev1.EnvVar{
				Name: "POSTGRES_PWD",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						Key: "DATABASE_PASSWORD",
						LocalObjectReference: corev1.LocalObjectReference{
							Name: instance.GetNameWithSuffix("-airbyte-secrets"),
						},
					},
				},
			})

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
					ServiceAccountName: instance.GetNameWithSuffix("-admin"),
					SecurityContext:    instance.Spec.SecurityContext,
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
									MountPath: "/etc/temporal/config/dynamicconfig/",
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

func (r *AirbyteReconciler) makeSecret(instance *stackv1alpha1.Airbyte) []*corev1.Secret {
	var secrets []*corev1.Secret
	labels := instance.GetLabels()
	if instance.Spec.Global != nil && instance.Spec.Global.Logs != nil && instance.Spec.Global.Logs.Gcs != nil {
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      instance.GetNameWithSuffix("-gcs-log-creds"),
				Namespace: instance.Namespace,
				Labels:    labels,
			},
			Type: corev1.SecretTypeOpaque,
			Data: map[string][]byte{
				"gcp.json": []byte(instance.Spec.Global.Logs.Gcs.CredentialsJson),
			},
		}
		secrets = append(secrets, secret)

	}

	if instance.Spec.Global != nil && instance.Spec.Global.Logs != nil {
		// Create a Data map to hold the secret data
		data := make(map[string][]byte)

		data["AWS_ACCESS_KEY_ID"] = []byte("")
		if instance.Spec.Global.Logs.AccessKey != nil {
			data["AWS_ACCESS_KEY_ID"] = []byte(instance.Spec.Global.Logs.AccessKey.Password)
		}

		data["AWS_SECRET_ACCESS_KEY"] = []byte("")
		if instance.Spec.Global.Logs.SecretKey != nil {
			data["AWS_SECRET_ACCESS_KEY"] = []byte(instance.Spec.Global.Logs.SecretKey.Password)
		}

		data["DATABASE_PASSWORD"] = []byte("")
		if instance.Spec.Postgres != nil {
			data["DATABASE_PASSWORD"] = []byte(instance.Spec.Postgres.Password)
		}

		data["DATABASE_USER"] = []byte("")
		if instance.Spec.Postgres != nil {
			data["DATABASE_USER"] = []byte(instance.Spec.Postgres.UserName)
		}

		data["STATE_STORAGE_MINIO_ACCESS_KEY"] = []byte("")
		if instance.Spec.Minio != nil {
			data["STATE_STORAGE_MINIO_ACCESS_KEY"] = []byte(instance.Spec.Minio.RootUser)
		}

		data["STATE_STORAGE_MINIO_SECRET_ACCESS_KEY"] = []byte("")
		if instance.Spec.Minio != nil {
			data["STATE_STORAGE_MINIO_SECRET_ACCESS_KEY"] = []byte(instance.Spec.Minio.RootPassword)
		}

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

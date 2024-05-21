package controller

import (
	"context"
	"encoding/base64"
	stackv1alpha1 "github.com/zncdatadev/airbyte-operator/api/v1alpha1"
	opgo "github.com/zncdatadev/operator-go/pkg/apis/commons/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
)

// reconcile temporal service
func (r *AirbyteReconciler) reconcileTemporalService(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	roleGroups := convertRoleGroupToRoleConfigObject(instance, Temporal)
	reconcileParams := r.createReconcileParams(ctx, roleGroups, instance, r.extractTemporalServiceForRoleGroup)
	if err := reconcileParams.createOrUpdateResource(); err != nil {
		return err
	}
	return nil
}

// extract temporal service for role group
func (r *AirbyteReconciler) extractTemporalServiceForRoleGroup(params ExtractorParams) (client.Object, error) {
	instance := params.instance.(*stackv1alpha1.Airbyte)
	temporal := instance.Spec.Temporal
	groupCfg := params.roleGroup
	roleCfg := temporal.RoleConfig
	clusterCfg := params.cluster
	mergedLabels := r.mergeLabels(groupCfg, instance.GetLabels(), clusterCfg, Temporal)
	schema := params.scheme
	port, serviceType, annotations := getServiceInfo(groupCfg, roleCfg, clusterCfg)
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        createSvcNameForRoleGroup(instance, Temporal, params.roleGroupName),
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
		return nil, err
	}
	return svc, nil
}

// reconcile temporal deployment
func (r *AirbyteReconciler) reconcileTemporalDeployment(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	roleGroups := convertRoleGroupToRoleConfigObject(instance, Temporal)
	reconcileParams := r.createReconcileParams(ctx, roleGroups, instance, r.extractTemporalDeploymentForRoleGroup)
	if err := reconcileParams.createOrUpdateResource(); err != nil {
		return err
	}
	return nil
}

// extract temporal deployment for role group
func (r *AirbyteReconciler) extractTemporalDeploymentForRoleGroup(params ExtractorParams) (client.Object, error) {
	instance := params.instance.(*stackv1alpha1.Airbyte)
	temporal := instance.Spec.Temporal
	groupCfg := params.roleGroup
	roleCfg := temporal.RoleConfig
	clusterCfg := params.cluster
	mergedLabels := r.mergeLabels(groupCfg, instance.GetLabels(), clusterCfg, Temporal)
	schema := params.scheme
	image, securityContext, replicas, resources, containerPorts := getDeploymentInfo(groupCfg, roleCfg, clusterCfg)
	envVars, err := r.mergeEnvVarsForTemporalDeployment(instance, params.ctx)
	if err != nil {
		return nil, err
	}
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      createNameForRoleGroup(instance, "temporal", params.roleGroupName),
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
										Name: createNameForRoleGroup(instance, "tempora-dynamicconfig", params.roleGroupName),
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
	ConnectorBuilderServerScheduler(instance, dep, groupCfg)

	err = ctrl.SetControllerReference(instance, dep, schema)
	if err != nil {
		r.Log.Error(err, "Failed to set controller reference for deployment")
		return nil, err
	}
	return dep, nil
}

// merge env vars for temporal deployment
func (r *AirbyteReconciler) mergeEnvVarsForTemporalDeployment(instance *stackv1alpha1.Airbyte,
	ctx context.Context) ([]corev1.EnvVar, error) {
	var envVars []corev1.EnvVar
	if instance != nil && instance.Spec.ClusterConfig != nil {
		if instance.Spec.ClusterConfig.DeploymentMode == "oss" {
			// db envs
			envVars = append(envVars, corev1.EnvVar{
				Name:  "AUTO_SETUP",
				Value: "true",
			}, corev1.EnvVar{
				Name:  "DB",
				Value: "postgresql",
			})

			db := instance.Spec.ClusterConfig.Config.Database
			if db != nil {
				dbrsc := &opgo.Database{
					ObjectMeta: metav1.ObjectMeta{Name: db.Reference},
				}
				resourceReq := &ResourceRequest{
					Ctx:       ctx,
					Client:    r.Client,
					Namespace: instance.Namespace,
					Log:       r.Log,
				}
				dbConnection, err := resourceReq.FetchDb(dbrsc)
				if err != nil {
					return nil, err
				}
				pg := dbConnection.Spec.Provider.Postgres
				pgEnvs := []corev1.EnvVar{
					{
						Name:  "DB_PORT",
						Value: strconv.Itoa(pg.Port),
					},
					{
						Name:  "POSTGRES_USER",
						Value: dbrsc.Spec.Credential.Username,
					},
					{
						Name:  "POSTGRES_PWD",
						Value: dbrsc.Spec.Credential.Password,
					},
					{
						Name:  "POSTGRES_SEEDS",
						Value: pg.Host,
					},
				}
				envVars = append(envVars, pgEnvs...)
			}
			envVars = append(envVars, corev1.EnvVar{
				Name:  "DYNAMIC_CONFIG_FILE_PATH",
				Value: "config/dynamicconfig/development.yaml",
			})
		}
	}

	// temporal secret envs
	if instance != nil && instance.Spec.Temporal != nil && instance.Spec.Temporal.RoleConfig.Secret != nil {
		for key := range instance.Spec.Temporal.RoleConfig.Secret {
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
	// merge env vars
	if instance != nil && (instance.Spec.Temporal != nil || instance.Spec.ClusterConfig != nil) {
		envVarsMap := make(map[string]string)
		if instance.Spec.Temporal != nil && instance.Spec.Temporal.RoleConfig.EnvVars != nil {
			for key, value := range instance.Spec.Temporal.RoleConfig.EnvVars {
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
	// extra env
	if instance != nil && instance.Spec.Temporal != nil && instance.Spec.Temporal.RoleConfig.ExtraEnv != (corev1.EnvVar{}) {
		envVars = append(envVars, instance.Spec.Temporal.RoleConfig.ExtraEnv)
	}
	return envVars, nil
}

// reconcile temporal configmap
func (r *AirbyteReconciler) reconcileTemporalConfigMap(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	roleGroups := convertRoleGroupToRoleConfigObject(instance, Temporal)
	reconcileParams := r.createReconcileParams(ctx, roleGroups, instance, r.extractTemporalDynamicConfigConfigMap)
	if err := reconcileParams.createOrUpdateResource(); err != nil {
		return err
	}
	return nil
}

// extract temporal dynamic config configmap
func (r *AirbyteReconciler) extractTemporalDynamicConfigConfigMap(params ExtractorParams) (client.Object, error) {
	instance := params.instance.(*stackv1alpha1.Airbyte)
	mergedLabels := r.mergeLabels(params.roleGroup, instance.GetLabels(), params.cluster, Temporal)
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      createNameForRoleGroup(instance, "tempora-dynamicconfig", params.roleGroupName),
			Namespace: instance.Namespace,
			Labels:    mergedLabels,
		},
		Data: map[string]string{
			"development.yaml": temporalDynamicConfigDeploymentCfg,
		},
	}
	err := ctrl.SetControllerReference(instance, configMap, params.scheme)
	if err != nil {
		return nil, err
	}
	return configMap, nil
}

// reconcile temporal secrets
func (r *AirbyteReconciler) reconcileTemporalSecret(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	roleGroups := convertRoleGroupToRoleConfigObject(instance, Temporal)
	reconcileParams := r.createReconcileParams(ctx, roleGroups, instance, r.extractTemporalSecret)
	if err := reconcileParams.createOrUpdateResource(); err != nil {
		return err
	}
	return nil
}

// extract temporal secret
func (r *AirbyteReconciler) extractTemporalSecret(params ExtractorParams) (client.Object, error) {
	instance := params.instance.(*stackv1alpha1.Airbyte)
	mergedLabels := r.mergeLabels(params.roleGroup, instance.GetLabels(), params.cluster, Temporal)
	if secretField := instance.Spec.Temporal.RoleConfig.Secret; secretField != nil {
		// Create a Data map to hold the secret data
		Data := make(map[string][]byte)
		// Iterate over instance.Spec.Temporal.Secret and create Secret data items
		for k, v := range secretField {
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
				Labels:    mergedLabels,
			},
			Type: corev1.SecretTypeOpaque,
			Data: Data,
		}
		return secret, nil
	}
	return nil, nil
}

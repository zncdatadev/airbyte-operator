package controller

import (
	"context"
	stackv1alpha1 "github.com/zncdatadev/airbyte-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
)

// reconcile Airbyte pod sweeper deployment
func (r *AirbyteReconciler) reconcilePodSweeperDeployment(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	roleGroups := convertRoleGroupToRoleConfigObject(instance, PodSweeper)
	reconcileParams := r.createReconcileParams(ctx, roleGroups, instance, r.extractPodSweeperDeploymentForRoleGroup)
	if err := reconcileParams.createOrUpdateResource(); err != nil {
		return err
	}
	return nil
}

// extract pod sweeper deployment for role group
func (r *AirbyteReconciler) extractPodSweeperDeploymentForRoleGroup(params ExtractorParams) (client.Object, error) {
	instance := params.instance.(*stackv1alpha1.Airbyte)
	server := instance.Spec.ApiServer
	groupCfg := params.roleGroup
	roleCfg := server.RoleConfig
	clusterCfg := params.cluster
	roleGroupName := params.roleGroupName
	mergedLabels := r.mergeLabels(groupCfg, instance.GetLabels(), clusterCfg, PodSweeper)
	schema := params.scheme

	realGroupCfg := groupCfg.(*stackv1alpha1.PodSweeperRoleConfigSpec)
	image, securityContext, replicas, resources, _ := getDeploymentInfo(groupCfg, roleCfg, clusterCfg)
	envVars := mergeEnvVarsForPodSweeperDeployment(instance, realGroupCfg)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      createNameForRoleGroup(instance, "pod-sweeper", roleGroupName),
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
							Name:            "airbyte-pod-sweeper",
							Image:           image.Repository + ":" + image.Tag,
							ImagePullPolicy: image.PullPolicy,
							Resources:       *resources,
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
										Name: createNameForRoleGroup(instance, "sweep-pod-script", roleGroupName),
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

	CronScheduler(instance, dep, groupCfg)

	container := &dep.Spec.Template.Spec.Containers[0]
	if realGroupCfg.Config != nil && len(realGroupCfg.ExtraVolumeMounts) > 0 {
		container.VolumeMounts = append(container.VolumeMounts, realGroupCfg.ExtraVolumeMounts...)
	} else if instance.Spec.PodSweeper != nil && len(instance.Spec.PodSweeper.RoleConfig.ExtraVolumeMounts) > 0 {
		container.VolumeMounts = append(container.VolumeMounts, instance.Spec.PodSweeper.RoleConfig.ExtraVolumeMounts...)
	}

	podSpec := &dep.Spec.Template.Spec
	if realGroupCfg.Config != nil && len(realGroupCfg.ExtraVolumes) > 0 {
		podSpec.Volumes = append(podSpec.Volumes, realGroupCfg.ExtraVolumes...)
	} else if instance.Spec.PodSweeper.RoleConfig != nil && len(instance.Spec.PodSweeper.RoleConfig.ExtraVolumes) > 0 {
		podSpec.Volumes = append(podSpec.Volumes, instance.Spec.PodSweeper.RoleConfig.ExtraVolumes...)
	}

	if realGroupCfg.Config != nil && realGroupCfg.LivenessProbe != nil && realGroupCfg.LivenessProbe.Enabled {
		dep.Spec.Template.Spec.Containers[0].LivenessProbe = createProbe(realGroupCfg.LivenessProbe)
	} else if instance.Spec.PodSweeper.RoleConfig.LivenessProbe != nil && instance.Spec.PodSweeper.RoleConfig.LivenessProbe.Enabled {
		dep.Spec.Template.Spec.Containers[0].LivenessProbe = createProbe(instance.Spec.PodSweeper.RoleConfig.LivenessProbe)
	}

	if realGroupCfg.Config != nil && realGroupCfg.ReadinessProbe != nil && realGroupCfg.ReadinessProbe.Enabled {
		dep.Spec.Template.Spec.Containers[0].ReadinessProbe = createProbe(realGroupCfg.ReadinessProbe)
	} else if instance.Spec.PodSweeper.RoleConfig.ReadinessProbe != nil && instance.Spec.PodSweeper.RoleConfig.ReadinessProbe.Enabled {
		dep.Spec.Template.Spec.Containers[0].ReadinessProbe = createProbe(instance.Spec.PodSweeper.RoleConfig.ReadinessProbe)
	}

	err := ctrl.SetControllerReference(instance, dep, schema)

	if err != nil {
		r.Log.Error(err, "Failed to set controller reference for deployment")
		return nil, err
	}
	return dep, nil
}

// merge env vars for deployment
func mergeEnvVarsForPodSweeperDeployment(instance *stackv1alpha1.Airbyte,
	roleGroup *stackv1alpha1.PodSweeperRoleConfigSpec) []corev1.EnvVar {
	var envVars []corev1.EnvVar
	envVars = append(envVars, corev1.EnvVar{
		Name: "KUBE_NAMESPACE",
		ValueFrom: &corev1.EnvVarSource{
			FieldRef: &corev1.ObjectFieldSelector{
				FieldPath: "metadata.namespace",
			},
		},
	})

	if roleGroup != nil && roleGroup.Config != nil && roleGroup.TimeToDeletePods != nil {
		envVars = append(envVars, corev1.EnvVar{
			Name:  "RUNNING_TTL_MINUTES",
			Value: roleGroup.TimeToDeletePods.Running,
		})
		envVars = append(envVars, corev1.EnvVar{
			Name:  "SUCCEEDED_TTL_MINUTES",
			Value: strconv.Itoa(int(roleGroup.TimeToDeletePods.Succeeded)),
		})
		envVars = append(envVars, corev1.EnvVar{
			Name:  "UNSUCCESSFUL_TTL_MINUTES",
			Value: strconv.Itoa(int(roleGroup.TimeToDeletePods.Unsuccessful)),
		})
	} else if instance.Spec.PodSweeper != nil && instance.Spec.PodSweeper.RoleConfig.TimeToDeletePods != nil {
		envVars = append(envVars, corev1.EnvVar{
			Name:  "RUNNING_TTL_MINUTES",
			Value: instance.Spec.PodSweeper.RoleConfig.TimeToDeletePods.Running,
		})
		envVars = append(envVars, corev1.EnvVar{
			Name:  "SUCCEEDED_TTL_MINUTES",
			Value: strconv.Itoa(int(instance.Spec.PodSweeper.RoleConfig.TimeToDeletePods.Succeeded)),
		})
		envVars = append(envVars, corev1.EnvVar{
			Name:  "UNSUCCESSFUL_TTL_MINUTES",
			Value: strconv.Itoa(int(instance.Spec.PodSweeper.RoleConfig.TimeToDeletePods.Unsuccessful)),
		})
	}
	return envVars
}

// reconcile Airbyte pod sweeper config map
func (r *AirbyteReconciler) reconcilePodSweeperConfigMap(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	roleGroups := convertRoleGroupToRoleConfigObject(instance, PodSweeper)
	reconcileParams := r.createReconcileParams(ctx, roleGroups, instance, r.extractPodSweeperConfigMapForRoleGroup)
	if err := reconcileParams.createOrUpdateResource(); err != nil {
		return err
	}
	return nil
}

// extract pod sweeper config map for role group
func (r *AirbyteReconciler) extractPodSweeperConfigMapForRoleGroup(params ExtractorParams) (client.Object, error) {
	instance := params.instance.(*stackv1alpha1.Airbyte)
	clusterCfg := params.cluster
	mergedLabels := r.mergeLabels(params.roleGroup, instance.GetLabels(), clusterCfg, PodSweeper)
	schema := params.scheme
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      createNameForRoleGroup(instance, "sweep-pod-script", params.roleGroupName),
			Namespace: instance.Namespace,
			Labels:    mergedLabels,
		},
		Data: map[string]string{
			"sweep-pod.sh": sweepPodScript,
		},
	}
	err := ctrl.SetControllerReference(instance, configMap, schema)
	if err != nil {
		r.Log.Error(err, "Failed to set controller reference for airbyte-pod-sweeper-sweep-pod-script configmap")
		return nil, err
	}
	return configMap, nil
}

func createProbe(probeConfig *stackv1alpha1.ProbeSpec) *corev1.Probe {
	return &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			Exec: &corev1.ExecAction{
				Command: []string{
					"/bin/sh",
					"-c",
					"grep -q sweep-pod.sh /proc/1/cmdline",
				},
			},
		},
		InitialDelaySeconds: probeConfig.InitialDelaySeconds,
		PeriodSeconds:       probeConfig.PeriodSeconds,
		TimeoutSeconds:      probeConfig.TimeoutSeconds,
		SuccessThreshold:    probeConfig.SuccessThreshold,
		FailureThreshold:    probeConfig.FailureThreshold,
	}
}

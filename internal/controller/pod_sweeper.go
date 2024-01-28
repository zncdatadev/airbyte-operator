package controller

import (
	stackv1alpha1 "github.com/zncdata-labs/airbyte-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"strconv"
)

func (r *AirbyteReconciler) makePodSweeperDeployment(instance *stackv1alpha1.Airbyte) ([]*appsv1.Deployment, error) {
	var deployments []*appsv1.Deployment

	if instance.Spec.PodSweeper.RoleGroups != nil {
		for roleGroupName, roleGroup := range instance.Spec.PodSweeper.RoleGroups {
			dep := r.makePodSweeperDeploymentForRoleGroup(instance, roleGroupName, roleGroup, r.Scheme)
			if dep != nil {
				deployments = append(deployments, dep)
			}
		}
	}

	return deployments, nil
}

func (r *AirbyteReconciler) makePodSweeperDeploymentForRoleGroup(instance *stackv1alpha1.Airbyte, roleGroupName string, roleGroup *stackv1alpha1.RoleGroupPodSweeperSpec, schema *runtime.Scheme) *appsv1.Deployment {
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
		image = *instance.Spec.PodSweeper.Image
	}

	if roleGroup != nil && roleGroup.Config != nil && roleGroup.Config.SecurityContext != nil {
		securityContext = roleGroup.Config.SecurityContext
	} else {
		securityContext = instance.Spec.PodSweeper.SecurityContext
	}

	var envVars []corev1.EnvVar

	envVars = append(envVars, corev1.EnvVar{
		Name: "KUBE_NAMESPACE",
		ValueFrom: &corev1.EnvVarSource{
			FieldRef: &corev1.ObjectFieldSelector{
				FieldPath: "metadata.namespace",
			},
		},
	})

	if roleGroup != nil && roleGroup.Config != nil && roleGroup.Config.TimeToDeletePods != nil {
		envVars = append(envVars, corev1.EnvVar{
			Name:  "RUNNING_TTL_MINUTES",
			Value: roleGroup.Config.TimeToDeletePods.Running,
		})
		envVars = append(envVars, corev1.EnvVar{
			Name:  "SUCCEEDED_TTL_MINUTES",
			Value: strconv.Itoa(int(roleGroup.Config.TimeToDeletePods.Succeeded)),
		})
		envVars = append(envVars, corev1.EnvVar{
			Name:  "UNSUCCESSFUL_TTL_MINUTES",
			Value: strconv.Itoa(int(roleGroup.Config.TimeToDeletePods.Unsuccessful)),
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

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.GetNameWithSuffix("-pod-sweeper" + "-" + roleGroupName),
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
							Name:            "airbyte-pod-sweeper",
							Image:           image.Repository + ":" + image.Tag,
							ImagePullPolicy: image.PullPolicy,
							Resources:       *roleGroup.Config.Resources,
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

	container := &dep.Spec.Template.Spec.Containers[0]
	if roleGroup.Config != nil && len(roleGroup.Config.ExtraVolumeMounts) > 0 {
		for _, volumeMount := range roleGroup.Config.ExtraVolumeMounts {
			container.VolumeMounts = append(container.VolumeMounts, volumeMount)
		}
	} else if instance.Spec.PodSweeper != nil && len(instance.Spec.PodSweeper.ExtraVolumeMounts) > 0 {
		for _, volumeMount := range instance.Spec.PodSweeper.ExtraVolumeMounts {
			container.VolumeMounts = append(container.VolumeMounts, volumeMount)
		}
	}

	podSpec := &dep.Spec.Template.Spec
	if roleGroup.Config != nil && len(roleGroup.Config.ExtraVolumes) > 0 {
		for _, volume := range roleGroup.Config.ExtraVolumes {
			podSpec.Volumes = append(podSpec.Volumes, volume)
		}
	} else if instance.Spec.PodSweeper.RoleConfig != nil && len(instance.Spec.PodSweeper.ExtraVolumes) > 0 {
		for _, volume := range instance.Spec.PodSweeper.ExtraVolumes {
			podSpec.Volumes = append(podSpec.Volumes, volume)
		}
	}

	if roleGroup.Config != nil && roleGroup.Config.LivenessProbe != nil && roleGroup.Config.LivenessProbe.Enabled {
		dep.Spec.Template.Spec.Containers[0].LivenessProbe = createProbe(roleGroup.Config.LivenessProbe)
	} else if instance.Spec.PodSweeper.LivenessProbe != nil && instance.Spec.PodSweeper.LivenessProbe.Enabled {
		dep.Spec.Template.Spec.Containers[0].LivenessProbe = createProbe(instance.Spec.PodSweeper.LivenessProbe)
	}

	if roleGroup.Config != nil && roleGroup.Config.ReadinessProbe != nil && roleGroup.Config.ReadinessProbe.Enabled {
		dep.Spec.Template.Spec.Containers[0].ReadinessProbe = createProbe(roleGroup.Config.ReadinessProbe)
	} else if instance.Spec.PodSweeper.ReadinessProbe != nil && instance.Spec.PodSweeper.ReadinessProbe.Enabled {
		dep.Spec.Template.Spec.Containers[0].ReadinessProbe = createProbe(instance.Spec.PodSweeper.ReadinessProbe)
	}

	err := ctrl.SetControllerReference(instance, dep, schema)

	if err != nil {
		r.Log.Error(err, "Failed to set controller reference for deployment")
		return nil
	}
	return dep
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

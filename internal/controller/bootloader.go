package controller

import (
	"context"
	stackv1alpha1 "github.com/zncdatadev/airbyte-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// reconcile Bootloader pod
func (r *AirbyteReconciler) reconcileBootloaderPod(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	reconcileParams := r.createReconcileParams(ctx, nil, instance, r.extractBootloaderPod)
	if err := reconcileParams.createOrUpdateResourceNoGroup(); err != nil {
		return err
	}
	return nil
}

// extract bootloader pod for pre-install and pre-update
func (r *AirbyteReconciler) extractBootloaderPod(params ExtractorParams) (client.Object, error) {
	instance := params.instance.(*stackv1alpha1.Airbyte)
	instanceLabels := instance.Labels

	//copy labels of instance for pod
	var labels = make(map[string]string)
	for k, v := range instanceLabels {
		labels[k] = v
	}

	labels["app.kubernetes.io/name"] = labels["app.kubernetes.io/name"] + "-" + roleNameMapper[Bootloader]
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      createNameForRoleGroup(instance, "bootloader", ""),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            createNameForRoleGroup(instance, "bootloader-container", ""),
					Image:           "airbyte/bootloader:0.50.45",
					ImagePullPolicy: corev1.PullIfNotPresent,
					Env:             r.mergeEnvVars(instance),
				},
			},
		},
	}
	return pod, nil
}

func (r *AirbyteReconciler) mergeEnvVars(instance *stackv1alpha1.Airbyte) []corev1.EnvVar {
	var envVars []corev1.EnvVar
	// AIRBYTE_VERSION
	envVars = append(envVars, corev1.EnvVar{
		Name: "AIRBYTE_VERSION",
		ValueFrom: &corev1.EnvVarSource{
			ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
				Key: "AIRBYTE_VERSION",
				LocalObjectReference: corev1.LocalObjectReference{
					Name: createNameForRoleGroup(instance, "env", ""),
				},
			},
		},
	})
	//DATABASE_HOST
	envVars = append(envVars, corev1.EnvVar{
		Name: "DATABASE_HOST",
		ValueFrom: &corev1.EnvVarSource{
			ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
				Key: "DATABASE_HOST",
				LocalObjectReference: corev1.LocalObjectReference{
					Name: createNameForRoleGroup(instance, "env", ""),
				},
			},
		},
	})
	//DATABASE_PORT
	envVars = append(envVars, corev1.EnvVar{
		Name: "DATABASE_PORT",
		ValueFrom: &corev1.EnvVarSource{
			ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
				Key: "DATABASE_PORT",
				LocalObjectReference: corev1.LocalObjectReference{
					Name: createNameForRoleGroup(instance, "env", ""),
				},
			},
		},
	})
	//DATABASE_USER
	envVars = append(envVars, corev1.EnvVar{
		Name: "DATABASE_USER",
		ValueFrom: &corev1.EnvVarSource{
			SecretKeyRef: &corev1.SecretKeySelector{
				Key: "DATABASE_USER",
				LocalObjectReference: corev1.LocalObjectReference{
					Name: createNameForRoleGroup(instance, "secrets", ""),
				},
			},
		},
	})
	//DATABASE_PASSWORD
	envVars = append(envVars, corev1.EnvVar{
		Name: "DATABASE_PASSWORD",
		ValueFrom: &corev1.EnvVarSource{
			SecretKeyRef: &corev1.SecretKeySelector{
				Key: "DATABASE_PASSWORD",
				LocalObjectReference: corev1.LocalObjectReference{
					Name: createNameForRoleGroup(instance, "secrets", ""),
				},
			},
		},
	})
	//DATABASE_URL
	envVars = append(envVars, corev1.EnvVar{
		Name: "DATABASE_URL",
		ValueFrom: &corev1.EnvVarSource{
			ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
				Key: "DATABASE_URL",
				LocalObjectReference: corev1.LocalObjectReference{
					Name: createNameForRoleGroup(instance, "env", ""),
				},
			},
		},
	})
	return envVars

}

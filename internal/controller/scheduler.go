package controller

import (
	stackv1alpha1 "github.com/zncdatadev/airbyte-operator/api/v1alpha1"
	"github.com/zncdatadev/airbyte-operator/api/v1alpha1/rolegroup"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

//server

func ServerScheduler(instance *stackv1alpha1.Airbyte, dep *appsv1.Deployment, group rolegroup.RoleConfigObject) {
	scheduleForRole(instance, group, dep, func() rolegroup.RoleConfigObject { return instance.Spec.Server.RoleConfig })
}

//worker

func WorkerScheduler(instance *stackv1alpha1.Airbyte, dep *appsv1.Deployment, group rolegroup.RoleConfigObject) {
	scheduleForRole(instance, group, dep, func() rolegroup.RoleConfigObject { return instance.Spec.Worker.RoleConfig })
}

//api server

func ApiServerScheduler(instance *stackv1alpha1.Airbyte, dep *appsv1.Deployment, group rolegroup.RoleConfigObject) {
	scheduleForRole(instance, group, dep, func() rolegroup.RoleConfigObject { return instance.Spec.ApiServer.RoleConfig })
}

//connector builder server

func ConnectorBuilderServerScheduler(instance *stackv1alpha1.Airbyte, dep *appsv1.Deployment, group rolegroup.RoleConfigObject) {
	scheduleForRole(instance, group, dep, func() rolegroup.RoleConfigObject { return instance.Spec.ConnectorBuilderServer.RoleConfig })
}

//cron

func CronScheduler(instance *stackv1alpha1.Airbyte, dep *appsv1.Deployment, group rolegroup.RoleConfigObject) {
	scheduleForRole(instance, group, dep, func() rolegroup.RoleConfigObject { return instance.Spec.Cron.RoleConfig })
}

func scheduleForRole(instance *stackv1alpha1.Airbyte, group rolegroup.RoleConfigObject, dep *appsv1.Deployment,
	roleExtractor func() rolegroup.RoleConfigObject) {
	cluster := instance.Spec.ClusterConfig
	role := roleExtractor()
	affinity, toleration, nodeSelector := getScheduleInfo(group, role, cluster)

	if nodeSelector != nil {
		dep.Spec.Template.Spec.NodeSelector = nodeSelector
	}

	if toleration != nil {
		dep.Spec.Template.Spec.Tolerations = []corev1.Toleration{
			*toleration.DeepCopy(),
		}
	}

	if affinity != nil {
		dep.Spec.Template.Spec.Affinity = affinity.DeepCopy()
	}
}

package rolegroup

import (
	"github.com/zncdatadev/airbyte-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

type RoleConfigObject interface {
	GetReplicas() int32
	GetImage() *v1alpha1.ImageSpec
	GetLabels() map[string]string
	GetSecurityContext() *corev1.PodSecurityContext
	GetMatchLabels() map[string]string
	GetAffinity() *corev1.Affinity
	GetNodeSelector() map[string]string
	GetTolerations() *corev1.Toleration
	GetService() *v1alpha1.ServiceSpec
}

type ConfigObject interface {
	GetConfig() any
	GetResource() *v1alpha1.ResourcesSpec
}

// todo: add BaseRoleConfigSpec and BaseConfigSpec into here that will occur a circle import

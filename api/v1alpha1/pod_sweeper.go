package v1alpha1

import corev1 "k8s.io/api/core/v1"

type PodSweeperSpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:default=true
	Enabled bool `json:"enabled,omitempty"`

	// +kubebuilder:validation:Optional
	RoleConfig *PodSweeperRoleConfigSpec `json:"roleConfig,omitempty"`

	// +kubebuilder:validation:Optional
	RoleGroups map[string]*PodSweeperRoleConfigSpec `json:"roleGroups,omitempty"`
}

type PodSweeperRoleConfigSpec struct {
	GenericRoleConfigSpec `json:",inline"`

	// +kubebuilder:validation:Optional
	LivenessProbe *ProbeSpec `json:"livenessProbe,omitempty"`

	// +kubebuilder:validation:Optional
	ReadinessProbe *ProbeSpec `json:"readinessProbe,omitempty"`

	// +kubebuilder:validation:Optional
	ExtraVolumes []corev1.Volume `json:"extraVolumes,omitempty"`

	// +kubebuilder:validation:Optional
	ExtraVolumeMounts []corev1.VolumeMount `json:"extraVolumeMounts,omitempty"`

	// +kubebuilder:validation:Optional
	TimeToDeletePods *PodSweeperTimeToDeletePodsSpec `json:"timeToDeletePods,omitempty"`
}

type PodSweeperTimeToDeletePodsSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	Running string `json:"running,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=120
	Succeeded int32 `json:"succeeded,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1440
	Unsuccessful int32 `json:"unsuccessful,omitempty"`
}

type ProbeSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=true
	Enabled bool `json:"enabled,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=5
	InitialDelaySeconds int32 `json:"initialDelaySeconds,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=30
	PeriodSeconds int32 `json:"periodSeconds,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	TimeoutSeconds int32 `json:"timeoutSeconds,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=3
	FailureThreshold int32 `json:"failureThreshold,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	SuccessThreshold int32 `json:"successThreshold,omitempty"`
}

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
)

type TemporalSpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:default=true
	Enabled bool `json:"enabled,omitempty"`

	// +kubebuilder:validation:Optional
	RoleConfig *TemporalRoleConfigSpec `json:"roleConfig,omitempty"`

	// +kubebuilder:validation:Optional
	RoleGroups map[string]*TemporalRoleConfigSpec `json:"roleGroups,omitempty"`
}

type TemporalRoleConfigSpec struct {
	BaseRoleConfigSpec `json:",inline"`

	// +kubebuilder:validation:Optional
	Config TemporalConfigSpec `json:"config"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default={}
	Secret map[string]string `json:"secret,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=INFO
	LogLevel string `json:"logLevel,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default={}
	EnvVars map[string]string `json:"envVars,omitempty"`

	// +kubebuilder:validation:Optional
	ExtraEnv corev1.EnvVar `json:"extraEnv,omitempty"`
}

type TemporalConfigSpec struct {
	BaseConfigSpec `json:",inline"`

	// +kubebuilder:validation:Optional
	S3 *S3Spec `json:"minio,omitempty"`

	Database *DatabaseSpec `json:"database,omitempty"`
}

func (t *TemporalRoleConfigSpec) GetConfig() any {
	return t.Config
}

func (t *TemporalRoleConfigSpec) GetResource() *ResourcesSpec {
	return t.Config.Resources
}

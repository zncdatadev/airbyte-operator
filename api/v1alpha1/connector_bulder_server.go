package v1alpha1

import corev1 "k8s.io/api/core/v1"

type ConnectorBuilderServerSpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:default=true
	Enabled bool `json:"enabled,omitempty"`

	// +kubebuilder:validation:Optional
	RoleConfig *ConnectorBuilderServerRoleConfigSpec `json:"roleConfig,omitempty"`

	// +kubebuilder:validation:Optional
	RoleGroups map[string]*ConnectorBuilderServerRoleConfigSpec `json:"roleGroups,omitempty"`
}

type ConnectorBuilderServerRoleConfigSpec struct {
	GenericRoleConfigSpec `json:",inline"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default={}
	Secret map[string]string `json:"secret,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=INFO
	LogLevel string `json:"logLevel,omitempty"`

	// +kubebuilder:validation:Optional
	Debug *DebugSpec `json:"debug,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default={}
	EnvVars map[string]string `json:"envVars,omitempty"`

	// +kubebuilder:validation:Optional
	ExtraEnv *corev1.EnvVar `json:"extraEnv,omitempty"`
}

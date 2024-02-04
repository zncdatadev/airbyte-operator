package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
)

type WebAppSpec struct {

	// +kubebuilder:validation:Required
	// +kubebuilder:default=true
	Enabled bool `json:"enabled,omitempty"`

	// +kubebuilder:validation:Optional
	RoleConfig *WebAppRoleConfigSpec `json:"roleConfig,omitempty"`

	// +kubebuilder:validation:Optional
	RoleGroups map[string]*WebAppRoleConfigSpec `json:"roleGroups,omitempty"`
}

type WebAppRoleConfigSpec struct {
	IngresRoleConfigSpec `json:",inline"`

	// +kubebuilder:validation:Optional
	Config *BaseConfigSpec `json:"config"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default={}
	Secret *map[string]string `json:"secret,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="/api/v1/"
	ApiUrl string `json:"apiUrl,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="/connector-builder-api"
	ConnectorBuilderServerUrl string `json:"connectorBuilderServerUrl,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	Url string `json:"url,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=INFO
	LogLevel string `json:"logLevel,omitempty"`

	// +kubebuilder:validation:Optional
	EnvVars map[string]string `json:"envVars,omitempty"`

	// +kubebuilder:validation:Optional
	ExtraEnv corev1.EnvVar `json:"extraEnv,omitempty"`
}

func (w *WebAppRoleConfigSpec) GetConfig() any {
	return w.Config
}

func (w *WebAppRoleConfigSpec) GetResource() *ResourcesSpec {
	return w.Config.Resources
}

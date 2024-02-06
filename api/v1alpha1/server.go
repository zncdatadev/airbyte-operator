package v1alpha1

type ServerSpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:default=true
	Enabled bool `json:"enabled,omitempty"`

	// +kubebuilder:validation:Optional
	RoleConfig *ServerRoleConfigSpec `json:"roleConfig,omitempty"`

	// +kubebuilder:validation:Optional
	RoleGroups map[string]*ServerRoleConfigSpec `json:"roleGroups,omitempty"`
}

type ServerRoleConfigSpec struct {
	GenericRoleConfigSpec `json:",inline"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default={}
	Secret map[string]string `json:"secret,omitempty"`
}

package v1alpha1

type KeycloakSpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:default=false
	Enabled bool `json:"enabled,omitempty"`

	// +kubebuilder:validation:Optional
	RoleConfig *KeycloakRoleConfigSpec `json:"roleConfig,omitempty"`

	// +kubebuilder:validation:Optional
	RoleGroups map[string]*KeycloakRoleConfigSpec `json:"roleGroups,omitempty"`
}

type KeycloakRoleConfigSpec struct {
	GenericRoleConfigSpec `json:",inline"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=INFO
	LogLevel string `json:"logLevel,omitempty"`
}

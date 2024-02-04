package v1alpha1

const (
	AirbyteBootloaderRepo     = "airbyte/airbyte-bootloader"
	AirbyteBootloaderTag      = "0.50.30"
	AirbyteBootloaderPort     = 80
	AirbyteBootloaderReplicas = 1
)

type AirbyteBootloaderSpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:default=true
	Enabled bool `json:"enabled,omitempty"`

	// +kubebuilder:validation:Optional
	RoleConfig *BootloaderRoleConfigSpec `json:"roleConfig,omitempty"`

	// +kubebuilder:validation:Optional
	RoleGroups map[string]*BootloaderRoleConfigSpec `json:"roleGroups,omitempty"`
}

type BootloaderRoleConfigSpec struct {
	GenericRoleConfigSpec `json:",inline"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=INFO
	LogLevel string `json:"logLevel,omitempty"`

	// +kubebuilder:validation:Optional
	RunDatabaseMigrationsOnStartup *bool `json:"runDatabaseMigrationsOnStartup,omitempty"`
}

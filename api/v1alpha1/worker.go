package v1alpha1

type WorkerSpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:default=true
	Enabled bool `json:"enabled,omitempty"`

	// +kubebuilder:validation:Optional
	RoleConfig *WorkerRoleConfigSpec `json:"roleConfig,omitempty"`

	// +kubebuilder:validation:Optional
	RoleGroups map[string]*WorkerRoleConfigSpec `json:"roleGroups,omitempty"`
}

type WorkerRoleConfigSpec struct {
	GenericRoleConfigSpec `json:",inline"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=INFO
	LogLevel string `json:"logLevel,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	Debug bool `json:"debug,omitempty"`

	// +kubebuilder:validation:Optional
	ContainerOrchestrator *ContainerOrchestratorSpec `json:"containerOrchestrator,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	ActivityMaxAttempt string `json:"activityMaxAttempt,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	ActivityInitialDelayBetweenAttemptsSeconds string `json:"activityInitialDelayBetweenAttemptsSeconds,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	ActivityMaxDelayBetweenAttemptsSeconds string `json:"activityMaxDelayBetweenAttemptsSeconds,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="5"
	MaxNotifyWorkers string `json:"maxNotifyWorkers,omitempty"`
}

type ContainerOrchestratorSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=true
	Enabled bool `json:"enabled,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	Image string `json:"image,omitempty"`
}

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
)

type ClusterRoleConfigSpec struct {
	BaseRoleConfigSpec `json:",inline"`

	// +kubebuilder:validation:Optional
	Config ClusterConfigSpec `json:"config"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="airbyte-admin"
	ServiceAccountName string `json:"serviceAccountName,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="oss"
	DeploymentMode string `json:"deploymentMode,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default={}
	EnvVars map[string]string `json:"envVars,omitempty"`

	//// +kubebuilder:validation:Optional
	//Database *DatabaseClusterConfigSpec `json:"database,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="community"
	Edition string `json:"edition,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="S3"
	StateStorageType string `json:"stateStorageType,omitempty"`

	// +kubebuilder:validation:Optional
	StateStorage *StateStorageClusterConfigSpec `json:"stateStorage,omitempty"`

	// +kubebuilder:validation:Optional
	Logs *LogsClusterConfigSpec `json:"logs,omitempty"`

	// +kubebuilder:validation:Optional
	Metrics *MetricsClusterConfigSpec `json:"metrics,omitempty"`

	// +kubebuilder:validation:Optional
	Jobs *JobsClusterConfigSpec `json:"jobs,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=INFO
	LogLevel string `json:"logLevel,omitempty"`

	// +kubebuilder:validation:Optional
	RunDatabaseMigrationsOnStartup *bool `json:"runDatabaseMigrationsOnStartup,omitempty"`

	// +kubebuilder:validation:Optional
	ExtraEnv *corev1.EnvVar `json:"extraEnv,omitempty"`
}

type StateStorageClusterConfigSpec struct {
	// +kubebuilder:validation:Optional
	S3 *S3Spec `json:"s3,omitempty"`
}

type ClusterConfigSpec struct {
	GenericRoleConfigSpec `json:",inline"`

	// +kubebuilder:validation:Required
	Database *DatabaseSpec `json:"database,omitempty"`
}

type MetricsClusterConfigSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	MetricClient string `json:"metricClient,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	OtelCollectorEndpoint string `json:"otelCollectorEndpoint,omitempty"`
}

type JobsClusterConfigSpec struct {
	// +kubebuilder:validation:Optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// +kubebuilder:validation:Optional
	Kube *JobsKubeClusterConfigSpec `json:"kube,omitempty"`
}

type JobsKubeClusterConfigSpec struct {
	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels,omitempty"`

	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// +kubebuilder:validation:Optional
	Tolerations *corev1.Toleration `json:"tolerations,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	MainContainerImagePullSecret string `json:"mainContainerImagePullSecret,omitempty"`

	// +kubebuilder:validation:Optional
	Images *JobsKubeImagesClusterConfigSpec `json:"images,omitempty"`
}

type JobsKubeImagesClusterConfigSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	Busybox string `json:"busybox,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	Socat string `json:"socat,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	Curl string `json:"curl,omitempty"`
}

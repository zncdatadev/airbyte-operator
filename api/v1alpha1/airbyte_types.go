/*
Copyright 2023 zncdata-labs.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AirbyteSpec defines the desired state of Airbyte
type AirbyteSpec struct {
	// +kubebuilder:validation:Required
	SecurityContext *corev1.PodSecurityContext `json:"securityContext"`

	// +kubebuilder:validation:Optional
	Secret SecretParam `json:"secret"`

	// +kubebuilder:validation:Optional
	Global *GlobalSpec `json:"global"`

	// +kubebuilder:validation:Optional
	Postgres *PostgresSpec `json:"postgres"`

	// +kubebuilder:validation:Optional
	Server *ServerSpec `json:"server"`

	// +kubebuilder:validation:Optional
	Worker *WorkerSpec `json:"worker"`

	// +kubebuilder:validation:Optional
	AirbyteApiServer *AirbyteApiServerSpec `json:"airbyteApiServer"`

	// +kubebuilder:validation:Optional
	WebApp *WebAppSpec `json:"webApp"`

	// +kubebuilder:validation:Optional
	PodSweeper *PodSweeperSpec `json:"podSweeper"`

	// +kubebuilder:validation:Optional
	ConnectorBuilderServer *ConnectorBuilderServerSpec `json:"connectorBuilderServer"`

	// +kubebuilder:validation:Optional
	AirbyteBootloader *AirbyteBootloaderSpec `json:"airbyteBootloader"`

	// +kubebuilder:validation:Optional
	Temporal *TemporalSpec `json:"temporal"`

	// +kubebuilder:validation:Optional
	Keycloak *KeycloakSpec `json:"keycloak"`

	// +kubebuilder:validation:Optional
	Cron *CronSpec `json:"cron"`
}

type GlobalSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default="airbyte-admin"
	ServiceAccountName string `json:"serviceAccountName"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="oss"
	DeploymentMode string `json:"deploymentMode"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default={}
	EnvVars map[string]string `json:"envVars"`

	// +kubebuilder:validation:Optional
	Database *GlobalDatabaseSpec `json:"database"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="community"
	Edition string `json:"edition"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="MINIO"
	StateStorageType string `json:"stateStorageType"`

	// +kubebuilder:validation:Optional
	Log *GlobalLogSpec `json:"log"`

	// +kubebuilder:validation:Optional
	Metrics *GlobalMetricsSpec `json:"metrics"`

	// +kubebuilder:validation:Optional
	Jobs *GlobalJobsSpec `json:"jobs"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	ConfigMapName string `json:"configMapName"`
}

type GlobalDatabaseSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	SecretName string `json:"secretName"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	SecretValue string `json:"secretValue"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="example.com"
	Host string `json:"host"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=5432
	Port int32 `json:"port"`
}

type GlobalJobsSpec struct {
	// +kubebuilder:validation:Optional
	Resources *corev1.ResourceRequirements `json:"resources"`

	// +kubebuilder:validation:Optional
	Kube *GlobalJobsKubeSpec `json:"kube"`
}

type GlobalJobsKubeSpec struct {
	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations"`

	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels"`

	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector"`

	// +kubebuilder:validation:Optional
	Tolerations *corev1.Toleration `json:"tolerations"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	MainContainerImagePullSecret string `json:"mainContainerImagePullSecret"`

	// +kubebuilder:validation:Optional
	Images *GlobalJobsKubeImagesSpec `json:"images"`
}

type GlobalJobsKubeImagesSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	Busybox string `json:"busybox"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	Socat string `json:"socat"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	Curl string `json:"curl"`
}

type GlobalLogSpec struct {
	// +kubebuilder:validation:Optional
	Minio *GlobalLogMinioSpec `json:"minio"`

	// +kubebuilder:validation:Optional
	ExternalMinio *GlobalLogExternalMinioSpec `json:"externalMinio"`

	// +kubebuilder:validation:Optional
	S3 *GlobalLogS3Spec `json:"s3"`

	// +kubebuilder:validation:Optional
	Gcs *GlobalLogGcsSpec `json:"gcs"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="MINIO"
	StorageType string `json:"storageType"`
}

type GlobalMetricsSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	MetricClient string `json:"metricClient"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	OtelCollectorEndpoint string `json:"otelCollectorEndpoint"`
}

type GlobalLogMinioSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=true
	Enabled bool `json:"enabled"`
}

type GlobalLogExternalMinioSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	Enabled bool `json:"enabled"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	Endpoint string `json:"endpoint"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="localhost"
	Host string `json:"host"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=9000
	Port int32 `json:"port"`
}

type GlobalLogS3Spec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:default=false
	Enabled bool `json:"enabled"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="airbyte-dev-logs"
	Bucket string `json:"bucket"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	BucketRegion string `json:"bucketRegion"`
}

type GlobalLogGcsSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	Bucket string `json:"bucket"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	Credentials string `json:"credentials"`
}

type PostgresSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default="postgresql"
	Host string `json:"host"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default="5432"
	Port string `json:"port"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default="airbyte"
	UserName string `json:"username"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default="airbyte"
	Password string `json:"password"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default="db-airbyte"
	DataBase string `json:"database"`
}

type ServerSpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:default=true
	Enabled bool `json:"enabled"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas"`

	// +kubebuilder:validation:Required
	Image *ServerImageSpec `json:"image"`

	// +kubebuilder:validation:Required
	Resources *corev1.ResourceRequirements `json:"resources"`

	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations"`

	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels"`

	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector"`

	// +kubebuilder:validation:Optional
	Tolerations *corev1.Toleration `json:"tolerations"`

	// +kubebuilder:validation:Optional
	Affinity *corev1.Affinity `json:"affinity"`

	// +kubebuilder:validation:Optional
	Secret SecretParam `json:"secret,omitempty"`

	// +kubebuilder:validation:Optional
	Service *ServerServiceSpec `json:"service"`
}

type WorkerSpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:default=true
	Enabled bool `json:"enabled"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas"`

	// +kubebuilder:validation:Required
	Image *WorkerImageSpec `json:"image"`

	// +kubebuilder:validation:Required
	Resources *corev1.ResourceRequirements `json:"resources"`

	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations"`

	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels"`

	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector"`

	// +kubebuilder:validation:Optional
	Tolerations *corev1.Toleration `json:"tolerations"`

	// +kubebuilder:validation:Optional
	Affinity *corev1.Affinity `json:"affinity"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=INFO
	LogLevel string `json:"logLevel"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	Debug bool `json:"debug"`

	// +kubebuilder:validation:Optional
	ContainerOrchestrator *ContainerOrchestratorSpec `json:"containerOrchestrator"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	ActivityMaxAttempt string `json:"activityMaxAttempt"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	ActivityInitialDelayBetweenAttemptsSeconds string `json:"activityInitialDelayBetweenAttemptsSeconds"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	ActivityMaxDelayBetweenAttemptsSeconds string `json:"activityMaxDelayBetweenAttemptsSeconds"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="5"
	MaxNotifyWorkers string `json:"maxNotifyWorkers"`
}

type ContainerOrchestratorSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=true
	Enabled bool `json:"enabled"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	Image string `json:"image"`
}

type AirbyteApiServerSpec struct {

	// +kubebuilder:validation:Optional
	EnvVars map[string]string `json:"envVars"`

	// +kubebuilder:validation:Required
	// +kubebuilder:default=true
	Enabled bool `json:"enabled"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas"`

	// +kubebuilder:validation:Required
	Image *AirbyteApiServerImageSpec `json:"image"`

	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations"`

	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels"`

	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector"`

	// +kubebuilder:validation:Optional
	Tolerations *corev1.Toleration `json:"tolerations"`

	// +kubebuilder:validation:Optional
	Affinity *corev1.Affinity `json:"affinity"`

	// +kubebuilder:validation:Optional
	Resources *corev1.ResourceRequirements `json:"resources"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=INFO
	LogLevel string `json:"logLevel"`

	// +kubebuilder:validation:Optional
	Service *AirbyteApiServerServiceSpec `json:"service"`

	// +kubebuilder:validation:Optional
	Debug *AirbyteApiServerDebugSpec `json:"debug"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default={}
	Secret map[string]string `json:"secret"`

	// +kubebuilder:validation:Optional
	ExtraEnv corev1.EnvVar `json:"extraEnv"`
}

type WebAppSpec struct {

	// +kubebuilder:validation:Required
	// +kubebuilder:default=true
	Enabled bool `json:"enabled"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas"`

	// +kubebuilder:validation:Required
	Image *WebAppImageSpec `json:"image"`

	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations"`

	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels"`

	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector"`

	// +kubebuilder:validation:Optional
	Tolerations *corev1.Toleration `json:"tolerations"`

	// +kubebuilder:validation:Optional
	Affinity *corev1.Affinity `json:"affinity"`

	// +kubebuilder:validation:Optional
	Resources *corev1.ResourceRequirements `json:"resources"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=INFO
	LogLevel string `json:"logLevel"`

	// +kubebuilder:validation:Optional
	Service *WebAppServiceSpec `json:"service"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="/api/v1/"
	ApiUrl string `json:"apiUrl"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="/connector-builder-api"
	ConnectorBuilderServerUrl string `json:"connectorBuilderServerUrl"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	Url string `json:"url"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default={}
	Secret map[string]string `json:"secret"`

	// +kubebuilder:validation:Optional
	EnvVars map[string]string `json:"envVars"`

	// +kubebuilder:validation:Optional
	ExtraEnv corev1.EnvVar `json:"extraEnv"`
}

type PodSweeperSpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:default=true
	Enabled bool `json:"enabled"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas"`

	// +kubebuilder:validation:Required
	Image *PodSweeperImageSpec `json:"image"`

	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations"`

	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels"`

	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector"`

	// +kubebuilder:validation:Optional
	Tolerations *corev1.Toleration `json:"tolerations"`

	// +kubebuilder:validation:Optional
	Affinity *corev1.Affinity `json:"affinity"`

	// +kubebuilder:validation:Optional
	Resources *corev1.ResourceRequirements `json:"resources"`

	// +kubebuilder:validation:Optional
	TimeToDeletePods *PodSweeperTimeToDeletePodsSpec `json:"timeToDeletePods"`

	// +kubebuilder:validation:Optional
	LivenessProbe *PodSweeperLivenessProbeSpec `json:"livenessProbe"`

	// +kubebuilder:validation:Optional
	ReadinessProbe *PodSweeperReadinessProbeSpec `json:"readinessProbe"`
}

type PodSweeperLivenessProbeSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=true
	Enabled bool `json:"enabled"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=5
	InitialDelaySeconds int32 `json:"initialDelaySeconds"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=30
	PeriodSeconds int32 `json:"periodSeconds"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	TimeoutSeconds int32 `json:"timeoutSeconds"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=3
	FailureThreshold int32 `json:"failureThreshold"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	SuccessThreshold int32 `json:"successThreshold"`
}

type PodSweeperReadinessProbeSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=true
	Enabled bool `json:"enabled"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=5
	InitialDelaySeconds int32 `json:"initialDelaySeconds"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=30
	PeriodSeconds int32 `json:"periodSeconds"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	TimeoutSeconds int32 `json:"timeoutSeconds"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=3
	FailureThreshold int32 `json:"failureThreshold"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	SuccessThreshold int32 `json:"successThreshold"`
}

type PodSweeperTimeToDeletePodsSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	Running string `json:"running"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=120
	Succeeded int32 `json:"succeeded"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1440
	Unsuccessful int32 `json:"unsuccessful"`
}

type ConnectorBuilderServerSpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:default=true
	Enabled bool `json:"enabled"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas"`

	// +kubebuilder:validation:Required
	Image *ConnectorBuilderServerImageSpec `json:"image"`

	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations"`

	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels"`

	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector"`

	// +kubebuilder:validation:Optional
	Tolerations *corev1.Toleration `json:"tolerations"`

	// +kubebuilder:validation:Optional
	Affinity *corev1.Affinity `json:"affinity"`

	// +kubebuilder:validation:Optional
	Resources *corev1.ResourceRequirements `json:"resources"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=INFO
	LogLevel string `json:"logLevel"`

	// +kubebuilder:validation:Optional
	Service *ConnectorBuilderServerServiceSpec `json:"service"`

	// +kubebuilder:validation:Optional
	Debug *ConnectorBuilderServerDebugSpec `json:"debug"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default={}
	Secret map[string]string `json:"secret"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default={}
	EnvVars map[string]string `json:"envVars"`

	// +kubebuilder:validation:Optional
	ExtraEnv corev1.EnvVar `json:"extraEnv"`
}

type AirbyteBootloaderSpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:default=true
	Enabled bool `json:"enabled"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas"`

	// +kubebuilder:validation:Required
	Image *AirbyteBootloaderImageSpec `json:"image"`

	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations"`

	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels"`

	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector"`

	// +kubebuilder:validation:Optional
	Tolerations *corev1.Toleration `json:"tolerations"`

	// +kubebuilder:validation:Optional
	Affinity *corev1.Affinity `json:"affinity"`

	// +kubebuilder:validation:Optional
	Resources *corev1.ResourceRequirements `json:"resources"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=INFO
	LogLevel string `json:"logLevel"`

	// +kubebuilder:validation:Optional
	Service *AirbyteBootloaderServiceSpec `json:"service"`

	// +kubebuilder:validation:Optional
	RunDatabaseMigrationsOnStartup *bool `json:"runDatabaseMigrationsOnStartup"`
}

type TemporalSpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:default=true
	Enabled bool `json:"enabled"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas"`

	// +kubebuilder:validation:Required
	Image *TemporalImageSpec `json:"image"`

	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations"`

	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels"`

	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector"`

	// +kubebuilder:validation:Optional
	Tolerations *corev1.Toleration `json:"tolerations"`

	// +kubebuilder:validation:Optional
	Affinity *corev1.Affinity `json:"affinity"`

	// +kubebuilder:validation:Optional
	Resources *corev1.ResourceRequirements `json:"resources"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=INFO
	LogLevel string `json:"logLevel"`

	// +kubebuilder:validation:Optional
	Service *TemporalServiceSpec `json:"service"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default={}
	Secret map[string]string `json:"secret"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default={}
	EnvVars map[string]string `json:"envVars"`

	// +kubebuilder:validation:Optional
	ExtraEnv corev1.EnvVar `json:"extraEnv"`
}

type KeycloakSpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:default=false
	Enabled bool `json:"enabled"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas"`

	// +kubebuilder:validation:Required
	Image *KeycloakImageSpec `json:"image"`

	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations"`

	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels"`

	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector"`

	// +kubebuilder:validation:Optional
	Tolerations *corev1.Toleration `json:"tolerations"`

	// +kubebuilder:validation:Optional
	Affinity *corev1.Affinity `json:"affinity"`

	// +kubebuilder:validation:Optional
	Resources *corev1.ResourceRequirements `json:"resources"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=INFO
	LogLevel string `json:"logLevel"`

	// +kubebuilder:validation:Optional
	Service *KeycloakServiceSpec `json:"service"`
}

type CronSpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:default=true
	Enabled bool `json:"enabled"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas"`

	// +kubebuilder:validation:Required
	Image *CronImageSpec `json:"image"`

	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations"`

	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels"`

	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector"`

	// +kubebuilder:validation:Optional
	Tolerations *corev1.Toleration `json:"tolerations"`

	// +kubebuilder:validation:Optional
	Affinity *corev1.Affinity `json:"affinity"`

	// +kubebuilder:validation:Optional
	Resources *corev1.ResourceRequirements `json:"resources"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=INFO
	LogLevel string `json:"logLevel"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default={}
	Secret map[string]string `json:"secret"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default={}
	EnvVars map[string]string `json:"envVars"`

	// +kubebuilder:validation:Optional
	ExtraEnv corev1.EnvVar `json:"extraEnv"`
}

type ServerImageSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=airbyte/server
	Repository string `json:"repository,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="0.50.30"
	Tag string `json:"tag,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=IfNotPresent
	PullPolicy corev1.PullPolicy `json:"pullPolicy,omitempty"`
}

type WorkerImageSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=airbyte/worker
	Repository string `json:"repository,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="0.50.30"
	Tag string `json:"tag,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=IfNotPresent
	PullPolicy corev1.PullPolicy `json:"pullPolicy,omitempty"`
}

type AirbyteApiServerImageSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=airbyte/airbyte-api-server
	Repository string `json:"repository,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="0.50.30"
	Tag string `json:"tag,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=IfNotPresent
	PullPolicy corev1.PullPolicy `json:"pullPolicy,omitempty"`
}

type WebAppImageSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=airbyte/webapp
	Repository string `json:"repository,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="0.50.30"
	Tag string `json:"tag,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=IfNotPresent
	PullPolicy corev1.PullPolicy `json:"pullPolicy,omitempty"`
}

type PodSweeperImageSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=bitnami/kubectl
	Repository string `json:"repository,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="latest"
	Tag string `json:"tag,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=IfNotPresent
	PullPolicy corev1.PullPolicy `json:"pullPolicy,omitempty"`
}

type ConnectorBuilderServerImageSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=airbyte/connector-builder-server
	Repository string `json:"repository,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="0.50.30"
	Tag string `json:"tag,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=IfNotPresent
	PullPolicy corev1.PullPolicy `json:"pullPolicy,omitempty"`
}

type AirbyteBootloaderImageSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=airbyte/bootloader
	Repository string `json:"repository,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="0.50.30"
	Tag string `json:"tag,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=IfNotPresent
	PullPolicy corev1.PullPolicy `json:"pullPolicy,omitempty"`
}

type TemporalImageSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=temporalio/auto-setup
	Repository string `json:"repository,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="1.20.1"
	Tag string `json:"tag,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=IfNotPresent
	PullPolicy corev1.PullPolicy `json:"pullPolicy,omitempty"`
}

type KeycloakImageSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=airbyte/keycloak
	Repository string `json:"repository,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="1.20.1"
	Tag string `json:"tag,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=IfNotPresent
	PullPolicy corev1.PullPolicy `json:"pullPolicy,omitempty"`
}

type CronImageSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=airbyte/cron
	Repository string `json:"repository,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="0.50.30"
	Tag string `json:"tag,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=IfNotPresent
	PullPolicy corev1.PullPolicy `json:"pullPolicy,omitempty"`
}

func (Airbyte *Airbyte) GetNameWithSuffix(name string) string {
	return Airbyte.GetName() + "" + name
}

type ServerServiceSpec struct {
	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:enum=ClusterIP;NodePort;LoadBalancer;ExternalName
	// +kubebuilder:default=ClusterIP
	Type corev1.ServiceType `json:"type"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:default=8001
	Port int32 `json:"port"`
}

type AirbyteApiServerServiceSpec struct {
	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:enum=ClusterIP;NodePort;LoadBalancer;ExternalName
	// +kubebuilder:default=ClusterIP
	Type corev1.ServiceType `json:"type"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:default=80
	Port int32 `json:"port"`
}

type WebAppServiceSpec struct {
	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:enum=ClusterIP;NodePort;LoadBalancer;ExternalName
	// +kubebuilder:default=ClusterIP
	Type corev1.ServiceType `json:"type"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:default=80
	Port int32 `json:"port"`
}

type ConnectorBuilderServerServiceSpec struct {
	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:enum=ClusterIP;NodePort;LoadBalancer;ExternalName
	// +kubebuilder:default=ClusterIP
	Type corev1.ServiceType `json:"type"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:default=80
	Port int32 `json:"port"`
}

type TemporalServiceSpec struct {
	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:enum=ClusterIP;NodePort;LoadBalancer;ExternalName
	// +kubebuilder:default=ClusterIP
	Type corev1.ServiceType `json:"type"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:default=7233
	Port int32 `json:"port"`
}

type AirbyteBootloaderServiceSpec struct {
	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// +kubebuilder:validation:enum=ClusterIP;NodePort;LoadBalancer;ExternalName
	// +kubebuilder:default=ClusterIP
	Type corev1.ServiceType `json:"type"`

	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:default=80
	Port int32 `json:"port"`
}

type KeycloakServiceSpec struct {
	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:enum=ClusterIP;NodePort;LoadBalancer;ExternalName
	// +kubebuilder:default=ClusterIP
	Type corev1.ServiceType `json:"type"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:default=8180
	Port int32 `json:"port"`
}

type AirbyteApiServerDebugSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	Enabled bool `json:"enabled"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=5005
	RemoteDebugPort int32 `json:"RemoteDebugPort"`
}

type ConnectorBuilderServerDebugSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	Enabled bool `json:"enabled"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=5005
	RemoteDebugPort int32 `json:"RemoteDebugPort"`
}

type SecretParam struct {
	AirbyteSecrets map[string]string `json:"airbyte-secrets,omitempty"`
	GcsLogCreds    map[string]string `json:"gcs-log-creds,omitempty"`
	// airbyte-secrets
	// Suffix    string            `json:"suffix,omitempty"`
	// GcpJson   string            `json:"gcpJson,omitempty"`
	// SecretMap map[string]string `json:"secretMap,omitempty"`
}

// SetStatusCondition updates the status condition using the provided arguments.
// If the condition already exists, it updates the condition; otherwise, it appends the condition.
// If the condition status has changed, it updates the condition's LastTransitionTime.
func (r *Airbyte) SetStatusCondition(condition metav1.Condition) {
	// if the condition already exists, update it
	existingCondition := apimeta.FindStatusCondition(r.Status.Conditions, condition.Type)
	if existingCondition == nil {
		condition.ObservedGeneration = r.GetGeneration()
		condition.LastTransitionTime = metav1.Now()
		r.Status.Conditions = append(r.Status.Conditions, condition)
	} else if existingCondition.Status != condition.Status || existingCondition.Reason != condition.Reason || existingCondition.Message != condition.Message {
		existingCondition.Status = condition.Status
		existingCondition.Reason = condition.Reason
		existingCondition.Message = condition.Message
		existingCondition.ObservedGeneration = r.GetGeneration()
		existingCondition.LastTransitionTime = metav1.Now()
	}
}

// InitStatusConditions initializes the status conditions to the provided conditions.
func (r *Airbyte) InitStatusConditions() {
	r.Status.Conditions = []metav1.Condition{}
	r.SetStatusCondition(metav1.Condition{
		Type:               ConditionTypeProgressing,
		Status:             metav1.ConditionTrue,
		Reason:             ConditionReasonPreparing,
		Message:            "Airbyte is preparing",
		ObservedGeneration: r.GetGeneration(),
		LastTransitionTime: metav1.Now(),
	})
	r.SetStatusCondition(metav1.Condition{
		Type:               ConditionTypeAvailable,
		Status:             metav1.ConditionFalse,
		Reason:             ConditionReasonPreparing,
		Message:            "Airbyte is preparing",
		ObservedGeneration: r.GetGeneration(),
		LastTransitionTime: metav1.Now(),
	})
}

// AirbyteStatus defines the observed state of Airbyte
type AirbyteStatus struct {
	// +kubebuilder:validation:Optional
	Conditions []metav1.Condition `json:"condition,omitempty"`
}

//+kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Image",type="string",JSONPath=".spec.image.repository",description="The Docker Image of Airbyte"
// +kubebuilder:printcolumn:name="Tag",type="string",JSONPath=".spec.image.tag",description="The Docker Tag of Airbyte"
//+kubebuilder:subresource:status

// Airbyte is the Schema for the airbytes API
type Airbyte struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AirbyteSpec   `json:"spec,omitempty"`
	Status AirbyteStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AirbyteList contains a list of Airbyte
type AirbyteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Airbyte `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Airbyte{}, &AirbyteList{})
}

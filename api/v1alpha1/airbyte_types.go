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
	"github.com/zncdata-labs/operator-go/pkg/status"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

var (
	PullPolicy                 = "IfNotPresent"
	ServerRepo                 = "airbyte/server"
	ServerTag                  = "0.50.30"
	ServerPort                 = 8001
	WorkerRepo                 = "airbyte/worker"
	WorkerTag                  = "0.50.30"
	AirbyteApiServerRepo       = "airbyte/airbyte-api-server"
	AirbyteApiServerTag        = "0.50.30"
	AirbyteApiServerPort       = 8006
	WebAppRepo                 = "airbyte/webapp"
	WebAppTag                  = "0.50.30"
	WebAppPort                 = 80
	PodSweeperRepo             = "bitnami/kubectl"
	PodSweeperTag              = "latest"
	ConnectorBuilderServerRepo = "airbyte/connector-builder-server"
	ConnectorBuilderServerTag  = "0.50.30"
	ConnectorBuilderServerPort = 80
	AirbyteBootloaderRepo      = "airbyte/airbyte-bootloader"
	AirbyteBootloaderTag       = "0.50.30"
	AirbyteBootloaderPort      = 80
	TemporalRepo               = "temporalio/auto-setup"
	TemporalTag                = "1.20.1"
	TemporalPort               = 7233
	KeycloakRepo               = "airbyte/keycloak"
	KeycloakTag                = "0.50.30"
	KeycloakPort               = 8180
	CronRepo                   = "airbyte/cron"
	CronTag                    = "0.50.30"
)

// AirbyteSpec defines the desired state of Airbyte
type AirbyteSpec struct {
	// +kubebuilder:validation:Required
	SecurityContext *corev1.PodSecurityContext `json:"securityContext"`

	// +kubebuilder:validation:Optional
	ClusterConfig *ClusterConfigSpec `json:"clusterConfig,omitempty"`

	// +kubebuilder:validation:Optional
	Postgres *PostgresSpec `json:"postgres,omitempty"`

	// +kubebuilder:validation:Optional
	Minio *MinioSpec `json:"minio,omitempty"`

	// +kubebuilder:validation:Required
	Server *ServerSpec `json:"server"`

	// +kubebuilder:validation:Required
	Worker *WorkerSpec `json:"worker"`

	// +kubebuilder:validation:Required
	AirbyteApiServer *AirbyteApiServerSpec `json:"airbyteApiServer"`

	// +kubebuilder:validation:Required
	WebApp *WebAppSpec `json:"webApp"`

	// +kubebuilder:validation:Required
	PodSweeper *PodSweeperSpec `json:"podSweeper"`

	// +kubebuilder:validation:Required
	ConnectorBuilderServer *ConnectorBuilderServerSpec `json:"connectorBuilderServer"`

	// +kubebuilder:validation:Required
	AirbyteBootloader *AirbyteBootloaderSpec `json:"airbyteBootloader"`

	// +kubebuilder:validation:Required
	Temporal *TemporalSpec `json:"temporal"`

	// +kubebuilder:validation:Required
	Keycloak *KeycloakSpec `json:"keycloak"`

	// +kubebuilder:validation:Required
	Cron *CronSpec `json:"cron"`
}

type MinioSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default="minio"
	RootUser string `json:"rootUser,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="minio123"
	RootPassword string `json:"rootPassword,omitempty"`
}

type ClusterConfigSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default="airbyte-admin"
	ServiceAccountName string `json:"serviceAccountName,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="oss"
	DeploymentMode string `json:"deploymentMode,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default={}
	EnvVars map[string]string `json:"envVars,omitempty"`

	// +kubebuilder:validation:Optional
	Database *DatabaseClusterConfigSpec `json:"database,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="community"
	Edition string `json:"edition,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="MINIO"
	StateStorageType string `json:"stateStorageType,omitempty"`

	// +kubebuilder:validation:Optional
	Logs *LogsClusterConfigSpec `json:"logs,omitempty"`

	// +kubebuilder:validation:Optional
	Metrics *MetricsClusterConfigSpec `json:"metrics,omitempty"`

	// +kubebuilder:validation:Optional
	Jobs *JobsClusterConfigSpec `json:"jobs,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	ConfigMapName string `json:"configMapName,omitempty"`
}

type DatabaseClusterConfigSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	SecretName string `json:"secretName,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	SecretValue string `json:"secretValue,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="example.com"
	Host string `json:"host,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=5432
	Port int32 `json:"port,omitempty"`
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

type LogsClusterConfigSpec struct {
	// +kubebuilder:validation:Optional
	AccessKey *LogsAccessKeyClusterConfigSpec `json:"accessKey,omitempty"`

	// +kubebuilder:validation:Optional
	SecretKey *LogsSecretKeyClusterConfigSpec `json:"secretKey,omitempty"`

	// +kubebuilder:validation:Optional
	Minio *LogMinioClusterConfigSpec `json:"minio,omitempty"`

	// +kubebuilder:validation:Optional
	ExternalMinio *LogExternalMinioClusterConfigSpec `json:"externalMinio,omitempty"`

	// +kubebuilder:validation:Optional
	S3 *LogS3ClusterConfigSpec `json:"s3,omitempty"`

	// +kubebuilder:validation:Optional
	Gcs *LogGcsClusterConfigSpec `json:"gcs,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="MINIO"
	StorageType string `json:"storageType,omitempty"`
}

type LogsAccessKeyClusterConfigSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	Password string `json:"password,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	ExistingSecret string `json:"existingSecret,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	ExistingSecretKey string `json:"existingSecretKey,omitempty"`
}

type LogsSecretKeyClusterConfigSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	Password string `json:"password,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	ExistingSecret string `json:"existingSecret,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	ExistingSecretKey string `json:"existingSecretKey,omitempty"`
}

type MetricsClusterConfigSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	MetricClient string `json:"metricClient,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	OtelCollectorEndpoint string `json:"otelCollectorEndpoint,omitempty"`
}

type LogMinioClusterConfigSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=true
	Enabled bool `json:"enabled,omitempty"`
}

type LogExternalMinioClusterConfigSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	Enabled bool `json:"enabled,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	Endpoint string `json:"endpoint,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="localhost"
	Host string `json:"host,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=9000
	Port int32 `json:"port,omitempty"`
}

type LogS3ClusterConfigSpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:default=false
	Enabled bool `json:"enabled,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="airbyte-dev-logs"
	Bucket string `json:"bucket,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	BucketRegion string `json:"bucketRegion,omitempty"`
}

type LogGcsClusterConfigSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	Bucket string `json:"bucket,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	Credentials string `json:"credentials,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	CredentialsJson string `json:"credentialsJson,omitempty"`
}

type PostgresSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default="postgresql"
	Host string `json:"host,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default="5432"
	Port string `json:"port,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default="airbyte"
	UserName string `json:"username,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default="airbyte"
	Password string `json:"password,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default="db-airbyte"
	DataBase string `json:"database,omitempty"`
}

type ServerSpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:default=true
	Enabled bool `json:"enabled,omitempty"`

	// +kubebuilder:validation:Optional
	RoleGroups map[string]*RoleGroupServerSpec `json:"roleGroups,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Optional
	Image *ImageSpec `json:"image,omitempty"`

	// +kubebuilder:validation:Required
	Resources *corev1.ResourceRequirements `json:"resources"`

	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels,omitempty"`

	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// +kubebuilder:validation:Optional
	Tolerations *corev1.Toleration `json:"tolerations,omitempty"`

	// +kubebuilder:validation:Optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default={}
	Secret map[string]string `json:"secret,omitempty"`

	// +kubebuilder:validation:Optional
	Service *ServiceSpec `json:"service,omitempty"`
}

type RoleGroupServerSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=1
	Replicas int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Optional
	Config *ConfigRoleGroupServerSpec `json:"config,omitempty"`
}

type ConfigRoleGroupServerSpec struct {
	// +kubebuilder:validation:Optional
	Image *ImageSpec `json:"image,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// +kubebuilder:validation:Optional
	SecurityContext *corev1.PodSecurityContext `json:"securityContext,omitempty"`

	// +kubebuilder:validation:Optional
	MatchLabels map[string]string `json:"matchLabels,omitempty"`

	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// +kubebuilder:validation:Optional
	Tolerations *corev1.Toleration `json:"tolerations,omitempty"`

	// +kubebuilder:validation:Optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default={}
	Secret map[string]string `json:"secret,omitempty"`

	// +kubebuilder:validation:Optional
	Service *ServiceSpec `json:"service,omitempty"`
}

func (server *ServerSpec) GetImage() *ImageSpec {
	if server.Image == nil {
		return &ImageSpec{
			Repository: ServerRepo,
			Tag:        ServerTag,
			PullPolicy: corev1.PullPolicy(PullPolicy),
		}
	}
	return server.Image
}

func (server *ServerSpec) GetService() *ServiceSpec {
	if server.Service == nil {
		return &ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Port: int32(ServerPort),
		}
	}
	return server.Service
}

type WorkerSpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:default=true
	Enabled bool `json:"enabled"`

	// +kubebuilder:validation:Optional
	RoleConfig *RoleConfigWorkerSpec `json:"roleConfig,omitempty"`

	// +kubebuilder:validation:Optional
	RoleGroups map[string]*RoleGroupWorkerSpec `json:"roleGroups,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Optional
	Image *ImageSpec `json:"image,omitempty"`

	// +kubebuilder:validation:Required
	Resources *corev1.ResourceRequirements `json:"resources"`

	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels,omitempty"`

	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// +kubebuilder:validation:Optional
	Tolerations *corev1.Toleration `json:"tolerations,omitempty"`

	// +kubebuilder:validation:Optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`
}

type RoleConfigWorkerSpec struct {
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

type RoleGroupWorkerSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=1
	Replicas int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Optional
	Config *ConfigRoleGroupWorkerSpec `json:"config,omitempty"`
}

type ConfigRoleGroupWorkerSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Optional
	Image *ImageSpec `json:"image,omitempty"`

	// +kubebuilder:validation:Optional
	SecurityContext *corev1.PodSecurityContext `json:"securityContext"`

	// +kubebuilder:validation:Required
	Resources *corev1.ResourceRequirements `json:"resources"`

	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// +kubebuilder:validation:Optional
	MatchLabels map[string]string `json:"matchLabels,omitempty"`

	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// +kubebuilder:validation:Optional
	Tolerations *corev1.Toleration `json:"tolerations,omitempty"`

	// +kubebuilder:validation:Optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

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

func (worker *WorkerSpec) GetImage() *ImageSpec {
	if worker.Image == nil {
		return &ImageSpec{
			Repository: WorkerRepo,
			Tag:        WorkerTag,
			PullPolicy: corev1.PullPolicy(PullPolicy),
		}
	}
	return worker.Image
}

type ContainerOrchestratorSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=true
	Enabled bool `json:"enabled,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=""
	Image string `json:"image,omitempty"`
}

type AirbyteApiServerSpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:default=true
	Enabled bool `json:"enabled"`

	// +kubebuilder:validation:Optional
	RoleConfig *RoleConfigConnectorBuilderServerAndAirbyteApiServerSpec `json:"roleConfig,omitempty"`

	// +kubebuilder:validation:Optional
	RoleGroups map[string]*RoleGroupAirbyteApiServerSpec `json:"roleGroups,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Optional
	Image *ImageSpec `json:"image,omitempty"`

	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// +kubebuilder:validation:Optional
	SecurityContext *corev1.PodSecurityContext `json:"securityContext,omitempty"`

	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels,omitempty"`

	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// +kubebuilder:validation:Optional
	Tolerations *corev1.Toleration `json:"tolerations,omitempty"`

	// +kubebuilder:validation:Optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// +kubebuilder:validation:Optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// +kubebuilder:validation:Optional
	Service *ServiceSpec `json:"service,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default={}
	Secret map[string]string `json:"secret,omitempty"`
}

type RoleGroupAirbyteApiServerSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=1
	Replicas int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Optional
	Config *ConfigRoleGroupAirbyteApiServerSpec `json:"config,omitempty"`
}

type ConfigRoleGroupAirbyteApiServerSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Optional
	Image *ImageSpec `json:"image,omitempty"`

	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// +kubebuilder:validation:Optional
	SecurityContext *corev1.PodSecurityContext `json:"securityContext,omitempty"`

	// +kubebuilder:validation:Optional
	MatchLabels map[string]string `json:"matchLabels,omitempty"`

	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// +kubebuilder:validation:Optional
	Tolerations *corev1.Toleration `json:"tolerations,omitempty,omitempty"`

	// +kubebuilder:validation:Optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// +kubebuilder:validation:Optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// +kubebuilder:validation:Optional
	Service *ServiceSpec `json:"service,omitempty"`

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
	ExtraEnv corev1.EnvVar `json:"extraEnv,omitempty"`
}

func (airbyteApiServer *AirbyteApiServerSpec) GetImage() *ImageSpec {
	if airbyteApiServer.Image == nil {
		return &ImageSpec{
			Repository: AirbyteApiServerRepo,
			Tag:        AirbyteApiServerTag,
			PullPolicy: corev1.PullPolicy(PullPolicy),
		}
	}
	return airbyteApiServer.Image
}

func (airbyteApiServer *AirbyteApiServerSpec) GetService() *ServiceSpec {
	if airbyteApiServer.Service == nil {
		return &ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Port: int32(AirbyteApiServerPort),
		}
	}
	return airbyteApiServer.Service
}

type WebAppSpec struct {

	// +kubebuilder:validation:Required
	// +kubebuilder:default=true
	Enabled bool `json:"enabled"`

	// +kubebuilder:validation:Optional
	RoleConfig *RoleConfigWebAppSpec `json:"roleConfig,omitempty"`

	// +kubebuilder:validation:Optional
	RoleGroups map[string]*RoleGroupWebAppSpec `json:"roleGroups,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Optional
	Image *ImageSpec `json:"image,omitempty"`

	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// +kubebuilder:validation:Optional
	SecurityContext *corev1.PodSecurityContext `json:"securityContext,omitempty"`

	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels,omitempty"`

	// +kubebuilder:validation:Required
	Ingress *WebAppIngressSpec `json:"ingress"`

	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// +kubebuilder:validation:Optional
	Tolerations *corev1.Toleration `json:"tolerations,omitempty"`

	// +kubebuilder:validation:Optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// +kubebuilder:validation:Optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// +kubebuilder:validation:Optional
	Service *ServiceSpec `json:"service,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default={}
	Secret map[string]string `json:"secret,omitempty"`
}

type RoleConfigWebAppSpec struct {
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
	EnvVars map[string]string `json:"envVars,omitempty"`

	// +kubebuilder:validation:Optional
	ExtraEnv corev1.EnvVar `json:"extraEnv,omitempty"`
}

type RoleGroupWebAppSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=1
	Replicas int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Optional
	Config *ConfigRoleGroupWebAppSpec `json:"config,omitempty"`
}

type ConfigRoleGroupWebAppSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Optional
	Image *ImageSpec `json:"image,omitempty"`

	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// +kubebuilder:validation:Optional
	SecurityContext *corev1.PodSecurityContext `json:"securityContext,omitempty"`

	// +kubebuilder:validation:Optional
	MatchLabels map[string]string `json:"matchLabels,omitempty"`

	// +kubebuilder:validation:Optional
	Ingress *WebAppIngressSpec `json:"ingress,omitempty"`

	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// +kubebuilder:validation:Optional
	Tolerations *corev1.Toleration `json:"tolerations,omitempty"`

	// +kubebuilder:validation:Optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// +kubebuilder:validation:Optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// +kubebuilder:validation:Optional
	Service *ServiceSpec `json:"service,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default={}
	Secret map[string]string `json:"secret,omitempty"`

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

func (webApp *WebAppSpec) GetImage() *ImageSpec {
	if webApp.Image == nil {
		return &ImageSpec{
			Repository: WebAppRepo,
			Tag:        WebAppTag,
			PullPolicy: corev1.PullPolicy(PullPolicy),
		}
	}
	return webApp.Image
}

func (webApp *WebAppSpec) GetService() *ServiceSpec {
	if webApp.Service == nil {
		return &ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Port: int32(WebAppPort),
		}
	}
	return webApp.Service
}

type PodSweeperSpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:default=true
	Enabled bool `json:"enabled"`

	// +kubebuilder:validation:Optional
	RoleConfig *RoleConfigPodSweeperSpec `json:"roleConfig,omitempty"`

	// +kubebuilder:validation:Optional
	RoleGroups map[string]*RoleGroupPodSweeperSpec `json:"roleGroups,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Optional
	Image *ImageSpec `json:"image,omitempty"`

	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// +kubebuilder:validation:Optional
	SecurityContext *corev1.PodSecurityContext `json:"securityContext,omitempty"`

	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels,omitempty"`

	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// +kubebuilder:validation:Optional
	Tolerations *corev1.Toleration `json:"tolerations,omitempty"`

	// +kubebuilder:validation:Optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// +kubebuilder:validation:Optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// +kubebuilder:validation:Optional
	LivenessProbe *ProbeSpec `json:"livenessProbe,omitempty"`

	// +kubebuilder:validation:Optional
	ReadinessProbe *ProbeSpec `json:"readinessProbe,omitempty"`

	// +kubebuilder:validation:Optional
	ExtraVolumes []corev1.Volume `json:"extraVolumes,omitempty"`

	// +kubebuilder:validation:Optional
	ExtraVolumeMounts []corev1.VolumeMount `json:"extraVolumeMounts,omitempty"`
}

type RoleConfigPodSweeperSpec struct {
	// +kubebuilder:validation:Optional
	TimeToDeletePods *PodSweeperTimeToDeletePodsSpec `json:"timeToDeletePods,omitempty"`
}

type RoleGroupPodSweeperSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=1
	Replicas int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Optional
	Config *ConfigRoleGroupPodSweeperSpec `json:"config,omitempty"`
}

type ConfigRoleGroupPodSweeperSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Optional
	Image *ImageSpec `json:"image,omitempty"`

	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// +kubebuilder:validation:Optional
	SecurityContext *corev1.PodSecurityContext `json:"securityContext,omitempty"`

	// +kubebuilder:validation:Optional
	MatchLabels map[string]string `json:"matchLabels,omitempty"`

	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// +kubebuilder:validation:Optional
	Tolerations *corev1.Toleration `json:"tolerations,omitempty"`

	// +kubebuilder:validation:Optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// +kubebuilder:validation:Optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

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

func (podSweeper *PodSweeperSpec) GetImage() *ImageSpec {
	if podSweeper.Image == nil {
		return &ImageSpec{
			Repository: PodSweeperRepo,
			Tag:        PodSweeperTag,
			PullPolicy: corev1.PullPolicy(PullPolicy),
		}
	}
	return podSweeper.Image
}

type ConnectorBuilderServerSpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:default=true
	Enabled bool `json:"enabled"`

	// +kubebuilder:validation:Optional
	RoleConfig *RoleConfigConnectorBuilderServerAndAirbyteApiServerSpec `json:"roleConfig,omitempty"`

	// +kubebuilder:validation:Optional
	RoleGroups map[string]*RoleGroupConnectorBuilderServerSpec `json:"roleGroups,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Optional
	Image *ImageSpec `json:"image,omitempty"`

	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// +kubebuilder:validation:Optional
	SecurityContext *corev1.PodSecurityContext `json:"securityContext,omitempty"`

	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels,omitempty"`

	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// +kubebuilder:validation:Optional
	Tolerations *corev1.Toleration `json:"tolerations,omitempty"`

	// +kubebuilder:validation:Optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// +kubebuilder:validation:Optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// +kubebuilder:validation:Optional
	Service *ServiceSpec `json:"service,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default={}
	Secret map[string]string `json:"secret,omitempty"`
}

type RoleConfigConnectorBuilderServerAndAirbyteApiServerSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=INFO
	LogLevel string `json:"logLevel,omitempty"`

	// +kubebuilder:validation:Optional
	Debug *DebugSpec `json:"debug,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default={}
	EnvVars map[string]string `json:"envVars,omitempty"`

	// +kubebuilder:validation:Optional
	ExtraEnv corev1.EnvVar `json:"extraEnv,omitempty"`
}

type RoleGroupConnectorBuilderServerSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=1
	Replicas int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Optional
	Config *ConfigRoleGroupConnectorBuilderServerSpec `json:"config,omitempty"`
}

type ConfigRoleGroupConnectorBuilderServerSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Optional
	Image *ImageSpec `json:"image,omitempty"`

	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// +kubebuilder:validation:Optional
	SecurityContext *corev1.PodSecurityContext `json:"securityContext,omitempty"`

	// +kubebuilder:validation:Optional
	MatchLabels map[string]string `json:"matchLabels,omitempty"`

	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// +kubebuilder:validation:Optional
	Tolerations *corev1.Toleration `json:"tolerations,omitempty"`

	// +kubebuilder:validation:Optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// +kubebuilder:validation:Optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// +kubebuilder:validation:Optional
	Service *ServiceSpec `json:"service,omitempty"`

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
	ExtraEnv corev1.EnvVar `json:"extraEnv,omitempty"`
}

func (connectorBuilderServer *ConnectorBuilderServerSpec) GetImage() *ImageSpec {
	if connectorBuilderServer.Image == nil {
		return &ImageSpec{
			Repository: ConnectorBuilderServerRepo,
			Tag:        ConnectorBuilderServerTag,
			PullPolicy: corev1.PullPolicy(PullPolicy),
		}
	}
	return connectorBuilderServer.Image
}

func (connectorBuilderServer *ConnectorBuilderServerSpec) GetService() *ServiceSpec {
	if connectorBuilderServer.Service == nil {
		return &ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Port: int32(ConnectorBuilderServerPort),
		}
	}
	return connectorBuilderServer.Service
}

type AirbyteBootloaderSpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:default=true
	Enabled bool `json:"enabled"`

	// +kubebuilder:validation:Optional
	RoleConfig *RoleConfigAirbyteBootloaderSpec `json:"roleConfig,omitempty"`

	// +kubebuilder:validation:Optional
	RoleGroups map[string]*RoleGroupAirbyteBootloaderSpec `json:"roleGroups,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Optional
	Image *ImageSpec `json:"image,omitempty"`

	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels,omitempty"`

	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// +kubebuilder:validation:Optional
	Tolerations *corev1.Toleration `json:"tolerations,omitempty"`

	// +kubebuilder:validation:Optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// +kubebuilder:validation:Optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// +kubebuilder:validation:Optional
	Service *ServiceSpec `json:"service,omitempty"`
}

type RoleConfigAirbyteBootloaderSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=INFO
	LogLevel string `json:"logLevel,omitempty"`

	// +kubebuilder:validation:Optional
	RunDatabaseMigrationsOnStartup *bool `json:"runDatabaseMigrationsOnStartup,omitempty"`
}

type RoleGroupAirbyteBootloaderSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=1
	Replicas int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Optional
	Config *ConfigRoleGroupAirbyteBootloaderSpec `json:"config,omitempty"`
}

type ConfigRoleGroupAirbyteBootloaderSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Optional
	Image *ImageSpec `json:"image,omitempty"`

	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// +kubebuilder:validation:Optional
	MatchLabels map[string]string `json:"matchLabels,omitempty"`

	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// +kubebuilder:validation:Optional
	Tolerations *corev1.Toleration `json:"tolerations,omitempty"`

	// +kubebuilder:validation:Optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// +kubebuilder:validation:Optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// +kubebuilder:validation:Optional
	Service *ServiceSpec `json:"service,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=INFO
	LogLevel string `json:"logLevel,omitempty"`

	// +kubebuilder:validation:Optional
	RunDatabaseMigrationsOnStartup *bool `json:"runDatabaseMigrationsOnStartup,omitempty"`
}

func (airbyteBootloader *AirbyteBootloaderSpec) GetImage() *ImageSpec {
	if airbyteBootloader.Image == nil {
		return &ImageSpec{
			Repository: AirbyteBootloaderRepo,
			Tag:        AirbyteBootloaderTag,
			PullPolicy: corev1.PullPolicy(PullPolicy),
		}
	}
	return airbyteBootloader.Image
}

func (airbyteBootloader *AirbyteBootloaderSpec) GetService() *ServiceSpec {
	if airbyteBootloader.Service == nil {
		return &ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Port: int32(AirbyteBootloaderPort),
		}
	}
	return airbyteBootloader.Service
}

type TemporalSpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:default=true
	Enabled bool `json:"enabled"`

	// +kubebuilder:validation:Optional
	RoleConfig *RoleConfigCronAndTemporalSpec `json:"roleConfig,omitempty"`

	// +kubebuilder:validation:Optional
	RoleGroups map[string]*RoleGroupTemporalSpec `json:"roleGroups,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Optional
	Image *ImageSpec `json:"image,omitempty"`

	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels,omitempty"`

	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// +kubebuilder:validation:Optional
	Tolerations *corev1.Toleration `json:"tolerations,omitempty"`

	// +kubebuilder:validation:Optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// +kubebuilder:validation:Optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// +kubebuilder:validation:Optional
	Service *ServiceSpec `json:"service,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default={}
	Secret map[string]string `json:"secret,omitempty"`
}

type RoleGroupTemporalSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=1
	Replicas int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Optional
	Config *ConfigRoleGroupTemporalSpec `json:"config,omitempty"`
}

type ConfigRoleGroupTemporalSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Optional
	Image *ImageSpec `json:"image,omitempty"`

	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// +kubebuilder:validation:Optional
	MatchLabels map[string]string `json:"matchLabels,omitempty"`

	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// +kubebuilder:validation:Optional
	Tolerations *corev1.Toleration `json:"tolerations,omitempty"`

	// +kubebuilder:validation:Optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// +kubebuilder:validation:Optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// +kubebuilder:validation:Optional
	Service *ServiceSpec `json:"service,omitempty"`

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

func (temporal *TemporalSpec) GetImage() *ImageSpec {
	if temporal.Image == nil {
		return &ImageSpec{
			Repository: TemporalRepo,
			Tag:        TemporalTag,
			PullPolicy: corev1.PullPolicy(PullPolicy),
		}
	}
	return temporal.Image
}

func (temporal *TemporalSpec) GetService() *ServiceSpec {
	if temporal.Service == nil {
		return &ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Port: int32(TemporalPort),
		}
	}
	return temporal.Service
}

type KeycloakSpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:default=false
	Enabled bool `json:"enabled"`

	// +kubebuilder:validation:Optional
	RoleConfig *RoleConfigKeycloakSpec `json:"roleConfig,omitempty"`

	// +kubebuilder:validation:Optional
	RoleGroups map[string]*RoleGroupKeycloakSpec `json:"roleGroups,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Optional
	Image *ImageSpec `json:"image,omitempty"`

	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels,omitempty"`

	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// +kubebuilder:validation:Optional
	Tolerations *corev1.Toleration `json:"tolerations,omitempty"`

	// +kubebuilder:validation:Optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// +kubebuilder:validation:Optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// +kubebuilder:validation:Optional
	Service *ServiceSpec `json:"service,omitempty"`
}

type RoleConfigKeycloakSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=INFO
	LogLevel string `json:"logLevel,omitempty"`
}

type RoleGroupKeycloakSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=1
	Replicas int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Optional
	Config *ConfigRoleGroupKeycloakSpec `json:"config,omitempty"`
}

type ConfigRoleGroupKeycloakSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Optional
	Image *ImageSpec `json:"image,omitempty"`

	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// +kubebuilder:validation:Optional
	MatchLabels map[string]string `json:"matchLabels,omitempty"`

	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// +kubebuilder:validation:Optional
	Tolerations *corev1.Toleration `json:"tolerations,omitempty"`

	// +kubebuilder:validation:Optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// +kubebuilder:validation:Optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// +kubebuilder:validation:Optional
	Service *ServiceSpec `json:"service,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=INFO
	LogLevel string `json:"logLevel,omitempty"`
}

func (keycloak *KeycloakSpec) GetImage() *ImageSpec {
	if keycloak.Image == nil {
		return &ImageSpec{
			Repository: KeycloakRepo,
			Tag:        KeycloakTag,
			PullPolicy: corev1.PullPolicy(PullPolicy),
		}
	}
	return keycloak.Image
}

func (keycloak *KeycloakSpec) GetService() *ServiceSpec {
	if keycloak.Service == nil {
		return &ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Port: int32(KeycloakPort),
		}
	}
	return keycloak.Service
}

type CronSpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:default=true
	Enabled bool `json:"enabled"`

	// +kubebuilder:validation:Optional
	RoleConfig *RoleConfigCronAndTemporalSpec `json:"roleConfig,omitempty"`

	// +kubebuilder:validation:Optional
	RoleGroups map[string]*RoleGroupCronSpec `json:"roleGroups,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Optional
	Image *ImageSpec `json:"image,omitempty"`

	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels,omitempty"`

	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// +kubebuilder:validation:Optional
	Tolerations *corev1.Toleration `json:"tolerations,omitempty"`

	// +kubebuilder:validation:Optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// +kubebuilder:validation:Optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default={}
	Secret map[string]string `json:"secret,omitempty"`
}

type RoleConfigCronAndTemporalSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=INFO
	LogLevel string `json:"logLevel,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default={}
	EnvVars map[string]string `json:"envVars,omitempty"`

	// +kubebuilder:validation:Optional
	ExtraEnv corev1.EnvVar `json:"extraEnv,omitempty"`
}

type RoleGroupCronSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=1
	Replicas int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Optional
	Config *ConfigRoleGroupCronSpec `json:"config,omitempty"`
}

type ConfigRoleGroupCronSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Optional
	Image *ImageSpec `json:"image,omitempty"`

	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// +kubebuilder:validation:Optional
	MatchLabels map[string]string `json:"matchLabels,omitempty"`

	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// +kubebuilder:validation:Optional
	Tolerations *corev1.Toleration `json:"tolerations,omitempty"`

	// +kubebuilder:validation:Optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// +kubebuilder:validation:Optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

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

func (cron *CronSpec) GetImage() *ImageSpec {
	if cron.Image == nil {
		return &ImageSpec{
			Repository: CronRepo,
			Tag:        CronTag,
			PullPolicy: corev1.PullPolicy(PullPolicy),
		}
	}
	return cron.Image
}

type ImageSpec struct {
	// +kubebuilder:validation:Optional
	Repository string `json:"repository,omitempty"`

	// +kubebuilder:validation:Optional
	Tag string `json:"tag,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=IfNotPresent
	PullPolicy corev1.PullPolicy `json:"pullPolicy,omitempty"`
}

func (Airbyte *Airbyte) GetNameWithSuffix(name string) string {
	return Airbyte.GetName() + "" + name
}

type ServiceSpec struct {
	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:enum=ClusterIP;NodePort;LoadBalancer;ExternalName
	// +kubebuilder:default=ClusterIP
	Type corev1.ServiceType `json:"type,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:default=8001
	Port int32 `json:"port,omitempty"`
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

type DebugSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	Enabled bool `json:"enabled,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=5005
	RemoteDebugPort int32 `json:"RemoteDebugPort,omitempty"`
}

type WebAppIngressSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=true
	Enabled bool `json:"enabled,omitempty"`
	// +kubebuilder:validation:Optional
	TLS *networkingv1.IngressTLS `json:"tls,omitempty"`
	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:="webapp.example.com"
	Host string `json:"host,omitempty"`
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
	// +kubebuilder:validation:Optional
	URLs []StatusURL `json:"urls,omitempty"`
}

type StatusURL struct {
	Name string `json:"name"`
	URL  string `json:"url"`
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
	Status status.Status `json:"status,omitempty"`
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

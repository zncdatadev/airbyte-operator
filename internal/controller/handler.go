package controller

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"

	stackv1alpha1 "github.com/zncdata-labs/airbyte-operator/api/v1alpha1"
)

func (r *AirbyteReconciler) reconcileIngress(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	obj := r.makeIngress(instance, r.Scheme)
	if obj == nil {
		return nil
	}

	if err := CreateOrUpdate(ctx, r.Client, obj); err != nil {
		r.Log.Error(err, "Failed to create or update ingress")
		return err
	}

	if instance.Spec.WebApp.Ingress.Enabled {
		url := fmt.Sprintf("http://%s", instance.Spec.WebApp.Ingress.Host)
		if instance.Status.URLs == nil {
			instance.Status.URLs = []stackv1alpha1.StatusURL{
				{
					Name: "webui",
					URL:  url,
				},
			}
			if err := r.UpdateStatus(ctx, instance); err != nil {
				return err
			}

		} else if instance.Spec.WebApp.Ingress.Host != instance.Status.URLs[0].Name {
			instance.Status.URLs[0].URL = url
			if err := r.UpdateStatus(ctx, instance); err != nil {
				return err
			}

		}
	}

	return nil
}

func (r *AirbyteReconciler) reconcileService(ctx context.Context, instance *stackv1alpha1.Airbyte) error {

	ServerService := r.makeServerService(instance, r.Scheme)
	if ServerService == nil {
		return nil
	}

	AirbyteApiServerService := r.makeAirbyteApiServerService(instance, r.Scheme)
	if AirbyteApiServerService == nil {
		return nil
	}

	ConnectorBuilderServerService := r.makeConnectorBuilderServerService(instance, r.Scheme)
	if ConnectorBuilderServerService == nil {
		return nil
	}

	TemporalService := r.makeTemporalService(instance, r.Scheme)
	if TemporalService == nil {
		return nil
	}

	WebAppService := r.makeWebAppService(instance, r.Scheme)
	if WebAppService == nil {
		return nil
	}

	if instance.Spec.Server.Enabled {
		serverService := r.makeServerService(instance, r.Scheme)
		if serverService != nil {
			if err := CreateOrUpdate(ctx, r.Client, serverService); err != nil {
				r.Log.Error(err, "Failed to create or update service")
				return err
			}
		}
	}

	if instance.Spec.AirbyteApiServer.Enabled {
		airbyteApiServerService := r.makeAirbyteApiServerService(instance, r.Scheme)
		if airbyteApiServerService != nil {
			if err := CreateOrUpdate(ctx, r.Client, airbyteApiServerService); err != nil {
				r.Log.Error(err, "Failed to create or update airbyteApiServer service")
				return err
			}
		}
	}

	if instance.Spec.ConnectorBuilderServer.Enabled {
		connectorBuilderServerService := r.makeConnectorBuilderServerService(instance, r.Scheme)
		if connectorBuilderServerService != nil {
			if err := CreateOrUpdate(ctx, r.Client, connectorBuilderServerService); err != nil {
				r.Log.Error(err, "Failed to create or update connectorBuilderServerService service")
				return err
			}
		}
	}

	if instance.Spec.Temporal.Enabled {
		temporalService := r.makeTemporalService(instance, r.Scheme)
		if temporalService != nil {
			if err := CreateOrUpdate(ctx, r.Client, temporalService); err != nil {
				r.Log.Error(err, "Failed to create or update temporalService service")
				return err
			}
		}
	}

	if instance.Spec.WebApp.Enabled {
		webAppService := r.makeWebAppService(instance, r.Scheme)
		if webAppService != nil {
			if err := CreateOrUpdate(ctx, r.Client, webAppService); err != nil {
				r.Log.Error(err, "Failed to create or update webAppService service")
				return err
			}
		}
	}
	return nil
}

func (r *AirbyteReconciler) reconcileDeployment(ctx context.Context, instance *stackv1alpha1.Airbyte) error {

	ServerDeployment := r.makeServerDeployment(instance, r.Scheme)
	if ServerDeployment == nil {
		return nil
	}

	WorkerDeployment := r.makeWorkerDeployment(instance, r.Scheme)
	if WorkerDeployment == nil {
		return nil
	}

	AirbyteApiServerDeployment := r.makeAirbyteApiServerDeployment(instance, r.Scheme)
	if AirbyteApiServerDeployment == nil {
		return nil
	}

	ConnectorBuilderServerDeployment := r.makeConnectorBuilderServerDeployment(instance, r.Scheme)
	if ConnectorBuilderServerDeployment == nil {
		return nil
	}

	CronDeployment := r.makeCronDeployment(instance, r.Scheme)
	if CronDeployment == nil {
		return nil
	}

	PodSweeperDeployment := r.makePodSweeperDeployment(instance, r.Scheme)
	if PodSweeperDeployment == nil {
		return nil
	}

	TemporalDeployment := r.makeTemporalDeployment(instance, r.Scheme)
	if TemporalDeployment == nil {
		return nil
	}

	WebAppDeployment := r.makeWebAppDeployment(instance, r.Scheme)
	if WebAppDeployment == nil {
		return nil
	}

	if instance.Spec.Server.Enabled {
		serverDeployment := r.makeServerDeployment(instance, r.Scheme)
		if serverDeployment != nil {
			if err := CreateOrUpdate(ctx, r.Client, serverDeployment); err != nil {
				r.Log.Error(err, "Failed to create or update Serverdeployment")
				return err
			}
		}
	}

	if instance.Spec.Worker.Enabled {
		workerDeployment := r.makeWorkerDeployment(instance, r.Scheme)
		if workerDeployment != nil {
			if err := CreateOrUpdate(ctx, r.Client, workerDeployment); err != nil {
				r.Log.Error(err, "Failed to create or update Workerdeployment")
				return err
			}
		}
	}

	if instance.Spec.AirbyteApiServer.Enabled {
		airbyteApiServerDeployment := r.makeAirbyteApiServerDeployment(instance, r.Scheme)
		if airbyteApiServerDeployment != nil {
			if err := CreateOrUpdate(ctx, r.Client, airbyteApiServerDeployment); err != nil {
				r.Log.Error(err, "Failed to create or update AirbyteApiServer")
				return err
			}
		}
	}

	if instance.Spec.ConnectorBuilderServer.Enabled {
		connectorBuilderServerDeployment := r.makeConnectorBuilderServerDeployment(instance, r.Scheme)
		if connectorBuilderServerDeployment != nil {
			if err := CreateOrUpdate(ctx, r.Client, connectorBuilderServerDeployment); err != nil {
				r.Log.Error(err, "Failed to create or update ConnectorBuilderServerDeployment")
				return err
			}
		}
	}

	if instance.Spec.Cron.Enabled {
		cronDeployment := r.makeCronDeployment(instance, r.Scheme)
		if cronDeployment != nil {
			if err := CreateOrUpdate(ctx, r.Client, cronDeployment); err != nil {
				r.Log.Error(err, "Failed to create or update cronDeployment")
				return err
			}
		}
	}

	if instance.Spec.PodSweeper.Enabled {
		podSweeperDeployment := r.makePodSweeperDeployment(instance, r.Scheme)
		if podSweeperDeployment != nil {
			if err := CreateOrUpdate(ctx, r.Client, podSweeperDeployment); err != nil {
				r.Log.Error(err, "Failed to create or update podSweeperDeployment")
				return err
			}
		}
	}

	if instance.Spec.Temporal.Enabled {
		temporalDeployment := r.makeTemporalDeployment(instance, r.Scheme)
		if temporalDeployment != nil {
			if err := CreateOrUpdate(ctx, r.Client, temporalDeployment); err != nil {
				r.Log.Error(err, "Failed to create or update temporalDeployment")
				return err
			}
		}
	}

	if instance.Spec.WebApp.Enabled {
		webAppDeployment := r.makeWebAppDeployment(instance, r.Scheme)
		if webAppDeployment != nil {
			if err := CreateOrUpdate(ctx, r.Client, webAppDeployment); err != nil {
				r.Log.Error(err, "Failed to create or update webAppDeployment")
				return err
			}
		}
	}
	return nil
}

func (r *AirbyteReconciler) reconcileSecret(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	logger := log.FromContext(ctx)
	objs := r.makeSecret(instance)

	if objs == nil {
		return nil
	}
	for _, obj := range objs {
		if err := CreateOrUpdate(ctx, r.Client, obj); err != nil {
			logger.Error(err, "Failed to create or update secret")
			return err
		}
	}
	// Create or update the secret

	return nil
}

func (r *AirbyteReconciler) makeYmlConfigMap(instance *stackv1alpha1.Airbyte, schema *runtime.Scheme) *corev1.ConfigMap {
	labels := instance.GetLabels()
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.GetNameWithSuffix("-airbyte-yml"),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Data: map[string]string{
			"fileContents": "",
		},
	}
	err := ctrl.SetControllerReference(instance, configMap, schema)
	if err != nil {
		r.Log.Error(err, "Failed to set controller reference for configmap")
		return nil
	}
	return configMap
}

func addMapToData(data map[string]string, key string, value map[string]string) {
	if value != nil {
		var strList []string
		for k, v := range value {
			strList = append(strList, k+"="+v)
		}
		data[key] = strings.Join(strList, ",")
	}
}

func (r *AirbyteReconciler) makeEnvConfigMap(instance *stackv1alpha1.Airbyte, schema *runtime.Scheme) *corev1.ConfigMap {
	labels := instance.GetLabels()

	data := map[string]string{
		"AIRBYTE_VERSION": instance.Spec.Server.Image.Tag,
		"AIRBYTE_EDITION": instance.Spec.Global.Edition,
		"CONFIG_ROOT":     "/configs",
		"CONFIGS_DATABASE_MINIMUM_FLYWAY_MIGRATION_VERSION": "0.35.15.001",
		"DATA_DOCKER_MOUNT":          "airbyte_data",
		"DATABASE_DB":                instance.Spec.Postgres.DataBase,
		"DATABASE_HOST":              instance.Spec.Postgres.Host,
		"DATABASE_PORT":              instance.Spec.Postgres.Port,
		"DATABASE_URL":               "jdbc:postgresql://" + instance.Spec.Postgres.Host + ":" + instance.Spec.Postgres.Port + "/" + instance.Spec.Postgres.DataBase,
		"DB_DOCKER_MOUNT":            "airbyte_db",
		"INTERNAL_API_HOST":          instance.Name + "-airbyte-server-svc:" + strconv.FormatInt(int64(instance.Spec.Server.Service.Port), 10),
		"CONNECTOR_BUILDER_API_HOST": instance.Name + "-airbyte-connector-builder-server-svc:" + strconv.FormatInt(int64(instance.Spec.ConnectorBuilderServer.Service.Port), 10),
		"AIRBYTE_API_HOST":           instance.Name + "-airbyte-api-server-svc:" + strconv.FormatInt(int64(instance.Spec.AirbyteApiServer.Service.Port), 10),
		"JOBS_DATABASE_MINIMUM_FLYWAY_MIGRATION_VERSION": "0.29.15.001",
		"LOCAL_ROOT":                             "/tmp/airbyte_local",
		"STATE_STORAGE_MINIO_BUCKET_NAME":        "airbyte-state-storage",
		"TEMPORAL_HOST":                          instance.Name + "-temporal:" + strconv.FormatInt(int64(instance.Spec.Temporal.Service.Port), 10),
		"TEMPORAL_WORKER_PORTS":                  "9001,9002,9003,9004,9005,9006,9007,9008,9009,9010,9011,9012,9013,9014,9015,9016,9017,9018,9019,9020,9021,9022,9023,9024,9025,9026,9027,9028,9029,9030,9031,9032,9033,9034,9035,9036,9037,9038,9039,9040",
		"TRACKING_STRATEGY":                      "segment",
		"WORKER_ENVIRONMENT":                     "kubernetes",
		"WORKSPACE_DOCKER_MOUNT":                 "airbyte_workspace",
		"WORKSPACE_ROOT":                         "/workspace",
		"WORKFLOW_FAILURE_RESTART_DELAY_SECONDS": "",
		"USE_STREAM_CAPABLE_STATE":               "true",
		"AUTO_DETECT_SCHEMA":                     "true",
		"WORKERS_MICRONAUT_ENVIRONMENTS":         "control-plane",
		"CRON_MICRONAUT_ENVIRONMENTS":            "control-plane",
		"SHOULD_RUN_NOTIFY_WORKFLOWS":            "true",
	}

	stateStorageType := ""
	if instance.Spec.Global.StateStorageType != "" {
		stateStorageType = instance.Spec.Global.StateStorageType
	}
	data["WORKER_STATE_STORAGE_TYPE"] = stateStorageType

	containerOrchestratorImage := ""
	if instance != nil && instance.Spec.Worker != nil && instance.Spec.Worker.ContainerOrchestrator != nil {
		if instance.Spec.Worker.ContainerOrchestrator.Image != "" {
			containerOrchestratorImage = instance.Spec.Worker.ContainerOrchestrator.Image
		}
	}
	data["CONTAINER_ORCHESTRATOR_IMAGE"] = containerOrchestratorImage

	containerOrchestratorEnabled := "true"
	if instance != nil && instance.Spec.Worker != nil && instance.Spec.Worker.ContainerOrchestrator != nil {
		if instance.Spec.Worker.ContainerOrchestrator.Enabled {
			containerOrchestratorEnabled = strconv.FormatBool(instance.Spec.Worker.ContainerOrchestrator.Enabled)
		}
	}
	data["CONTAINER_ORCHESTRATOR_ENABLED"] = containerOrchestratorEnabled

	logStorageType := ""
	if instance != nil && instance.Spec.Global != nil && instance.Spec.Global.Logs != nil {
		if instance.Spec.Global.Logs.StorageType != "" {
			logStorageType = instance.Spec.Global.Logs.StorageType
		}
	}
	data["WORKER_LOGS_STORAGE_TYPE"] = logStorageType

	apiUrl := ""
	if instance != nil && instance.Spec.WebApp != nil {
		if instance.Spec.WebApp.ApiUrl != "" {
			apiUrl = instance.Spec.WebApp.ApiUrl
		}
	}
	data["API_URL"] = apiUrl

	ServerUrl := ""
	if instance != nil && instance.Spec.WebApp != nil {
		if instance.Spec.WebApp.ConnectorBuilderServerUrl != "" {
			ServerUrl = instance.Spec.WebApp.ConnectorBuilderServerUrl
		}
	}
	data["CONNECTOR_BUILDER_API_URL"] = ServerUrl

	gcsBucket := ""
	if instance != nil && instance.Spec.Global != nil && instance.Spec.Global.Logs != nil && instance.Spec.Global.Logs.Gcs != nil {
		if instance.Spec.Global.Logs.Gcs.Bucket != "" {
			gcsBucket = instance.Spec.Global.Logs.Gcs.Bucket
		}
	}
	data["GCS_LOG_BUCKET"] = gcsBucket

	gcsCredentials := ""
	if instance != nil && instance.Spec.Global != nil && instance.Spec.Global.Logs != nil && instance.Spec.Global.Logs.Gcs != nil {
		if instance.Spec.Global.Logs.Gcs.Credentials != "" {
			gcsCredentials = instance.Spec.Global.Logs.Gcs.Credentials
		}
	}
	data["GOOGLE_APPLICATION_CREDENTIALS"] = gcsCredentials

	s3Bucket := ""
	if instance != nil && instance.Spec.Global != nil && instance.Spec.Global.Logs != nil && instance.Spec.Global.Logs.S3 != nil {
		if instance.Spec.Global.Logs.S3.Bucket != "" {
			s3Bucket = instance.Spec.Global.Logs.S3.Bucket
		}
	}
	data["S3_LOG_BUCKET"] = s3Bucket

	s3BucketRegion := ""
	if instance != nil && instance.Spec.Global != nil && instance.Spec.Global.Logs != nil && instance.Spec.Global.Logs.S3 != nil {
		if instance.Spec.Global.Logs.S3.BucketRegion != "" {
			s3BucketRegion = instance.Spec.Global.Logs.S3.BucketRegion
		}
	}
	data["S3_LOG_BUCKET_REGION"] = s3BucketRegion

	if instance != nil && instance.Spec.Global != nil {
		if instance.Spec.Global.Edition == "pro" {
			if instance.Spec.Keycloak != nil && instance.Spec.Keycloak.Service != nil {
				data["KEYCLOAK_INTERNAL_HOST"] = instance.Name + "-airbyte-keycloak-svc:" + strconv.FormatInt(int64(instance.Spec.Keycloak.Service.Port), 10)
				data["KEYCLOAK_PORT"] = strconv.FormatInt(int64(instance.Spec.Keycloak.Service.Port), 10)
				data["KEYCLOAK_HOSTNAME_URL"] = "webapp-url/auth"
			}
		} else {
			data["KEYCLOAK_INTERNAL_HOST"] = "localhost"
		}
	}

	if instance != nil && instance.Spec.Global != nil && instance.Spec.Global.Jobs != nil && instance.Spec.Global.Jobs.Kube != nil {
		if instance.Spec.Global.Jobs.Kube.Annotations != nil {
			addMapToData(data, "JOB_KUBE_ANNOTATIONS", instance.Spec.Global.Jobs.Kube.Annotations)
		}

		if instance.Spec.Global.Jobs.Kube.Labels != nil {
			addMapToData(data, "JOB_KUBE_LABELS", instance.Spec.Global.Jobs.Kube.Labels)
		}

		if instance.Spec.Global.Jobs.Kube.NodeSelector != nil {
			addMapToData(data, "JOB_KUBE_NODE_SELECTOR", instance.Spec.Global.Jobs.Kube.NodeSelector)
		}

		if instance.Spec.Global.Jobs.Kube.MainContainerImagePullSecret != "" {
			data["JOB_KUBE_MAIN_CONTAINER_IMAGE_PULL_SECRET"] = instance.Spec.Global.Jobs.Kube.MainContainerImagePullSecret
		}

		if instance.Spec.Global.Jobs.Kube.Tolerations != nil {
			tolerations, err := json.Marshal(instance.Spec.Global.Jobs.Kube.Tolerations)
			if err != nil {
				// 处理错误
			}
			data["JOB_KUBE_TOLERATIONS"] = string(tolerations)
		}

		if instance.Spec.Global.Jobs.Kube.Images.Busybox != "" {
			data["JOB_KUBE_BUSYBOX_IMAGE"] = instance.Spec.Global.Jobs.Kube.Images.Busybox
		}

		if instance.Spec.Global.Jobs.Kube.Images.Socat != "" {
			data["JOB_KUBE_SOCAT_IMAGE"] = instance.Spec.Global.Jobs.Kube.Images.Socat
		}

		if instance.Spec.Global.Jobs.Kube.Images.Curl != "" {
			data["JOB_KUBE_CURL_IMAGE"] = instance.Spec.Global.Jobs.Kube.Images.Curl
		}
	}

	cpuLimit := ""
	if instance != nil && instance.Spec.Global != nil && instance.Spec.Global.Jobs != nil && instance.Spec.Global.Jobs.Resources != nil && instance.Spec.Global.Jobs.Resources.Limits != nil {
		if cpu := instance.Spec.Global.Jobs.Resources.Limits.Cpu(); cpu != nil {
			cpuLimit = cpu.String()
		}
	}
	data["JOB_MAIN_CONTAINER_CPU_LIMIT"] = cpuLimit

	cpuRequest := ""
	if instance != nil && instance.Spec.Global != nil && instance.Spec.Global.Jobs != nil && instance.Spec.Global.Jobs.Resources != nil && instance.Spec.Global.Jobs.Resources.Requests != nil {
		if cpu := instance.Spec.Global.Jobs.Resources.Requests.Cpu(); cpu != nil {
			cpuRequest = cpu.String()
		}
	}
	data["JOB_MAIN_CONTAINER_CPU_REQUEST"] = cpuRequest

	memoryLimit := ""
	if instance != nil && instance.Spec.Global != nil && instance.Spec.Global.Jobs != nil && instance.Spec.Global.Jobs.Resources != nil && instance.Spec.Global.Jobs.Resources.Limits != nil {
		if mem := instance.Spec.Global.Jobs.Resources.Limits.Memory(); mem != nil {
			memoryLimit = mem.String()
		}
	}
	data["JOB_MAIN_CONTAINER_MEMORY_LIMIT"] = memoryLimit

	memoryRequest := ""
	if instance != nil && instance.Spec.Global != nil && instance.Spec.Global.Jobs != nil && instance.Spec.Global.Jobs.Resources != nil && instance.Spec.Global.Jobs.Resources.Requests != nil {
		if mem := instance.Spec.Global.Jobs.Resources.Requests.Memory(); mem != nil {
			memoryRequest = mem.String()
		}
	}
	data["JOB_MAIN_CONTAINER_MEMORY_REQUEST"] = memoryRequest

	runDatabaseMigrationOnStartup := "true"
	if instance.Spec.AirbyteBootloader.RunDatabaseMigrationsOnStartup != nil {
		runDatabaseMigrationOnStartup = strconv.FormatBool(*instance.Spec.AirbyteBootloader.RunDatabaseMigrationsOnStartup)
	}
	data["RUN_DATABASE_MIGRATION_ON_STARTUP"] = runDatabaseMigrationOnStartup

	s3MinioEndpoint := ""
	s3PathStyleAccess := false
	if instance != nil && instance.Spec.Global != nil && instance.Spec.Global.Logs != nil {
		if instance.Spec.Global.Logs.Minio != nil && instance.Spec.Global.Logs.Minio.Enabled {
			s3MinioEndpoint = "http://airbyte-minio-svc:9000"
			s3PathStyleAccess = true
		} else if instance.Spec.Global.Logs.ExternalMinio != nil {
			if instance.Spec.Global.Logs.ExternalMinio.Endpoint != "" {
				s3MinioEndpoint = instance.Spec.Global.Logs.ExternalMinio.Endpoint
			} else if instance.Spec.Global.Logs.ExternalMinio.Enabled {
				s3MinioEndpoint = "http://" + instance.Spec.Global.Logs.ExternalMinio.Host + ":" + strconv.Itoa(int(instance.Spec.Global.Logs.ExternalMinio.Port))
				s3PathStyleAccess = true
			}
		}
	}
	data["S3_MINIO_ENDPOINT"] = s3MinioEndpoint
	data["STATE_STORAGE_MINIO_ENDPOINT"] = s3MinioEndpoint
	data["S3_PATH_STYLE_ACCESS"] = strconv.FormatBool(s3PathStyleAccess)

	webAppUrl := ""
	if instance != nil && instance.Spec.WebApp != nil {
		if instance.Spec.WebApp.Url != "" {
			webAppUrl = instance.Spec.WebApp.Url
		} else if instance.Spec.WebApp.Service != nil {
			webAppUrl = fmt.Sprintf("http://%s-airbyte-webapp-svc:%d", instance.Name, instance.Spec.WebApp.Service.Port)
		}
	}
	data["WEBAPP_URL"] = webAppUrl

	metricClient := ""
	if instance != nil && instance.Spec.Global != nil && instance.Spec.Global.Metrics != nil {
		if instance.Spec.Global.Metrics.MetricClient != "" {
			metricClient = instance.Spec.Global.Metrics.MetricClient
		}
	}
	data["METRIC_CLIENT"] = metricClient

	otelCollectorEndpoint := ""
	if instance != nil && instance.Spec.Global != nil && instance.Spec.Global.Metrics != nil {
		if instance.Spec.Global.Metrics.OtelCollectorEndpoint != "" {
			otelCollectorEndpoint = instance.Spec.Global.Metrics.OtelCollectorEndpoint
		}
	}
	data["OTEL_COLLECTOR_ENDPOINT"] = otelCollectorEndpoint

	activityMaxAttempt := ""
	activityInitialDelayBetweenAttemptsSeconds := ""
	activityMaxDelayBetweenAttemptsSeconds := ""
	maxNotifyWorkers := "5"
	if instance != nil && instance.Spec.Worker != nil {
		if instance.Spec.Worker.ActivityMaxAttempt != "" {
			activityMaxAttempt = instance.Spec.Worker.ActivityMaxAttempt
		}

		if instance.Spec.Worker.ActivityInitialDelayBetweenAttemptsSeconds != "" {
			activityInitialDelayBetweenAttemptsSeconds = instance.Spec.Worker.ActivityInitialDelayBetweenAttemptsSeconds
		}

		if instance.Spec.Worker.ActivityMaxDelayBetweenAttemptsSeconds != "" {
			activityMaxDelayBetweenAttemptsSeconds = instance.Spec.Worker.ActivityMaxDelayBetweenAttemptsSeconds
		}

		if instance.Spec.Worker.MaxNotifyWorkers != "" {
			maxNotifyWorkers = instance.Spec.Worker.MaxNotifyWorkers
		}
	}
	data["ACTIVITY_MAX_ATTEMPT"] = activityMaxAttempt
	data["ACTIVITY_INITIAL_DELAY_BETWEEN_ATTEMPTS_SECONDS"] = activityInitialDelayBetweenAttemptsSeconds
	data["ACTIVITY_MAX_DELAY_BETWEEN_ATTEMPTS_SECONDS"] = activityMaxDelayBetweenAttemptsSeconds
	data["MAX_NOTIFY_WORKERS"] = maxNotifyWorkers

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.GetNameWithSuffix("-airbyte-env"),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Data: data,
	}
	err := ctrl.SetControllerReference(instance, configMap, schema)
	if err != nil {
		r.Log.Error(err, "Failed to set controller reference for configmap")
		return nil
	}
	return configMap
}

func (r *AirbyteReconciler) reconcileConfigMap(ctx context.Context, instance *stackv1alpha1.Airbyte) error {

	YmlConfigMap := r.makeYmlConfigMap(instance, r.Scheme)
	if YmlConfigMap == nil {
		return nil
	}

	EnvConfigMap := r.makeEnvConfigMap(instance, r.Scheme)
	if EnvConfigMap == nil {
		return nil
	}

	SweepPodScriptConfigMap := r.makeSweepPodScriptConfigMap(instance, r.Scheme)
	if SweepPodScriptConfigMap == nil {
		return nil
	}

	TemporalDynamicconfigConfigMap := r.makeTemporalDynamicconfigConfigMap(instance, r.Scheme)
	if TemporalDynamicconfigConfigMap == nil {
		return nil
	}

	ymlConfigMap := r.makeYmlConfigMap(instance, r.Scheme)
	if ymlConfigMap != nil {
		if err := CreateOrUpdate(ctx, r.Client, ymlConfigMap); err != nil {
			r.Log.Error(err, "Failed to create or update airbyte-yml configmap")
			return err
		}
	}

	envConfigMap := r.makeEnvConfigMap(instance, r.Scheme)
	if envConfigMap != nil {
		if err := CreateOrUpdate(ctx, r.Client, envConfigMap); err != nil {
			r.Log.Error(err, "Failed to create or update airbyte-env configmap")
			return err
		}
	}

	if instance.Spec.PodSweeper.Enabled {
		sweepPodScriptConfigMap := r.makeSweepPodScriptConfigMap(instance, r.Scheme)
		if sweepPodScriptConfigMap != nil {
			if err := CreateOrUpdate(ctx, r.Client, sweepPodScriptConfigMap); err != nil {
				r.Log.Error(err, "Failed to create or update airbyte-sweep-pod-script configmap")
				return err
			}
		}
	}

	if instance.Spec.Temporal.Enabled {
		temporalDynamicconfigConfigMap := r.makeTemporalDynamicconfigConfigMap(instance, r.Scheme)
		if temporalDynamicconfigConfigMap != nil {
			if err := CreateOrUpdate(ctx, r.Client, temporalDynamicconfigConfigMap); err != nil {
				r.Log.Error(err, "Failed to create or update airbyte-temporal-dynamicconfig configmap")
				return err
			}
		}
	}
	return nil
}

func (r *AirbyteReconciler) makeAdminServiceAccount(instance *stackv1alpha1.Airbyte, schema *runtime.Scheme) *corev1.ServiceAccount {
	labels := instance.GetLabels()
	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.GetNameWithSuffix("-admin"),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
	}
	err := ctrl.SetControllerReference(instance, serviceAccount, schema)
	if err != nil {
		r.Log.Error(err, "Failed to set controller reference for serviceaccount")
		return nil
	}
	return serviceAccount

}

func (r *AirbyteReconciler) reconcileServiceAccount(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	serviceAccount := r.makeAdminServiceAccount(instance, r.Scheme)
	if serviceAccount != nil {
		if err := CreateOrUpdate(ctx, r.Client, serviceAccount); err != nil {
			r.Log.Error(err, "Failed to create or update serviceaccount")
			return err
		}
	}
	return nil
}

func (r *AirbyteReconciler) makeAdminRole(instance *stackv1alpha1.Airbyte, schema *runtime.Scheme) *rbacv1.Role {
	labels := instance.GetLabels()
	role := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.GetNameWithSuffix("-admin-role"),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"*"},
				Resources: []string{"jobs", "pods", "pods/log", "pods/exec", "pods/attach"},
				Verbs:     []string{"get", "list", "watch", "create", "update", "patch", "delete"},
			},
		},
	}
	err := ctrl.SetControllerReference(instance, role, schema)
	if err != nil {
		r.Log.Error(err, "Failed to set controller reference for role")
		return nil
	}
	return role

}

func (r *AirbyteReconciler) reconcileRole(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	role := r.makeAdminRole(instance, r.Scheme)
	if role != nil {
		if err := CreateOrUpdate(ctx, r.Client, role); err != nil {
			r.Log.Error(err, "Failed to create or update role")
			return err
		}
	}
	return nil
}

func (r *AirbyteReconciler) makeAdminRoleBinding(instance *stackv1alpha1.Airbyte, schema *runtime.Scheme) *rbacv1.RoleBinding {
	labels := instance.GetLabels()
	roleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.GetNameWithSuffix("-admin-binding"),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      instance.GetNameWithSuffix("-admin"),
				Namespace: instance.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "Role",
			Name:     instance.GetNameWithSuffix("-admin-role"),
			APIGroup: "",
		},
	}
	err := ctrl.SetControllerReference(instance, roleBinding, schema)
	if err != nil {
		r.Log.Error(err, "Failed to set controller reference for rolebinding")
		return nil
	}
	return roleBinding
}

func (r *AirbyteReconciler) reconcileRoleBinding(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	roleBinding := r.makeAdminRoleBinding(instance, r.Scheme)
	if roleBinding != nil {
		if err := CreateOrUpdate(ctx, r.Client, roleBinding); err != nil {
			r.Log.Error(err, "Failed to create or update rolebinding")
			return err
		}
	}
	return nil
}

func createSecretData(secretMap map[string]string) map[string][]byte {
	data := make(map[string][]byte)
	for k, v := range secretMap {
		value := ""
		if v != "" {
			value = base64.StdEncoding.EncodeToString([]byte(v))
		}
		data[k] = []byte(value)
	}
	return data
}

func appendEnvVarsFromSecret(envVars []corev1.EnvVar, secretMap map[string]string, secretPrefix string, roleGroupName string) []corev1.EnvVar {
	for key := range secretMap {
		envVars = append(envVars, corev1.EnvVar{
			Name: key,
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: secretPrefix + "-" + roleGroupName,
					},
					Key: key,
				},
			},
		})
	}
	return envVars
}

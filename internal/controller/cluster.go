package controller

import (
	"context"
	"encoding/base64"
	"fmt"
	stackv1alpha1 "github.com/zncdata-labs/airbyte-operator/api/v1alpha1"
	opgo "github.com/zncdata-labs/operator-go/pkg/apis/commons/v1alpha1"
	"github.com/zncdata-labs/operator-go/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/json"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
	"strings"
)

// reconcile cluster config map
func (r *AirbyteReconciler) reconcileClusterConfigMap(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	reconcileParams := r.createReconcileParams(ctx, nil, instance, r.extractClusterEnvConfigMap)
	if err := reconcileParams.createOrUpdateResourceNoGroup(); err != nil {
		return err
	}
	return nil
}

// extract cluster env config map
func (r *AirbyteReconciler) extractClusterEnvConfigMap(params ExtractorParams) (client.Object, error) {
	instance := params.instance.(*stackv1alpha1.Airbyte)
	ctx := params.ctx
	schema := params.scheme

	labels := instance.GetLabels()

	data := map[string]string{
		"AIRBYTE_VERSION": instance.Spec.Server.RoleConfig.Image.Tag,
		"AIRBYTE_EDITION": instance.Spec.ClusterConfig.Edition,
		"CONFIG_ROOT":     "/configs",
		"CONFIGS_DATABASE_MINIMUM_FLYWAY_MIGRATION_VERSION": "0.35.15.001",
		"DATA_DOCKER_MOUNT":                              "airbyte_data",
		"DB_DOCKER_MOUNT":                                "airbyte_db",
		"INTERNAL_API_HOST":                              instance.Name + "-airbyte-server-svc:" + strconv.FormatInt(int64(instance.Spec.Server.RoleConfig.Service.Port), 10),
		"CONNECTOR_BUILDER_API_HOST":                     instance.Name + "-airbyte-connector-builder-server-svc:" + strconv.FormatInt(int64(instance.Spec.ConnectorBuilderServer.RoleConfig.Service.Port), 10),
		"AIRBYTE_API_HOST":                               instance.Name + "-airbyte-api-server-svc:" + strconv.FormatInt(int64(instance.Spec.ApiServer.RoleConfig.Service.Port), 10),
		"JOBS_DATABASE_MINIMUM_FLYWAY_MIGRATION_VERSION": "0.29.15.001",
		"LOCAL_ROOT":                                     "/tmp/airbyte_local",
		"STATE_STORAGE_MINIO_BUCKET_NAME":                "airbyte-state-storage",
		"TEMPORAL_HOST":                                  instance.Name + "-temporal:" + strconv.FormatInt(int64(instance.Spec.Temporal.RoleConfig.Service.Port), 10),
		"TEMPORAL_WORKER_PORTS":                          "9001,9002,9003,9004,9005,9006,9007,9008,9009,9010,9011,9012,9013,9014,9015,9016,9017,9018,9019,9020,9021,9022,9023,9024,9025,9026,9027,9028,9029,9030,9031,9032,9033,9034,9035,9036,9037,9038,9039,9040",
		"TRACKING_STRATEGY":                              "segment",
		"WORKER_ENVIRONMENT":                             "kubernetes",
		"WORKSPACE_DOCKER_MOUNT":                         "airbyte_workspace",
		"WORKSPACE_ROOT":                                 "/workspace",
		"WORKFLOW_FAILURE_RESTART_DELAY_SECONDS":         "",
		"USE_STREAM_CAPABLE_STATE":                       "true",
		"AUTO_DETECT_SCHEMA":                             "true",
		"WORKERS_MICRONAUT_ENVIRONMENTS":                 "control-plane",
		"CRON_MICRONAUT_ENVIRONMENTS":                    "control-plane",
		"SHOULD_RUN_NOTIFY_WORKFLOWS":                    "true",
	}

	var clusterConfig = instance.Spec.ClusterConfig
	resourceReq := &ResourceRequest{
		Client:    r.Client,
		Log:       r.Log,
		Ctx:       ctx,
		Namespace: instance.Namespace,
	}
	// s3
	logS3Ref := clusterConfig.Logs.S3.Reference
	s3Bucket, s3Connection, err := resourceReq.FetchS3ByReference(logS3Ref)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Failed to fetch s3 log bucket by reference %s", logS3Ref))
	}
	data["S3_LOG_BUCKET"] = s3Bucket.Spec.BucketName
	data["S3_LOG_BUCKET_REGION"] = s3Connection.Spec.Region
	data["S3_ENDPOINT"] = s3Connection.Spec.Endpoint
	data["S3_PATH_STYLE_ACCESS"] = strconv.FormatBool(s3Connection.Spec.PathStyle)
	//db
	if db := clusterConfig.Config.Database; db != nil {
		dbRef := db.Reference
		database, databaseConnection, err := resourceReq.FetchDbByReference(dbRef)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Failed to fetch db by reference %s", dbRef))
		}
		data["DATABASE_DB"] = database.Spec.DatabaseName
		pg := databaseConnection.Spec.Provider.Postgres
		data["DATABASE_HOST"] = pg.Host
		data["DATABASE_PORT"] = strconv.FormatInt(int64(pg.Port), 10)
		data["DATABASE_URL"] = "jdbc:postgresql://" + pg.Host + ":" + strconv.FormatInt(int64(pg.Port), 10) + "/" + database.Spec.DatabaseName

	}
	// state storage
	stateStorageType := ""
	if instance.Spec.ClusterConfig.StateStorageType != "" {
		stateStorageType = instance.Spec.ClusterConfig.StateStorageType
	}
	if stateStorage := instance.Spec.ClusterConfig.StateStorage; stateStorage != nil {
		if stateS3 := stateStorage.S3; stateS3 != nil {
			stateS3Bucket, _, err := resourceReq.FetchS3ByReference(stateS3.Reference)
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("Failed to fetch s3 state storage by reference %s", logS3Ref))
			}
			data["STATE_STORAGE_MINIO_BUCKET_NAME"] = stateS3Bucket.Spec.BucketName
		}
	}
	data["WORKER_STATE_STORAGE_TYPE"] = stateStorageType

	containerOrchestratorImage := ""
	if instance != nil && instance.Spec.Worker != nil && instance.Spec.Worker.RoleConfig.ContainerOrchestrator != nil {
		if instance.Spec.Worker.RoleConfig.ContainerOrchestrator.Image != "" {
			containerOrchestratorImage = instance.Spec.Worker.RoleConfig.ContainerOrchestrator.Image
		}
	}
	data["CONTAINER_ORCHESTRATOR_IMAGE"] = containerOrchestratorImage

	containerOrchestratorEnabled := "true"
	if instance != nil && instance.Spec.Worker != nil && instance.Spec.Worker.RoleConfig.ContainerOrchestrator != nil {
		if instance.Spec.Worker.RoleConfig.ContainerOrchestrator.Enabled {
			containerOrchestratorEnabled = strconv.FormatBool(instance.Spec.Worker.RoleConfig.ContainerOrchestrator.Enabled)
		}
	}
	data["CONTAINER_ORCHESTRATOR_ENABLED"] = containerOrchestratorEnabled
	data["WORKER_LOGS_STORAGE_TYPE"] = "S3"

	apiUrl := ""
	if instance != nil && instance.Spec.WebApp != nil {
		if instance.Spec.WebApp.RoleConfig.ApiUrl != "" {
			apiUrl = instance.Spec.WebApp.RoleConfig.ApiUrl
		}
	}
	data["API_URL"] = apiUrl

	ServerUrl := ""
	if instance != nil && instance.Spec.WebApp != nil {
		if instance.Spec.WebApp.RoleConfig.ConnectorBuilderServerUrl != "" {
			ServerUrl = instance.Spec.WebApp.RoleConfig.ConnectorBuilderServerUrl
		}
	}
	data["CONNECTOR_BUILDER_API_URL"] = ServerUrl
	data["GCS_LOG_BUCKET"] = ""
	data["GOOGLE_APPLICATION_CREDENTIALS"] = ""

	if instance != nil && instance.Spec.ClusterConfig != nil {
		if instance.Spec.ClusterConfig.Edition == "pro" {
			if instance.Spec.Keycloak != nil && instance.Spec.Keycloak.RoleConfig.Service != nil {
				data["KEYCLOAK_INTERNAL_HOST"] = instance.Name + "-airbyte-keycloak-svc:" + strconv.FormatInt(int64(instance.Spec.Keycloak.RoleConfig.Service.Port), 10)
				data["KEYCLOAK_PORT"] = strconv.FormatInt(int64(instance.Spec.Keycloak.RoleConfig.Service.Port), 10)
				data["KEYCLOAK_HOSTNAME_URL"] = "webapp-url/auth"
			}
		} else {
			data["KEYCLOAK_INTERNAL_HOST"] = "localhost"
		}
	}

	if instance != nil && instance.Spec.ClusterConfig != nil && instance.Spec.ClusterConfig.Jobs != nil && instance.Spec.ClusterConfig.Jobs.Kube != nil {
		if instance.Spec.ClusterConfig.Jobs.Kube.Annotations != nil {
			addMapToData(data, "JOB_KUBE_ANNOTATIONS", instance.Spec.ClusterConfig.Jobs.Kube.Annotations)
		}

		if instance.Spec.ClusterConfig.Jobs.Kube.Labels != nil {
			addMapToData(data, "JOB_KUBE_LABELS", instance.Spec.ClusterConfig.Jobs.Kube.Labels)
		}

		if instance.Spec.ClusterConfig.Jobs.Kube.NodeSelector != nil {
			addMapToData(data, "JOB_KUBE_NODE_SELECTOR", instance.Spec.ClusterConfig.Jobs.Kube.NodeSelector)
		}

		if instance.Spec.ClusterConfig.Jobs.Kube.MainContainerImagePullSecret != "" {
			data["JOB_KUBE_MAIN_CONTAINER_IMAGE_PULL_SECRET"] = instance.Spec.ClusterConfig.Jobs.Kube.MainContainerImagePullSecret
		}

		if instance.Spec.ClusterConfig.Jobs.Kube.Tolerations != nil {
			tolerations, err := json.Marshal(instance.Spec.ClusterConfig.Jobs.Kube.Tolerations)
			if err != nil {
				// 处理错误
				r.Log.Error(err, "Failed to json-marshal tolerations")
				return nil, err
			}
			data["JOB_KUBE_TOLERATIONS"] = string(tolerations)
		}

		if instance.Spec.ClusterConfig.Jobs.Kube.Images.Busybox != "" {
			data["JOB_KUBE_BUSYBOX_IMAGE"] = instance.Spec.ClusterConfig.Jobs.Kube.Images.Busybox
		}

		if instance.Spec.ClusterConfig.Jobs.Kube.Images.Socat != "" {
			data["JOB_KUBE_SOCAT_IMAGE"] = instance.Spec.ClusterConfig.Jobs.Kube.Images.Socat
		}

		if instance.Spec.ClusterConfig.Jobs.Kube.Images.Curl != "" {
			data["JOB_KUBE_CURL_IMAGE"] = instance.Spec.ClusterConfig.Jobs.Kube.Images.Curl
		}
	}

	cpuLimit := ""
	if instance != nil && instance.Spec.ClusterConfig != nil && instance.Spec.ClusterConfig.Jobs != nil && instance.Spec.ClusterConfig.Jobs.Resources != nil && instance.Spec.ClusterConfig.Jobs.Resources.Limits != nil {
		if cpu := instance.Spec.ClusterConfig.Jobs.Resources.Limits.Cpu(); cpu != nil {
			cpuLimit = cpu.String()
		}
	}
	data["JOB_MAIN_CONTAINER_CPU_LIMIT"] = cpuLimit

	cpuRequest := ""
	if instance != nil && instance.Spec.ClusterConfig != nil && instance.Spec.ClusterConfig.Jobs != nil && instance.Spec.ClusterConfig.Jobs.Resources != nil && instance.Spec.ClusterConfig.Jobs.Resources.Requests != nil {
		if cpu := instance.Spec.ClusterConfig.Jobs.Resources.Requests.Cpu(); cpu != nil {
			cpuRequest = cpu.String()
		}
	}
	data["JOB_MAIN_CONTAINER_CPU_REQUEST"] = cpuRequest

	memoryLimit := ""
	if instance != nil && instance.Spec.ClusterConfig != nil && instance.Spec.ClusterConfig.Jobs != nil && instance.Spec.ClusterConfig.Jobs.Resources != nil && instance.Spec.ClusterConfig.Jobs.Resources.Limits != nil {
		if mem := instance.Spec.ClusterConfig.Jobs.Resources.Limits.Memory(); mem != nil {
			memoryLimit = mem.String()
		}
	}
	data["JOB_MAIN_CONTAINER_MEMORY_LIMIT"] = memoryLimit

	memoryRequest := ""
	if instance != nil && instance.Spec.ClusterConfig != nil && instance.Spec.ClusterConfig.Jobs != nil && instance.Spec.ClusterConfig.Jobs.Resources != nil && instance.Spec.ClusterConfig.Jobs.Resources.Requests != nil {
		if mem := instance.Spec.ClusterConfig.Jobs.Resources.Requests.Memory(); mem != nil {
			memoryRequest = mem.String()
		}
	}
	data["JOB_MAIN_CONTAINER_MEMORY_REQUEST"] = memoryRequest

	runDatabaseMigrationOnStartup := "true"
	if instance.Spec.AirbyteBootloader.RoleConfig.RunDatabaseMigrationsOnStartup != nil {
		runDatabaseMigrationOnStartup = strconv.FormatBool(*instance.Spec.AirbyteBootloader.RoleConfig.RunDatabaseMigrationsOnStartup)
	}
	data["RUN_DATABASE_MIGRATION_ON_STARTUP"] = runDatabaseMigrationOnStartup

	webAppUrl := ""
	if instance != nil && instance.Spec.WebApp != nil {
		if instance.Spec.WebApp.RoleConfig.Url != "" {
			webAppUrl = instance.Spec.WebApp.RoleConfig.Url
		} else if instance.Spec.WebApp.RoleConfig.Service != nil {
			webAppUrl = fmt.Sprintf("http://%s-airbyte-webapp-svc:%d", instance.Name, instance.Spec.WebApp.RoleConfig.Service.Port)
		}
	}
	data["WEBAPP_URL"] = webAppUrl

	metricClient := ""
	if instance != nil && instance.Spec.ClusterConfig != nil && instance.Spec.ClusterConfig.Metrics != nil {
		if instance.Spec.ClusterConfig.Metrics.MetricClient != "" {
			metricClient = instance.Spec.ClusterConfig.Metrics.MetricClient
		}
	}
	data["METRIC_CLIENT"] = metricClient

	otelCollectorEndpoint := ""
	if instance != nil && instance.Spec.ClusterConfig != nil && instance.Spec.ClusterConfig.Metrics != nil {
		if instance.Spec.ClusterConfig.Metrics.OtelCollectorEndpoint != "" {
			otelCollectorEndpoint = instance.Spec.ClusterConfig.Metrics.OtelCollectorEndpoint
		}
	}
	data["OTEL_COLLECTOR_ENDPOINT"] = otelCollectorEndpoint

	activityMaxAttempt := ""
	activityInitialDelayBetweenAttemptsSeconds := ""
	activityMaxDelayBetweenAttemptsSeconds := ""
	maxNotifyWorkers := "5"
	if instance != nil && instance.Spec.Worker != nil {
		if instance.Spec.Worker.RoleConfig.ActivityMaxAttempt != "" {
			activityMaxAttempt = instance.Spec.Worker.RoleConfig.ActivityMaxAttempt
		}

		if instance.Spec.Worker.RoleConfig.ActivityInitialDelayBetweenAttemptsSeconds != "" {
			activityInitialDelayBetweenAttemptsSeconds = instance.Spec.Worker.RoleConfig.ActivityInitialDelayBetweenAttemptsSeconds
		}

		if instance.Spec.Worker.RoleConfig.ActivityMaxDelayBetweenAttemptsSeconds != "" {
			activityMaxDelayBetweenAttemptsSeconds = instance.Spec.Worker.RoleConfig.ActivityMaxDelayBetweenAttemptsSeconds
		}

		if instance.Spec.Worker.RoleConfig.MaxNotifyWorkers != "" {
			maxNotifyWorkers = instance.Spec.Worker.RoleConfig.MaxNotifyWorkers
		}
	}
	data["ACTIVITY_MAX_ATTEMPT"] = activityMaxAttempt
	data["ACTIVITY_INITIAL_DELAY_BETWEEN_ATTEMPTS_SECONDS"] = activityInitialDelayBetweenAttemptsSeconds
	data["ACTIVITY_MAX_DELAY_BETWEEN_ATTEMPTS_SECONDS"] = activityMaxDelayBetweenAttemptsSeconds
	data["MAX_NOTIFY_WORKERS"] = maxNotifyWorkers

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.GetNameWithSuffix("env"),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Data: data,
	}
	err = ctrl.SetControllerReference(instance, configMap, schema)
	if err != nil {
		r.Log.Error(err, "Failed to set controller reference for configmap")
		return nil, err
	}
	return configMap, nil
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
func (r *AirbyteReconciler) makeAdminServiceAccount(instance *stackv1alpha1.Airbyte, schema *runtime.Scheme) *corev1.ServiceAccount {
	labels := instance.GetLabels()
	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.GetNameWithSuffix("admin"),
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
			Name:      instance.GetNameWithSuffix("admin-role"),
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
			Name:      instance.GetNameWithSuffix("admin-binding"),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      instance.GetNameWithSuffix("admin"),
				Namespace: instance.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "Role",
			Name:     instance.GetNameWithSuffix("admin-role"),
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
	if secretMap == nil {
		return nil
	}
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

// reconcile cluster yml content config map
func (r *AirbyteReconciler) reconcileClusterYmlContentConfigMap(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	reconcileParams := r.createReconcileParams(ctx, nil, instance, r.extractClusterYmlContentConfigMap)
	if err := reconcileParams.createOrUpdateResourceNoGroup(); err != nil {
		return err
	}
	return nil
}

// extract cluster yml content config map
func (r *AirbyteReconciler) extractClusterYmlContentConfigMap(params ExtractorParams) (client.Object, error) {
	instance := params.instance.(*stackv1alpha1.Airbyte)
	labels := instance.GetLabels()
	schema := params.scheme
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.GetNameWithSuffix("yml"),
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
		return nil, err
	}
	return configMap, nil
}

// reconcile airbyte log storage secret
func (r *AirbyteReconciler) reconcileClusterSecret(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	roleGroups := convertRoleGroupToRoleConfigObject(instance, Temporal)
	reconcileParams := r.createReconcileParams(ctx, roleGroups, instance, r.extractClusterSecret)
	if err := reconcileParams.createOrUpdateResource(); err != nil {
		return err
	}
	return nil
}

// extract airbyte global secret
func (r *AirbyteReconciler) extractClusterSecret(params ExtractorParams) (client.Object, error) {
	instance := params.instance.(*stackv1alpha1.Airbyte)
	mergedLabels := r.mergeLabels(params.roleGroup, instance.GetLabels(), params.cluster)
	if logs3 := instance.Spec.ClusterConfig.Logs.S3; logs3 != nil {
		s3rsc := &opgo.S3Bucket{
			ObjectMeta: metav1.ObjectMeta{
				Name: logs3.Reference,
			},
		}
		resourceReq := &ResourceRequest{
			Ctx:       params.ctx,
			Client:    r.Client,
			Namespace: instance.Namespace,
			Log:       r.Log,
		}
		s3Connection, err := resourceReq.FetchS3(s3rsc)
		if err != nil {
			return nil, err
		}
		// Create a Data map to hold the secret data
		data := make(map[string][]byte)
		r.makeS3Data(s3Connection, &data)

		if db := instance.Spec.ClusterConfig.Config.Database; db != nil {
			dbRef := db.Reference
			database, databaseConnection, err := resourceReq.FetchDbByReference(dbRef)
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("Failed to fetch db by reference %s", dbRef))
			}
			data["DATABASE_DB"] = []byte(database.Spec.DatabaseName)
			data["DATABASE_USER"] = []byte(database.Spec.Credential.Username)
			data["DATABASE_PASSWORD"] = []byte(database.Spec.Credential.Password)
			pg := databaseConnection.Spec.Provider.Postgres
			data["DATABASE_HOST"] = []byte(pg.Host)
			data["DATABASE_PORT"] = []byte(strconv.FormatInt(int64(pg.Port), 10))
			data["DATABASE_URL"] = []byte("jdbc:postgresql://" + pg.Host + ":" + strconv.FormatInt(int64(pg.Port), 10) + "/" + database.Spec.DatabaseName)
		}

		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      createNameForRoleGroup(instance, "secrets", ""),
				Namespace: instance.Namespace,
				Labels:    mergedLabels,
			},
			Type: corev1.SecretTypeOpaque,
			Data: data,
		}
		return secret, nil
	}
	return nil, nil
}

func (r *AirbyteReconciler) makeS3Data(s3Connection *opgo.S3Connection, data *map[string][]byte) {
	if s3Connection != nil {
		s3Credential := s3Connection.Spec.S3Credential
		(*data)["AWS_ACCESS_KEY_ID"] = []byte(s3Credential.AccessKey)
		(*data)["AWS_SECRET_ACCESS_KEY"] = []byte(s3Credential.SecretKey)
		(*data)["AWS_DEFAULT_REGION"] = []byte(s3Connection.Spec.Region)
	}
}

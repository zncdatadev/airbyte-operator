package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/zncdata-labs/airbyte-operator/api/v1alpha1"
	"github.com/zncdata-labs/airbyte-operator/api/v1alpha1/rolegroup"
	opgo "github.com/zncdata-labs/operator-go/pkg/apis/commons/v1alpha1"
	"github.com/zncdata-labs/operator-go/pkg/util"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ExtractorParams Params of extract resource function for role group
type ExtractorParams struct {
	instance      any
	ctx           context.Context
	roleGroupName string
	cluster       rolegroup.RoleConfigObject
	roleGroup     rolegroup.RoleConfigObject
	scheme        *runtime.Scheme
}

// ReconcileParams Params of reconcile resource
type ReconcileParams struct {
	instance client.Object
	client.Client
	roleGroups         map[string]rolegroup.RoleConfigObject
	ctx                context.Context
	scheme             *runtime.Scheme
	log                logr.Logger
	roleGroupExtractor func(params ExtractorParams) (client.Object, error)
}

// create Reconcile params
func (r *AirbyteReconciler) createReconcileParams(ctx context.Context, roleGroup map[string]rolegroup.RoleConfigObject,
	instance *v1alpha1.Airbyte, method func(params ExtractorParams) (client.Object, error)) *ReconcileParams {
	return &ReconcileParams{
		instance:           instance,
		Client:             r.Client,
		roleGroups:         roleGroup,
		ctx:                ctx,
		scheme:             r.Scheme,
		log:                r.Log,
		roleGroupExtractor: method,
	}
}

// Extract resource from k8s cluster for role group ,and after collect them to a slice
// see extractResourcesByGroupName
func (r *ReconcileParams) extractResources() ([]client.Object, error) {
	var resources []client.Object
	if r.roleGroups != nil {
		for roleGroupName, roleGroup := range r.roleGroups {
			rsc, err := r.roleGroupExtractor(ExtractorParams{
				instance:      r.instance,
				ctx:           r.ctx,
				roleGroupName: roleGroupName,
				roleGroup:     roleGroup,
				scheme:        r.scheme,
			})
			if err != nil {
				return nil, err
			}
			resources = append(resources, rsc)
		}
	}
	return resources, nil
}

// Extract resource from k8s cluster for role group ,and after collect them by group name
// see extractResources
// it's return map[string]client.Object, the key is role group name, not same as extractResources, it's return []client.Object
// this is used for update instance status, such as URL of ingress
func (r *ReconcileParams) extractResourcesByGroupName() (map[string]client.Object, error) {
	resources := make(map[string]client.Object)
	if r.roleGroups != nil {
		for roleGroupName, roleGroup := range r.roleGroups {
			rsc, err := r.roleGroupExtractor(ExtractorParams{
				instance:      r.instance,
				ctx:           r.ctx,
				roleGroupName: roleGroupName,
				roleGroup:     roleGroup,
				scheme:        r.scheme,
			})
			if err != nil {
				return nil, err
			}
			resources[roleGroupName] = rsc
		}
	}
	return resources, nil
}

func (r *ReconcileParams) createOrUpdateResourceNoGroup() error {
	rsc, err := r.roleGroupExtractor(ExtractorParams{
		instance:      r.instance,
		ctx:           r.ctx,
		roleGroupName: "",
		roleGroup:     nil,
		scheme:        r.scheme,
	})
	if err != nil {
		return err
	}
	if err = CreateOrUpdate(r.ctx, r, rsc); err != nil {
		return err
	}
	return nil
}

// create or update resource
func (r *ReconcileParams) createOrUpdateResource() error {
	resources, err := r.extractResources()
	if err != nil {
		return err
	}

	for _, rsc := range resources {
		if rsc == nil {
			continue
		}

		if err = CreateOrUpdate(r.ctx, r, rsc); err != nil {
			r.log.Error(err, "Failed to create or update Resource", "resource", rsc)
			return err
		}
	}
	return nil
}

func (r *ReconcileParams) updateInstanceStatus(params []interface{}, groupName string, rsc any,
	function func(groupName string, rsc any, client client.Client, instance client.Object, ctx context.Context,
		object []interface{}) error) error {
	if err := function(groupName, rsc, r.Client, r.instance, r.ctx, params); err != nil {
		return err
	}
	if err := util.UpdateStatus(r.ctx, r.Client, r.instance); err != nil {
		return err
	}
	return nil
}

// create or update resource and update instance status
func (r *ReconcileParams) createOrUpdateResourceAndUpdateInstanceStatus(params []interface{},
	function func(groupName string, rsc any, client client.Client, instance client.Object, ctx context.Context,
		object []interface{}) error) error {
	resources, err := r.extractResourcesByGroupName()
	if err != nil {
		return err
	}

	for groupName, rsc := range resources {
		if rsc == nil {
			continue
		}

		if err = CreateOrUpdate(r.ctx, r, rsc); err != nil {
			r.log.Error(err, "Failed to create or update Resource", "resource", rsc)
			return err
		}
		// update instance status, such as URL of ingres
		err = r.updateInstanceStatus(params, groupName, rsc, function)
		if err != nil {
			r.log.Error(err, "Failed to update instance status", "resource", rsc)
			return err
		}
	}
	return nil
}

// ResourceRequest Params of fetch resource  from k8s cluster
type ResourceRequest struct {
	Ctx context.Context
	client.Client
	Namespace string
	Log       logr.Logger
}

func (r *ResourceRequest) fetchResource(obj client.Object) error {
	name := obj.GetName()
	kind := obj.GetObjectKind()
	if err := r.Get(r.Ctx, client.ObjectKey{Namespace: r.Namespace, Name: name}, obj); err != nil {
		opt := []any{"ns", r.Namespace, "name", name, "kind", kind}
		if apierrors.IsNotFound(err) {
			r.Log.Error(err, "Fetch resource NotFound", opt...)
		} else {
			r.Log.Error(err, "Fetch resource occur some unknown err", opt...)
		}
		return err
	}
	return nil
}

// FetchS3ByReference fetch s3 s3Bucket and  s3 connection  by s3 reference
func (r *ResourceRequest) FetchS3ByReference(reference string) (*opgo.S3Bucket, *opgo.S3Connection, error) {
	if reference == "" {
		return nil, nil, errors.New("s3 reference is empty")
	}
	s3rsc := &opgo.S3Bucket{
		ObjectMeta: metav1.ObjectMeta{
			Name: reference,
		},
	}
	// Fetch S3 by s3Bucket reference
	if s3Connection, err := r.FetchS3(s3rsc); err != nil {
		return nil, nil, err
	} else {
		return s3rsc, s3Connection, nil
	}
}

// FetchDbByReference fetch database and database connection by database reference
func (r *ResourceRequest) FetchDbByReference(reference string) (*opgo.Database, *opgo.DatabaseConnection, error) {
	database := &opgo.Database{
		ObjectMeta: metav1.ObjectMeta{
			Name: reference,
		},
	}
	// Fetch Database by database reference
	if dbConnection, err := r.FetchDb(database); err != nil {
		return nil, nil, err
	} else {
		return database, dbConnection, nil
	}
}

// todo: db and s3 reference can be implemented by generic way, such as interface

// FetchS3 fetch s3 connection by s3 reference, while s3 is fetched by resource request
func (r *ResourceRequest) FetchS3(s3 *opgo.S3Bucket) (*opgo.S3Connection, error) {
	// 1 - fetch exists S3 by reference
	if err := r.fetchResource(s3); err != nil {
		return nil, err
	}
	// if exist secret reference, fetch secret by reference,and override accessKey and secretKey
	if secret, err := r.FetchSecretByReference(s3.Spec.Credential.ExistSecret); err != nil {
		if secret != nil {
			s3.Spec.Credential.AccessKey = string(secret.Data["ACCESS_KEY"])
			s3.Spec.Credential.SecretKey = string(secret.Data["SECRET_KEY"])
		}
		return nil, err
	}
	//2 - fetch exist s3 connection by pre-fetch 's3.spec.name'
	s3Connection := &opgo.S3Connection{
		ObjectMeta: metav1.ObjectMeta{Name: s3.Spec.Reference},
	}
	if err := r.fetchResource(s3Connection); err != nil {
		return nil, err
	}
	return s3Connection, nil
}

// FetchDb fetch database connection by database reference, while database is fetched by resource request
func (r *ResourceRequest) FetchDb(database *opgo.Database) (*opgo.DatabaseConnection, error) {
	// 1 - fetch exists Database by reference
	if err := r.fetchResource(database); err != nil {
		return nil, err
	}
	// if exist secret reference, fetch secret by reference, and override username and password
	if secret, err := r.FetchSecretByReference(database.Spec.Credential.ExistSecret); err != nil {
		if secret != nil {
			database.Spec.Credential.Username = string(secret.Data["USERNAME"])
			database.Spec.Credential.Password = string(secret.Data["PASSWORD"])
		}
		return nil, err
	}

	//2 - fetch exist database connection by pre-fetch 'database.spec.name'
	dbConnection := &opgo.DatabaseConnection{
		ObjectMeta: metav1.ObjectMeta{Name: database.Spec.Reference},
	}
	if err := r.fetchResource(dbConnection); err != nil {
		return nil, err
	}
	return dbConnection, nil
}

// FetchSecretByReference fetch secret by reference
func (r *ResourceRequest) FetchSecretByReference(secretRef string) (*corev1.Secret, error) {
	if secretRef == "" {
		return nil, nil
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: secretRef,
		},
	}
	if err := r.fetchResource(secret); err != nil {
		return nil, err
	}
	data := secret.Data
	if len(data) == 0 {
		return nil, fmt.Errorf("data of secret reference '%s' is empty", secretRef)
	}
	return secret, nil
}

// air byte role reconcile task

type AirByteRole string

const (
	Cluster                AirByteRole = "Cluster"
	Server                 AirByteRole = "Server"
	Cron                   AirByteRole = "Cron"
	ApiServer              AirByteRole = "ApiServer"
	ConnectorBuilderServer AirByteRole = "ConnectorBuilderServer"
	PodSweeper             AirByteRole = "PodSweeper"
	Temporal               AirByteRole = "Temporal"
	Webapp                 AirByteRole = "Webapp"
	Worker                 AirByteRole = "Worker"
	LogStorage             AirByteRole = "LogStorage"
)

const (
	ResourceReconcileErrTemp = "Reconcile %s %s failed"
)

type AirByteRoleReconcileTask struct {
	Role     AirByteRole
	function func(ctx context.Context, instance *v1alpha1.Airbyte) error
}

func reconcileRoleTask(tasks []AirByteRoleReconcileTask, ctx context.Context, instance *v1alpha1.Airbyte,
	log logr.Logger, errLogTemp string, resourceType ResourceType) error {
	for _, t := range tasks {
		if err := t.function(ctx, instance); err != nil {
			errMsg := fmt.Sprintf(errLogTemp, t.Role, resourceType)
			log.Error(err, errMsg)
			return err
		}
	}
	return nil
}

// convert origin role group to role group object
func convertRoleGroupToRoleConfigObject(instance *v1alpha1.Airbyte, roleType AirByteRole) map[string]rolegroup.RoleConfigObject {
	var roleGroups = make(map[string]rolegroup.RoleConfigObject)
	switch roleType {
	case ApiServer:
		if origin := instance.Spec.ApiServer.RoleGroups; origin != nil {
			for name, group := range origin {
				roleGroups[name] = group
			}
		}
	case ConnectorBuilderServer:
		if origin := instance.Spec.ConnectorBuilderServer.RoleGroups; origin != nil {
			for name, group := range origin {
				roleGroups[name] = group
			}
		}
	case Cron:
		if origin := instance.Spec.Cron.RoleGroups; origin != nil {
			for name, group := range origin {
				roleGroups[name] = group
			}
		}
	case PodSweeper:
		if origin := instance.Spec.PodSweeper.RoleGroups; origin != nil {
			for name, group := range origin {
				roleGroups[name] = group
			}
		}
	case Server:
		if origin := instance.Spec.Server.RoleGroups; origin != nil {
			for name, group := range origin {
				roleGroups[name] = group
			}
		}
	case Temporal:
		if origin := instance.Spec.Temporal.RoleGroups; origin != nil {
			for name, group := range origin {
				roleGroups[name] = group
			}
		}
	case Webapp:
		if origin := instance.Spec.WebApp.RoleGroups; origin != nil {
			for name, group := range origin {
				roleGroups[name] = group
			}
		}
	case Worker:
		if origin := instance.Spec.Worker.RoleGroups; origin != nil {
			for name, group := range origin {
				roleGroups[name] = group
			}
		}
	}
	return roleGroups
}

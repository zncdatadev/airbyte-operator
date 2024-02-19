package controller

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/zncdata-labs/operator-go/pkg/status"
	"github.com/zncdata-labs/operator-go/pkg/util"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ResourceType string

const (
	Pod            ResourceType = "Pod"
	Deployment     ResourceType = "Deployment"
	Service        ResourceType = "Service"
	Secret         ResourceType = "Secret"
	ConfigMap      ResourceType = "ConfigMap"
	Pvc            ResourceType = "Pvc"
	Ingress        ResourceType = "Ingress"
	ServiceAccount ResourceType = "ServiceAccount"
	Role           ResourceType = "Role"
	RoleBinding                 = "RoleBinding"
)

var (
	ResourceMapper = map[ResourceType]string{
		Pod:            "ReconcilePod",
		Deployment:     status.ConditionTypeReconcileDeployment,
		Service:        status.ConditionTypeReconcileService,
		Secret:         status.ConditionTypeReconcileSecret,
		ConfigMap:      status.ConditionTypeReconcileConfigMap,
		Pvc:            status.ConditionTypeReconcilePVC,
		Ingress:        status.ConditionTypeReconcileIngress,
		Role:           "ReconcileRole",
		ServiceAccount: "ReconcileServiceAccount", //todo: update opgo
		RoleBinding:    "ReconcileRoleBinding",
	}
)

type ReconcileTask[T client.Object] struct {
	resourceName  ResourceType
	reconcileFunc func(ctx context.Context, instance T) error
}

type TaskReconcilePara[T client.Object] struct {
	ctx            context.Context
	instance       T
	log            logr.Logger
	client         client.Client
	instanceStatus *status.Status
	serverName     string
}

func ReconcileTasks[T client.Object](tasks *[]ReconcileTask[T], params TaskReconcilePara[T]) error {
	log := params.log
	for _, task := range *tasks {
		jobFunc := task.reconcileFunc
		if err := jobFunc(params.ctx, params.instance); err != nil {
			log.Error(err, fmt.Sprintf("unable to reconcile %s", task.resourceName))
			return err
		}
		instance := params.instance
		if updated := params.instanceStatus.SetStatusCondition(v1.Condition{
			Type:               ResourceMapper[ResourceType(task.resourceName)],
			Status:             v1.ConditionTrue,
			Reason:             status.ConditionReasonRunning,
			Message:            createSuccessMessage(params.serverName, task.resourceName),
			ObservedGeneration: instance.GetGeneration(),
		}); updated {
			err := util.UpdateStatus(params.ctx, params.client, instance)
			if err != nil {
				log.Error(err, createUpdateErrorMessage(task.resourceName))
				return err
			}
		}
	}
	return nil
}

func createSuccessMessage(serverName string, resourceName ResourceType) string {
	return fmt.Sprintf("%sServer's %s is running", serverName, resourceName)
}

func createUpdateErrorMessage(resourceName ResourceType) string {
	return fmt.Sprintf("unable to update status for %s", resourceName)
}

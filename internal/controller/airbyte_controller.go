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

package controller

import (
	"context"
	"github.com/go-logr/logr"
	stackv1alpha1 "github.com/zncdata-labs/airbyte-operator/api/v1alpha1"
	"github.com/zncdata-labs/operator-go/pkg/status"
	"github.com/zncdata-labs/operator-go/pkg/util"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// AirbyteReconciler reconciles an Airbyte object
type AirbyteReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

//+kubebuilder:rbac:groups=stack.zncdata.net,resources=airbytes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=stack.zncdata.net,resources=airbytes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=stack.zncdata.net,resources=airbytes/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extensions,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Airbyte object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *AirbyteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.Log.Info("Reconciling instance")

	existInstance, err := r.getAirbyteInstance(ctx, req)
	if err != nil {
		return ctrl.Result{}, err
	}
	// if instance not found, return
	if existInstance == nil {
		return ctrl.Result{}, nil
	}

	instance := existInstance.DeepCopy()

	if err := r.handleStatusConditions(ctx, instance); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.executeReconcileTasks(ctx, instance); err != nil {
		return ctrl.Result{}, err
	}

	r.Log.Info("Successfully reconciled Airbyte")
	return ctrl.Result{}, nil
}

func (r *AirbyteReconciler) getAirbyteInstance(ctx context.Context, req ctrl.Request) (*stackv1alpha1.Airbyte, error) {
	instance := &stackv1alpha1.Airbyte{}

	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		if client.IgnoreNotFound(err) != nil {
			r.Log.Error(err, "unable to fetch Airbyte")
			return nil, err
		}
		r.Log.Info("Airbyte resource not found. Ignoring since object must be deleted")
		return nil, nil
	}

	r.Log.Info("Airbyte found", "Name", instance.Name)
	return instance, nil
}

func (r *AirbyteReconciler) handleStatusConditions(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	readCondition := apimeta.FindStatusCondition(instance.Status.Conditions, status.ConditionTypeProgressing)
	if readCondition == nil || readCondition.ObservedGeneration != instance.GetGeneration() {
		instance.InitStatusConditions()
		if err := util.UpdateStatus(ctx, r.Client, instance); err != nil {
			return err
		}
	}
	return nil
}

func (r *AirbyteReconciler) executeReconcileTasks(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	tasks := []ReconcileTask[*stackv1alpha1.Airbyte]{
		{resourceName: "Pod", reconcileFunc: r.reconcileBootloaderPod},
		{resourceName: "Secret", reconcileFunc: r.reconcileSecret},
		{resourceName: "ConfigMap", reconcileFunc: r.reconcileConfigMap},
		{resourceName: "ServiceAccount", reconcileFunc: r.reconcileServiceAccount},
		{resourceName: "Role", reconcileFunc: r.reconcileRole},
		{resourceName: "RoleBinding", reconcileFunc: r.reconcileRoleBinding},
		{resourceName: "Deployment", reconcileFunc: r.reconcileDeployment},
		{resourceName: "Service", reconcileFunc: r.reconcileService},
		{resourceName: "Ingress", reconcileFunc: r.reconcileIngress},
	}
	reconcileParams := TaskReconcilePara[*stackv1alpha1.Airbyte]{
		ctx:            ctx,
		instance:       instance,
		log:            r.Log,
		client:         r.Client,
		instanceStatus: &instance.Status,
		serverName:     "Airbyte",
	}
	return ReconcileTasks(&tasks, reconcileParams)
}

// SetupWithManager sets up the controller with the Manager.
func (r *AirbyteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stackv1alpha1.Airbyte{}).
		Complete(r)
}

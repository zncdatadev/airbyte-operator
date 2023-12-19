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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"

	stackv1alpha1 "github.com/zncdata-labs/airbyte-operator/api/v1alpha1"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// AirbyteReconciler reconciles a Airbyte object
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

	airbyte := &stackv1alpha1.Airbyte{}

	if err := r.Get(ctx, req.NamespacedName, airbyte); err != nil {
		if client.IgnoreNotFound(err) != nil {
			r.Log.Error(err, "unable to fetch Airbyte")
			return ctrl.Result{}, err
		}
		r.Log.Info("Airbyte resource not found. Ignoring since object must be deleted")
		return ctrl.Result{}, nil
	}

	// Get the status condition, if it exists and its generation is not the
	//same as the Airbyte's generation, reset the status conditions
	readCondition := apimeta.FindStatusCondition(airbyte.Status.Conditions, stackv1alpha1.ConditionTypeProgressing)
	if readCondition == nil || readCondition.ObservedGeneration != airbyte.GetGeneration() {
		airbyte.InitStatusConditions()

		if err := r.UpdateStatus(ctx, airbyte); err != nil {
			return ctrl.Result{}, err
		}
	}

	r.Log.Info("Airbyte found", "Name", airbyte.Name)

	if err := r.reconcileDeployment(ctx, airbyte); err != nil {
		r.Log.Error(err, "unable to reconcile Deployment")
		return ctrl.Result{}, err
	}

	if err := r.reconcileService(ctx, airbyte); err != nil {
		r.Log.Error(err, "unable to reconcile Service")
		return ctrl.Result{}, err
	}

	if err := r.reconcileSecret(ctx, airbyte); err != nil {
		r.Log.Error(err, "unable to reconcile Secret")
		return ctrl.Result{}, err
	}

	if err := r.reconcileConfigMap(ctx, airbyte); err != nil {
		r.Log.Error(err, "unable to reconcile ConfigMap")
		return ctrl.Result{}, err
	}

	if err := r.reconcileServiceAccount(ctx, airbyte); err != nil {
		r.Log.Error(err, "unable to reconcile ServiceAccount")
		return ctrl.Result{}, err
	}

	if err := r.reconcileRole(ctx, airbyte); err != nil {
		r.Log.Error(err, "unable to reconcile Role")
		return ctrl.Result{}, err
	}

	if err := r.reconcileRoleBinding(ctx, airbyte); err != nil {
		r.Log.Error(err, "unable to reconcile RoleBinding")
		return ctrl.Result{}, err
	}

	airbyte.SetStatusCondition(metav1.Condition{
		Type:               stackv1alpha1.ConditionTypeAvailable,
		Status:             metav1.ConditionTrue,
		Reason:             stackv1alpha1.ConditionReasonRunning,
		Message:            "Airbyte is running",
		ObservedGeneration: airbyte.GetGeneration(),
	})

	if err := r.UpdateStatus(ctx, airbyte); err != nil {
		return ctrl.Result{}, err
	}

	r.Log.Info("Successfully reconciled Airbyte")
	return ctrl.Result{}, nil

}

// UpdateStatus updates the status of the Airbyte resource
// https://stackoverflow.com/questions/76388004/k8s-controller-update-status-and-condition
func (r *AirbyteReconciler) UpdateStatus(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		return r.Status().Update(ctx, instance)
		//return r.Status().Patch(ctx, instance, client.MergeFrom(instance))
	})

	if retryErr != nil {
		r.Log.Error(retryErr, "Failed to update vfm status after retries")
		return retryErr
	}

	r.Log.V(1).Info("Successfully patched object status")
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AirbyteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stackv1alpha1.Airbyte{}).
		Complete(r)
}

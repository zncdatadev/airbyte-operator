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
	"reflect"

	"github.com/go-logr/logr"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	stackv1alpha1 "github.com/zncdata-labs/airbyte-operator/api/v1alpha1"
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
	logger := log.FromContext(ctx)

	logger.Info("Reconciling instance")

	airbyte := &stackv1alpha1.Airbyte{}
	if err := r.Get(ctx, req.NamespacedName, airbyte); err != nil {
		if client.IgnoreNotFound(err) != nil {
			logger.Error(err, "unable to fetch instance")
			return ctrl.Result{}, err
		}
		logger.Info("Instance deleted")
		return ctrl.Result{}, nil
	}
	print(airbyte.Spec.Replicas)

	logger.Info("Instance found", "Name", airbyte.Name)

	if len(airbyte.Status.Conditions) == 0 {
		airbyte.Status.Nodes = []string{}
		airbyte.Status.Conditions = append(airbyte.Status.Conditions, corev1.ComponentCondition{
			Type:   corev1.ComponentHealthy,
			Status: corev1.ConditionFalse,
		})
		err := r.Status().Update(ctx, airbyte)
		if err != nil {
			logger.Error(err, "unable to update instance status")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	} else if airbyte.Status.Conditions[0].Status == corev1.ConditionTrue {
		airbyte.Status.Conditions[0].Status = corev1.ConditionFalse
		err := r.Status().Update(ctx, airbyte)
		if err != nil {
			logger.Error(err, "unable to update instance status")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	if err := r.reconcilePVC(ctx, airbyte); err != nil {
		logger.Error(err, "unable to reconcile PVC")
		return ctrl.Result{}, err
	}

	if err := r.reconcileDeployment(ctx, airbyte); err != nil {
		logger.Error(err, "unable to reconcile Deployment")
		return ctrl.Result{}, err
	}

	if err := r.reconcileService(ctx, airbyte); err != nil {
		logger.Error(err, "unable to reconcile Service")
		return ctrl.Result{}, err
	}

	if err := r.reconcileSecret(ctx, airbyte); err != nil {
		logger.Error(err, "unable to reconcile Secret")
		return ctrl.Result{}, err
	}

	if err := r.reconcileConfigMap(ctx, airbyte); err != nil {
		logger.Error(err, "unable to reconcile ConfigMap")
		return ctrl.Result{}, err
	}

	podList := &corev1.PodList{}
	if err := r.List(ctx, podList, &client.ListOptions{Namespace: airbyte.Namespace, LabelSelector: labels.SelectorFromSet(airbyte.GetLabels())}); err != nil {
		logger.Error(err, "unable to list pods")
		return ctrl.Result{}, err
	}

	podNames := getPodNames(podList.Items)

	if !reflect.DeepEqual(podNames, airbyte.Status.Nodes) {
		logger.Info("Updating instance status", "nodes", podNames)
		airbyte.Status.Nodes = podNames
		airbyte.Status.Conditions[0].Status = corev1.ConditionTrue
		err := r.Status().Update(ctx, airbyte)
		if err != nil {
			logger.Error(err, "unable to update instance status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil

}

func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

// SetupWithManager sets up the controller with the Manager.
func (r *AirbyteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stackv1alpha1.Airbyte{}).
		Complete(r)
}

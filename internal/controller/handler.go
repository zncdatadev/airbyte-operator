package controller

import (
	"context"
	stackv1alpha1 "github.com/zncdata-labs/airbyte-operator/api/v1alpha1"
)

func (r *AirbyteReconciler) ReconcilePod(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	// list tasks of reconcile pod
	tasks := []AirByteRoleReconcileTask{
		{
			Role:     Bootloader,
			function: r.reconcileBootloaderPod,
		},
	}
	// reconcile pod tasks
	if err := reconcileRoleTask(tasks, ctx, instance, r.Log, ResourceReconcileErrTemp, Pod); err != nil {
		return err
	}
	return nil
}

func (r *AirbyteReconciler) reconcileIngress(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	// list tasks of reconcile ingress
	tasks := []AirByteRoleReconcileTask{
		{
			Role:     Webapp,
			function: r.reconcileWebAppIngress,
		},
	}
	// reconcile ingress tasks
	if err := reconcileRoleTask(tasks, ctx, instance, r.Log, ResourceReconcileErrTemp, Ingress); err != nil {
		return err
	}
	return nil
}

func (r *AirbyteReconciler) reconcileService(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	// list tasks of reconcile service
	tasks := []AirByteRoleReconcileTask{
		{
			Role:     Server,
			function: r.reconcileServerService,
		},
		{
			Role:     ConnectorBuilderServer,
			function: r.reconcileConnectorBuilderServerService,
		},
		{
			Role:     Temporal,
			function: r.reconcileTemporalService,
		},
		{
			Role:     Webapp,
			function: r.reconcileWebAppService,
		},
		{
			Role:     ApiServer,
			function: r.reconcileApiServerService,
		},
	}
	// reconcile service tasks
	if err := reconcileRoleTask(tasks, ctx, instance, r.Log, ResourceReconcileErrTemp, Service); err != nil {
		return err
	}
	return nil
}

func (r *AirbyteReconciler) reconcileDeployment(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	// list tasks of reconcile deployment
	tasks := []AirByteRoleReconcileTask{
		{
			Role:     Server,
			function: r.reconcileServerDeployment,
		},
		{
			Role:     Worker,
			function: r.reconcileWorkerDeployment,
		},
		{
			Role:     ApiServer,
			function: r.reconcileAirbyteApiServerDeployment,
		},
		{
			Role:     ConnectorBuilderServer,
			function: r.reconcileConnectorBuilderServerDeployment,
		},
		{
			Role:     Cron,
			function: r.reconcileCronDeployment,
		},
		{
			Role:     PodSweeper,
			function: r.reconcilePodSweeperDeployment,
		},
		{
			Role:     Temporal,
			function: r.reconcileTemporalDeployment,
		},
		{
			Role:     Webapp,
			function: r.reconcileWebAppDeployment,
		},
	}
	// reconcile deployment tasks
	if err := reconcileRoleTask(tasks, ctx, instance, r.Log, ResourceReconcileErrTemp, Deployment); err != nil {
		return err
	}
	return nil
}

func (r *AirbyteReconciler) reconcileSecret(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	//list tasks of reconcile secret
	tasks := []AirByteRoleReconcileTask{
		{
			Role:     Webapp,
			function: r.reconcileWebAppSecret,
		},
		{
			Role:     ConnectorBuilderServer,
			function: r.reconcileConnectorBuilderServerSecret,
		},
		{
			Role:     ApiServer,
			function: r.reconcileAirbyteApiServerSecret,
		},
		{
			Role:     Temporal,
			function: r.reconcileTemporalSecret,
		},
		{
			Role:     LogStorage,
			function: r.reconcileClusterSecret,
		},
	}
	// reconcile secret tasks
	if err := reconcileRoleTask(tasks, ctx, instance, r.Log, ResourceReconcileErrTemp, Secret); err != nil {
		return err
	}
	return nil
}

func (r *AirbyteReconciler) reconcileConfigMap(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	// list tasks of reconcile configmap
	tasks := []AirByteRoleReconcileTask{
		// cluster config
		{
			Role:     Cluster,
			function: r.reconcileClusterConfigMap,
		},
		{
			Role:     PodSweeper,
			function: r.reconcilePodSweeperConfigMap,
		},
		{
			Role:     Temporal,
			function: r.reconcileTemporalConfigMap,
		},

		// cluster yml content config
		{
			Role:     Cluster,
			function: r.reconcileClusterYmlContentConfigMap,
		},
	}
	// reconcile configmap tasks
	if err := reconcileRoleTask(tasks, ctx, instance, r.Log, ResourceReconcileErrTemp, ConfigMap); err != nil {
		return err
	}
	return nil
}

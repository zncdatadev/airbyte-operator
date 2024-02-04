package controller

import (
	"context"
	"fmt"
	stackv1alpha1 "github.com/zncdata-labs/airbyte-operator/api/v1alpha1"
	"github.com/zncdata-labs/operator-go/pkg/errors"
	"github.com/zncdata-labs/operator-go/pkg/status"
	"github.com/zncdata-labs/operator-go/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// reconcile Airbyte webapp service
func (r *AirbyteReconciler) reconcileWebAppService(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	roleGroups := convertRoleGroupToRoleConfigObject(instance, Webapp)
	reconcileParams := r.createReconcileParams(ctx, roleGroups, instance, r.extractWebAppServiceForRoleGroup)
	if err := reconcileParams.createOrUpdateResource(); err != nil {
		return err
	}
	return nil
}

// extract airbyte webapp service for role group
func (r *AirbyteReconciler) extractWebAppServiceForRoleGroup(params ExtractorParams) (client.Object, error) {
	instance := params.instance.(*stackv1alpha1.Airbyte)
	webapp := instance.Spec.WebApp
	groupCfg := params.roleGroup
	roleCfg := webapp.RoleConfig
	clusterCfg := params.cluster
	roleGroupName := params.roleGroupName
	mergedLabels := r.mergeLabels(groupCfg, instance.GetLabels(), clusterCfg)
	schema := params.scheme

	port, serviceType, annotations := getServiceInfo(groupCfg, roleCfg, clusterCfg)
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        createNameForRoleGroup(instance, "webapp", roleGroupName),
			Namespace:   instance.Namespace,
			Labels:      mergedLabels,
			Annotations: annotations,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port:     port,
					Name:     "http",
					Protocol: "TCP",
				},
			},
			Selector: mergedLabels,
			Type:     serviceType,
		},
	}
	err := ctrl.SetControllerReference(instance, svc, schema)
	if err != nil {
		r.Log.Error(err, "Failed to set controller reference for service")
		return nil, errors.Wrap(err, "Failed to set controller reference for service")
	}
	return svc, nil
}

// reconcile webapp deployment
func (r *AirbyteReconciler) reconcileWebAppDeployment(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	roleGroups := convertRoleGroupToRoleConfigObject(instance, Webapp)
	reconcileParams := r.createReconcileParams(ctx, roleGroups, instance, r.extractWebAppDeploymentForRoleGroup)
	if err := reconcileParams.createOrUpdateResource(); err != nil {
		return err
	}
	return nil
}

// extract webapp deployment for role group
func (r *AirbyteReconciler) extractWebAppDeploymentForRoleGroup(params ExtractorParams) (client.Object, error) {
	instance := params.instance.(*stackv1alpha1.Airbyte)
	server := instance.Spec.ApiServer
	groupCfg := params.roleGroup
	roleCfg := server.RoleConfig
	clusterCfg := params.cluster
	roleGroupName := params.roleGroupName
	mergedLabels := r.mergeLabels(groupCfg, instance.GetLabels(), clusterCfg)
	schema := params.scheme

	realGroupCfg := groupCfg.(*stackv1alpha1.WebAppRoleConfigSpec)
	image, securityContext, replicas, resources, containerPorts := getDeploymentInfo(groupCfg, roleCfg, clusterCfg)
	envVars := mergeEnvVarsForWebAppDeployment(instance, realGroupCfg, roleGroupName)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      createNameForRoleGroup(instance, "webapp", roleGroupName),
			Namespace: instance.Namespace,
			Labels:    mergedLabels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: mergedLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: mergedLabels,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: createNameForRoleGroup(instance, "admin", ""),
					SecurityContext:    securityContext,
					Containers: []corev1.Container{
						{
							Name:            instance.Name,
							Image:           image.Repository + ":" + image.Tag,
							ImagePullPolicy: image.PullPolicy,
							Resources:       *resources,
							Env:             envVars,
							Ports:           containerPorts,
						},
					},
				},
			},
		},
	}
	ConnectorBuilderServerScheduler(instance, dep, groupCfg)

	err := ctrl.SetControllerReference(instance, dep, schema)
	if err != nil {
		r.Log.Error(err, "Failed to set controller reference for deployment")
		return nil, err
	}
	return dep, nil
}

// merge env vars for deployment
func mergeEnvVarsForWebAppDeployment(instance *stackv1alpha1.Airbyte, roleGroup *stackv1alpha1.WebAppRoleConfigSpec,
	roleGroupName string) []corev1.EnvVar {

	envVarNames := []string{"TRACKING_STRATEGY", "INTERNAL_API_HOST", "KEYCLOAK_INTERNAL_HOST", "CONNECTOR_BUILDER_API_HOST",
		"AIRBYTE_VERSION", "API_URL", "CONNECTOR_BUILDER_API_URL"}
	var envVars []corev1.EnvVar

	if instance != nil && instance.Spec.ClusterConfig != nil {
		if instance.Spec.ClusterConfig.DeploymentMode == "oss" {
			for _, envVarName := range envVarNames {
				envVar := corev1.EnvVar{
					Name: envVarName,
					ValueFrom: &corev1.EnvVarSource{
						ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
							Key: envVarName,
							LocalObjectReference: corev1.LocalObjectReference{
								Name: createNameForRoleGroup(instance, "env", ""),
							},
						},
					},
				}
				envVars = append(envVars, envVar)
			}
		}
	}

	if roleGroup.Config != nil && roleGroup.Secret != nil {
		envVars = appendEnvVarsFromSecret(envVars, *roleGroup.Secret, "webapp-secrets", roleGroupName)
	} else if instance.Spec.WebApp != nil && instance.Spec.WebApp.RoleConfig.Secret != nil {
		envVars = appendEnvVarsFromSecret(envVars, *instance.Spec.WebApp.RoleConfig.Secret, "webapp-secrets", roleGroupName)
	}

	if instance != nil && (instance.Spec.WebApp != nil || instance.Spec.ClusterConfig != nil) {
		envVarsMap := make(map[string]string)

		if roleGroup.Config != nil && roleGroup.EnvVars != nil {
			for key, value := range roleGroup.EnvVars {
				envVarsMap[key] = value
			}
		} else if instance.Spec.WebApp != nil && instance.Spec.WebApp.RoleConfig.EnvVars != nil {
			for key, value := range instance.Spec.WebApp.RoleConfig.EnvVars {
				envVarsMap[key] = value
			}
		}

		if instance.Spec.ClusterConfig != nil && instance.Spec.ClusterConfig.EnvVars != nil {
			for key, value := range instance.Spec.ClusterConfig.EnvVars {
				envVarsMap[key] = value
			}
		}

		for key, value := range envVarsMap {
			envVars = append(envVars, corev1.EnvVar{
				Name:  key,
				Value: value,
			})
		}
	}

	if roleGroup.Config != nil && roleGroup.ExtraEnv != (corev1.EnvVar{}) {
		envVars = append(envVars, roleGroup.ExtraEnv)
	} else if instance != nil && instance.Spec.WebApp != nil && instance.Spec.WebApp.RoleConfig.ExtraEnv != (corev1.EnvVar{}) {
		envVars = append(envVars, instance.Spec.WebApp.RoleConfig.ExtraEnv)
	}
	return envVars
}

// reconcile webapp ingress
func (r *AirbyteReconciler) reconcileWebAppIngress(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	roleGroups := convertRoleGroupToRoleConfigObject(instance, Webapp)
	reconcileParams := r.createReconcileParams(ctx, roleGroups, instance, r.extractWebAppIngressForRoleGroup)
	var params []interface{}
	// Note!!! this is not a common way to do this, Only in this case, the ingress of the webapp is enabled
	if err := reconcileParams.createOrUpdateResourceAndUpdateInstanceStatus(params, r.enableIngresForWebapp); err != nil {
		return err
	}
	return nil
}

func (r *AirbyteReconciler) enableIngresForWebapp(groupName string, rsc any, client client.Client,
	instanceObj client.Object, ctx context.Context, _ []interface{}) error {
	ingress := rsc.(*v1.Ingress)
	host := ingress.Spec.Rules[0].Host
	url := fmt.Sprintf("http://%s", host)
	instance := instanceObj.(*stackv1alpha1.Airbyte)
	if instance.Status.URLs == nil {
		instance.Status.URLs = []status.URL{
			{
				Name: "webui-" + groupName,
				URL:  url,
			},
		}
	} else if host != instance.Status.URLs[0].Name {
		instance.Status.URLs[0].URL = url
	}
	if err := util.UpdateStatus(ctx, client, instance); err != nil {
		return err
	}
	return nil
}

// extract webapp ingress for role group
func (r *AirbyteReconciler) extractWebAppIngressForRoleGroup(params ExtractorParams) (client.Object, error) {
	instance := params.instance.(*stackv1alpha1.Airbyte)
	mergedLabels := r.mergeLabels(params.roleGroup, instance.GetLabels(), params.cluster)
	realRoleCfg := instance.Spec.WebApp.RoleConfig
	realGroupCfg := params.roleGroup.(*stackv1alpha1.WebAppRoleConfigSpec)
	groupName := params.roleGroupName
	pt := v1.PathTypeImplementationSpecific

	var (
		host = getIngressInfo(realGroupCfg.Ingress, realRoleCfg.Ingress, nil)
		port = getServicePort(realGroupCfg, realRoleCfg, params.cluster)
	)
	ing := &v1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      createNameForRoleGroup(instance, "webapp", groupName),
			Namespace: instance.Namespace,
			Labels:    mergedLabels,
		},
		Spec: v1.IngressSpec{
			Rules: []v1.IngressRule{
				{
					Host: host,
					IngressRuleValue: v1.IngressRuleValue{
						HTTP: &v1.HTTPIngressRuleValue{
							Paths: []v1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pt,
									Backend: v1.IngressBackend{
										Service: &v1.IngressServiceBackend{
											Name: createNameForRoleGroup(instance, "webapp", groupName),
											Port: v1.ServiceBackendPort{
												Number: port,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	err := ctrl.SetControllerReference(instance, ing, params.scheme)
	if err != nil {
		r.Log.Error(err, "Failed to set controller reference for ingress")
		return nil, errors.Wrap(err, "Failed to set controller reference for ingress")
	}
	return ing, nil
}

// reconcile webapp secret
func (r *AirbyteReconciler) reconcileWebAppSecret(ctx context.Context, instance *stackv1alpha1.Airbyte) error {
	roleGroups := convertRoleGroupToRoleConfigObject(instance, Webapp)
	reconcileParams := r.createReconcileParams(ctx, roleGroups, instance, r.extractWebappSecretForRoleGroup)
	if err := reconcileParams.createOrUpdateResource(); err != nil {
		return err
	}
	return nil
}

// extract webapp secret for role group
func (r *AirbyteReconciler) extractWebappSecretForRoleGroup(params ExtractorParams) (client.Object, error) {
	instance := params.instance.(*stackv1alpha1.Airbyte)
	mergedLabels := r.mergeLabels(params.roleGroup, instance.GetLabels(), params.cluster)
	groupName := params.roleGroupName
	realRoleCfg := instance.Spec.WebApp.RoleConfig
	realGroupCfg := params.roleGroup.(*stackv1alpha1.WebAppRoleConfigSpec)
	secretInfo := getSecretInfo(realGroupCfg.Secret, realRoleCfg.Secret, nil)

	if secretInfo == nil {
		return nil, nil
	}
	data := createSecretData(*secretInfo)
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "webapp-secrets" + "-" + groupName,
			Namespace: instance.Namespace,
			Labels:    mergedLabels,
		},
		Type: corev1.SecretTypeOpaque,
		Data: data,
	}

	if err := ctrl.SetControllerReference(instance, secret, params.scheme); err != nil {
		return nil, errors.Wrap(err, "Failed to set controller reference for secret")
	}
	return secret, nil
}

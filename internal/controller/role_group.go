package controller

import (
	"fmt"
	stackv1alpha1 "github.com/zncdata-labs/airbyte-operator/api/v1alpha1"
	"github.com/zncdata-labs/airbyte-operator/api/v1alpha1/rolegroup"
	corev1 "k8s.io/api/core/v1"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *AirbyteReconciler) getRoleGroupLabels(config rolegroup.RoleConfigObject) map[string]string {
	if config == nil {
		return make(map[string]string)
	}
	additionalLabels := make(map[string]string)
	if configLabels := config.GetMatchLabels(); configLabels != nil {
		for k, v := range configLabels {
			additionalLabels[k] = v
		}
	}
	return additionalLabels
}

// merge labels
func (r *AirbyteReconciler) mergeLabels(group rolegroup.RoleConfigObject, roleLabels map[string]string,
	cluster rolegroup.RoleConfigObject, roleType AirByteRole) Map {
	var mergeLabels = make(Map)
	mergeLabels.MapMerge(r.getRoleGroupLabels(cluster), true)
	mergeLabels.MapMerge(roleLabels, true)
	mergeLabels.MapMerge(r.getRoleGroupLabels(group), true)
	if roleType != "" {
		mergeLabels["app.kubernetes.io/name"] = mergeLabels["app.kubernetes.io/name"] + "-" + roleNameMapper[roleType]
	}
	return mergeLabels
}

// get ingress info
func getIngressInfo(group, role, cluster *stackv1alpha1.IngressSpec) string {
	if group != nil {
		return group.Host
	}
	if role != nil {
		return role.Host
	}
	if cluster != nil {
		return cluster.Host
	}
	return ""
}

// get service Port
func getServicePort(group, role, cluster rolegroup.RoleConfigObject) int32 {
	var groupSvc, roleSvc, clusterSvc *stackv1alpha1.ServiceSpec
	if group != nil {
		groupSvc = group.GetService()
	}
	if role != nil {
		roleSvc = role.GetService()
	}
	if cluster != nil {
		clusterSvc = cluster.GetService()
	}

	if groupSvc != nil {
		return groupSvc.Port
	}
	if roleSvc != nil {
		return roleSvc.Port
	}
	if clusterSvc != nil {
		return clusterSvc.Port
	}
	return 0
}

// get service info
func getServiceInfo(group, role, cluster rolegroup.RoleConfigObject) (int32, corev1.ServiceType, map[string]string) {
	var groupSvc, roleSvc, clusterSvc *stackv1alpha1.ServiceSpec
	if group != nil {
		groupSvc = group.GetService()
	}
	if role != nil {
		roleSvc = role.GetService()
	}
	if cluster != nil {
		clusterSvc = cluster.GetService()
	}
	if groupSvc != nil {
		return groupSvc.Port, groupSvc.Type, groupSvc.Annotations
	}
	if roleSvc != nil {
		return roleSvc.Port, roleSvc.Type, roleSvc.Annotations
	}
	if clusterSvc != nil {
		return clusterSvc.Port, clusterSvc.Type, clusterSvc.Annotations
	}
	return 0, "", nil
}

// get deployment info
func getDeploymentInfo(group, role, cluster rolegroup.RoleConfigObject) (*stackv1alpha1.ImageSpec,
	*corev1.PodSecurityContext, int32, *corev1.ResourceRequirements, []corev1.ContainerPort) {
	var (
		image           *stackv1alpha1.ImageSpec
		securityContext *corev1.PodSecurityContext
		replicas        *int32
		resources       *stackv1alpha1.ResourcesSpec
		containerPorts  []corev1.ContainerPort
	)
	if cluster != nil {
		image = cluster.GetImage()
		securityContext = cluster.GetSecurityContext()
		clusterReplicas := cluster.GetReplicas()
		replicas = &clusterReplicas
		if configSpec, ok := interface{}(cluster).(rolegroup.ConfigObject); ok {
			resources = configSpec.GetResource()
		}
	}
	// if image of role is not nil, use it, if securityContext of role is not nil, use it; if replicas of role is not nil, use it;...
	overrideFunc := func(roleCfg rolegroup.RoleConfigObject) {
		if roleCfg.GetImage() != nil {
			image = roleCfg.GetImage()
		}
		if roleCfg.GetSecurityContext() != nil {
			securityContext = roleCfg.GetSecurityContext()
		}
		if roleCfg.GetReplicas() != 0 {
			roleReplicas := roleCfg.GetReplicas()
			replicas = &roleReplicas
		}
		if configSpec, ok := interface{}(roleCfg).(rolegroup.ConfigObject); ok {
			if currentResources := configSpec.GetResource(); currentResources != nil {
				resources = currentResources
			}
		}
		if roleCfg.GetService() != nil {
			if rolePort := roleCfg.GetService().Port; rolePort != 0 {
				containerPorts = append(containerPorts, corev1.ContainerPort{
					ContainerPort: rolePort,
				})
			}
		}
	}
	if role != nil {
		overrideFunc(role)
	}
	if group != nil {
		overrideFunc(group)
	}
	var corev1Resource *corev1.ResourceRequirements
	if resources != nil {
		corev1Resource = convertToResourceRequirements(resources)
	}
	return image, securityContext, *replicas, corev1Resource, containerPorts
}

func convertToResourceRequirements(resources *stackv1alpha1.ResourcesSpec) *corev1.ResourceRequirements {
	if resources == nil {
		return nil
	}
	return &corev1.ResourceRequirements{
		Limits:   corev1.ResourceList{"cpu": *resources.CPU.Max, "memory": *resources.Memory.Limit},
		Requests: corev1.ResourceList{"cpu": *resources.CPU.Min, "memory": *resources.Memory.Limit},
	}
}

// get secret info
func getSecretInfo(group, role, cluster *map[string]string) *map[string]string {
	if group != nil {
		return group
	}
	if role != nil {
		return role
	}
	if cluster != nil {
		return cluster
	}
	return nil
}

// get affinity, tolerations, nodeSelector
func getScheduleInfo(group, role, cluster rolegroup.RoleConfigObject) (*corev1.Affinity, *corev1.Toleration,
	map[string]string) {
	var (
		affinity     *corev1.Affinity
		tolerations  *corev1.Toleration
		nodeSelector map[string]string
	)
	if cluster != nil {
		affinity = cluster.GetAffinity()
		tolerations = cluster.GetTolerations()
		nodeSelector = cluster.GetNodeSelector()
	}
	// if affinity of role is not nil, use it; if tolerations of role is not nil, use it; if nodeSelector of role is not nil, use it;...
	overrideFunc := func(roleCfg rolegroup.RoleConfigObject) {
		if roleCfg.GetAffinity() != nil {
			affinity = roleCfg.GetAffinity()
		}
		if roleCfg.GetTolerations() != nil {
			tolerations = roleCfg.GetTolerations()
		}
		if roleCfg.GetNodeSelector() != nil && len(roleCfg.GetNodeSelector()) != 0 {
			nodeSelector = roleCfg.GetNodeSelector()
		}
	}
	if role != nil {
		overrideFunc(role)
	}
	if group != nil {
		overrideFunc(group)
	}
	return affinity, tolerations, nodeSelector
}

// todo(WIP):  :construction:

func GetRoleGroupEx(clusterCfg any, roleCfg any, groupCfg any, roleFields []string,
	configFields []string) client.Object {
	var mergedRoleGroup client.Object

	selectFields := make(map[string]bool)
	for _, field := range roleFields {
		selectFields[field] = true
	}

	roleGroupValue := reflect.ValueOf(groupCfg)
	roleConfigValue := reflect.ValueOf(roleCfg)
	clusterConfigValue := reflect.ValueOf(clusterCfg)
	// can edit use reflect.ValueOf(&mergedRoleGroup)
	mergedRoleGroupValue := reflect.ValueOf(mergedRoleGroup).Elem()

	for i := 0; i < roleGroupValue.NumField(); i++ {
		field := roleGroupValue.Field(i)
		if !selectFields[field.Type().Name()] {
			continue
		}
		if field.IsNil() {
			roleConfigField := roleConfigValue.Field(i)
			if !roleConfigField.IsNil() {
				mergedRoleGroupValue.Field(i).Set(roleConfigField)
			} else {
				clusterConfigField := clusterConfigValue.Field(i)
				if !clusterConfigField.IsNil() {
					mergedRoleGroupValue.Field(i).Set(clusterConfigField)
				}
			}
		} else {
			mergedRoleGroupValue.Field(i).Set(field)
		}
	}
	return mergedRoleGroup
}

var roleNameMapper = map[AirByteRole]string{
	Bootloader:             "bootloader",
	Server:                 "server",
	Cron:                   "cron",
	ApiServer:              "api-server",
	ConnectorBuilderServer: "connector-builder-server",
	PodSweeper:             "pod-sweeper",
	Temporal:               "temporal",
	Worker:                 "worker",
	Webapp:                 "webapp",
}

func createNameForRoleGroup(instance *stackv1alpha1.Airbyte, resourceName string, groupName string) string {
	var name = instance.GetNameWithSuffix(resourceName)
	if groupName != "" {
		name = name + "-" + groupName
	}
	return name
}

func createSvcNameForRoleGroup(instance *stackv1alpha1.Airbyte, roleType AirByteRole, groupName string) string {
	name := instance.GetNameWithSuffix(roleNameMapper[roleType])
	return name + "-" + groupName + "-svc"
}

func getServicePortByRole(roleType AirByteRole, groupName string, instance *stackv1alpha1.Airbyte) (int32, error) {
	clusterCfg := instance.Spec.ClusterConfig.BaseRoleConfigSpec
	var roleCfg rolegroup.RoleConfigObject
	var groupCfg rolegroup.RoleConfigObject
	switch roleType {
	case Server:
		roleCfg = instance.Spec.Server.RoleConfig
		groups := instance.Spec.Server.RoleGroups
		groupCfg = groups[groupName]
	case ApiServer:
		roleCfg = instance.Spec.ApiServer.RoleConfig
		groups := instance.Spec.ApiServer.RoleGroups
		groupCfg = groups[groupName]
	case ConnectorBuilderServer:
		roleCfg = instance.Spec.ConnectorBuilderServer.RoleConfig
		groups := instance.Spec.ConnectorBuilderServer.RoleGroups
		groupCfg = groups[groupName]
	case Webapp:
		roleCfg = instance.Spec.WebApp.RoleConfig
		groups := instance.Spec.WebApp.RoleGroups
		groupCfg = groups[groupName]
	case Cron:
		roleCfg = instance.Spec.Cron.RoleConfig
		groups := instance.Spec.Cron.RoleGroups
		groupCfg = groups[groupName]
	case Temporal:
		roleCfg = instance.Spec.Temporal.RoleConfig
		groups := instance.Spec.Temporal.RoleGroups
		groupCfg = groups[groupName]
	case Worker:
		roleCfg = instance.Spec.Worker.RoleConfig
		groups := instance.Spec.Worker.RoleGroups
		groupCfg = groups[groupName]
	case PodSweeper:
		roleCfg = instance.Spec.PodSweeper.RoleConfig
		groups := instance.Spec.PodSweeper.RoleGroups
		groupCfg = groups[groupName]
	default:
		return 0, fmt.Errorf("roleType %s is not supported", roleType)
	}
	return getServicePort(groupCfg, roleCfg, &clusterCfg), nil

}

func getSvcNameAndPort(instance *stackv1alpha1.Airbyte, roleType AirByteRole, groupName string) (string, int32) {
	handler := func(roleType AirByteRole, groupName string, instance *stackv1alpha1.Airbyte) (string, int32) {
		name := createSvcNameForRoleGroup(instance, roleType, groupName)
		var sp int32
		if port, err := getServicePortByRole(roleType, groupName, instance); err != nil {
			panic(err)
		} else {
			sp = port
		}
		return name, sp
	}

	switch roleType {
	case Server:
		return handler(Server, groupName, instance)
	case ApiServer:
		return handler(ApiServer, groupName, instance)
	case ConnectorBuilderServer:
		return handler(ConnectorBuilderServer, groupName, instance)
	case Webapp:
		return handler(Webapp, groupName, instance)
	case Temporal:
		return handler(Temporal, groupName, instance)
	default:
		return "", 0
	}
}

func getSvcHost(instance *stackv1alpha1.Airbyte, roleType AirByteRole, groupName string) string {
	name, port := getSvcNameAndPort(instance, roleType, groupName)
	return fmt.Sprintf("%s:%d", name, port)
}

func getSvcUrl(instance *stackv1alpha1.Airbyte, roleType AirByteRole, groupName string) string {
	return fmt.Sprintf("http://%s", getSvcHost(instance, roleType, groupName))
}

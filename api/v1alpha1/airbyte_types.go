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

package v1alpha1

import (
	"github.com/zncdata-labs/operator-go/pkg/status"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Image",type="string",JSONPath=".spec.image.repository",description="The Docker Image of Airbyte"
// +kubebuilder:printcolumn:name="Tag",type="string",JSONPath=".spec.image.tag",description="The Docker Tag of Airbyte"
// +kubebuilder:subresource:status

// Airbyte is the Schema for the airbytes API
type Airbyte struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AirbyteSpec   `json:"spec,omitempty"`
	Status status.Status `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AirbyteList contains a list of Airbyte
type AirbyteList struct {
	metav1.TypeMeta `json:",inline,omitempty"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Airbyte `json:"items,omitempty"`
}

// AirbyteSpec defines the desired state of Airbyte
type AirbyteSpec struct {
	// +kubebuilder:validation:Required
	SecurityContext *corev1.PodSecurityContext `json:"securityContext,omitempty"`

	// +kubebuilder:validation:Optional
	ClusterConfig *ClusterRoleConfigSpec `json:"clusterConfig,omitempty"`

	// +kubebuilder:validation:Required
	Server *ServerSpec `json:"server"`

	// +kubebuilder:validation:Required
	Worker *WorkerSpec `json:"worker"`

	// +kubebuilder:validation:Required
	ApiServer *AirbyteApiServerSpec `json:"airbyteApiServer"`

	// +kubebuilder:validation:Required
	WebApp *WebAppSpec `json:"webApp"`

	// +kubebuilder:validation:Required
	PodSweeper *PodSweeperSpec `json:"podSweeper"`

	// +kubebuilder:validation:Required
	ConnectorBuilderServer *ConnectorBuilderServerSpec `json:"connectorBuilderServer"`

	// +kubebuilder:validation:Required
	AirbyteBootloader *AirbyteBootloaderSpec `json:"airbyteBootloader"`

	// +kubebuilder:validation:Required
	Temporal *TemporalSpec `json:"temporal"`

	// +kubebuilder:validation:Required
	Keycloak *KeycloakSpec `json:"keycloak"`

	// +kubebuilder:validation:Required
	Cron *CronSpec `json:"cron"`
}

type IngresRoleConfigSpec struct {
	BaseRoleConfigSpec `json:",inline"`
	// +kubebuilder:validation=Optional
	Ingress *IngressSpec `json:"ingress,omitempty"`
}

type GenericRoleConfigSpec struct {
	BaseRoleConfigSpec `json:",inline"`
	Config             *BaseConfigSpec `json:"config,omitempty"`
}

type DatabaseSpec struct {
	// +kubebuilder:validation:Required
	Reference string `json:"reference"`
}

type S3Spec struct {
	S3BucketSpec `json:"bucket"`

	// +kubebuilder:validation=Optional
	// +kubebuilder:default=20
	MaxConnect int `json:"maxConnect,omitempty"`

	// +kubebuilder:validation=Optional
	PathStyleAccess bool `json:"pathStyle_access,omitempty"`
}

type S3BucketSpec struct {
	Reference string `json:"reference"`
}

type LogsClusterConfigSpec struct {
	// +kubebuilder:validation:Optional
	S3 *S3Spec `json:"s3,omitempty"`
}

type IngressSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=true
	Enabled bool `json:"enabled,omitempty"`

	// +kubebuilder:validation:Optional
	TLS *networkingv1.IngressTLS `json:"tls,omitempty"`

	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default:="webapp.example.com"
	Host string `json:"host,omitempty"`
}

type ServiceSpec struct {
	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:enum=ClusterIP;NodePort;LoadBalancer;ExternalName
	// +kubebuilder:default=ClusterIP
	Type corev1.ServiceType `json:"type,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:default=8001
	Port int32 `json:"port,omitempty"`
}

type DebugSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	Enabled bool `json:"enabled,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=5005
	RemoteDebugPort int32 `json:"RemoteDebugPort,omitempty"`
}

type BaseRoleConfigSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=1
	Replicas int32 `json:"replicas"`

	// +kubebuilder:validation:Optional
	Image *ImageSpec `json:"image"`

	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels,omitempty"`

	// +kubebuilder:validation:Optional
	SecurityContext *corev1.PodSecurityContext `json:"securityContext"`

	// +kubebuilder:validation:Optional
	MatchLabels map[string]string `json:"matchLabels,omitempty"`

	// +kubebuilder:validation:Optional
	Affinity *corev1.Affinity `json:"affinity"`

	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector"`

	// +kubebuilder:validation:Optional
	Tolerations *corev1.Toleration `json:"tolerations"`

	// +kubebuilder:validation:Optional
	Service *ServiceSpec `json:"service"`
}

type BaseConfigSpec struct {
	// +kubebuilder:validation:Optional
	Resources *ResourcesSpec `json:"resources"`
}

type ResourcesSpec struct {
	// +kubebuilder:validation:Optional
	CPU *CPUResource `json:"cpu,omitempty"`

	// +kubebuilder:validation:Optional
	Memory *MemoryResource `json:"memory,omitempty"`

	// +kubebuilder:validation:Optional
	Storage map[string]*StorageResource `json:"storage,omitempty"`
}

type CPUResource struct {
	// +kubebuilder:validation:Optional
	Max *resource.Quantity `json:"max,omitempty"`

	// +kubebuilder:validation:Optional
	Min *resource.Quantity `json:"min,omitempty"`
}

type MemoryResource struct {
	// +kubebuilder:validation:Optional
	Limit *resource.Quantity `json:"limit,omitempty"`
}

type StorageResource struct {
	Capacity *resource.Quantity `json:"capacity,omitempty"`
}

// BaseRoleConfigSpec implements RoleConfigObject interface

func safeAccess(ptr interface{}, fun func() interface{}) interface{} {
	if ptr == nil {
		return nil
	}
	return fun()
}

func (b *BaseRoleConfigSpec) GetReplicas() int32 {
	return safeAccess(b, func() interface{} { return b.Replicas }).(int32)
}
func (b *BaseRoleConfigSpec) GetImage() *ImageSpec {
	return safeAccess(b, func() interface{} { return b.Image }).(*ImageSpec)
}
func (b *BaseRoleConfigSpec) GetLabels() map[string]string {
	return safeAccess(b, func() interface{} { return b.Labels }).(map[string]string)
}
func (b *BaseRoleConfigSpec) GetSecurityContext() *corev1.PodSecurityContext {
	return safeAccess(b, func() interface{} { return b.SecurityContext }).(*corev1.PodSecurityContext)
}
func (b *BaseRoleConfigSpec) GetMatchLabels() map[string]string {
	return safeAccess(b, func() interface{} { return b.MatchLabels }).(map[string]string)
}
func (b *BaseRoleConfigSpec) GetAffinity() *corev1.Affinity {
	return safeAccess(b, func() interface{} { return b.Affinity }).(*corev1.Affinity)
}
func (b *BaseRoleConfigSpec) GetNodeSelector() map[string]string {
	return safeAccess(b, func() interface{} { return b.NodeSelector }).(map[string]string)
}
func (b *BaseRoleConfigSpec) GetTolerations() *corev1.Toleration {
	return safeAccess(b, func() interface{} { return b.Tolerations }).(*corev1.Toleration)
}
func (b *BaseRoleConfigSpec) GetService() *ServiceSpec {
	return safeAccess(b, func() interface{} { return b.Service }).(*ServiceSpec)
}

func (g *GenericRoleConfigSpec) GetConfig() any {
	return safeAccess(g, func() interface{} { return g.Config })
}
func (g *GenericRoleConfigSpec) GetResource() *ResourcesSpec {
	if safeAccess(g, func() interface{} { return g.Config }) == nil {
		return nil
	}
	return g.Config.Resources
}

func (a *Airbyte) GetNameWithSuffix(name string) string {
	return a.GetName() + "-" + name
}

func (a *Airbyte) InitStatusConditions() {
	a.Status.InitStatus(a)
	a.Status.InitStatusConditions()
}

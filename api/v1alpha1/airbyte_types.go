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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	// DefaultImageRepository is the default docker image repository for the operator
	DefaultImageRepository = "airbyte/server"
	// DefaultImageTag is the default docker image tag for the operator
	DefaultImageTag = "0.50.30"
	// DefaultImagePullPolicy is the default docker image pull policy for the operator
	DefaultImagePullPolicy = corev1.PullIfNotPresent
	// DefaultServiceType is the default service type for the operator
	DefaultServiceType = corev1.ServiceTypeClusterIP
	// DefaultSize is the default size for the operator
	DefaultSize = "10Gi"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AirbyteSpec defines the desired state of Airbyte
type AirbyteSpec struct {
	Image ImageSpec `json:"image"`

	// +kubebuilder:validation=Optional
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas"`

	Resource *corev1.ResourceRequirements `json:"resource"`

	// +kubebuilder:validation:Optional
	SecurityContext *corev1.SecurityContext `json:"securityContext"`

	// +kubebuilder:validation:Optional
	PodSecurityContext *corev1.PodSecurityContext `json:"podSecurityContext"`

	// +kubebuilder:validation=Optional
	Service *ServiceSpec `json:"service"`

	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector"`

	// +kubebuilder:validation:Optional
	Tolerations []corev1.Toleration `json:"tolerations"`

	// +kubebuilder:validation=Optional
	Persistence *PersistenceSpec `json:"persistence"`

	// +kubebuilder:validation=Optional
	Secret SecretParam `json:"secret,omitempty"`
}

func (Airbyte *Airbyte) GetLabels() map[string]string {
	return map[string]string{
		"app": Airbyte.Name,
	}
}

type ImageSpec struct {
	// +kubebuilder:validation=Optional
	// +kubebuilder:default=airbyte/server
	Repository string `json:"repository,omitempty"`

	// +kubebuilder:validation=Optional
	// +kubebuilder:default=latest
	Tag string `json:"tag,omitempty"`

	// +kubebuilder:validation=Optional
	// +kubebuilder:default=IfNotPresent
	PullPolicy corev1.PullPolicy `json:"pullPolicy,omitempty"`
}

func (Airbyte *Airbyte) GetImageTag() string {
	image := Airbyte.Spec.Image.Repository
	if image == "" {
		image = DefaultImageRepository
	}
	tag := Airbyte.Spec.Image.Tag
	if tag == "" {
		tag = DefaultImageTag
	}
	return image + ":" + tag
}

func (Airbyte *Airbyte) GetImagePullPolicy() corev1.PullPolicy {
	pullPolicy := Airbyte.Spec.Image.PullPolicy
	if pullPolicy == "" {
		pullPolicy = DefaultImagePullPolicy
	}
	return pullPolicy
}

func (Airbyte *Airbyte) GetNameWithSuffix(name string) string {
	return Airbyte.GetName() + "" + name
}

func (Airbyte *Airbyte) GetPvcName() string {
	if Airbyte.Spec.Persistence.Existing() {
		return *Airbyte.Spec.Persistence.ExistingClaim
	}
	return Airbyte.GetNameWithSuffix("-pvc")
}

func (Airbyte *Airbyte) GetSvcName() string {
	return Airbyte.GetNameWithSuffix("-svc")
}

type ServiceSpec struct {
	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// +kubebuilder:validation:enum=ClusterIP;NodePort;LoadBalancer;ExternalName
	// +kubebuilder:default=ClusterIP
	Type corev1.ServiceType `json:"type"`

	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:default=10000
	Port int32 `json:"port"`
}

func (Airbyte *Airbyte) GetServiceType() corev1.ServiceType {
	serviceType := Airbyte.Spec.Service.Type
	if serviceType == "" {
		serviceType = DefaultServiceType
	}
	return serviceType
}

type PersistenceSpec struct {
	// +kubebuilder:validation:Optional
	StorageClassName *string `json:"storageClassName,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=ReadWriteOnce
	AccessMode string `json:"accessMode,omitempty"`

	// +kubebuilder:default="10Gi"
	StorageSize string `json:"storageSize,omitempty"`

	// +kubebuilder:validation:Optional
	ExistingClaim *string `json:"existingClaim,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=Filesystem
	VolumeMode *corev1.PersistentVolumeMode `json:"volumeMode,omitempty"`

	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels,omitempty"`

	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`
}

func (p *PersistenceSpec) Existing() bool {
	return p.ExistingClaim != nil
}

type SecretParam struct {
	AirbyteSecrets map[string]string `json:"airbyte-secrets,omitempty"`
	GcsLogCreds    map[string]string `json:"gcs-log-creds,omitempty"`
	// airbyte-secrets
	// Suffix    string            `json:"suffix,omitempty"`
	// GcpJson   string            `json:"gcpJson,omitempty"`
	// SecretMap map[string]string `json:"secretMap,omitempty"`
}

// AirbyteStatus defines the observed state of Airbyte
type AirbyteStatus struct {
	Nodes      []string                    `json:"nodes"`
	Conditions []corev1.ComponentCondition `json:"conditons"`
}

//+kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Image",type="string",JSONPath=".spec.image.repository",description="The Docker Image of Airbyte"
// +kubebuilder:printcolumn:name="Tag",type="string",JSONPath=".spec.image.tag",description="The Docker Tag of Airbyte"
//+kubebuilder:subresource:status

// Airbyte is the Schema for the airbytes API
type Airbyte struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AirbyteSpec   `json:"spec,omitempty"`
	Status AirbyteStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AirbyteList contains a list of Airbyte
type AirbyteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Airbyte `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Airbyte{}, &AirbyteList{})
}

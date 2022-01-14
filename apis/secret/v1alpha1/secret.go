/*
Copyright 2021 The Crossplane Authors.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

// A SecretReference is a reference to a secret in an arbitrary namespace.
type SecretReference struct {
	// Name of the secret.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Namespace of the secret.
	// +kubebuilder:validation:Required
	Namespace string `json:"namespace"`

	// Key whose value will be used. If not given, the whole map in the Secret
	// data will be used.
	Key *string `json:"key,omitempty"`
}

// A SecretParameters defines desired state of a Secret
type SecretParameters struct {
	// Name of this secret.
	// *Note:* The secret ID is defined in `metadata.annotations[crossplane.io/external-name]`.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Description of the secret
	// +kubebuilder:validation:Required
	Description string `json:"description"`

	// Reference to the K8s secret that holds the secret value
	// +kubebuilder:validation:Required
	SecretRef SecretReference `json:"secretRef"`

	// ChecksumSecretRef to the K8s secret that stores the checksum for the external secret.
	// This field and the secret will be autogenerated by controller during reconcile.
	// +optional
	ChecksumSecretRef *SecretReference `json:"checksumSecretRef,omitempty"`
}

// A SecretSpec defines the desired state of a Secret.
type SecretSpec struct {
	xpv1.ResourceSpec `json:",inline"`

	// ForProvider contains secret paramerers. Secret ID is defined in `metadata.annotations[crossplane.io/external-name]`.
	ForProvider SecretParameters `json:"forProvider"`
}

// A SecretObservation defines the desired state of a Secret
type SecretObservation struct {
	// LastModifiedAt the time the external resource was last modified.
	LastModifiedAt *metav1.Time `json:"lastUpdated,omitempty"`
}

// A SecretStatus represents the status of a Secret.
type SecretStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          SecretObservation `json:"atProvider"`
}

// +kubebuilder:object:root=true

// A Secret is the schema for Styra Secrets API
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,styra}
type Secret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SecretSpec   `json:"spec"`
	Status SecretStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SecretList contains a list of Secret
type SecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Secret `json:"items"`
}

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
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

// CustomSystemParameters that are not part of the Styra API spec.
type CustomSystemParameters struct {
	// Labels for this systems
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
}

// A SystemParameters defines desired state of a System
type SystemParameters struct {
	CustomSystemParameters `json:",inline"`

	// configuration settings to be used by the system agents
	// +optional
	DeploymentParameters *V1SystemDeploymentParameters `json:"deploymentParameters,omitempty"`

	// description for the system
	// +optional
	Description *string `json:"description,omitempty"`

	// external system ID
	// +optional
	ExternalID *string `json:"externalId,omitempty"`

	// prevents users from modifying policies using Styra UIs
	// +optional
	ReadOnly *bool `json:"readOnly,omitempty"`

	// system type e.g. kubernetes
	// +kubebuilder:validation:Required
	// +immutable
	Type string `json:"type"`
}

// A SystemSpec defines the desired state of a System.
type SystemSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       SystemParameters `json:"forProvider"`
}

// A SystemObservation defines the desired state of a System
type SystemObservation struct {
	// datasources created for the system
	Datasources []*V1DatasourceConfig `json:"datasources,omitempty"`
}

// A SystemStatus represents the status of a System.
type SystemStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          SystemObservation `json:"atProvider"`
}

// +kubebuilder:object:root=true

// A System is the schema for Styra Systems API
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,styra}
type System struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SystemSpec   `json:"spec"`
	Status SystemStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SystemList contains a list of System
type SystemList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []System `json:"items"`
}

// HasLabels whether the system has labels
func (in *SystemParameters) HasLabels() bool {
	// so far only handle labels for kubernetes*
	return in.Type == "kubernetes:v2"
}

// GetAssetTypes gets available asset types
func (in *SystemParameters) GetAssetTypes() []string {
	switch {
	case strings.HasPrefix(in.Type, "kubernetes"):
		return []string{"helm-values"}
	case in.Type == "custom":
		return []string{"opa-config"}
	}

	return []string{}
}

// HasAssets whether the system has available assets
func (in *SystemParameters) HasAssets() bool {
	return len(in.GetAssetTypes()) > 0
}

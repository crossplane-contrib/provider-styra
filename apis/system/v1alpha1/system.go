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

// CustomSystemParameters that are not part of the Styra API spec.
type CustomSystemParameters struct {
	Labels map[string]string `json:"labels,omitempty"`
}

// A SystemParameters defines desired state of a System
type SystemParameters struct {
	CustomSystemParameters `json:",inline"`

	// bundle registry configuration
	BundleRegistry *V1BundleRegistryConfig `json:"bundleRegistry,omitempty"`

	// location of key attributes and additional columns in the decisions grouped by policy entry point path
	DecisionMappings map[string]V1RuleDecisionMappings `json:"decisionMappings,omitempty"`

	// configuration settings to be used by the system agents
	DeploymentParameters *V1SystemDeploymentParameters `json:"deploymentParameters,omitempty"`

	// description for the system
	Description string `json:"description,omitempty"`

	// external system ID
	ExternalID string `json:"externalId,omitempty"`

	// prevents users from modifying policies using Styra UIs
	ReadOnly *bool `json:"readOnly,omitempty"`

	// source control system configuration
	SourceControl *V1SourceControlConfig `json:"sourceControl,omitempty"`

	// system type e.g. kubernetes
	// +kubebuilder:validation:Required
	Type *string `json:"type"`
}

// A SystemSpec defines the desired state of a System.
type SystemSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       SystemParameters `json:"forProvider"`
}

// A SystemObservation defines the desired state of a System
type SystemObservation struct {
	// authorization config
	Authz *V1AuthzConfig `json:"authz,omitempty"`

	// datasources created for the system
	Datasources []*V1DatasourceConfig `json:"datasources,omitempty"`

	// current deployment errors
	Errors map[string]V1AgentErrors `json:"errors,omitempty"`

	// system ID
	ID string `json:"id,omitempty"`

	// // installation instructions by installation method and asset type (deprecated)
	// Install map[string]V1SystemConfigInstall `json:"install,omitempty"`

	// system object metadata
	Metadata *V1ObjectMeta `json:"metadata,omitempty"`

	// policies created for the system
	Policies []*V1PolicyConfig `json:"policies,omitempty"`

	// tokens created for the system
	Tokens []*V1Token `json:"tokens,omitempty"`

	// uninstallation instructions by installation method (deprecated)
	Uninstall map[string]string `json:"uninstall,omitempty"`

	// current deployment warnings
	Warnings map[string]V1SystemConfigWarnings `json:"warnings,omitempty"`
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

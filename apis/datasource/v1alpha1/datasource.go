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

//go:generate go run github.com/crossplane/crossplane-tools/cmd/angryjet generate-methodsets ./...

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

// A DataSourceParameters defines desired state of a DataSource
type DataSourceParameters struct {
	DatasourcesV1Common `json:",inline"`

	AWSECR *DatasourcesV1AWSECR `json:"awsECR,omitempty"`

	BundleS3 *DatasourcesV1BundleS3 `json:"bundleS3,omitempty"`

	GitBlame *DatasourcesV1GitBlame `json:"gitBlame,omitempty"`

	GitContent *DatasourcesV1GitContent `json:"gitContent,omitempty"`

	GitRego *DatasourcesV1GitRego `json:"gitRego,omitempty"`

	HTTP *DatasourcesV1HTTP `json:"http,omitempty"`

	KubernetesResources *DatasourcesV1KubernetesResources `json:"kubernetesResources,omitempty"`

	LDAP *DatasourcesV1LDAP `json:"ldap,omitempty"`

	PolicyLibrary *DatasourcesV1PolicyLibrary `json:"policyLibrary,omitempty"`

	Rest *DatasourcesV1Rest `json:"rest,omitempty"`
}

// A DataSourceSpec defines the desired state of a DataSource.
type DataSourceSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       DataSourceParameters `json:"forProvider"`
}

// A DataSourceObservation defines the observed state of a DataSource.
type DataSourceObservation struct {
	// The last time the data source was executed
	Executed *string `json:"executed,omitempty"`

	// The last observed status of the datasource
	Status *DataSourceExternalStatus `json:"status,omitempty"`
}

// A DataSourceStatus represents the status of a DataSource.
type DataSourceStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          DataSourceObservation `json:"atProvider"`
}

// +kubebuilder:object:root=true

// A DataSource is the schema for Styra DataSources API
// +kubebuilder:printcolumn:name="CATEGORY",type="string",JSONPath=".spec.forProvider.category"
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,styra}
type DataSource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DataSourceSpec   `json:"spec"`
	Status DataSourceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DataSourceList contains a list of DataSource
type DataSourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DataSource `json:"items"`
}

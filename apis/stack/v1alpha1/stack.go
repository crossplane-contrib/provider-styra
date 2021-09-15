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

// CustomStackParameters that are not part of the Styra API
type CustomStackParameters struct {
	// SelectorInclude used to identify systems to apply stack to
	SelectorInclude map[string][]string `json:"selectorInclude,omitempty"`

	// SelectorExclude used to exclude systems from stack application
	SelectorExclude map[string][]string `json:"selectorExclude,omitempty"`
}

// A StackParameters defines desired state of a Stack
type StackParameters struct {
	CustomStackParameters `json:",inline"`

	// description
	// +kubebuilder:validation:Required
	Description string `json:"description"`

	// read only
	// +kubebuilder:validation:Required
	ReadOnly bool `json:"readOnly"`

	// source control
	// +optional
	SourceControl *V1SourceControlConfig `json:"sourceControl,omitempty"`

	// type
	// +kubebuilder:validation:Required
	Type string `json:"type"`
}

// A StackSpec defines the desired state of a Stack.
type StackSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       StackParameters `json:"forProvider"`
}

// A StackStatus represents the status of a Stack.
type StackStatus struct {
	xpv1.ResourceStatus `json:",inline"`
}

// +kubebuilder:object:root=true

// A Stack is the schema for Styra Stacks API
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,styra}
type Stack struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StackSpec   `json:"spec"`
	Status StackStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// StackList contains a list of Stack
type StackList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Stack `json:"items"`
}

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
)

// V1DatasourceConfig v1 datasource config
type V1DatasourceConfig struct {

	// datasource category
	// +kubebuilder:validation:Required
	Category string `json:"category"`

	// datasource ID
	// +kubebuilder:validation:Required
	ID string `json:"id"`

	// optional datasources can be deleted without being recreated automatically
	Optional bool `json:"optional,omitempty"`

	// datasource status
	Status *V1Status `json:"status,omitempty"`

	// pull or push
	// +kubebuilder:validation:Required
	Type string `json:"type"`
}

// V1Status v1 status
type V1Status struct {

	// code
	// +kubebuilder:validation:Required
	Code *string `json:"code"`

	// message
	// +kubebuilder:validation:Required
	Message *string `json:"message"`

	// timestamp
	// +kubebuilder:validation:Required
	Timestamp *metav1.Time `json:"timestamp"`
}

// V1SystemDeploymentParameters v1 system deployment parameters
type V1SystemDeploymentParameters struct {

	// true to fail close
	DenyOnOpaFail *bool `json:"denyOnOpaFail,omitempty"`

	// // extra deployment settings
	// Extra interface{} `json:"extra,omitempty"`

	// HTTP proxy URL
	HTTPProxy *string `json:"httpProxy,omitempty"`

	// HTTPS proxy URL
	HTTPSProxy *string `json:"httpsProxy,omitempty"`

	// minimum Kubernetes version expected (where applicable)
	KubernetesVersion *string `json:"kubernetesVersion,omitempty"`

	// Kubernetes namespace the system is deployed to
	Namespace *string `json:"namespace,omitempty"`

	// URLs that should be excluded from proxying
	NoProxy *string `json:"noProxy,omitempty"`

	// Kubernetes webhook timeout (where applicable)
	TimeoutSeconds *int32 `json:"timeoutSeconds,omitempty"`

	// trusted CA certificates
	TrustedCaCerts []string `json:"trustedCaCerts,omitempty"`

	// trusted container registry
	TrustedContainerRegistry *string `json:"trustedContainerRegistry,omitempty"`
}

// V1AgentErrors v1 agent errors
type V1AgentErrors struct {

	// list of system errors
	// +kubebuilder:validation:Required
	Errors []*V1Status `json:"errors"`

	// true if the the system is waiting for error to be resolved
	// +kubebuilder:validation:Required
	Waiting *bool `json:"waiting"`
}

// V1SystemConfigWarnings v1 system config warnings
type V1SystemConfigWarnings struct {

	// code
	// +kubebuilder:validation:Required
	Code *string `json:"code"`

	// message
	// +kubebuilder:validation:Required
	Message *string `json:"message"`

	// timestamp
	// +kubebuilder:validation:Required
	Timestamp *metav1.Time `json:"timestamp"`
}

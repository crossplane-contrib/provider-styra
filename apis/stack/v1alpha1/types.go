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

// V1ObjectMeta v1 object meta
//
type V1ObjectMeta struct {

	// created at
	CreatedAt metav1.Time `json:"created_at,omitempty"`

	// created by
	CreatedBy string `json:"created_by,omitempty"`

	// created through
	CreatedThrough string `json:"created_through,omitempty"`

	// last modified at
	LastModifiedAt metav1.Time `json:"last_modified_at,omitempty"`

	// last modified by
	LastModifiedBy string `json:"last_modified_by,omitempty"`

	// last modified through
	LastModifiedThrough string `json:"last_modified_through,omitempty"`
}

// V1PolicyConfig v1 policy config
//
type V1PolicyConfig struct {

	// policy on when to (re)generate the policy
	Created string `json:"created,omitempty"`

	// enforcement status of the policy
	// +kubebuilder:validation:Required
	Enforcement *V1EnforcementConfig `json:"enforcement"`

	// policy ID (path)
	// +kubebuilder:validation:Required
	ID *string `json:"id"`

	// rego modules policy consists of
	Modules []*V1Module `json:"modules"`

	// rule count
	Rules *V1RuleCounts `json:"rules,omitempty"`

	// policy type e.g. validating/rules
	// +kubebuilder:validation:Required
	Type *string `json:"type"`
}

// V1EnforcementConfig v1 enforcement config
//
type V1EnforcementConfig struct {

	// true if the policy is enforced
	// +kubebuilder:validation:Required
	Enforced bool `json:"enforced"`

	// enforcement type e.g. opa, test, mask
	// +kubebuilder:validation:Required
	Type string `json:"type"`
}

// V1Module v1 module
//
type V1Module struct {

	// module name
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// module is a placeholder
	Placeholder *bool `json:"placeholder,omitempty"`

	// true if module is read-only
	// +kubebuilder:validation:Required
	ReadOnly bool `json:"read_only"`

	// module rule count
	Rules *V1RuleCounts `json:"rules,omitempty"`
}

// V1RuleCounts v1 rule counts
//
type V1RuleCounts struct {

	// number of allow rules
	// +kubebuilder:validation:Required
	Allow int32 `json:"allow"`

	// number of deny rules
	// +kubebuilder:validation:Required
	Deny int32 `json:"deny"`

	// number of enforce rules
	// +kubebuilder:validation:Required
	Enforce int32 `json:"enforce"`

	// number of ignore rules
	// +kubebuilder:validation:Required
	Ignore int32 `json:"ignore"`

	// number of monitor rules
	// +kubebuilder:validation:Required
	Monitor int32 `json:"monitor"`

	// number of notify rules
	// +kubebuilder:validation:Required
	Notify int32 `json:"notify"`

	// number of unclassified rules
	// +kubebuilder:validation:Required
	Other int32 `json:"other"`

	// number of test rules
	// +kubebuilder:validation:Required
	Test int32 `json:"test"`

	// total number of rules
	// +kubebuilder:validation:Required
	Total int32 `json:"total"`
}

// V1SourceControlConfig v1 source control config
//
type V1SourceControlConfig struct {

	// origin
	// +kubebuilder:validation:Required
	Origin *V1GitRepoConfig `json:"origin"`
}

// V1GitRepoConfig v1 git repo config
//
type V1GitRepoConfig struct {

	// Credentials are looked under the key <name>/<creds>
	// +kubebuilder:validation:Required
	Credentials *string `json:"credentials"`

	// Path to limit the import to
	// +kubebuilder:validation:Required
	Path *string `json:"path"`

	// Remote reference, defaults to refs/heads/master
	// +kubebuilder:validation:Required
	Reference *string `json:"reference"`

	// Repository URL
	// +kubebuilder:validation:Required
	URL *string `json:"url"`
}

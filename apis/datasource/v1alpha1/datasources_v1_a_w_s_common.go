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
	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

// DatasourcesV1AWSCommon fields
type DatasourcesV1AWSCommon struct {

	// Secret ID with AWS credentials
	// +kubebuilder:validation:Required
	// +crossplane:generate:reference:type=github.com/crossplane-contrib/provider-styra/apis/secret/v1alpha1.Secret
	// +crossplane:generate:reference:refFieldName=CredentialsRef
	// +crossplane:generate:reference:selectorFieldName=CredentialsSelector
	Credentials string `json:"credentials"`

	// CredentialsRef is a reference to a Secret used to set Credentials.
	// +optional
	CredentialsRef *xpv1.Reference `json:"credentialsRef,omitempty"`

	// CredentialsSelector selects references to a Secret used to set Credentials.
	// +optional
	CredentialsSelector *xpv1.Selector `json:"credentialsSelector,omitempty"`

	// AWS region
	// +kubebuilder:validation:Required
	Region string `json:"region"`
}

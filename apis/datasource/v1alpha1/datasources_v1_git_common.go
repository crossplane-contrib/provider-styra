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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DatasourcesV1GitCommon fields
type DatasourcesV1GitCommon struct {
	DatasourcesV1Poller `json:",inline"`

	DatasourcesV1RateLimiter `json:",inline"`

	// Secret ID with credentials
	// +kubebuilder:validation:Required
	// +crossplane:generate:reference:type=github.com/crossplane-contrib/provider-styra/apis/secret/v1alpha1.Secret
	// +crossplane:generate:reference:refFieldName=CredentialsRef
	// +crossplane:generate:reference:selectorFieldName=CredentialsSelector
	Credentials *string `json:"credentials,omitempty"`

	// CredentialsRef is a reference to a Secret used to set Credentials.
	// +optional
	CredentialsRef *xpv1.Reference `json:"credentialsRef,omitempty"`

	// CredentialsSelector selects references to a Secret used to set Credentials.
	// +optional
	CredentialsSelector *xpv1.Selector `json:"credentialsSelector,omitempty"`

	// reference
	Reference *string `json:"reference,omitempty"`

	// ssh credentials
	SSHCredentials *DatasourcesV1GitCommonAO3SSHCredentials `json:"sshCredentials,omitempty"`

	// timeout
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// Git URL
	// +kubebuilder:validation:Required
	URL string `json:"url"`
}

// DatasourcesV1GitCommonAO3SSHCredentials fields
type DatasourcesV1GitCommonAO3SSHCredentials struct {

	// Secret ID with passphrase
	// +kubebuilder:validation:Required
	// +crossplane:generate:reference:type=github.com/crossplane-contrib/provider-styra/apis/secret/v1alpha1.Secret
	// +crossplane:generate:reference:refFieldName=PassphraseRef
	// +crossplane:generate:reference:selectorFieldName=PassphraseSelector
	Passphrase *string `json:"passphrase,omitempty"`

	// PassphraseRef is a reference to a Secret used to set Passphrase.
	// +optional
	PassphraseRef *xpv1.Reference `json:"passphraseRef,omitempty"`

	// PassphraseSelector selects references to a Secret used to set Passphrase.
	// +optional
	PassphraseSelector *xpv1.Selector `json:"passphraseSelector,omitempty"`

	// Secret ID with private key
	// +kubebuilder:validation:Required
	// +crossplane:generate:reference:type=github.com/crossplane-contrib/provider-styra/apis/secret/v1alpha1.Secret
	// +crossplane:generate:reference:refFieldName=PrivateKeyRef
	// +crossplane:generate:reference:selectorFieldName=PrivateKeySelector
	// +kubebuilder:validation:Required
	PrivateKey string `json:"privateKey"`

	// PassphraseRef is a reference to a Secret used to set Passphrase.
	// +optional
	PrivateKeyRef *xpv1.Reference `json:"privateKeyRef,omitempty"`

	// PassphraseSelector selects references to a Secret used to set Passphrase.
	// +optional
	PrivateKeySelector *xpv1.Selector `json:"privateKeySelector,omitempty"`
}

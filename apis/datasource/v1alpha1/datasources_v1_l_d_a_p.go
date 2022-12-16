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

// DatasourcesV1LDAP properties specific to LDAP datasources
type DatasourcesV1LDAP struct {
	DatasourcesV1RateLimiter `json:",inline"`

	DatasourcesV1Poller `json:",inline"`

	DatasourcesV1RegoFiltering `json:",inline"`

	DatasourcesV1TLSSettings `json:",inline"`

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

	// search
	Search *DatasourcesV1LDAPAO5Search `json:"search,omitempty"`

	// List of URLs: main + replicas
	// +kubebuilder:validation:Required
	Urls []string `json:"urls"`
}

// DatasourcesV1LDAPAO5Search Search Request.
// Documentation: https://ldapwiki.com/wiki/SearchRequest
//
// swagger:model DatasourcesV1LDAPAO5Search
type DatasourcesV1LDAPAO5Search struct {

	// Search attribute selection.
	// Documentation: https://ldapwiki.com/wiki/AttributeSelection
	//
	Attributes []string `json:"attributes,omitempty"`

	// Search Base DN.
	// Documentation: https://ldapwiki.com/wiki/BaseDN
	//
	// +kubebuilder:validation:Required
	BaseDN string `json:"baseDN"`

	// Search dereference policy.
	// Documentation: https://ldapwiki.com/wiki/Dereference%20Policy
	//
	// Enum: [never searching finding always]
	Deref *string `json:"deref,omitempty"`

	// Search filter.
	// Documentation: https://ldapwiki.com/wiki/LDAP%20SearchFilters
	// Examples: https://ldapwiki.com/wiki/LDAP%20Query%20Examples
	//
	// +kubebuilder:validation:Required
	Filter string `json:"filter"`

	// Search page size.
	// Documentation: https://ldapwiki.com/wiki/MaxPageSize
	//
	// Maximum: 4.294967295e+09
	// Minimum: 0
	PageSize *int64 `json:"pageSize,omitempty"`

	// Search scope.
	// Documentation: https://ldapwiki.com/wiki/LDAP%20Search%20Scopes
	//
	// Enum: [base-object single-level whole-subtree]
	Scope *string `json:"scope,omitempty"`

	// Search page limit.
	// Documentation: https://ldapwiki.com/wiki/SizeLimit
	//
	// Minimum: 0
	SizeLimit *int64 `json:"sizeLimit,omitempty"`
}

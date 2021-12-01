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

// DatasourcesV1RegoFiltering fields
type DatasourcesV1RegoFiltering struct {

	// Policy Filter (if set, then policyQuery must be set as well)
	PolicyFilter *string `json:"policyFilter,omitempty"`

	// Policy Query (if set, then policyFilter must be set as well)
	PolicyQuery *string `json:"policyQuery,omitempty"`
}

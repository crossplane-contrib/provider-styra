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

// Supported datasource categories
const (
	DataSourceCategoryAWSECR              = "aws/ecr"
	DataSourceCategoryBundleS3            = "bundle/s3"
	DataSourceCategoryGitBlame            = "git/blame"
	DataSourceCategoryGitContent          = "git/content"
	DataSourceCategoryGitRego             = "git/rego"
	DataSourceCategoryHTTP                = "http"
	DataSourceCategoryKubernetesResources = "kubernetes/resources"
	DataSourceCategoryLDAP                = "ldap"
	DataSourceCategoryPolicyLibrary       = "policy-library"
	DataSourceCategoryRest                = "rest"

	// DataSourceStatusFailed describes the status code when the external resource has failed
	DataSourceStatusFailed = "failed"
)

// DataSourceExternalStatus represents the external status of a datasource tracked by Styra.
type DataSourceExternalStatus struct {
	Code *string `json:"code,omitempty"`

	Message *string `json:"message,omitempty"`

	Timestamp *string `json:"timestamp,omitempty"`
}

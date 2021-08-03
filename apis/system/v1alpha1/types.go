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

// V1AuthzConfig v1 authz config
//
type V1AuthzConfig struct {

	// a list of role binding configs
	RoleBindings []*V1RoleBindingConfig `json:"roleBindings"`
}

// V1RoleBindingConfig v1 role binding config
//
type V1RoleBindingConfig struct {

	// role binding ID
	// +kubebuilder:validation:Required
	ID *string `json:"id"`

	// role name
	// +kubebuilder:validation:Required
	RoleName *string `json:"roleName"`
}

// V1BundleRegistryConfig v1 bundle registry config
//
type V1BundleRegistryConfig struct {

	// configuration for external S3 bucket to use for bundle distribution
	DistributionS3 *V1BundleDistributionS3Config `json:"distributionS3,omitempty"`

	// manual deployment mode to prevent automatic deployment of new bundles
	ManualDeployment bool `json:"manualDeployment,omitempty"`

	// maximum number of all bundles to store (default 100)
	MaxBundles int64 `json:"maxBundles,omitempty"`

	// maximum number of previously deployed bundles to store (default 10)
	MaxDeployedBundles int64 `json:"maxDeployedBundles,omitempty"`
}

// V1BundleDistributionS3Config v1 bundle distribution s3 config
//
type V1BundleDistributionS3Config struct {

	// access key id and secret access key are looked under the key <name>/<access_keys>
	AccessKeys string `json:"accessKeys,omitempty"`

	// bucket name
	// +kubebuilder:validation:Required
	Bucket *string `json:"bucket"`

	// discovery bundle path
	// +kubebuilder:validation:Required
	DiscoveryPath *string `json:"discoveryPath"`

	// if provided, OPA uses this 'services[_].credentials.s3_signing' config to connect directly to S3 for bundle downloads
	// +kubebuilder:validation:Required
	OpaCredentials *V1BundleDistributionS3ConfigOpaCredentials `json:"opaCredentials"`

	// bundle path
	// +kubebuilder:validation:Required
	PolicyPath *string `json:"policyPath"`

	// AWS region
	// +kubebuilder:validation:Required
	Region *string `json:"region"`
}

// V1BundleDistributionS3ConfigOpaCredentials v1 bundle distribution s3 config opa credentials
type V1BundleDistributionS3ConfigOpaCredentials struct {

	// // environment credentials
	// EnvironmentCredentials V1BundleDistributionS3ConfigOpaCredentialsEnvironmentCredentials `json:"environment_credentials,omitempty"`

	// metadata credentials
	MetadataCredentials *V1BundleDistributionS3ConfigOpaCredentialsMetadataCredentials `json:"metadataCredentials,omitempty"`

	// web identity credentials
	WebIdentityCredentials *V1BundleDistributionS3ConfigOpaCredentialsWebIdentityCredentials `json:"webIdentityCredentials,omitempty"`
}

// // V1BundleDistributionS3ConfigOpaCredentialsEnvironmentCredentials v1 bundle distribution s3 config opa credentials environment credentials
// //
// type V1BundleDistributionS3ConfigOpaCredentialsEnvironmentCredentials interface{}

// V1BundleDistributionS3ConfigOpaCredentialsMetadataCredentials v1 bundle distribution s3 config opa credentials metadata credentials
//
type V1BundleDistributionS3ConfigOpaCredentialsMetadataCredentials struct {

	// aws region
	// +kubebuilder:validation:Required
	AwsRegion *string `json:"awsRegion"`

	// iam role
	// +kubebuilder:validation:Required
	IamRole *string `json:"iamRole"`
}

// V1BundleDistributionS3ConfigOpaCredentialsWebIdentityCredentials v1 bundle distribution s3 config opa credentials web identity credentials
//
type V1BundleDistributionS3ConfigOpaCredentialsWebIdentityCredentials struct {

	// aws region
	// +kubebuilder:validation:Required
	AwsRegion *string `json:"awsRegion"`

	// session name
	// +kubebuilder:validation:Required
	SessionName *string `json:"sessionName"`
}

// V1DatasourceConfig v1 datasource config
//
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
//
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

// V1RuleDecisionMappings v1 rule decision mappings
//
type V1RuleDecisionMappings struct {

	// rules to determine decision type (allowed, denied)
	Allowed *V1AllowedMapping `json:"allowed,omitempty"`

	// decision mappings for additional columns
	Columns []*V1ColumnMapping `json:"columns"`

	// decision mapping for the reason field
	Reason *V1ReasonMapping `json:"reason,omitempty"`
}

// V1AllowedMapping v1 allowed mapping
//
type V1AllowedMapping struct {

	// // expected value of the decision property
	// Expected V1AllowedMappingExpected `json:"expected,omitempty"`

	// when set to true, decision is Allowed when the mapped property IS NOT equal to the expected value
	Negated *bool `json:"negated,omitempty"`

	// dot-separated decision property path
	// +kubebuilder:validation:Required
	Path *string `json:"path"`
}

// // V1AllowedMappingExpected v1 allowed mapping expected
// //
// type V1AllowedMappingExpected interface{}

// V1ColumnMapping v1 column mapping
//
type V1ColumnMapping struct {

	// column key (also the search key)
	// +kubebuilder:validation:Required
	Key *string `json:"key"`

	// dot-separated decision property path
	// +kubebuilder:validation:Required
	Path *string `json:"path"`

	// column type: one of "string", "boolean", "date", "integer", "float"
	Type *string `json:"type,omitempty"`
}

// V1ReasonMapping v1 reason mapping
//
type V1ReasonMapping struct {

	// dot-separated decision property path
	// +kubebuilder:validation:Required
	Path *string `json:"path"`
}

// V1SystemDeploymentParameters v1 system deployment parameters
//
type V1SystemDeploymentParameters struct {

	// true to fail close
	DenyOnOpaFail *bool `json:"denyOnOpaFail,omitempty"`

	// // extra deployment settings
	// Extra interface{} `json:"extra,omitempty"`

	// HTTP proxy URL
	HTTPProxy string `json:"httpProxy,omitempty"`

	// HTTPS proxy URL
	HTTPSProxy string `json:"httpsProxy,omitempty"`

	// minimum Kubernetes version expected (where applicable)
	KubernetesVersion string `json:"kubernetesVersion,omitempty"`

	// Kubernetes namespace the system is deployed to
	Namespace string `json:"namespace,omitempty"`

	// URLs that should be excluded from proxying
	NoProxy string `json:"noProxy,omitempty"`

	// Kubernetes webhook timeout (where applicable)
	TimeoutSeconds int32 `json:"timeoutSeconds,omitempty"`

	// trusted CA certificates
	TrustedCaCerts []string `json:"trustedCaCerts,omitempty"`

	// trusted container registry
	TrustedContainerRegistry string `json:"trustedContainerRegistry,omitempty"`
}

// V1AgentErrors v1 agent errors
//
type V1AgentErrors struct {

	// list of system errors
	// +kubebuilder:validation:Required
	Errors []*V1Status `json:"errors"`

	// true if the the system is waiting for error to be resolved
	// +kubebuilder:validation:Required
	Waiting *bool `json:"waiting"`
}

// // V1SystemConfigInstall v1 system config install
// //
// type V1SystemConfigInstall interface{}

// V1ObjectMeta v1 object meta
//
type V1ObjectMeta struct {

	// created at
	CreatedAt metav1.Time `json:"createdAt,omitempty"`

	// created by
	CreatedBy string `json:"createdBy,omitempty"`

	// created through
	CreatedThrough string `json:"createdThrough,omitempty"`

	// last modified at
	LastModifiedAt metav1.Time `json:"lastModifiedAt,omitempty"`

	// last modified by
	LastModifiedBy string `json:"lastModifiedBy,omitempty"`

	// last modified through
	LastModifiedThrough string `json:"lastModifiedThrough,omitempty"`
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
	ReadOnly bool `json:"readOnly"`

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

// V1Token v1 token
//
type V1Token struct {

	// allow path patterns
	// +kubebuilder:validation:Required
	AllowPathPatterns []string `json:"allowPathPatterns"`

	// description
	// +kubebuilder:validation:Required
	Description *string `json:"description"`

	// expires
	Expires *metav1.Time `json:"expires,omitempty"`

	// id
	// +kubebuilder:validation:Required
	ID *string `json:"id"`

	// metadata
	Metadata *V1ObjectMeta `json:"metadata,omitempty"`

	// token
	Token string `json:"token,omitempty"`

	// ttl
	// +kubebuilder:validation:Required
	TTL *string `json:"ttl"`

	// uses
	// +kubebuilder:validation:Required
	Uses *int64 `json:"uses"`
}

// V1SystemConfigWarnings v1 system config warnings
//
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

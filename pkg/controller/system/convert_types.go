package system

import (
	"time"

	"github.com/go-openapi/strfmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mistermx/styra-go-client/pkg/client/systems"
	"github.com/mistermx/styra-go-client/pkg/models"

	"github.com/crossplane-contrib/provider-styra/apis/system/v1alpha1"
	styraclient "github.com/crossplane-contrib/provider-styra/pkg/client"
)

// GenerateSystem generates a V1SystemConfig from System
func GenerateSystem(resp *models.V1SystemConfig) (cr *v1alpha1.System) { // nolint:gocyclo
	cr = &v1alpha1.System{}

	if resp.Authz != nil {
		cr.Status.AtProvider.Authz = &v1alpha1.V1AuthzConfig{}
		if resp.Authz.RoleBindings != nil {
			cr.Status.AtProvider.Authz.RoleBindings = make([]*v1alpha1.V1RoleBindingConfig, len(resp.Authz.RoleBindings))
			for i, v := range resp.Authz.RoleBindings {
				n := &v1alpha1.V1RoleBindingConfig{}
				if v.ID != nil {
					n.ID = v.ID
				}
				if v.RoleName != nil {
					n.RoleName = v.RoleName
				}

				cr.Status.AtProvider.Authz.RoleBindings[i] = n
			}
		}
	}
	if resp.BundleRegistry != nil {
		cr.Spec.ForProvider.BundleRegistry = &v1alpha1.V1BundleRegistryConfig{}
		cr.Spec.ForProvider.BundleRegistry.ManualDeployment = resp.BundleRegistry.ManualDeployment
		cr.Spec.ForProvider.BundleRegistry.MaxBundles = resp.BundleRegistry.MaxBundles
		cr.Spec.ForProvider.BundleRegistry.MaxDeployedBundles = resp.BundleRegistry.MaxDeployedBundles
		if resp.BundleRegistry.DistributionS3 != nil {
			cr.Spec.ForProvider.BundleRegistry.DistributionS3 = &v1alpha1.V1BundleDistributionS3Config{}
			cr.Spec.ForProvider.BundleRegistry.DistributionS3.AccessKeys = resp.BundleRegistry.DistributionS3.AccessKeys
			if resp.BundleRegistry.DistributionS3.Bucket != nil {
				cr.Spec.ForProvider.BundleRegistry.DistributionS3.Bucket = resp.BundleRegistry.DistributionS3.Bucket
			}
			if resp.BundleRegistry.DistributionS3.DiscoveryPath != nil {
				cr.Spec.ForProvider.BundleRegistry.DistributionS3.DiscoveryPath = resp.BundleRegistry.DistributionS3.DiscoveryPath
			}
			if resp.BundleRegistry.DistributionS3.PolicyPath != nil {
				cr.Spec.ForProvider.BundleRegistry.DistributionS3.PolicyPath = resp.BundleRegistry.DistributionS3.PolicyPath
			}
			if resp.BundleRegistry.DistributionS3.Region != nil {
				cr.Spec.ForProvider.BundleRegistry.DistributionS3.Region = resp.BundleRegistry.DistributionS3.Region
			}
			if resp.BundleRegistry.DistributionS3.OpaCredentials != nil {
				cr.Spec.ForProvider.BundleRegistry.DistributionS3.OpaCredentials = &v1alpha1.V1BundleDistributionS3ConfigOpaCredentials{}
				if resp.BundleRegistry.DistributionS3.OpaCredentials.MetadataCredentials != nil {
					cr.Spec.ForProvider.BundleRegistry.DistributionS3.OpaCredentials.MetadataCredentials = &v1alpha1.V1BundleDistributionS3ConfigOpaCredentialsMetadataCredentials{}
					if resp.BundleRegistry.DistributionS3.OpaCredentials.MetadataCredentials.AwsRegion != nil {
						cr.Spec.ForProvider.BundleRegistry.DistributionS3.OpaCredentials.MetadataCredentials.AwsRegion = resp.BundleRegistry.DistributionS3.OpaCredentials.MetadataCredentials.AwsRegion
					}
					if resp.BundleRegistry.DistributionS3.OpaCredentials.MetadataCredentials.IamRole != nil {
						cr.Spec.ForProvider.BundleRegistry.DistributionS3.OpaCredentials.MetadataCredentials.IamRole = resp.BundleRegistry.DistributionS3.OpaCredentials.MetadataCredentials.IamRole
					}
				}
				if resp.BundleRegistry.DistributionS3.OpaCredentials.WebIdentityCredentials != nil {
					cr.Spec.ForProvider.BundleRegistry.DistributionS3.OpaCredentials.WebIdentityCredentials = &v1alpha1.V1BundleDistributionS3ConfigOpaCredentialsWebIdentityCredentials{}
					if resp.BundleRegistry.DistributionS3.OpaCredentials.WebIdentityCredentials.AwsRegion != nil {
						cr.Spec.ForProvider.BundleRegistry.DistributionS3.OpaCredentials.WebIdentityCredentials.AwsRegion = resp.BundleRegistry.DistributionS3.OpaCredentials.WebIdentityCredentials.AwsRegion
					}
					if resp.BundleRegistry.DistributionS3.OpaCredentials.WebIdentityCredentials.SessionName != nil {
						cr.Spec.ForProvider.BundleRegistry.DistributionS3.OpaCredentials.WebIdentityCredentials.SessionName = resp.BundleRegistry.DistributionS3.OpaCredentials.WebIdentityCredentials.SessionName
					}
				}
			}
		}
	}
	if resp.Datasources != nil {
		cr.Status.AtProvider.Datasources = make([]*v1alpha1.V1DatasourceConfig, len(resp.Datasources))
		for i, v := range resp.Datasources {
			n := &v1alpha1.V1DatasourceConfig{}
			n.ID = v.ID
			n.Optional = v.Optional
			n.Category = v.Category
			n.Type = v.Type
			if v.Status != nil {
				n.Status = &v1alpha1.V1Status{}
				if v.Status.Code != nil {
					n.Status.Code = v.Status.Code
				}
				if v.Status.Message != nil {
					n.Status.Message = v.Status.Message
				}
				// if v.Status.Timestamp != nil {
				// 	n.Status.Timestamp = v.Status.Timestamp
				// }
			}
			cr.Status.AtProvider.Datasources[i] = n
		}
	}
	// if resp.DecisionMappings != nil {
	// 	cr.Spec.ForProvider.DecisionMappings = make(map[string]v1alpha1.V1RuleDecisionMappings, len(resp.DecisionMappings))
	// 	for k, v := range resp.DecisionMappings {
	// 		n := &v1alpha1.V1RuleDecisionMappings{}
	// 	}
	// }

	if resp.DeploymentParameters != nil {
		cr.Spec.ForProvider.DeploymentParameters = &v1alpha1.V1SystemDeploymentParameters{
			DenyOnOpaFail:            resp.DeploymentParameters.DenyOnOpaFail,
			HTTPProxy:                resp.DeploymentParameters.HTTPProxy,
			HTTPSProxy:               resp.DeploymentParameters.HTTPSProxy,
			KubernetesVersion:        resp.DeploymentParameters.KubernetesVersion,
			Namespace:                resp.DeploymentParameters.Namespace,
			NoProxy:                  resp.DeploymentParameters.NoProxy,
			TimeoutSeconds:           resp.DeploymentParameters.TimeoutSeconds,
			TrustedCaCerts:           resp.DeploymentParameters.TrustedCaCerts,
			TrustedContainerRegistry: resp.DeploymentParameters.TrustedContainerRegistry,
		}
	}

	cr.Spec.ForProvider.Description = resp.Description

	// Errors

	cr.Spec.ForProvider.ExternalID = resp.ExternalID

	// Install

	if resp.Metadata != nil {
		cr.Status.AtProvider.Metadata = &v1alpha1.V1ObjectMeta{}
		cr.Status.AtProvider.Metadata.CreatedAt = GenerateTime(resp.Metadata.CreatedAt)
		cr.Status.AtProvider.Metadata.CreatedBy = resp.Metadata.CreatedBy
		cr.Status.AtProvider.Metadata.CreatedThrough = resp.Metadata.CreatedThrough
		cr.Status.AtProvider.Metadata.LastModifiedAt = GenerateTime(resp.Metadata.LastModifiedAt)
		cr.Status.AtProvider.Metadata.LastModifiedBy = resp.Metadata.LastModifiedBy
		cr.Status.AtProvider.Metadata.LastModifiedThrough = resp.Metadata.LastModifiedThrough
	}

	if resp.Policies != nil {
		cr.Status.AtProvider.Policies = make([]*v1alpha1.V1PolicyConfig, len(resp.Policies))
		for i, v := range resp.Policies {
			n := &v1alpha1.V1PolicyConfig{}
			n.Created = v.Created
			if v.Enforcement != nil {
				n.Enforcement = &v1alpha1.V1EnforcementConfig{}
				n.Enforcement.Enforced = v.Enforcement.Enforced
				n.Enforcement.Type = v.Enforcement.Type
			}
			if v.ID != nil {
				n.ID = v.ID
			}

			if v.Modules != nil {
				n.Modules = make([]*v1alpha1.V1Module, len(v.Modules))
				for im, vm := range v.Modules {
					nm := &v1alpha1.V1Module{}
					nm.Name = vm.Name
					nm.Placeholder = vm.Placeholder
					nm.ReadOnly = vm.ReadOnly
					if vm.Rules != nil {
						nm.Rules = &v1alpha1.V1RuleCounts{}
						nm.Rules.Allow = vm.Rules.Allow
						nm.Rules.Deny = vm.Rules.Deny
						nm.Rules.Enforce = vm.Rules.Enforce
						nm.Rules.Ignore = vm.Rules.Ignore
						nm.Rules.Monitor = vm.Rules.Monitor
						nm.Rules.Notify = vm.Rules.Notify
						nm.Rules.Other = vm.Rules.Other
						nm.Rules.Test = vm.Rules.Test
						nm.Rules.Total = vm.Rules.Total
					}
					n.Modules[im] = nm
				}
			}

			// Rules

			if v.Type != nil {
				n.Type = v.Type
			}

			cr.Status.AtProvider.Policies[i] = n
		}
	}

	if resp.ReadOnly != nil {
		cr.Spec.ForProvider.ReadOnly = resp.ReadOnly
	}

	// SourceControl

	if resp.Tokens != nil {
		cr.Status.AtProvider.Tokens = make([]*v1alpha1.V1Token, len(resp.Tokens))
		for i, v := range resp.Tokens {
			n := &v1alpha1.V1Token{}
			n.AllowPathPatterns = v.AllowPathPatterns
			n.Description = v.Description
			n.ID = v.ID
			n.Expires = GenerateTimePtr(v.Expires)

			if v.Metadata != nil {
				n.Metadata = &v1alpha1.V1ObjectMeta{}
				n.Metadata.CreatedBy = v.Metadata.CreatedBy
				n.Metadata.CreatedThrough = v.Metadata.CreatedThrough
				n.Metadata.CreatedAt = GenerateTime(v.Metadata.CreatedAt)
				n.Metadata.LastModifiedBy = v.Metadata.LastModifiedBy
				n.Metadata.LastModifiedThrough = v.Metadata.LastModifiedThrough
				n.Metadata.LastModifiedAt = GenerateTime(v.Metadata.LastModifiedAt)
			}

			n.Token = v.Token
			n.TTL = v.TTL
			n.Uses = v.Uses

			cr.Status.AtProvider.Tokens[i] = n
		}
	}

	if resp.Type != nil {
		cr.Spec.ForProvider.Type = resp.Type
	}

	// Uninstall

	// Warnings

	return cr
}

// GenerateSystemPostRequest generates models.V1SystemsPostRequest from v1alpha1.System
func GenerateSystemPostRequest(cr *v1alpha1.System) *models.V1SystemsPostRequest {
	decisionMappings := make(map[string]models.V1RuleDecisionMappings, len(cr.Spec.ForProvider.DecisionMappings))

	return &models.V1SystemsPostRequest{
		BundleRegistry:       cr.Spec.ForProvider.BundleRegistry.DeepCopyToModel(),
		DecisionMappings:     decisionMappings,
		DeploymentParameters: cr.Spec.ForProvider.DeploymentParameters.DeepCopyToModel(),
		Description:          cr.Spec.ForProvider.Description,
		ExternalID:           cr.Spec.ForProvider.ExternalID,
		Name:                 styraclient.String(cr.ObjectMeta.Name),
		ReadOnly:             cr.Spec.ForProvider.ReadOnly,
		SourceControl:        cr.Spec.ForProvider.SourceControl.DeepCopyToModel(),
		Type:                 cr.Spec.ForProvider.Type,
	}
}

// GenerateSystemPutRequest generates models.V1SystemsPutRequest from v1alpha1.System
func GenerateSystemPutRequest(cr *v1alpha1.System) *models.V1SystemsPutRequest {
	decisionMappings := make(map[string]models.V1RuleDecisionMappings, len(cr.Spec.ForProvider.DecisionMappings))

	return &models.V1SystemsPutRequest{
		BundleRegistry:       cr.Spec.ForProvider.BundleRegistry.DeepCopyToModel(),
		DecisionMappings:     decisionMappings,
		DeploymentParameters: cr.Spec.ForProvider.DeploymentParameters.DeepCopyToModel(),
		Description:          cr.Spec.ForProvider.Description,
		ExternalID:           cr.Spec.ForProvider.ExternalID,
		Name:                 styraclient.String(cr.ObjectMeta.Name),
		ReadOnly:             cr.Spec.ForProvider.ReadOnly,
		SourceControl:        cr.Spec.ForProvider.SourceControl.DeepCopyToModel(),
		Type:                 cr.Spec.ForProvider.Type,
	}
}

// GenerateTimePtr generates v1.Time from strfmt.DateTime
func GenerateTimePtr(d *strfmt.DateTime) *v1.Time {
	if d == nil || d.Equal(strfmt.DateTime{}) {
		return nil
	}

	t := time.Time(*d)
	vt := v1.NewTime(t)
	return &vt
}

// GenerateTime generates v1.Time from strfmt.DateTime
func GenerateTime(d strfmt.DateTime) v1.Time {
	t := time.Time(d)
	vt := v1.NewTime(t)
	return vt
}

// IsNotFound returns whether the given error is of type NotFound or not.
func IsNotFound(err error) bool {
	_, ok := err.(*systems.GetSystemNotFound)
	return ok
}

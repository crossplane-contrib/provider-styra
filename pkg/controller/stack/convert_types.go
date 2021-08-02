package stack

import (
	"time"

	"github.com/go-openapi/strfmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mistermx/styra-go-client/pkg/client/stacks"
	"github.com/mistermx/styra-go-client/pkg/models"

	"github.com/crossplane-contrib/provider-styra/apis/stack/v1alpha1"
	styraclient "github.com/crossplane-contrib/provider-styra/pkg/client"
)

// GenerateStack generates a V1StackConfig from Stack
func GenerateStack(resp *models.V1StackConfig) (cr *v1alpha1.Stack) { // nolint:gocyclo
	cr = &v1alpha1.Stack{}

	// if resp.DecisionMappings != nil {
	// 	cr.Spec.ForProvider.DecisionMappings = make(map[string]v1alpha1.V1RuleDecisionMappings, len(resp.DecisionMappings))
	// 	for k, v := range resp.DecisionMappings {
	// 		n := &v1alpha1.V1RuleDecisionMappings{}
	// 	}
	// }

	// DeploymentParameters

	cr.Spec.ForProvider.Description = styraclient.StringValue(resp.Description)

	// Errors

	// Install

	if resp.Metadata != nil {
		cr.Status.AtProvider.Metadata = &v1alpha1.V1ObjectMeta{}
		cr.Status.AtProvider.Metadata.CreatedBy = resp.Metadata.CreatedBy
		cr.Status.AtProvider.Metadata.CreatedThrough = resp.Metadata.CreatedThrough
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

	cr.Spec.ForProvider.ReadOnly = styraclient.BoolValue(resp.ReadOnly)

	// SourceControl

	cr.Spec.ForProvider.Type = styraclient.StringValue(resp.Type)

	// Uninstall

	// Warnings

	return cr
}

// GenerateStackPostRequest generates models.V1StacksPostRequest from v1alpha1.Stack
func GenerateStackPostRequest(cr *v1alpha1.Stack) *models.V1StacksPostRequest {
	return &models.V1StacksPostRequest{
		Description:   styraclient.String(cr.Spec.ForProvider.Description),
		Name:          styraclient.String(cr.ObjectMeta.Name),
		ReadOnly:      styraclient.Bool(cr.Spec.ForProvider.ReadOnly),
		SourceControl: cr.Spec.ForProvider.SourceControl.DeepCopyToModel(),
		Type:          styraclient.String(cr.Spec.ForProvider.Type),
	}
}

// GenerateStackPutRequest generates models.V1StacksPutRequest from v1alpha1.Stack
func GenerateStackPutRequest(cr *v1alpha1.Stack) *models.V1StacksPutRequest {
	return &models.V1StacksPutRequest{
		Description:   styraclient.String(cr.Spec.ForProvider.Description),
		Name:          styraclient.String(cr.ObjectMeta.Name),
		ReadOnly:      styraclient.Bool(cr.Spec.ForProvider.ReadOnly),
		SourceControl: cr.Spec.ForProvider.SourceControl.DeepCopyToModel(),
		Type:          styraclient.String(cr.Spec.ForProvider.Type),
	}
}

// GenerateTime generates v1.Time from strfmt.DateTime
func GenerateTime(d *strfmt.DateTime) *v1.Time {
	if d == nil {
		return nil
	}

	t := time.Time(*d)
	vt := v1.NewTime(t)
	return &vt
}

// IsNotFound returns whether the given error is of type NotFound or not.
func IsNotFound(err error) bool {
	_, ok := err.(*stacks.GetStackNotFound)
	return ok
}

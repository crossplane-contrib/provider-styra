package stack

import (
	"github.com/mistermx/styra-go-client/pkg/client/stacks"
	"github.com/mistermx/styra-go-client/pkg/models"

	"github.com/crossplane-contrib/provider-styra/apis/stack/v1alpha1"
	styraclient "github.com/crossplane-contrib/provider-styra/pkg/client"
)

// GenerateStack generates a V1StackConfig from Stack
func generateStack(resp *models.V1StackConfig) (cr *v1alpha1.Stack) {
	cr = &v1alpha1.Stack{}

	if resp.SourceControl != nil && resp.SourceControl.Origin != nil {
		cr.Spec.ForProvider.SourceControl = &v1alpha1.V1SourceControlConfig{}
		cr.Spec.ForProvider.SourceControl.Origin = v1alpha1.V1GitRepoConfig{
			Credentials: styraclient.StringValue(resp.SourceControl.Origin.Credentials),
			Path:        styraclient.StringValue(resp.SourceControl.Origin.Path),
			Reference:   styraclient.StringValue(resp.SourceControl.Origin.Reference),
			URL:         styraclient.StringValue(resp.SourceControl.Origin.URL),
		}
	}

	return cr
}

// GenerateStackPostRequest generates models.V1StacksPostRequest from v1alpha1.Stack
func generateStackPostRequest(cr *v1alpha1.Stack) *models.V1StacksPostRequest {
	return &models.V1StacksPostRequest{
		Description:   styraclient.String(cr.Spec.ForProvider.Description),
		Name:          styraclient.String(cr.ObjectMeta.Name),
		ReadOnly:      styraclient.Bool(cr.Spec.ForProvider.ReadOnly),
		SourceControl: generateModelSourceControlConfig(cr.Spec.ForProvider.SourceControl),
		Type:          styraclient.String(cr.Spec.ForProvider.Type),
	}
}

// GenerateStackPutRequest generates models.V1StacksPutRequest from v1alpha1.Stack
func generateStackPutRequest(cr *v1alpha1.Stack) *models.V1StacksPutRequest {
	return &models.V1StacksPutRequest{
		Description:   styraclient.String(cr.Spec.ForProvider.Description),
		Name:          styraclient.String(cr.ObjectMeta.Name),
		ReadOnly:      styraclient.Bool(cr.Spec.ForProvider.ReadOnly),
		SourceControl: generateModelSourceControlConfig(cr.Spec.ForProvider.SourceControl),
		Type:          styraclient.String(cr.Spec.ForProvider.Type),
	}
}

func generateModelSourceControlConfig(spec *v1alpha1.V1SourceControlConfig) *models.V1SourceControlConfig {
	if spec == nil {
		return nil
	}
	return &models.V1SourceControlConfig{
		Origin: &models.V1GitRepoConfig{
			Credentials: styraclient.String(spec.Origin.Credentials),
			Path:        styraclient.String(spec.Origin.Path),
			Reference:   styraclient.String(spec.Origin.Reference),
			URL:         styraclient.String(spec.Origin.URL),
		},
	}
}

func lateInitializeSourceControlConfig(spec, current *v1alpha1.V1SourceControlConfig) *v1alpha1.V1SourceControlConfig {
	if spec != nil {
		return spec
	}
	return current
}

// IsNotFound returns whether the given error is of type NotFound or not.
func IsNotFound(err error) bool {
	_, ok := err.(*stacks.GetStackNotFound)
	return ok
}

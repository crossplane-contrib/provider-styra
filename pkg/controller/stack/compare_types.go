package stack

import (
	models "github.com/mistermx/styra-go-client/pkg/models"

	v1alpha1 "github.com/crossplane-contrib/provider-styra/apis/stack/v1alpha1"
	styraclient "github.com/crossplane-contrib/provider-styra/pkg/client"
)

func isEqualSourceControlConfig(spec *v1alpha1.V1SourceControlConfig, current *models.V1SourceControlConfig) bool {
	if spec == nil {
		return current == nil
	}
	if current == nil || current.Origin == nil {
		return false
	}
	return spec.Origin.Credentials == styraclient.StringValue(current.Origin.Credentials) &&
		spec.Origin.Path == styraclient.StringValue(current.Origin.Path) &&
		spec.Origin.Reference == styraclient.StringValue(current.Origin.Reference) &&
		spec.Origin.URL == styraclient.StringValue(current.Origin.URL)
}

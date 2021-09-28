package system

import (
	models "github.com/mistermx/styra-go-client/pkg/models"

	v1alpha1 "github.com/crossplane-contrib/provider-styra/apis/system/v1alpha1"
	styraclient "github.com/crossplane-contrib/provider-styra/pkg/client"
)

func isEqualSystemDeploymentParameters(spec *v1alpha1.V1SystemDeploymentParameters, current *models.V1SystemDeploymentParameters) bool { // nolint:gocyclo
	if spec == nil {
		return current == nil
	}

	if current == nil {
		return false
	}

	return styraclient.IsEqualBool(spec.DenyOnOpaFail, current.DenyOnOpaFail) &&
		styraclient.StringValue(spec.HTTPProxy) == current.HTTPProxy &&
		styraclient.StringValue(spec.HTTPSProxy) == current.HTTPSProxy &&
		styraclient.StringValue(spec.KubernetesVersion) == current.KubernetesVersion &&
		styraclient.StringValue(spec.Namespace) == current.Namespace &&
		styraclient.StringValue(spec.NoProxy) == current.NoProxy &&
		styraclient.Int32Value(spec.TimeoutSeconds) == current.TimeoutSeconds &&
		styraclient.IsEqualStringArrayContent(spec.TrustedCaCerts, current.TrustedCaCerts) &&
		styraclient.StringValue(spec.TrustedContainerRegistry) == current.TrustedContainerRegistry
}

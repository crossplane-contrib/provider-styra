package system

import (
	models "github.com/mistermx/styra-go-client/pkg/models"

	v1alpha1 "github.com/crossplane-contrib/provider-styra/apis/system/v1alpha1"
	styraclient "github.com/crossplane-contrib/provider-styra/pkg/client"
)

func isEqualBundleRegistry(spec *v1alpha1.V1BundleRegistryConfig, current *models.V1BundleRegistryConfig) bool { // nolint:gocyclo
	if spec == nil {
		return current == nil
	}

	if current == nil {
		return false
	}

	if spec.ManualDeployment != current.ManualDeployment ||
		spec.MaxBundles != current.MaxBundles ||
		spec.MaxDeployedBundles != current.MaxDeployedBundles ||
		spec.DistributionS3.AccessKeys != current.DistributionS3.AccessKeys ||
		!styraclient.IsEqualString(spec.DistributionS3.Bucket, current.DistributionS3.Bucket) ||
		!styraclient.IsEqualString(spec.DistributionS3.DiscoveryPath, current.DistributionS3.DiscoveryPath) ||
		!styraclient.IsEqualString(spec.DistributionS3.PolicyPath, current.DistributionS3.PolicyPath) ||
		!styraclient.IsEqualString(spec.DistributionS3.Region, current.DistributionS3.Region) ||
		!styraclient.IsEqualString(spec.DistributionS3.OpaCredentials.MetadataCredentials.AwsRegion, current.DistributionS3.OpaCredentials.MetadataCredentials.AwsRegion) ||
		!styraclient.IsEqualString(spec.DistributionS3.OpaCredentials.MetadataCredentials.IamRole, current.DistributionS3.OpaCredentials.MetadataCredentials.IamRole) ||
		!styraclient.IsEqualString(spec.DistributionS3.OpaCredentials.WebIdentityCredentials.AwsRegion, current.DistributionS3.OpaCredentials.WebIdentityCredentials.AwsRegion) ||
		!styraclient.IsEqualString(spec.DistributionS3.OpaCredentials.WebIdentityCredentials.SessionName, current.DistributionS3.OpaCredentials.WebIdentityCredentials.SessionName) {
		return false
	}

	return true
}

func isEqualDecisionMapping(spec map[string]v1alpha1.V1RuleDecisionMappings, current map[string]models.V1RuleDecisionMappings) bool {
	if len(spec) != len(current) {
		return false
	}

	for k, specMapping := range spec {
		currentMapping, exists := current[k]
		if !exists {
			return false
		}

		if !isEqualAllowedMapping(specMapping.Allowed, currentMapping.Allowed) ||
			!isEqualReasonMapping(specMapping.Reason, currentMapping.Reason) ||
			!isEqualColumnMapping(specMapping.Columns, currentMapping.Columns) {
			return false
		}

	}

	return true
}

func isEqualAllowedMapping(spec *v1alpha1.V1AllowedMapping, current *models.V1AllowedMapping) bool {
	if spec == nil {
		return current == nil
	}

	if current == nil {
		return false
	}

	return styraclient.IsEqualString(spec.Path, current.Path) &&
		styraclient.IsEqualBool(spec.Negated, current.Negated)
}

func isEqualReasonMapping(spec *v1alpha1.V1ReasonMapping, current *models.V1ReasonMapping) bool {
	if spec == nil {
		return current == nil
	}

	if current == nil {
		return false
	}

	return styraclient.IsEqualString(spec.Path, current.Path)
}

func isEqualColumnMapping(spec []*v1alpha1.V1ColumnMapping, current []*models.V1ColumnMapping) bool {
	if len(spec) != len(current) {
		return false
	}

	currentMap := make(map[string]*models.V1ColumnMapping, len(current))
	for _, c := range currentMap {
		cKeyVal := styraclient.StringValue(c.Key)
		currentMap[cKeyVal] = c
	}

	for _, s := range currentMap {
		sKeyVal := styraclient.StringValue(s.Key)
		c, exists := currentMap[sKeyVal]

		if !exists ||
			!styraclient.IsEqualString(s.Path, c.Path) ||
			!styraclient.IsEqualString(s.Type, c.Type) {
			return false
		}
	}

	return true
}

func isEqualSystemDeploymentParameters(spec *v1alpha1.V1SystemDeploymentParameters, current *models.V1SystemDeploymentParameters) bool { // nolint:gocyclo
	if spec == nil {
		return current == nil
	}

	if current == nil {
		return false
	}

	return styraclient.IsEqualBool(spec.DenyOnOpaFail, current.DenyOnOpaFail) &&
		spec.HTTPProxy == current.HTTPProxy &&
		spec.HTTPSProxy == current.HTTPSProxy &&
		spec.KubernetesVersion == current.KubernetesVersion &&
		spec.Namespace == current.Namespace &&
		spec.NoProxy == current.NoProxy &&
		// spec.TimeoutSeconds == current.TimeoutSeconds  // default value is 1 not 0, rework for lateinitialize required
		styraclient.IsEqualStringArrayContent(spec.TrustedCaCerts, current.TrustedCaCerts) &&
		spec.TrustedContainerRegistry == current.TrustedContainerRegistry
}

func isEqualSourceControlConfig(spec *v1alpha1.V1SourceControlConfig, current *models.V1SourceControlConfig) bool {
	if spec == nil {
		return current == nil
	}

	if current == nil {
		return false
	}

	return styraclient.IsEqualString(spec.Origin.Credentials, current.Origin.Credentials) &&
		styraclient.IsEqualString(spec.Origin.Path, current.Origin.Path) &&
		styraclient.IsEqualString(spec.Origin.Reference, current.Origin.Reference) &&
		styraclient.IsEqualString(spec.Origin.URL, current.Origin.URL)
}

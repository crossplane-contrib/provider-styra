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

// generateSystem generates a V1SystemConfig from System
func generateSystem(resp *models.V1SystemConfig) (cr *v1alpha1.System) { // nolint:gocyclo
	cr = &v1alpha1.System{}

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
				if v.Status.Timestamp != nil {
					n.Status.Timestamp = generateTimePtr(v.Status.Timestamp)
				}
			}
			cr.Status.AtProvider.Datasources[i] = n
		}
	}
	if resp.DeploymentParameters != nil {
		cr.Spec.ForProvider.DeploymentParameters = &v1alpha1.V1SystemDeploymentParameters{
			DenyOnOpaFail:            resp.DeploymentParameters.DenyOnOpaFail,
			HTTPProxy:                styraclient.String(resp.DeploymentParameters.HTTPProxy),
			HTTPSProxy:               styraclient.String(resp.DeploymentParameters.HTTPSProxy),
			KubernetesVersion:        styraclient.String(resp.DeploymentParameters.KubernetesVersion),
			Namespace:                styraclient.String(resp.DeploymentParameters.Namespace),
			NoProxy:                  styraclient.String(resp.DeploymentParameters.NoProxy),
			TimeoutSeconds:           styraclient.Int32(resp.DeploymentParameters.TimeoutSeconds),
			TrustedCaCerts:           resp.DeploymentParameters.TrustedCaCerts,
			TrustedContainerRegistry: styraclient.String(resp.DeploymentParameters.TrustedContainerRegistry),
		}
	}

	cr.Spec.ForProvider.Description = &resp.Description

	if resp.Errors != nil {
		cr.Status.AtProvider.Errors = map[string]v1alpha1.V1AgentErrors{}
		for k, e := range resp.Errors {
			errors := v1alpha1.V1AgentErrors{}
			errors.Waiting = e.Waiting
			errors.Errors = make([]*v1alpha1.V1Status, len(e.Errors))
			for i, ee := range e.Errors {
				errors.Errors[i] = &v1alpha1.V1Status{
					Code:      ee.Code,
					Message:   ee.Message,
					Timestamp: generateTimePtr(ee.Timestamp),
				}
			}
			cr.Status.AtProvider.Errors[k] = errors
		}
	}

	cr.Spec.ForProvider.ExternalID = &resp.ExternalID
	cr.Spec.ForProvider.ReadOnly = resp.ReadOnly
	cr.Spec.ForProvider.Type = styraclient.StringValue(resp.Type)

	if resp.Warnings != nil {
		cr.Status.AtProvider.Warnings = map[string]v1alpha1.V1SystemConfigWarnings{}
		for k, w := range resp.Warnings {
			cr.Status.AtProvider.Warnings[k] = v1alpha1.V1SystemConfigWarnings{
				Code:      w.Code,
				Message:   w.Message,
				Timestamp: generateTimePtr(w.Timestamp),
			}
		}
	}

	return cr
}

// generateSystemPostRequest generates models.V1SystemsPostRequest from v1alpha1.System
func generateSystemPostRequest(cr *v1alpha1.System) *models.V1SystemsPostRequest {
	return &models.V1SystemsPostRequest{
		DeploymentParameters: generateDeploymentParameters(cr.Spec.ForProvider.DeploymentParameters),
		Description:          styraclient.StringValue(cr.Spec.ForProvider.Description),
		ExternalID:           styraclient.StringValue(cr.Spec.ForProvider.ExternalID),
		Name:                 styraclient.String(cr.ObjectMeta.Name),
		ReadOnly:             cr.Spec.ForProvider.ReadOnly,
		Type:                 styraclient.String(cr.Spec.ForProvider.Type),
	}
}

// generateSystemPutRequest generates models.V1SystemsPutRequest from v1alpha1.System
func generateSystemPutRequest(cr *v1alpha1.System) *models.V1SystemsPutRequest {
	return &models.V1SystemsPutRequest{
		DeploymentParameters: generateDeploymentParameters(cr.Spec.ForProvider.DeploymentParameters),
		Description:          styraclient.StringValue(cr.Spec.ForProvider.Description),
		ExternalID:           styraclient.StringValue(cr.Spec.ForProvider.ExternalID),
		Name:                 styraclient.String(cr.ObjectMeta.Name),
		ReadOnly:             cr.Spec.ForProvider.ReadOnly,
		Type:                 styraclient.String(cr.Spec.ForProvider.Type),
	}
}

func generateDeploymentParameters(spec *v1alpha1.V1SystemDeploymentParameters) *models.V1SystemDeploymentParameters {
	if spec != nil {
		return &models.V1SystemDeploymentParameters{
			DenyOnOpaFail:            spec.DenyOnOpaFail,
			HTTPProxy:                styraclient.StringValue(spec.HTTPProxy),
			HTTPSProxy:               styraclient.StringValue(spec.HTTPSProxy),
			KubernetesVersion:        styraclient.StringValue(spec.KubernetesVersion),
			Namespace:                styraclient.StringValue(spec.Namespace),
			NoProxy:                  styraclient.StringValue(spec.NoProxy),
			TimeoutSeconds:           styraclient.Int32Value(spec.TimeoutSeconds),
			TrustedCaCerts:           spec.TrustedCaCerts,
			TrustedContainerRegistry: styraclient.StringValue(spec.TrustedContainerRegistry),
		}
	}
	return nil
}

func lateInitializeDeploymentParameters(spec *v1alpha1.V1SystemDeploymentParameters, current *v1alpha1.V1SystemDeploymentParameters) *v1alpha1.V1SystemDeploymentParameters {
	if current == nil {
		return spec
	}

	if spec == nil {
		return current
	}

	spec.DenyOnOpaFail = styraclient.LateInitializeBoolPtr(spec.DenyOnOpaFail, current.DenyOnOpaFail)
	spec.HTTPProxy = styraclient.LateInitializeStringPtr(spec.HTTPProxy, current.HTTPProxy)
	spec.HTTPSProxy = styraclient.LateInitializeStringPtr(spec.HTTPSProxy, current.HTTPSProxy)
	spec.KubernetesVersion = styraclient.LateInitializeStringPtr(spec.KubernetesVersion, current.KubernetesVersion)
	spec.Namespace = styraclient.LateInitializeStringPtr(spec.Namespace, current.Namespace)
	spec.NoProxy = styraclient.LateInitializeStringPtr(spec.NoProxy, current.NoProxy)
	spec.TimeoutSeconds = styraclient.LateInitializeInt32Ptr(spec.TimeoutSeconds, current.TimeoutSeconds)
	if spec.TrustedCaCerts == nil {
		spec.TrustedCaCerts = current.TrustedCaCerts
	}
	spec.TrustedContainerRegistry = styraclient.LateInitializeStringPtr(spec.TrustedContainerRegistry, current.TrustedContainerRegistry)
	return spec
}

// generateTimePtr generates v1.Time from strfmt.DateTime
func generateTimePtr(d *strfmt.DateTime) *v1.Time {
	if d == nil || d.Equal(strfmt.DateTime{}) {
		return nil
	}

	t := time.Time(*d)
	vt := v1.NewTime(t)
	return &vt
}

// isNotFound returns whether the given error is of type NotFound or not.
func isNotFound(err error) bool {
	_, ok := err.(*systems.GetSystemNotFound)
	return ok
}

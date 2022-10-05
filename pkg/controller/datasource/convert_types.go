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

package datasource

import (
	"context"
	"time"

	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/mistermx/styra-go-client/pkg/client/datasources"
	"github.com/mistermx/styra-go-client/pkg/models"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/crossplane-contrib/provider-styra/apis/datasource/v1alpha1"
	styraclient "github.com/crossplane-contrib/provider-styra/pkg/client"
)

func generateDataSource(resp *models.DatasourcesV1DatasourcesGetResponseResult) *v1alpha1.DataSource { // nolint:gocyclo
	if resp == nil {
		return &v1alpha1.DataSource{}
	}
	cr := &v1alpha1.DataSource{}

	cr.Spec.ForProvider.DatasourcesV1Common = v1alpha1.DatasourcesV1Common{
		Category:    resp.Category,
		Description: &resp.Description,
		Enabled:     &resp.Enabled,
		OnPremises:  resp.OnPremises,
		Type:        resp.Type,
	}

	rateLimiter := v1alpha1.DatasourcesV1RateLimiter{
		RateLimit: styraclient.Float64ToQuantity(resp.RateLimit),
	}
	poller := v1alpha1.DatasourcesV1Poller{
		PollingInterval: generateDurationFromSeconds(resp.PollingInterval),
	}

	switch cr.Spec.ForProvider.Category {
	case v1alpha1.DataSourceCategoryAWSECR:
		cr.Spec.ForProvider.AWSECR = &v1alpha1.DatasourcesV1AWSECR{
			DatasourcesV1AWSCommon: v1alpha1.DatasourcesV1AWSCommon{
				Credentials: resp.Credentials,
				Region:      resp.Region,
			},
			RegistryID:               &resp.RegistryID,
			DatasourcesV1RateLimiter: rateLimiter,
			DatasourcesV1Poller:      poller,
		}
	case v1alpha1.DataSourceCategoryBundleS3:
		cr.Spec.ForProvider.BundleS3 = &v1alpha1.DatasourcesV1BundleS3{
			DatasourcesV1AWSCommon: v1alpha1.DatasourcesV1AWSCommon{
				Credentials: resp.Credentials,
				Region:      resp.Region,
			},
			Bucket:              resp.Bucket,
			Endpoint:            &resp.Endpoint,
			Path:                resp.Path,
			DatasourcesV1Poller: poller,
		}
	case v1alpha1.DataSourceCategoryGitBlame:
		cr.Spec.ForProvider.GitBlame = &v1alpha1.DatasourcesV1GitBlame{
			DatasourcesV1GitCommon: generateDatasourcesV1GitCommon(resp),
			PathRegexp:             &resp.PathRegexp,
		}
	case v1alpha1.DataSourceCategoryGitContent:
		cr.Spec.ForProvider.GitContent = &v1alpha1.DatasourcesV1GitContent{
			DatasourcesV1GitCommon: generateDatasourcesV1GitCommon(resp),
		}
	case v1alpha1.DataSourceCategoryGitRego:
		cr.Spec.ForProvider.GitRego = &v1alpha1.DatasourcesV1GitRego{
			DatasourcesV1GitCommon: generateDatasourcesV1GitCommon(resp),
			Path:                   &resp.Path,
		}
	case v1alpha1.DataSourceCategoryHTTP:
		cr.Spec.ForProvider.HTTP = &v1alpha1.DatasourcesV1HTTP{
			DatasourcesV1RegoFiltering: v1alpha1.DatasourcesV1RegoFiltering{
				PolicyFilter: &resp.PolicyFilter,
				PolicyQuery:  &resp.PolicyQuery,
			},
			DatasourcesV1TLSSettings: v1alpha1.DatasourcesV1TLSSettings{
				CaCertificate:       &resp.CaCertificate,
				SkipTLSVerification: &resp.SkipTLSVerification,
			},
			Headers:             generateDatasourcesV1HTTPHeader(resp.Headers),
			URL:                 resp.URL,
			DatasourcesV1Poller: poller,
		}
	case v1alpha1.DataSourceCategoryKubernetesResources:
		cr.Spec.ForProvider.KubernetesResources = &v1alpha1.DatasourcesV1KubernetesResources{
			Masks:                    resp.Masks,
			Namespaces:               resp.Namespaces,
			Selectors:                resp.Selectors,
			DatasourcesV1RateLimiter: rateLimiter,
			DatasourcesV1Poller:      poller,
		}
	case v1alpha1.DataSourceCategoryLDAP:
		cr.Spec.ForProvider.LDAP = &v1alpha1.DatasourcesV1LDAP{
			DatasourcesV1RegoFiltering: v1alpha1.DatasourcesV1RegoFiltering{
				PolicyFilter: &resp.PolicyFilter,
				PolicyQuery:  &resp.PolicyQuery,
			},
			DatasourcesV1TLSSettings: v1alpha1.DatasourcesV1TLSSettings{
				CaCertificate:       &resp.CaCertificate,
				SkipTLSVerification: &resp.SkipTLSVerification,
			},
			Credentials:              &resp.Credentials,
			Search:                   generateLDAPSearch(resp.Search),
			Urls:                     resp.Urls,
			DatasourcesV1RateLimiter: rateLimiter,
			DatasourcesV1Poller:      poller,
		}
	case v1alpha1.DataSourceCategoryPolicyLibrary:
		cr.Spec.ForProvider.PolicyLibrary = &v1alpha1.DatasourcesV1PolicyLibrary{
			DatasourcesV1Poller: poller,
		}
	case v1alpha1.DataSourceCategoryRest:
		cr.Spec.ForProvider.Rest = &v1alpha1.DatasourcesV1Rest{
			ContentType: &resp.ContentType,
		}
	}

	cr.Status.AtProvider.Executed = &resp.Executed

	if resp.Status != nil {
		cr.Status.AtProvider.Status = &v1alpha1.DataSourceExternalStatus{
			Code:      &resp.Status.Code,
			Message:   &resp.Status.Message,
			Timestamp: styraclient.String(resp.Status.Timestamp.String()),
		}
	}

	return cr
}

func generateDurationFromSeconds(seconds int64) *metav1.Duration {
	d := (time.Duration)(seconds) * time.Second //nolint:durationcheck
	return &metav1.Duration{Duration: d}
}

func generateDatasourcesV1GitCommon(resp *models.DatasourcesV1DatasourcesGetResponseResult) v1alpha1.DatasourcesV1GitCommon {
	gitCommon := v1alpha1.DatasourcesV1GitCommon{
		Credentials: &resp.Credentials,
		Reference:   &resp.Reference,

		Timeout: generateDurationFromSeconds(resp.Timeout),
		URL:     resp.URL,
		DatasourcesV1RateLimiter: v1alpha1.DatasourcesV1RateLimiter{
			RateLimit: styraclient.Float64ToQuantity(resp.RateLimit),
		},
		DatasourcesV1Poller: v1alpha1.DatasourcesV1Poller{
			PollingInterval: generateDurationFromSeconds(resp.PollingInterval),
		},
	}

	if resp.SSHCredentials != nil {
		gitCommon.SSHCredentials = &v1alpha1.DatasourcesV1GitCommonAO3SSHCredentials{
			Passphrase: &resp.SSHCredentials.Passphrase,
			PrivateKey: styraclient.StringValue(resp.SSHCredentials.PrivateKey),
		}
	}

	return gitCommon
}

func generateDatasourcesV1HTTPHeader(current []*models.DatasourcesV1HTTPHeader) []v1alpha1.DatasourcesV1HTTPHeader {
	if current == nil {
		return nil
	}

	res := make([]v1alpha1.DatasourcesV1HTTPHeader, len(current))
	for i, cur := range current {
		if cur != nil {
			res[i] = v1alpha1.DatasourcesV1HTTPHeader{
				Name:     styraclient.StringValue(cur.Name),
				SecretID: &cur.SecretID,
				Value:    &cur.Value,
			}
		}
	}
	return res
}

func generateLDAPSearch(current *models.DatasourcesV1DatasourcesGetResponseResultSearch) *v1alpha1.DatasourcesV1LDAPAO5Search {
	if current == nil {
		return nil
	}

	return &v1alpha1.DatasourcesV1LDAPAO5Search{
		Attributes: current.Attributes,
		BaseDN:     current.BaseDN,
		Deref:      current.Deref,
		Filter:     current.Filter,
		PageSize:   current.PageSize,
		Scope:      current.Scope,
		SizeLimit:  current.SizeLimit,
	}
}

const (
	errRequireFieldFormat = "category '%s' requires field 'spec.forProvider.%s"
)

func errorRequireField(category, fieldName string) error {
	return errors.Errorf(errRequireFieldFormat, category, fieldName)
}

func generateDataSourceUpsertParams(ctx context.Context, cr *v1alpha1.DataSource) (*datasources.UpsertDatasourceParams, error) { // nolint:gocyclo
	req := &datasources.UpsertDatasourceParams{
		Datasource: meta.GetExternalName(cr),
		Context:    ctx,
	}

	common := models.DatasourcesV1Common{
		Category:    styraclient.String(cr.Spec.ForProvider.Category),
		Description: styraclient.StringValue(cr.Spec.ForProvider.Description),
		Enabled:     cr.Spec.ForProvider.Enabled,
		OnPremises:  &cr.Spec.ForProvider.OnPremises,
		Type:        &cr.Spec.ForProvider.Type,
	}

	switch cr.Spec.ForProvider.Category {
	case v1alpha1.DataSourceCategoryAWSECR:
		if cr.Spec.ForProvider.AWSECR == nil {
			return nil, errorRequireField(cr.Spec.ForProvider.Category, "awsECR")
		}

		req.Body = &models.DatasourcesV1AWSECR{
			DatasourcesV1Common:      common,
			DatasourcesV1Poller:      generateModelPoller(cr.Spec.ForProvider.AWSECR.DatasourcesV1Poller),
			DatasourcesV1RateLimiter: generateModelRateLimiter(cr.Spec.ForProvider.AWSECR.DatasourcesV1RateLimiter),
			DatasourcesV1AWSCommon: models.DatasourcesV1AWSCommon{
				Credentials: &cr.Spec.ForProvider.AWSECR.Credentials,
				Region:      &cr.Spec.ForProvider.AWSECR.Region,
			},
			RegistryID: styraclient.StringValue(cr.Spec.ForProvider.AWSECR.RegistryID),
		}
	case v1alpha1.DataSourceCategoryBundleS3:
		if cr.Spec.ForProvider.BundleS3 == nil {
			return nil, errorRequireField(cr.Spec.ForProvider.Category, "bundleS3")
		}

		req.Body = &models.DatasourcesV1BundleS3{
			DatasourcesV1Common: common,
			DatasourcesV1Poller: generateModelPoller(cr.Spec.ForProvider.BundleS3.DatasourcesV1Poller),
			DatasourcesV1AWSCommon: models.DatasourcesV1AWSCommon{
				Credentials: &cr.Spec.ForProvider.BundleS3.Credentials,
				Region:      &cr.Spec.ForProvider.BundleS3.Region,
			},
			Bucket:   &cr.Spec.ForProvider.BundleS3.Bucket,
			Endpoint: styraclient.StringValue(cr.Spec.ForProvider.BundleS3.Endpoint),
			Path:     &cr.Spec.ForProvider.BundleS3.Path,
		}
	case v1alpha1.DataSourceCategoryGitBlame:
		if cr.Spec.ForProvider.GitBlame == nil {
			return nil, errorRequireField(cr.Spec.ForProvider.Category, "gitBlame")
		}

		req.Body = &models.DatasourcesV1GitBlame{
			DatasourcesV1GitCommon: models.DatasourcesV1GitCommon{
				DatasourcesV1Common:      common,
				DatasourcesV1Poller:      generateModelPoller(cr.Spec.ForProvider.GitBlame.DatasourcesV1Poller),
				DatasourcesV1RateLimiter: generateModelRateLimiter(cr.Spec.ForProvider.GitBlame.DatasourcesV1RateLimiter),
				SSHCredentials:           generateModelDatasourcesV1GitCommonAO3SSHCredentials(cr.Spec.ForProvider.GitBlame.SSHCredentials),
				Credentials:              styraclient.StringValue(cr.Spec.ForProvider.GitBlame.Credentials),
				Reference:                cr.Spec.ForProvider.GitBlame.Reference,
				Timeout:                  styraclient.DurationToString(cr.Spec.ForProvider.GitBlame.Timeout),
				URL:                      &cr.Spec.ForProvider.GitBlame.URL,
			},
			PathRegexp: cr.Spec.ForProvider.GitBlame.PathRegexp,
		}
	case v1alpha1.DataSourceCategoryGitContent:
		if cr.Spec.ForProvider.GitContent == nil {
			return nil, errorRequireField(cr.Spec.ForProvider.Category, "gitContent")
		}

		req.Body = &models.DatasourcesV1GitContent{
			DatasourcesV1GitCommon: models.DatasourcesV1GitCommon{
				DatasourcesV1Common:      common,
				DatasourcesV1Poller:      generateModelPoller(cr.Spec.ForProvider.GitContent.DatasourcesV1Poller),
				DatasourcesV1RateLimiter: generateModelRateLimiter(cr.Spec.ForProvider.GitContent.DatasourcesV1RateLimiter),
				SSHCredentials:           generateModelDatasourcesV1GitCommonAO3SSHCredentials(cr.Spec.ForProvider.GitContent.SSHCredentials),
				Credentials:              styraclient.StringValue(cr.Spec.ForProvider.GitContent.Credentials),
				Reference:                cr.Spec.ForProvider.GitContent.Reference,
				Timeout:                  styraclient.DurationToString(cr.Spec.ForProvider.GitContent.Timeout),
				URL:                      &cr.Spec.ForProvider.GitContent.URL,
			},
		}
	case v1alpha1.DataSourceCategoryGitRego:
		if cr.Spec.ForProvider.GitRego == nil {
			return nil, errorRequireField(cr.Spec.ForProvider.Category, "gitRego")
		}

		req.Body = &models.DatasourcesV1GitRego{
			DatasourcesV1GitCommon: models.DatasourcesV1GitCommon{
				DatasourcesV1Common:      common,
				DatasourcesV1Poller:      generateModelPoller(cr.Spec.ForProvider.GitRego.DatasourcesV1Poller),
				DatasourcesV1RateLimiter: generateModelRateLimiter(cr.Spec.ForProvider.GitRego.DatasourcesV1RateLimiter),
				SSHCredentials:           generateModelDatasourcesV1GitCommonAO3SSHCredentials(cr.Spec.ForProvider.GitRego.SSHCredentials),
				Credentials:              styraclient.StringValue(cr.Spec.ForProvider.GitRego.Credentials),
				Reference:                cr.Spec.ForProvider.GitRego.Reference,
				Timeout:                  styraclient.DurationToString(cr.Spec.ForProvider.GitRego.Timeout),
				URL:                      &cr.Spec.ForProvider.GitRego.URL,
			},
			Path: styraclient.StringValue(cr.Spec.ForProvider.GitRego.Path),
		}
	case v1alpha1.DataSourceCategoryHTTP:
		if cr.Spec.ForProvider.HTTP == nil {
			return nil, errorRequireField(cr.Spec.ForProvider.Category, "http")
		}

		req.Body = &models.DatasourcesV1HTTP{
			DatasourcesV1Common: common,
			DatasourcesV1Poller: generateModelPoller(cr.Spec.ForProvider.HTTP.DatasourcesV1Poller),
			DatasourcesV1RegoFiltering: models.DatasourcesV1RegoFiltering{
				PolicyFilter: styraclient.StringValue(cr.Spec.ForProvider.HTTP.PolicyFilter),
				PolicyQuery:  styraclient.StringValue(cr.Spec.ForProvider.HTTP.PolicyQuery),
			},
			DatasourcesV1TLSSettings: models.DatasourcesV1TLSSettings{
				CaCertificate:       styraclient.StringValue(cr.Spec.ForProvider.HTTP.CaCertificate),
				SkipTLSVerification: cr.Spec.ForProvider.HTTP.SkipTLSVerification,
			},
			Headers: generateModelDatasourcesV1HTTPHeader(cr.Spec.ForProvider.HTTP.Headers),
			URL:     &cr.Spec.ForProvider.HTTP.URL,
		}
	case v1alpha1.DataSourceCategoryKubernetesResources:
		if cr.Spec.ForProvider.KubernetesResources == nil {
			return nil, errorRequireField(cr.Spec.ForProvider.Category, "kubernetesResources")
		}

		req.Body = &models.DatasourcesV1KubernetesResources{
			DatasourcesV1Common:      common,
			DatasourcesV1Poller:      generateModelPoller(cr.Spec.ForProvider.KubernetesResources.DatasourcesV1Poller),
			DatasourcesV1RateLimiter: generateModelRateLimiter(cr.Spec.ForProvider.KubernetesResources.DatasourcesV1RateLimiter),
			Masks:                    cr.Spec.ForProvider.KubernetesResources.Masks,
			Namespaces:               cr.Spec.ForProvider.KubernetesResources.Namespaces,
			Selectors:                cr.Spec.ForProvider.KubernetesResources.Selectors,
		}
	case v1alpha1.DataSourceCategoryLDAP:
		if cr.Spec.ForProvider.LDAP == nil {
			return nil, errorRequireField(cr.Spec.ForProvider.Category, "ldap")
		}

		req.Body = &models.DatasourcesV1LDAP{
			DatasourcesV1Common:      common,
			DatasourcesV1Poller:      generateModelPoller(cr.Spec.ForProvider.LDAP.DatasourcesV1Poller),
			DatasourcesV1RateLimiter: generateModelRateLimiter(cr.Spec.ForProvider.LDAP.DatasourcesV1RateLimiter),
			DatasourcesV1RegoFiltering: models.DatasourcesV1RegoFiltering{
				PolicyFilter: styraclient.StringValue(cr.Spec.ForProvider.LDAP.PolicyFilter),
				PolicyQuery:  styraclient.StringValue(cr.Spec.ForProvider.LDAP.PolicyQuery),
			},
			DatasourcesV1TLSSettings: models.DatasourcesV1TLSSettings{
				CaCertificate:       styraclient.StringValue(cr.Spec.ForProvider.LDAP.CaCertificate),
				SkipTLSVerification: cr.Spec.ForProvider.LDAP.SkipTLSVerification,
			},
			Credentials: styraclient.StringValue(cr.Spec.ForProvider.LDAP.Credentials),
			Search:      generateModelDatasourcesV1LDAPAO5Search(cr.Spec.ForProvider.LDAP.Search),
			Urls:        cr.Spec.ForProvider.LDAP.Urls,
		}
	case v1alpha1.DataSourceCategoryPolicyLibrary:
		if cr.Spec.ForProvider.PolicyLibrary == nil {
			return nil, errorRequireField(cr.Spec.ForProvider.Category, "policyLibrary")
		}
		req.Body = &models.DatasourcesV1PolicyLibrary{
			DatasourcesV1Common: common,
			DatasourcesV1Poller: generateModelPoller(cr.Spec.ForProvider.PolicyLibrary.DatasourcesV1Poller),
		}
	case v1alpha1.DataSourceCategoryRest:
		if cr.Spec.ForProvider.Rest == nil {
			return nil, errorRequireField(cr.Spec.ForProvider.Category, "rest")
		}
		req.Body = &models.DatasourcesV1Rest{
			DatasourcesV1Common: common,
			ContentType:         styraclient.StringValue(cr.Spec.ForProvider.Rest.ContentType),
		}
	}

	return req, nil
}

func generateModelDatasourcesV1GitCommonAO3SSHCredentials(spec *v1alpha1.DatasourcesV1GitCommonAO3SSHCredentials) *models.DatasourcesV1GitCommonAO3SSHCredentials {
	if spec == nil {
		return nil
	}

	return &models.DatasourcesV1GitCommonAO3SSHCredentials{
		Passphrase: styraclient.StringValue(spec.Passphrase),
		PrivateKey: &spec.PrivateKey,
	}
}

func generateModelDatasourcesV1HTTPHeader(spec []v1alpha1.DatasourcesV1HTTPHeader) []*models.DatasourcesV1HTTPHeader {
	if spec == nil {
		return nil
	}

	res := make([]*models.DatasourcesV1HTTPHeader, len(spec))
	for i, h := range spec {
		res[i] = &models.DatasourcesV1HTTPHeader{
			Name:     styraclient.String(h.Name),
			SecretID: styraclient.StringValue(h.SecretID),
			Value:    styraclient.StringValue(h.Value),
		}
	}
	return res
}

func generateModelRateLimiter(spec v1alpha1.DatasourcesV1RateLimiter) models.DatasourcesV1RateLimiter {
	return models.DatasourcesV1RateLimiter{
		RateLimit: styraclient.QuantityToFloat64Ptr(spec.RateLimit),
	}
}

func generateModelPoller(spec v1alpha1.DatasourcesV1Poller) models.DatasourcesV1Poller {
	return models.DatasourcesV1Poller{
		PollingInterval: styraclient.DurationToString(spec.PollingInterval),
	}
}

func generateModelDatasourcesV1LDAPAO5Search(spec *v1alpha1.DatasourcesV1LDAPAO5Search) *models.DatasourcesV1LDAPAO5Search {
	if spec == nil {
		return nil
	}
	return &models.DatasourcesV1LDAPAO5Search{
		Attributes: spec.Attributes,
		BaseDN:     &spec.BaseDN,
		Deref:      spec.Deref,
		Filter:     &spec.Filter,
		PageSize:   spec.PageSize,
		Scope:      spec.Scope,
		SizeLimit:  spec.SizeLimit,
	}
}

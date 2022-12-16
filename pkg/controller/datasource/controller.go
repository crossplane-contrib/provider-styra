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

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/google/go-cmp/cmp"

	"github.com/pkg/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"

	styra "github.com/mistermx/styra-go-client/pkg/client"
	"github.com/mistermx/styra-go-client/pkg/client/datasources"
	"github.com/mistermx/styra-go-client/pkg/models"

	"github.com/crossplane-contrib/provider-styra/apis/datasource/v1alpha1"
	styraclient "github.com/crossplane-contrib/provider-styra/pkg/client"
	"github.com/crossplane-contrib/provider-styra/pkg/interface/controller"
)

const (
	errNotDataSource  = "managed resource is not an datasource custom resource"
	errUpdateFailed   = "cannot update datasource"
	errCreateFailed   = "cannot create datasource"
	errDeleteFailed   = "cannot delete datasource"
	errDescribeFailed = "cannot describe datasource"
)

// SetupDataSource adds a controller that reconciles DataSources.
func SetupDataSource(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.DataSourceGroupKind)

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o.ForControllerRuntime()).
		For(&v1alpha1.DataSource{}).
		Complete(managed.NewReconciler(mgr,
			resource.ManagedKind(v1alpha1.DataSourceGroupVersionKind),
			managed.WithExternalConnecter(&connector{kube: mgr.GetClient(), newClientFn: styra.New}),
			managed.WithInitializers(managed.NewDefaultProviderConfig(mgr.GetClient()), managed.NewNameAsExternalName(mgr.GetClient())),
			managed.WithReferenceResolver(managed.NewAPISimpleReferenceResolver(mgr.GetClient())),
			managed.WithPollInterval(o.PollInterval),
			managed.WithLogger(o.Logger.WithValues("controller", name)),
			managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))),
			managed.WithConnectionPublishers(o.ConnectionPublisher...)))
}

type connector struct {
	kube        client.Client
	newClientFn func(transport runtime.ClientTransport, formats strfmt.Registry) *styra.StyraAPI
}

type external struct {
	client *styra.StyraAPI
	kube   client.Client
}

func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	_, ok := mg.(*v1alpha1.DataSource)
	if !ok {
		return nil, errors.New(errNotDataSource)
	}

	cfg, err := styraclient.GetConfig(ctx, c.kube, mg)
	// cfg.Debug = true
	if err != nil {
		return nil, err
	}

	client := c.newClientFn(cfg, strfmt.Default)

	return &external{client, c.kube}, nil
}

func (e *external) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.DataSource)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotDataSource)
	}

	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{}, nil
	}

	req := &datasources.GetDatasourceParams{
		Context:    ctx,
		Datasource: meta.GetExternalName(cr),
	}
	resp, err := e.client.Datasources.GetDatasource(req)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(resource.Ignore(isNotFound, err), errDescribeFailed)
	}

	currentSpec := cr.Spec.ForProvider.DeepCopy()
	externalDataSource := generateDataSource(resp.Payload.Result)
	externalDataSource.Status.AtProvider.DeepCopyInto(&cr.Status.AtProvider)
	lateInitialize(cr, resp.Payload.Result)

	if cr.Status.AtProvider.Status != nil {
		switch styraclient.StringValue(cr.Status.AtProvider.Status.Code) {
		case v1alpha1.DataSourceStatusFailed:
			cr.Status.SetConditions(v1.Unavailable())
		default:
			cr.Status.SetConditions(v1.Available())
		}
	} else {
		cr.Status.SetConditions(v1.Available())
	}

	return managed.ExternalObservation{
		ResourceExists:          true,
		ResourceUpToDate:        isUpToDate(cr, externalDataSource),
		ResourceLateInitialized: !cmp.Equal(&cr.Spec.ForProvider, currentSpec),
	}, nil
}

func isUpToDate(spec, current *v1alpha1.DataSource) bool { // nolint:gocyclo
	if !cmp.Equal(spec.Spec.ForProvider.DatasourcesV1Common, current.Spec.ForProvider.DatasourcesV1Common) {
		return false
	}

	switch spec.Spec.ForProvider.Category {
	case v1alpha1.DataSourceCategoryAWSECR:
		return cmp.Equal(spec.Spec.ForProvider.AWSECR, current.Spec.ForProvider.AWSECR)
	case v1alpha1.DataSourceCategoryBundleS3:
		return cmp.Equal(spec.Spec.ForProvider.BundleS3, current.Spec.ForProvider.BundleS3)
	case v1alpha1.DataSourceCategoryGitBlame:
		return cmp.Equal(spec.Spec.ForProvider.GitBlame, current.Spec.ForProvider.GitBlame)
	case v1alpha1.DataSourceCategoryGitContent:
		return cmp.Equal(spec.Spec.ForProvider.GitContent, current.Spec.ForProvider.GitContent)
	case v1alpha1.DataSourceCategoryGitRego:
		return cmp.Equal(spec.Spec.ForProvider.GitRego, current.Spec.ForProvider.GitRego)
	case v1alpha1.DataSourceCategoryHTTP:
		return cmp.Equal(spec.Spec.ForProvider.HTTP, current.Spec.ForProvider.HTTP)
	case v1alpha1.DataSourceCategoryKubernetesResources:
		return cmp.Equal(spec.Spec.ForProvider.KubernetesResources, current.Spec.ForProvider.KubernetesResources)
	case v1alpha1.DataSourceCategoryLDAP:
		return cmp.Equal(spec.Spec.ForProvider.LDAP, current.Spec.ForProvider.LDAP)
	case v1alpha1.DataSourceCategoryPolicyLibrary:
		return cmp.Equal(spec.Spec.ForProvider.PolicyLibrary, current.Spec.ForProvider.PolicyLibrary)
	case v1alpha1.DataSourceCategoryRest:
		return cmp.Equal(spec.Spec.ForProvider.Rest, current.Spec.ForProvider.Rest)
	}

	return true
}

func (e *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.DataSource)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotDataSource)
	}

	req, err := generateDataSourceUpsertParams(ctx, cr)
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreateFailed)
	}
	if _, err := e.client.Datasources.UpsertDatasource(req); err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreateFailed)
	}

	return managed.ExternalCreation{}, nil
}

func (e *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.DataSource)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotDataSource)
	}

	req, err := generateDataSourceUpsertParams(ctx, cr)
	if err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, errCreateFailed)
	}
	if _, err := e.client.Datasources.UpsertDatasource(req); err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, errUpdateFailed)
	}

	return managed.ExternalUpdate{}, nil
}

func (e *external) Delete(ctx context.Context, mg resource.Managed) error {
	cr, ok := mg.(*v1alpha1.DataSource)
	if !ok {
		return errors.New(errNotDataSource)
	}

	req := &datasources.DeleteDatasourceParams{
		Context:    ctx,
		Datasource: meta.GetExternalName(cr),
	}

	if _, err := e.client.Datasources.DeleteDatasource(req); err != nil {
		return errors.Wrap(err, errDeleteFailed)
	}

	return nil
}

func lateInitialize(cr *v1alpha1.DataSource, resp *models.DatasourcesV1DatasourcesGetResponseResult) { // nolint:gocyclo
	current := generateDataSource(resp)

	cr.Spec.ForProvider.Description = styraclient.LateInitializeStringPtr(cr.Spec.ForProvider.Description, current.Spec.ForProvider.Description)
	cr.Spec.ForProvider.Enabled = styraclient.LateInitializeBoolPtr(cr.Spec.ForProvider.Enabled, current.Spec.ForProvider.Enabled)

	switch cr.Spec.ForProvider.Category {
	case v1alpha1.DataSourceCategoryAWSECR:
		if cr.Spec.ForProvider.AWSECR == nil {
			cr.Spec.ForProvider.AWSECR = current.Spec.ForProvider.AWSECR
		} else if current.Spec.ForProvider.AWSECR != nil {
			cr.Spec.ForProvider.AWSECR.RegistryID = styraclient.LateInitializeStringPtr(cr.Spec.ForProvider.AWSECR.RegistryID, &current.Spec.ForProvider.AWSECR.Credentials)
			cr.Spec.ForProvider.AWSECR.PollingInterval = styraclient.LateInitializeDuration(cr.Spec.ForProvider.AWSECR.PollingInterval, current.Spec.ForProvider.AWSECR.PollingInterval)
			cr.Spec.ForProvider.AWSECR.RateLimit = styraclient.LateInitializeQuantity(cr.Spec.ForProvider.AWSECR.RateLimit, current.Spec.ForProvider.AWSECR.RateLimit)
		}
	case v1alpha1.DataSourceCategoryBundleS3:
		if cr.Spec.ForProvider.BundleS3 == nil {
			cr.Spec.ForProvider.BundleS3 = current.Spec.ForProvider.BundleS3
		} else if current.Spec.ForProvider.BundleS3 != nil {
			cr.Spec.ForProvider.BundleS3.Endpoint = styraclient.LateInitializeStringPtr(cr.Spec.ForProvider.BundleS3.Endpoint, current.Spec.ForProvider.BundleS3.Endpoint)
			cr.Spec.ForProvider.BundleS3.PollingInterval = styraclient.LateInitializeDuration(cr.Spec.ForProvider.BundleS3.PollingInterval, current.Spec.ForProvider.BundleS3.PollingInterval)
		}
	case v1alpha1.DataSourceCategoryGitBlame:
		if cr.Spec.ForProvider.GitBlame == nil {
			cr.Spec.ForProvider.GitBlame = current.Spec.ForProvider.GitBlame
		} else if current.Spec.ForProvider.GitBlame != nil {
			cr.Spec.ForProvider.GitBlame.DatasourcesV1GitCommon = lateInitializeDatasourcesV1GitCommon(cr.Spec.ForProvider.GitBlame.DatasourcesV1GitCommon, current.Spec.ForProvider.GitBlame.DatasourcesV1GitCommon)
			cr.Spec.ForProvider.GitBlame.PathRegexp = styraclient.LateInitializeStringPtr(cr.Spec.ForProvider.GitBlame.PathRegexp, current.Spec.ForProvider.GitBlame.PathRegexp)
			cr.Spec.ForProvider.GitBlame.PollingInterval = styraclient.LateInitializeDuration(cr.Spec.ForProvider.GitBlame.PollingInterval, current.Spec.ForProvider.GitBlame.PollingInterval)
			cr.Spec.ForProvider.GitBlame.RateLimit = styraclient.LateInitializeQuantity(cr.Spec.ForProvider.GitBlame.RateLimit, current.Spec.ForProvider.GitBlame.RateLimit)
		}
	case v1alpha1.DataSourceCategoryGitContent:
		if cr.Spec.ForProvider.GitContent == nil {
			cr.Spec.ForProvider.GitContent = current.Spec.ForProvider.GitContent
		} else if current.Spec.ForProvider.GitContent != nil {
			cr.Spec.ForProvider.GitContent.DatasourcesV1GitCommon = lateInitializeDatasourcesV1GitCommon(cr.Spec.ForProvider.GitContent.DatasourcesV1GitCommon, current.Spec.ForProvider.GitContent.DatasourcesV1GitCommon)
			cr.Spec.ForProvider.GitContent.PollingInterval = styraclient.LateInitializeDuration(cr.Spec.ForProvider.GitContent.PollingInterval, current.Spec.ForProvider.GitContent.PollingInterval)
			cr.Spec.ForProvider.GitContent.RateLimit = styraclient.LateInitializeQuantity(cr.Spec.ForProvider.GitContent.RateLimit, current.Spec.ForProvider.GitContent.RateLimit)
		}
	case v1alpha1.DataSourceCategoryGitRego:
		if cr.Spec.ForProvider.GitRego == nil {
			cr.Spec.ForProvider.GitRego = current.Spec.ForProvider.GitRego
		} else if current.Spec.ForProvider.GitRego != nil {
			cr.Spec.ForProvider.GitRego.DatasourcesV1GitCommon = lateInitializeDatasourcesV1GitCommon(cr.Spec.ForProvider.GitRego.DatasourcesV1GitCommon, current.Spec.ForProvider.GitRego.DatasourcesV1GitCommon)
			cr.Spec.ForProvider.GitRego.Path = styraclient.LateInitializeStringPtr(cr.Spec.ForProvider.GitRego.Path, current.Spec.ForProvider.GitRego.Path)
			cr.Spec.ForProvider.GitRego.PollingInterval = styraclient.LateInitializeDuration(cr.Spec.ForProvider.GitRego.PollingInterval, current.Spec.ForProvider.GitRego.PollingInterval)
			cr.Spec.ForProvider.GitRego.RateLimit = styraclient.LateInitializeQuantity(cr.Spec.ForProvider.GitRego.RateLimit, current.Spec.ForProvider.GitRego.RateLimit)
		}
	case v1alpha1.DataSourceCategoryHTTP:
		if cr.Spec.ForProvider.HTTP == nil {
			cr.Spec.ForProvider.HTTP = current.Spec.ForProvider.HTTP
		} else if current.Spec.ForProvider.HTTP != nil {
			cr.Spec.ForProvider.HTTP.CaCertificate = styraclient.LateInitializeStringPtr(cr.Spec.ForProvider.HTTP.CaCertificate, current.Spec.ForProvider.HTTP.CaCertificate)
			if cr.Spec.ForProvider.HTTP.Headers == nil {
				cr.Spec.ForProvider.HTTP.Headers = current.Spec.ForProvider.HTTP.Headers
			}
			cr.Spec.ForProvider.HTTP.PolicyFilter = styraclient.LateInitializeStringPtr(cr.Spec.ForProvider.HTTP.PolicyFilter, current.Spec.ForProvider.HTTP.PolicyFilter)
			cr.Spec.ForProvider.HTTP.PolicyQuery = styraclient.LateInitializeStringPtr(cr.Spec.ForProvider.HTTP.PolicyQuery, current.Spec.ForProvider.HTTP.PolicyQuery)
			cr.Spec.ForProvider.HTTP.SkipTLSVerification = styraclient.LateInitializeBoolPtr(cr.Spec.ForProvider.HTTP.SkipTLSVerification, current.Spec.ForProvider.HTTP.SkipTLSVerification)
			cr.Spec.ForProvider.HTTP.PollingInterval = styraclient.LateInitializeDuration(cr.Spec.ForProvider.HTTP.PollingInterval, current.Spec.ForProvider.HTTP.PollingInterval)
		}
	case v1alpha1.DataSourceCategoryKubernetesResources:
		if cr.Spec.ForProvider.KubernetesResources == nil {
			cr.Spec.ForProvider.KubernetesResources = current.Spec.ForProvider.KubernetesResources
		} else if current.Spec.ForProvider.KubernetesResources != nil {
			if cr.Spec.ForProvider.KubernetesResources.Masks == nil {
				cr.Spec.ForProvider.KubernetesResources.Masks = current.Spec.ForProvider.KubernetesResources.Masks
			}
			if cr.Spec.ForProvider.KubernetesResources.Namespaces == nil {
				cr.Spec.ForProvider.KubernetesResources.Namespaces = current.Spec.ForProvider.KubernetesResources.Namespaces
			}
			if cr.Spec.ForProvider.KubernetesResources.Selectors == nil {
				cr.Spec.ForProvider.KubernetesResources.Selectors = current.Spec.ForProvider.KubernetesResources.Selectors
			}
			cr.Spec.ForProvider.KubernetesResources.PollingInterval = styraclient.LateInitializeDuration(cr.Spec.ForProvider.KubernetesResources.PollingInterval, current.Spec.ForProvider.KubernetesResources.PollingInterval)
			cr.Spec.ForProvider.KubernetesResources.RateLimit = styraclient.LateInitializeQuantity(cr.Spec.ForProvider.KubernetesResources.RateLimit, current.Spec.ForProvider.KubernetesResources.RateLimit)
		}
	case v1alpha1.DataSourceCategoryLDAP:
		if cr.Spec.ForProvider.LDAP == nil {
			cr.Spec.ForProvider.LDAP = current.Spec.ForProvider.LDAP
		} else if current.Spec.ForProvider.LDAP != nil {
			cr.Spec.ForProvider.LDAP.CaCertificate = styraclient.LateInitializeStringPtr(cr.Spec.ForProvider.LDAP.CaCertificate, current.Spec.ForProvider.LDAP.CaCertificate)
			cr.Spec.ForProvider.LDAP.PolicyFilter = styraclient.LateInitializeStringPtr(cr.Spec.ForProvider.LDAP.PolicyFilter, current.Spec.ForProvider.LDAP.PolicyFilter)
			cr.Spec.ForProvider.LDAP.PolicyQuery = styraclient.LateInitializeStringPtr(cr.Spec.ForProvider.LDAP.PolicyQuery, current.Spec.ForProvider.LDAP.PolicyQuery)
			cr.Spec.ForProvider.LDAP.SkipTLSVerification = styraclient.LateInitializeBoolPtr(cr.Spec.ForProvider.LDAP.SkipTLSVerification, current.Spec.ForProvider.LDAP.SkipTLSVerification)

			cr.Spec.ForProvider.LDAP.Credentials = styraclient.LateInitializeStringPtr(cr.Spec.ForProvider.LDAP.Credentials, current.Spec.ForProvider.LDAP.Credentials)
			if cr.Spec.ForProvider.LDAP.Search == nil {
				cr.Spec.ForProvider.LDAP.Search = current.Spec.ForProvider.LDAP.Search
			}
			if cr.Spec.ForProvider.LDAP.Urls == nil {
				cr.Spec.ForProvider.LDAP.Urls = current.Spec.ForProvider.LDAP.Urls
			}
			cr.Spec.ForProvider.LDAP.PollingInterval = styraclient.LateInitializeDuration(cr.Spec.ForProvider.LDAP.PollingInterval, current.Spec.ForProvider.LDAP.PollingInterval)
			cr.Spec.ForProvider.LDAP.RateLimit = styraclient.LateInitializeQuantity(cr.Spec.ForProvider.LDAP.RateLimit, current.Spec.ForProvider.LDAP.RateLimit)

		}
	case v1alpha1.DataSourceCategoryPolicyLibrary:
		if cr.Spec.ForProvider.PolicyLibrary == nil {
			cr.Spec.ForProvider.PolicyLibrary = current.Spec.ForProvider.PolicyLibrary
		} else if current.Spec.ForProvider.PolicyLibrary != nil {
			cr.Spec.ForProvider.PolicyLibrary.PollingInterval = styraclient.LateInitializeDuration(cr.Spec.ForProvider.PolicyLibrary.PollingInterval, current.Spec.ForProvider.PolicyLibrary.PollingInterval)
		}
	case v1alpha1.DataSourceCategoryRest:
		if cr.Spec.ForProvider.Rest == nil {
			cr.Spec.ForProvider.Rest = current.Spec.ForProvider.Rest
		} else if current.Spec.ForProvider.Rest != nil {
			cr.Spec.ForProvider.Rest.ContentType = styraclient.LateInitializeStringPtr(cr.Spec.ForProvider.Rest.ContentType, current.Spec.ForProvider.Rest.ContentType)
		}
	}
}

func lateInitializeDatasourcesV1GitCommon(in, from v1alpha1.DatasourcesV1GitCommon) v1alpha1.DatasourcesV1GitCommon {
	in.Credentials = styraclient.LateInitializeStringPtr(in.Credentials, from.Credentials)
	in.Reference = styraclient.LateInitializeStringPtr(in.Reference, from.Reference)
	in.Timeout = styraclient.LateInitializeDuration(in.Timeout, from.Timeout)
	if in.SSHCredentials == nil {
		in.SSHCredentials = from.SSHCredentials
	}
	return in
}

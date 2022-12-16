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

package system

import (
	"bytes"
	"context"
	"fmt"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/google/go-cmp/cmp"
	"github.com/iancoleman/strcase"

	"github.com/pkg/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"

	styra "github.com/mistermx/styra-go-client/pkg/client"
	"github.com/mistermx/styra-go-client/pkg/client/policies"
	"github.com/mistermx/styra-go-client/pkg/client/systems"
	"github.com/mistermx/styra-go-client/pkg/models"

	"github.com/crossplane-contrib/provider-styra/apis/system/v1alpha1"
	styraclient "github.com/crossplane-contrib/provider-styra/pkg/client"
	"github.com/crossplane-contrib/provider-styra/pkg/interface/controller"
)

const (
	errNotSystem                = "managed resource is not an system custom resource"
	errUpdateFailed             = "cannot update system"
	errCreateFailed             = "cannot create system"
	errDeleteFailed             = "cannot delete system"
	errDescribeFailed           = "cannot describe system"
	errGetConnectionDetails     = "cannot get connection details"
	errGetAsset                 = "cannot get asset"
	errIsUpToDateFailed         = "isUpToDate failed"
	errGetLabels                = "cannot get system labels"
	errGetLabelsInvalidResponse = "get system labels returned an unexpected response"
	errCompareLabels            = "cannot compare labels"
	errUpdateLabels             = "cannotUpdateLabels"
)

// SetupSystem adds a controller that reconciles Systems.
func SetupSystem(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.SystemGroupKind)

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o.ForControllerRuntime()).
		For(&v1alpha1.System{}).
		Complete(managed.NewReconciler(mgr,
			resource.ManagedKind(v1alpha1.SystemGroupVersionKind),
			managed.WithExternalConnecter(&connector{kube: mgr.GetClient(), newClientFn: styra.New}),
			managed.WithInitializers(managed.NewDefaultProviderConfig(mgr.GetClient())),
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
	_, ok := mg.(*v1alpha1.System)
	if !ok {
		return nil, errors.New(errNotSystem)
	}

	cfg, err := styraclient.GetConfig(ctx, c.kube, mg)
	if err != nil {
		return nil, err
	}

	client := c.newClientFn(cfg, strfmt.Default)

	return &external{client, c.kube}, nil
}

func (e *external) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.System)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotSystem)
	}

	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{}, nil
	}

	req := &systems.GetSystemParams{
		Context: ctx,
		System:  meta.GetExternalName(cr),
	}
	resp, err := e.client.Systems.GetSystem(req)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(resource.Ignore(isNotFound, err), errDescribeFailed)
	}

	currentSpec := cr.Spec.ForProvider.DeepCopy()
	generateSystem(resp.Payload.Result).Status.AtProvider.DeepCopyInto(&cr.Status.AtProvider)

	e.LateInitialize(cr, resp.Payload.Result)
	isUpToDate, err := e.isUpToDate(ctx, cr, resp)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errIsUpToDateFailed)
	}

	cr.Status.SetConditions(v1.Available())

	connectionDetails, err := e.GetConnectionDetails(ctx, cr)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errGetConnectionDetails)
	}

	return managed.ExternalObservation{
		ResourceExists:          true,
		ResourceUpToDate:        isUpToDate,
		ResourceLateInitialized: !cmp.Equal(&cr.Spec.ForProvider, currentSpec),
		ConnectionDetails:       connectionDetails,
	}, nil
}

func (e *external) isUpToDate(ctx context.Context, cr *v1alpha1.System, resp *systems.GetSystemOK) (bool, error) { // nolint:gocyclo
	if cr.ObjectMeta.Name != styraclient.StringValue(resp.Payload.Result.Name) {
		return false, nil
	}
	if cr.Spec.ForProvider.DeploymentParameters != nil && !isEqualSystemDeploymentParameters(cr.Spec.ForProvider.DeploymentParameters, resp.Payload.Result.DeploymentParameters) {
		return false, nil
	}
	if styraclient.StringValue(cr.Spec.ForProvider.Description) != resp.Payload.Result.Description {
		return false, nil
	}
	if styraclient.StringValue(cr.Spec.ForProvider.ExternalID) != resp.Payload.Result.ExternalID {
		return false, nil
	}
	if cr.Spec.ForProvider.ReadOnly != nil && !styraclient.IsEqualBool(cr.Spec.ForProvider.ReadOnly, resp.Payload.Result.ReadOnly) {
		return false, nil
	}
	if cr.Spec.ForProvider.Type != styraclient.StringValue(resp.Payload.Result.Type) {
		return false, nil
	}
	return e.areLabelsUpToDate(ctx, cr)
}

func (e *external) areLabelsUpToDate(ctx context.Context, cr *v1alpha1.System) (bool, error) {
	if !cr.Spec.ForProvider.HasLabels() {
		return true, nil
	}

	req := &policies.GetPolicyParams{
		Context: ctx,
		Policy:  fmt.Sprintf("metadata/%s/labels", meta.GetExternalName((cr))),
	}

	res, err := e.client.Policies.GetPolicy(req)
	if err != nil {
		return false, errors.Wrap(err, errGetLabels)
	}

	result, ok := res.Payload.Result.(map[string]interface{})
	if !ok {
		return false, errors.New(errGetLabelsInvalidResponse)
	}

	if _, exists := result["modules"]; !exists {
		return false, errors.New(errGetLabelsInvalidResponse)
	}

	modules, ok := result["modules"].(map[string]interface{})
	if !ok {
		return false, errors.New(errGetLabelsInvalidResponse)
	}

	if _, exists := modules["labels.rego"]; !exists {
		return false, errors.New(errGetLabelsInvalidResponse)
	}

	labelsModule, ok := modules["labels.rego"].(string)
	if !ok {
		return false, errors.New(errGetLabelsInvalidResponse)
	}

	labelsAreEqual, err := compareLabels(ctx, labelsModule, cr)
	if err != nil {
		return false, errors.Wrap(err, errCompareLabels)
	}

	return labelsAreEqual, nil
}

func (e *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.System)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotSystem)
	}

	req := &systems.CreateSystemParams{
		Context: ctx,
		Body:    generateSystemPostRequest(cr),
	}

	resp, err := e.client.Systems.CreateSystem(req)
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreateFailed)
	}

	meta.SetExternalName(cr, styraclient.StringValue(resp.Payload.Result.ID))

	// Do not create/update labels and connection details here because an error
	// will result in a recreation of the system.
	// This shall be handled in Update().

	return managed.ExternalCreation{
		ExternalNameAssigned: true,
	}, nil
}

func (e *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.System)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotSystem)
	}

	req := &systems.UpdateSystemParams{
		Context: ctx,
		System:  meta.GetExternalName(cr),
		Body:    generateSystemPutRequest(cr),
	}

	_, err := e.client.Systems.UpdateSystem(req)
	if err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, errUpdateFailed)
	}

	if err := e.updateLabels(ctx, cr); err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, errUpdateFailed)
	}

	return managed.ExternalUpdate{}, nil
}

func (e *external) Delete(ctx context.Context, mg resource.Managed) error {
	cr, ok := mg.(*v1alpha1.System)
	if !ok {
		return errors.New(errNotSystem)
	}

	req := &systems.DeleteSystemParams{
		Context: ctx,
		System:  meta.GetExternalName(cr),

		// Should we not delete everything?
		// Recursive: styraclient.String("false"),
	}

	_, err := e.client.Systems.DeleteSystem(req)
	if err != nil {
		return errors.Wrap(err, errDeleteFailed)
	}

	return nil
}

func (e *external) LateInitialize(cr *v1alpha1.System, resp *models.SystemsV1SystemConfig) {
	system := generateSystem(resp)

	cr.Spec.ForProvider.Description = styraclient.LateInitializeStringPtr(cr.Spec.ForProvider.Description, system.Spec.ForProvider.Description)
	cr.Spec.ForProvider.ExternalID = styraclient.LateInitializeStringPtr(cr.Spec.ForProvider.ExternalID, system.Spec.ForProvider.ExternalID)
	cr.Spec.ForProvider.ReadOnly = styraclient.LateInitializeBoolPtr(cr.Spec.ForProvider.ReadOnly, system.Spec.ForProvider.ReadOnly)
	cr.Spec.ForProvider.DeploymentParameters = lateInitializeDeploymentParameters(cr.Spec.ForProvider.DeploymentParameters, system.Spec.ForProvider.DeploymentParameters)
}

func (e *external) GetConnectionDetails(ctx context.Context, cr *v1alpha1.System) (managed.ConnectionDetails, error) {
	if !cr.Spec.ForProvider.HasAssets() {
		return nil, nil
	}

	details := managed.ConnectionDetails{}

	for _, assetType := range cr.Spec.ForProvider.GetAssetTypes() {
		resp, err := e.getAsset(ctx, cr, assetType)

		if err != nil {
			return nil, errors.Wrapf(err, "cannot get %s", assetType)
		}

		key := strcase.ToLowerCamel(assetType)
		details[key] = resp
	}

	return details, nil
}

func (e *external) getAsset(ctx context.Context, cr *v1alpha1.System, at string) ([]byte, error) {
	req := &systems.GetAssetParams{
		Context:   ctx,
		System:    meta.GetExternalName(cr),
		Assettype: at,
	}

	buffer := bytes.Buffer{}
	_, err := e.client.Systems.GetAsset(req, &buffer, styraclient.ReturnRawResponse)
	if err != nil {
		return nil, errors.Wrap(err, errGetAsset)
	}

	return buffer.Bytes(), nil
}

func (e *external) updateLabels(ctx context.Context, cr *v1alpha1.System) error {
	rego, err := generateRegoLabels(cr)
	if err != nil {
		return errors.Wrap(err, errUpdateLabels)
	}

	req := &policies.UpdatePolicyParams{
		Context: ctx,
		Policy:  fmt.Sprintf("metadata/%s/labels", meta.GetExternalName((cr))),
		Body: &models.PoliciesV1PoliciesPutRequest{
			Modules: map[string]string{
				"labels.rego": rego,
			},
		},
	}

	_, err = e.client.Policies.UpdatePolicy(req)
	return errors.Wrap(err, errUpdateLabels)
}

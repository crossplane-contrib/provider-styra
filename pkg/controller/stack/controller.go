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

package stack

import (
	"context"
	"fmt"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/google/go-cmp/cmp"

	"github.com/pkg/errors"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"

	styra "github.com/mistermx/styra-go-client/pkg/client"
	"github.com/mistermx/styra-go-client/pkg/client/policies"
	"github.com/mistermx/styra-go-client/pkg/client/stacks"
	"github.com/mistermx/styra-go-client/pkg/models"

	v1alpha1 "github.com/crossplane-contrib/provider-styra/apis/stack/v1alpha1"
	styraclient "github.com/crossplane-contrib/provider-styra/pkg/client"
)

const (
	errNotStack                    = "managed resource is not an Stack custom resource"
	errUpdateFailed                = "cannot update Stack custom resource"
	errCreateFailed                = "cannot create Stack"
	errDeleteFailed                = "cannot delete Stack"
	errDescribeFailed              = "cannot describe Stack"
	errIsUpToDateFailed            = "isUpToDate failed"
	errGetSelectors                = "cannot get system selectors"
	errGetSelectorsInvalidResponse = "get system selectors returned an unexpected response"
	errCompareSelectors            = "cannot compare selectors"
	errUpdateSelectors             = "cannotUpdateselectors"
)

// SetupStack adds a controller that reconciles Stacks.
func SetupStack(mgr ctrl.Manager, l logging.Logger, rl workqueue.RateLimiter) error {
	name := managed.ControllerName(v1alpha1.StackGroupKind)

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(controller.Options{
			RateLimiter: ratelimiter.NewDefaultManagedRateLimiter(rl),
		}).
		For(&v1alpha1.Stack{}).
		Complete(managed.NewReconciler(mgr,
			resource.ManagedKind(v1alpha1.StackGroupVersionKind),
			managed.WithExternalConnecter(&connector{kube: mgr.GetClient(), newClientFn: styra.New}),
			managed.WithInitializers(managed.NewDefaultProviderConfig(mgr.GetClient())),
			managed.WithReferenceResolver(managed.NewAPISimpleReferenceResolver(mgr.GetClient())),
			managed.WithLogger(l.WithValues("controller", name)),
			managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name)))))
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
	_, ok := mg.(*v1alpha1.Stack)
	if !ok {
		return nil, errors.New(errNotStack)
	}

	cfg, err := styraclient.GetConfig(ctx, c.kube, mg)
	if err != nil {
		return nil, err
	}

	client := c.newClientFn(cfg, strfmt.Default)
	return &external{client, c.kube}, nil
}

func (e *external) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.Stack)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotStack)
	}

	id := meta.GetExternalName(cr)
	if id == "" {
		return managed.ExternalObservation{
			ResourceExists:   false,
			ResourceUpToDate: false,
		}, nil
	}

	req := &stacks.GetStackParams{
		Context: ctx,
		Stack:   meta.GetExternalName(cr),
	}
	resp, reqErr := e.client.Stacks.GetStack(req)
	if reqErr != nil {
		return managed.ExternalObservation{ResourceExists: false}, errors.Wrap(resource.Ignore(IsNotFound, reqErr), errDescribeFailed)
	}

	currentSpec := cr.Spec.ForProvider.DeepCopy()
	e.LateInitialize(cr, resp)

	isUpToDate, err := e.isUpToDate(ctx, cr, resp)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errIsUpToDateFailed)
	}

	cr.Status.SetConditions(v1.Available())

	return managed.ExternalObservation{
		ResourceExists:          true,
		ResourceUpToDate:        isUpToDate,
		ResourceLateInitialized: !cmp.Equal(&cr.Spec.ForProvider, currentSpec),
	}, nil
}

func (e *external) isUpToDate(ctx context.Context, cr *v1alpha1.Stack, resp *stacks.GetStackOK) (bool, error) { // nolint:gocyclo
	if cr.Spec.ForProvider.Description != styraclient.StringValue(resp.Payload.Result.Description) {
		return false, nil
	}
	if cr.Spec.ForProvider.ReadOnly != styraclient.BoolValue(resp.Payload.Result.ReadOnly) {
		return false, nil
	}
	if cr.Spec.ForProvider.Type != styraclient.StringValue(resp.Payload.Result.Type) {
		return false, nil
	}
	if !isEqualSourceControlConfig(cr.Spec.ForProvider.SourceControl, resp.Payload.Result.SourceControl) {
		return false, nil
	}
	return e.areSelectorsUpToDate(ctx, cr)
}

func (e *external) areSelectorsUpToDate(ctx context.Context, cr *v1alpha1.Stack) (bool, error) {
	req := &policies.GetPolicyParams{
		Context: ctx,
		Policy:  fmt.Sprintf("stacks/%s/selectors", meta.GetExternalName((cr))),
	}

	res, err := e.client.Policies.GetPolicy(req)
	if err != nil {
		return false, errors.Wrap(err, errGetSelectors)
	}

	// Q&D solution to get raw rego. There may be a more elegant solution for this
	result, ok := res.Payload.Result.(map[string]interface{})
	if !ok {
		return false, errors.New(errGetSelectorsInvalidResponse)
	}

	if _, exists := result["modules"]; !exists {
		return false, errors.New(errGetSelectorsInvalidResponse)
	}

	modules, ok := result["modules"].(map[string]interface{})
	if !ok {
		return false, errors.New(errGetSelectorsInvalidResponse)
	}

	if _, exists := modules["selector.rego"]; !exists {
		return false, errors.New(errGetSelectorsInvalidResponse)
	}

	selectorsModule, ok := modules["selector.rego"].(string)
	if !ok {
		return false, errors.New(errGetSelectorsInvalidResponse)
	}

	selectorsAreEqual, err := compareSelectors(selectorsModule, cr)
	return selectorsAreEqual, errors.Wrap(err, errCompareSelectors)
}

func (e *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.Stack)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotStack)
	}

	req := &stacks.CreateStackParams{
		Context: ctx,
		Body:    generateStackPostRequest(cr),
	}

	resp, err := e.client.Stacks.CreateStack(req)
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreateFailed)
	}

	meta.SetExternalName(cr, styraclient.StringValue(resp.Payload.Result.ID))
	return managed.ExternalCreation{
		ExternalNameAssigned: true,
	}, nil
}

func (e *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.Stack)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotStack)
	}

	req := &stacks.UpdateStackParams{
		Context: ctx,
		Stack:   meta.GetExternalName(cr),
		Body:    generateStackPutRequest(cr),
	}

	if _, err := e.client.Stacks.UpdateStack(req); err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, errUpdateFailed)
	}

	if err := e.updateSelectors(ctx, cr); err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, errUpdateFailed)
	}

	return managed.ExternalUpdate{}, nil
}

func (e *external) Delete(ctx context.Context, mg resource.Managed) error {
	cr, ok := mg.(*v1alpha1.Stack)
	if !ok {
		return errors.New(errNotStack)
	}

	req := &stacks.DeleteStackParams{
		Context: ctx,
		Stack:   meta.GetExternalName(cr),

		// Should we not delete everything?
		// Recursive: styraclient.String("false"),
	}

	_, err := e.client.Stacks.DeleteStack(req)
	if err != nil {
		return errors.Wrap(err, errDeleteFailed)
	}

	return nil
}

func (e *external) LateInitialize(cr *v1alpha1.Stack, resp *stacks.GetStackOK) {
	stack := generateStack(resp.Payload.Result)
	cr.Spec.ForProvider.SourceControl = lateInitializeSourceControlConfig(cr.Spec.ForProvider.SourceControl, stack.Spec.ForProvider.SourceControl)
}

func (e *external) updateSelectors(ctx context.Context, cr *v1alpha1.Stack) error {
	rego, err := generateRegoSelectors(cr)
	if err != nil {
		return errors.Wrap(err, errUpdateSelectors)
	}

	req := &policies.UpdatePolicyParams{
		Context: ctx,
		Policy:  fmt.Sprintf("stacks/%s/selectors", meta.GetExternalName((cr))),
		Body: &models.V1PoliciesPutRequest{
			Modules: map[string]string{
				"selector.rego": rego,
			},
		},
	}

	_, err = e.client.Policies.UpdatePolicy(req)

	return errors.Wrap(err, errUpdateSelectors)
}

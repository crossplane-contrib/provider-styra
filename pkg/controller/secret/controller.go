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

package secret

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/google/go-cmp/cmp"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
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
	"github.com/mistermx/styra-go-client/pkg/client/secrets"
	"github.com/mistermx/styra-go-client/pkg/models"

	v1alpha1 "github.com/crossplane-contrib/provider-styra/apis/secret/v1alpha1"
	styraclient "github.com/crossplane-contrib/provider-styra/pkg/client"
)

const (
	errNotSecret              = "managed resource is not an secret custom resource"
	errUpdateFailed           = "cannot update secret"
	errCreateFailed           = "cannot create secret"
	errDeleteFailed           = "cannot delete secret"
	errIsUpToDateFailed       = "isUpToDate failed"
	errDescribeFailed         = "cannot describe secret"
	errUpsertSecretFailed     = "failed to upsert secret"
	errGetChecksumFailed      = "failed to get checksum"
	errUpdateChecksumFailed   = "failed to update checksum"
	errGetSecretValueFailed   = "failed to get secret value"
	errGetK8sSecretFailed     = "failed to get k8s secret"
	errDeleteK8sSecretFailed  = "failed to delete k8s secret"
	errFmtKeyNotFound         = "key %s is not found in referenced Kubernetes secret"
	errGenerateChecksumFailed = "failed to generate checksum"
	errNoChecksumRef          = "no checksum ref"

	checksumSecretDefaultKey = "checksum"
)

// SetupSecret adds a controller that reconciles Secrets.
func SetupSecret(mgr ctrl.Manager, l logging.Logger, rl workqueue.RateLimiter) error {
	name := managed.ControllerName(v1alpha1.SecretGroupKind)

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(controller.Options{
			RateLimiter: ratelimiter.NewDefaultManagedRateLimiter(rl),
		}).
		For(&v1alpha1.Secret{}).
		Complete(managed.NewReconciler(mgr,
			resource.ManagedKind(v1alpha1.SecretGroupVersionKind),
			managed.WithExternalConnecter(&connector{kube: mgr.GetClient(), newClientFn: styra.New}),
			managed.WithInitializers(managed.NewDefaultProviderConfig(mgr.GetClient()), managed.NewNameAsExternalName(mgr.GetClient())),
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
	kube   *resource.ClientApplicator
}

func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	_, ok := mg.(*v1alpha1.Secret)
	if !ok {
		return nil, errors.New(errNotSecret)
	}

	cfg, err := styraclient.GetConfig(ctx, c.kube, mg)
	if err != nil {
		return nil, err
	}

	// Workaround to make a request without Content-Type header for DELETE.
	// It actually sets the header value to "", however, this is treated as not-set by Styra API.
	// Can be removed once https://github.com/go-openapi/runtime/issues/231 is resolved.
	cfg.DefaultMediaType = ""
	cfg.Producers[""] = nil
	client := c.newClientFn(cfg, strfmt.Default)

	applicator := &resource.ClientApplicator{
		Client:     c.kube,
		Applicator: resource.NewAPIPatchingApplicator(c.kube),
	}

	return &external{client, applicator}, nil
}

func (e *external) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.Secret)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotSecret)
	}

	// LateInitialize has to happen earlier than usual to ensure spec.forProvider.checksumRef is set for create
	currentSpec := cr.Spec.ForProvider.DeepCopy()
	lateInitialize(cr)

	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{}, nil
	}

	resp, err := e.getSecret(ctx, cr)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(resource.Ignore(isNotFound, err), errDescribeFailed)
	}

	generateSecret(resp).Status.AtProvider.DeepCopyInto(&cr.Status.AtProvider)

	isUpToDate, err := e.isUpToDate(ctx, cr, resp)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errIsUpToDateFailed)
	}

	cr.Status.SetConditions(v1.Available())
	return managed.ExternalObservation{
		ResourceExists:          true,
		ResourceUpToDate:        isUpToDate,
		ResourceLateInitialized: !cmp.Equal(cr.Spec.ForProvider, *currentSpec),
	}, nil
}

func (e *external) isUpToDate(ctx context.Context, cr *v1alpha1.Secret, resp *secrets.GetSecretOK) (bool, error) {
	switch {
	case cr.Spec.ForProvider.Name != styraclient.StringValue(resp.Payload.Result.Name):
	case cr.Spec.ForProvider.Description != styraclient.StringValue(resp.Payload.Result.Description):
		return false, nil
	}

	// Styra does not provide an API to retrieve value of a secret.
	// In order to track changes between the in-cluster value with the upstream version, the controller generates the
	// checksum of the secret value and the lastModifiedAt timestamp from styra after an update and stores it next to the
	// referenced secret.
	// During Observe both checksum are calculated compared with each other. If they match, the secrets are considered equal.

	secretValue, err := e.getSecretValue(ctx, cr)
	if err != nil {
		// Don't report an error if the k8s secret was deleted before cr.
		if meta.WasDeleted(cr) {
			return false, nil
		}
		return false, errors.Wrap(err, errGetSecretValueFailed)
	}

	specCheckSum, err := e.getChecksum(ctx, cr)
	if err != nil {
		return false, errors.Wrap(err, errGetChecksumFailed)
	}

	currentChecksum, err := generateSecretChecksum(secretValue, time.Time(resp.Payload.Result.Metadata.LastModifiedAt))
	if err != nil {
		return false, errors.Wrap(err, errGenerateChecksumFailed)
	}

	return currentChecksum == specCheckSum, nil
}

func (e *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.Secret)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotSecret)
	}

	return managed.ExternalCreation{}, errors.Wrap(e.updateSecret(ctx, cr), errCreateFailed)
}

func (e *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.Secret)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotSecret)
	}

	return managed.ExternalUpdate{}, errors.Wrap(e.updateSecret(ctx, cr), errUpdateFailed)
}

func (e *external) Delete(ctx context.Context, mg resource.Managed) error {
	cr, ok := mg.(*v1alpha1.Secret)
	if !ok {
		return errors.New(errNotSecret)
	}

	if err := e.deleteChecksum(ctx, cr); err != nil {
		return errors.Wrap(err, errDeleteFailed)
	}

	req := &secrets.DeleteSecretParams{
		Context:  ctx,
		SecretID: meta.GetExternalName(cr),
	}

	// Workaround to make a request without Content-Type header for DELETE.
	// It actually sets the header value to "", however, this is treated as not-set by Styra API.
	// Can be removed once https://github.com/go-openapi/runtime/issues/231 is resolved.
	_, err := e.client.Secrets.DeleteSecret(req, styraclient.DropContentTypeHeader)

	return errors.Wrap(err, errDeleteFailed)
}

func (e *external) getSecret(ctx context.Context, cr *v1alpha1.Secret) (*secrets.GetSecretOK, error) {
	req := &secrets.GetSecretParams{
		Context:  ctx,
		SecretID: meta.GetExternalName(cr),
	}
	resp, err := e.client.Secrets.GetSecret(req)
	return resp, err
}

// updateSecret upserts the styra secret and updates the checksum.
func (e *external) updateSecret(ctx context.Context, cr *v1alpha1.Secret) error {
	secretValue, err := e.getSecretValue(ctx, cr)
	if err != nil {
		return errors.Wrap(err, errGetSecretValueFailed)
	}

	if err := e.upsertSecret(ctx, cr, secretValue); err != nil {
		return errors.Wrap(err, errUpsertSecretFailed)
	}

	resp, err := e.getSecret(ctx, cr)
	if err != nil {
		return errors.Wrap(err, errDescribeFailed)
	}

	currentChecksum, err := generateSecretChecksum(secretValue, time.Time(resp.Payload.Result.Metadata.LastModifiedAt))
	if err != nil {
		return errors.Wrap(err, errGenerateChecksumFailed)
	}

	return errors.Wrap(e.updateChecksum(ctx, cr, currentChecksum), errUpdateChecksumFailed)
}

// getSecretValue from the referenced K8s secret.
func (e *external) getSecretValue(ctx context.Context, cr *v1alpha1.Secret) (string, error) {
	nn := types.NamespacedName{
		Name:      cr.Spec.ForProvider.SecretRef.Name,
		Namespace: cr.Spec.ForProvider.SecretRef.Namespace,
	}
	sc := &corev1.Secret{}
	if err := e.kube.Get(ctx, nn, sc); err != nil {
		return "", errors.Wrap(err, errGetK8sSecretFailed)
	}

	if cr.Spec.ForProvider.SecretRef.Key != nil {
		val, ok := sc.Data[styraclient.StringValue(cr.Spec.ForProvider.SecretRef.Key)]
		if !ok {
			return "", errors.New(fmt.Sprintf(errFmtKeyNotFound, styraclient.StringValue(cr.Spec.ForProvider.SecretRef.Key)))
		}
		return string(val), nil
	}
	d := map[string]string{}
	for k, v := range sc.Data {
		d[k] = string(v)
	}
	payload, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(payload), nil
}

// getChecksum from the referenced checksum secret.
func (e *external) getChecksum(ctx context.Context, cr *v1alpha1.Secret) (string, error) {
	if cr.Spec.ForProvider.ChecksumSecretRef == nil || cr.Spec.ForProvider.ChecksumSecretRef.Key == nil {
		return "", errors.New(errNoChecksumRef)
	}

	nn := types.NamespacedName{
		Name:      cr.Spec.ForProvider.ChecksumSecretRef.Name,
		Namespace: cr.Spec.ForProvider.ChecksumSecretRef.Namespace,
	}
	sc := &corev1.Secret{}
	if err := e.kube.Get(ctx, nn, sc); err != nil {
		return "", errors.Wrap(client.IgnoreNotFound(err), errGetK8sSecretFailed)
	}

	checkSum := sc.Data[*cr.Spec.ForProvider.ChecksumSecretRef.Key]
	return string(checkSum), nil
}

// updateChecksum in referenced checksum secret with a new value.
func (e *external) updateChecksum(ctx context.Context, cr *v1alpha1.Secret, checksum string) error {
	if cr.Spec.ForProvider.ChecksumSecretRef == nil || cr.Spec.ForProvider.ChecksumSecretRef.Key == nil {
		return errors.New(errNoChecksumRef)
	}

	sc := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Spec.ForProvider.ChecksumSecretRef.Name,
			Namespace: cr.Spec.ForProvider.ChecksumSecretRef.Namespace,
		},
		Data: map[string][]byte{
			*cr.Spec.ForProvider.ChecksumSecretRef.Key: []byte(checksum),
		},
	}

	return e.kube.Apply(ctx, sc)
}

// deleteChecksum deletes the checksum secret.
func (e *external) deleteChecksum(ctx context.Context, cr *v1alpha1.Secret) error {
	if cr.Spec.ForProvider.ChecksumSecretRef == nil {
		return nil
	}

	sc := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Spec.ForProvider.ChecksumSecretRef.Name,
			Namespace: cr.Spec.ForProvider.ChecksumSecretRef.Namespace,
		},
	}
	return errors.Wrap(client.IgnoreNotFound(e.kube.Delete(ctx, sc)), errDeleteK8sSecretFailed)
}

// upsertSecret on styra
func (e *external) upsertSecret(ctx context.Context, cr *v1alpha1.Secret, secretValue string) error {
	req := &secrets.CreateUpdateSecretParams{
		Context:  ctx,
		SecretID: meta.GetExternalName(cr),
		Body: &models.SecretsV1SecretsPutRequest{
			Description: &cr.Spec.ForProvider.Description,
			Name:        &cr.Spec.ForProvider.Name,
			Secret:      &secretValue,
		},
	}

	_, err := e.client.Secrets.CreateUpdateSecret(req)
	return err
}

func lateInitialize(cr *v1alpha1.Secret) {
	if cr.Spec.ForProvider.ChecksumSecretRef == nil {
		cr.Spec.ForProvider.ChecksumSecretRef = &v1alpha1.SecretReference{
			Name:      generateChecksumSecretName(cr.ObjectMeta.Name),
			Namespace: cr.Spec.ForProvider.SecretRef.Namespace,
			Key:       styraclient.String(checksumSecretDefaultKey),
		}
	} else if cr.Spec.ForProvider.ChecksumSecretRef.Key == nil {
		cr.Spec.ForProvider.ChecksumSecretRef.Key = styraclient.String(checksumSecretDefaultKey)
	}
}

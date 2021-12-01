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
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/crossplane-runtime/pkg/test"

	styra "github.com/mistermx/styra-go-client/pkg/client"
	"github.com/mistermx/styra-go-client/pkg/client/secrets"
	"github.com/mistermx/styra-go-client/pkg/models"

	v1alpha1 "github.com/crossplane-contrib/provider-styra/apis/secret/v1alpha1"
	styraclient "github.com/crossplane-contrib/provider-styra/pkg/client"
	mockkube "github.com/crossplane-contrib/provider-styra/pkg/client/mock/kube"
	mockresource "github.com/crossplane-contrib/provider-styra/pkg/client/mock/resource"
	mocksecret "github.com/crossplane-contrib/provider-styra/pkg/client/mock/secrets"
)

var (
	errBoom = errors.New("boom")

	testSecretRefName      = "test-secret"
	testSecretRefNameSpace = "test-namespace"
	testSecretRefKey       = "test-key"
	testSecretValue        = "test-value"

	testSecretID    = "test-secret"
	testSecretName  = "test/secret/name"
	testDescription = "test-description"

	timeNow     = time.Now()
	timeBefore  = time.Now().Add(-time.Hour)
	metaTimeNow = metav1.NewTime(timeNow)
)

type args struct {
	kube  mockKubeFn
	styra mockStyraFn
	cr    *v1alpha1.Secret
}

type mockStyraFn func(t *testing.T) *styra.StyraAPI
type mockStyraModifier func(t *testing.T, s *styra.StyraAPI)

func checksum(val string, t time.Time) []byte {
	res, _ := generateSecretChecksum(val, t)
	return []byte(res)
}

func mockStyra(m ...mockStyraModifier) mockStyraFn {
	return func(t *testing.T) *styra.StyraAPI {
		s := &styra.StyraAPI{}
		for _, mod := range m {
			mod(t, s)
		}
		return s
	}
}

type mockSecretModifier func(*mocksecret.MockClientService)

func withMockSecret(mod mockSecretModifier) mockStyraModifier {
	return func(t *testing.T, s *styra.StyraAPI) {
		ctrl := gomock.NewController(t)
		mock := mocksecret.NewMockClientService(ctrl)
		mod(mock)
		s.Secrets = mock
	}
}

type mockKubeFn func(t *testing.T) *resource.ClientApplicator
type mockKubeModifier func(t *testing.T, c *resource.ClientApplicator)

func mockKube(m ...mockKubeModifier) mockKubeFn {
	return func(t *testing.T) *resource.ClientApplicator {
		c := &resource.ClientApplicator{}
		for _, mod := range m {
			mod(t, c)
		}
		return c
	}
}

type mockKubeClientModifier func(c *mockkube.MockClient)

func withMockKubeClient(mod mockKubeClientModifier) mockKubeModifier {
	return func(t *testing.T, c *resource.ClientApplicator) {
		ctrl := gomock.NewController(t)
		mock := mockkube.NewMockClient(ctrl)
		mod(mock)
		c.Client = mock
	}
}

type mockAppliatorModifier func(c *mockresource.MockApplicator)

func withMockApplicator(mod mockAppliatorModifier) mockKubeModifier {
	return func(t *testing.T, c *resource.ClientApplicator) {
		ctrl := gomock.NewController(t)
		mock := mockresource.NewMockApplicator(ctrl)
		mod(mock)
		c.Applicator = mock
	}
}

type SecretModifier func(*v1alpha1.Secret)

func withName(v string) SecretModifier {
	return func(s *v1alpha1.Secret) {
		s.ObjectMeta.Name = v
	}
}

func withExternalName(v string) SecretModifier {
	return func(s *v1alpha1.Secret) {
		meta.SetExternalName(s, v)
	}
}

func withConditions(c ...xpv1.Condition) SecretModifier {
	return func(r *v1alpha1.Secret) { r.Status.ConditionedStatus.Conditions = c }
}

func withSpec(p v1alpha1.SecretParameters) SecretModifier {
	return func(r *v1alpha1.Secret) { r.Spec.ForProvider = p }
}

func withStatus(s v1alpha1.SecretObservation) SecretModifier {
	return func(r *v1alpha1.Secret) { r.Status.AtProvider = s }
}

func Secret(m ...SecretModifier) *v1alpha1.Secret {
	cr := &v1alpha1.Secret{}
	for _, f := range m {
		f(cr)
	}
	return cr
}

func TestObserve(t *testing.T) {
	type want struct {
		cr     *v1alpha1.Secret
		result managed.ExternalObservation
		err    error
	}

	cases := map[string]struct {
		args
		want
	}{
		"SuccessfulAvailable": {
			args: args{
				kube: mockKube(
					withMockKubeClient(func(mc *mockkube.MockClient) {
						getSecret := mc.EXPECT().
							Get(
								context.Background(),
								types.NamespacedName{
									Name:      testSecretRefName,
									Namespace: testSecretRefNameSpace,
								},
								&corev1.Secret{},
							).
							Do(func(_ context.Context, _ types.NamespacedName, sc *corev1.Secret) {
								sc.Data = map[string][]byte{
									testSecretRefKey: []byte(testSecretValue),
								}
							})
						mc.EXPECT().
							Get(
								context.Background(),
								types.NamespacedName{
									Name:      generateChecksumSecretName(testSecretID),
									Namespace: testSecretRefNameSpace,
								},
								&corev1.Secret{},
							).
							Do(func(_ context.Context, _ types.NamespacedName, sc *corev1.Secret) {
								checksum, _ := generateSecretChecksum(testSecretValue, timeNow)
								sc.Data = map[string][]byte{
									checksumSecretDefaultKey: []byte(checksum),
								}
							}).
							After(getSecret)
					}),
				),
				styra: mockStyra(
					withMockSecret(func(mcs *mocksecret.MockClientService) {
						mcs.EXPECT().
							GetSecret(&secrets.GetSecretParams{
								SecretID: testSecretID,
								Context:  context.Background(),
							}).
							Return(&secrets.GetSecretOK{
								Payload: &models.SecretsV1SecretsGetResponse{
									Result: &models.SecretsV1Secret{
										Description: &testDescription,
										ID:          &testSecretID,
										Name:        &testSecretID,
										Metadata: &models.MetaV1ObjectMeta{
											LastModifiedAt: strfmt.DateTime(timeNow),
										},
									},
								},
							}, nil)
					}),
				),
				cr: Secret(
					withExternalName(testSecretID),
					withSpec(v1alpha1.SecretParameters{
						Description: testDescription,
						SecretRef: v1alpha1.SecretReference{
							Name:      testSecretRefName,
							Namespace: testSecretRefNameSpace,
							Key:       &testSecretRefKey,
						},
						ChecksumSecretRef: &v1alpha1.SecretReference{
							Name:      generateChecksumSecretName(testSecretID),
							Namespace: testSecretRefNameSpace,
							Key:       styraclient.String(checksumSecretDefaultKey),
						},
					}),
				),
			},
			want: want{
				cr: Secret(
					withExternalName(testSecretID),
					withSpec(v1alpha1.SecretParameters{
						Description: testDescription,
						SecretRef: v1alpha1.SecretReference{
							Name:      testSecretRefName,
							Namespace: testSecretRefNameSpace,
							Key:       &testSecretRefKey,
						},
						ChecksumSecretRef: &v1alpha1.SecretReference{
							Name:      generateChecksumSecretName(testSecretID),
							Namespace: testSecretRefNameSpace,
							Key:       styraclient.String(checksumSecretDefaultKey),
						},
					}),
					withStatus(v1alpha1.SecretObservation{
						LastModifiedAt: &metaTimeNow,
					}),
					withConditions(xpv1.Available()),
				),
				result: managed.ExternalObservation{
					ResourceExists:   true,
					ResourceUpToDate: true,
				},
			},
		},
		"ChecksumsNotEqual": {
			args: args{
				kube: mockKube(
					withMockKubeClient(func(mc *mockkube.MockClient) {
						getSecret := mc.EXPECT().
							Get(
								context.Background(),
								types.NamespacedName{
									Name:      testSecretRefName,
									Namespace: testSecretRefNameSpace,
								},
								&corev1.Secret{},
							).
							Do(func(_ context.Context, _ types.NamespacedName, sc *corev1.Secret) {
								sc.Data = map[string][]byte{
									testSecretRefKey: []byte(testSecretValue),
								}
							})
						mc.EXPECT().
							Get(
								context.Background(),
								types.NamespacedName{
									Name:      generateChecksumSecretName(testSecretID),
									Namespace: testSecretRefNameSpace,
								},
								&corev1.Secret{},
							).
							Do(func(_ context.Context, _ types.NamespacedName, sc *corev1.Secret) {
								checksum, _ := generateSecretChecksum(testSecretValue, timeBefore)
								sc.Data = map[string][]byte{
									checksumSecretDefaultKey: []byte(checksum),
								}
							}).
							After(getSecret)
					}),
				),
				styra: mockStyra(
					withMockSecret(func(mcs *mocksecret.MockClientService) {
						mcs.EXPECT().
							GetSecret(&secrets.GetSecretParams{
								SecretID: testSecretID,
								Context:  context.Background(),
							}).
							Return(&secrets.GetSecretOK{
								Payload: &models.SecretsV1SecretsGetResponse{
									Result: &models.SecretsV1Secret{
										Description: &testDescription,
										ID:          &testSecretID,
										Name:        &testSecretID,
										Metadata: &models.MetaV1ObjectMeta{
											LastModifiedAt: strfmt.DateTime(timeNow),
										},
									},
								},
							}, nil)
					}),
				),
				cr: Secret(
					withExternalName(testSecretID),
					withSpec(v1alpha1.SecretParameters{
						Description: testDescription,
						SecretRef: v1alpha1.SecretReference{
							Name:      testSecretRefName,
							Namespace: testSecretRefNameSpace,
							Key:       &testSecretRefKey,
						},
						ChecksumSecretRef: &v1alpha1.SecretReference{
							Name:      generateChecksumSecretName(testSecretID),
							Namespace: testSecretRefNameSpace,
							Key:       styraclient.String(checksumSecretDefaultKey),
						},
					}),
					withStatus(v1alpha1.SecretObservation{
						LastModifiedAt: &metaTimeNow,
					}),
				),
			},
			want: want{
				cr: Secret(
					withExternalName(testSecretID),
					withSpec(v1alpha1.SecretParameters{
						Description: testDescription,
						SecretRef: v1alpha1.SecretReference{
							Name:      testSecretRefName,
							Namespace: testSecretRefNameSpace,
							Key:       &testSecretRefKey,
						},
						ChecksumSecretRef: &v1alpha1.SecretReference{
							Name:      generateChecksumSecretName(testSecretID),
							Namespace: testSecretRefNameSpace,
							Key:       styraclient.String(checksumSecretDefaultKey),
						},
					}),
					withStatus(v1alpha1.SecretObservation{
						LastModifiedAt: &metaTimeNow,
					}),
					withConditions(xpv1.Available()),
				),
				result: managed.ExternalObservation{
					ResourceExists:   true,
					ResourceUpToDate: false,
				},
			},
		},
		"LateInitialize": {
			args: args{
				kube: mockKube(
					withMockKubeClient(func(mc *mockkube.MockClient) {
						getSecret := mc.EXPECT().
							Get(
								context.Background(),
								types.NamespacedName{
									Name:      testSecretRefName,
									Namespace: testSecretRefNameSpace,
								},
								&corev1.Secret{},
							).
							Do(func(_ context.Context, _ types.NamespacedName, sc *corev1.Secret) {
								sc.Data = map[string][]byte{
									testSecretRefKey: []byte(testSecretValue),
								}
							})
						mc.EXPECT().
							Get(
								context.Background(),
								types.NamespacedName{
									Name:      generateChecksumSecretName(testSecretID),
									Namespace: testSecretRefNameSpace,
								},
								&corev1.Secret{},
							).
							Do(func(_ context.Context, _ types.NamespacedName, sc *corev1.Secret) {
								checksum, _ := generateSecretChecksum(testSecretValue, timeNow)
								sc.Data = map[string][]byte{
									checksumSecretDefaultKey: []byte(checksum),
								}
							}).
							After(getSecret)
					}),
				),
				styra: mockStyra(
					withMockSecret(func(mcs *mocksecret.MockClientService) {
						mcs.EXPECT().
							GetSecret(&secrets.GetSecretParams{
								SecretID: testSecretID,
								Context:  context.Background(),
							}).
							Return(&secrets.GetSecretOK{
								Payload: &models.SecretsV1SecretsGetResponse{
									Result: &models.SecretsV1Secret{
										Description: &testDescription,
										ID:          &testSecretID,
										Name:        &testSecretID,
										Metadata: &models.MetaV1ObjectMeta{
											LastModifiedAt: strfmt.DateTime(timeNow),
										},
									},
								},
							}, nil)
					}),
				),
				cr: Secret(
					withName(testSecretID),
					withExternalName(testSecretID),
					withSpec(v1alpha1.SecretParameters{
						Description: testDescription,
						SecretRef: v1alpha1.SecretReference{
							Name:      testSecretRefName,
							Namespace: testSecretRefNameSpace,
							Key:       &testSecretRefKey,
						},
					}),
				),
			},
			want: want{
				cr: Secret(
					withName(testSecretID),
					withExternalName(testSecretID),
					withSpec(v1alpha1.SecretParameters{
						Description: testDescription,
						SecretRef: v1alpha1.SecretReference{
							Name:      testSecretRefName,
							Namespace: testSecretRefNameSpace,
							Key:       &testSecretRefKey,
						},
						ChecksumSecretRef: &v1alpha1.SecretReference{
							Name:      generateChecksumSecretName(testSecretID),
							Namespace: testSecretRefNameSpace,
							Key:       styraclient.String(checksumSecretDefaultKey),
						},
					}),
					withStatus(v1alpha1.SecretObservation{
						LastModifiedAt: &metaTimeNow,
					}),
					withConditions(xpv1.Available()),
				),
				result: managed.ExternalObservation{
					ResourceExists:          true,
					ResourceUpToDate:        true,
					ResourceLateInitialized: true,
				},
			},
		},
		"NotAvailable": {
			args: args{
				kube: mockKube(),
				styra: mockStyra(
					withMockSecret(func(mcs *mocksecret.MockClientService) {
						mcs.EXPECT().
							GetSecret(&secrets.GetSecretParams{
								SecretID: testSecretID,
								Context:  context.Background(),
							}).
							Return(nil, secrets.NewGetSecretNotFound())
					}),
				),
				cr: Secret(
					withExternalName(testSecretID),
					withSpec(v1alpha1.SecretParameters{
						Description: testDescription,
						SecretRef: v1alpha1.SecretReference{
							Name:      testSecretRefName,
							Namespace: testSecretRefNameSpace,
							Key:       &testSecretRefKey,
						},
						ChecksumSecretRef: &v1alpha1.SecretReference{
							Name:      generateChecksumSecretName(testSecretID),
							Namespace: testSecretRefNameSpace,
							Key:       styraclient.String(checksumSecretDefaultKey),
						},
					}),
					withStatus(v1alpha1.SecretObservation{
						LastModifiedAt: &metaTimeNow,
					}),
				),
			},
			want: want{
				cr: Secret(
					withExternalName(testSecretID),
					withSpec(v1alpha1.SecretParameters{
						Description: testDescription,
						SecretRef: v1alpha1.SecretReference{
							Name:      testSecretRefName,
							Namespace: testSecretRefNameSpace,
							Key:       &testSecretRefKey,
						},
						ChecksumSecretRef: &v1alpha1.SecretReference{
							Name:      generateChecksumSecretName(testSecretID),
							Namespace: testSecretRefNameSpace,
							Key:       styraclient.String(checksumSecretDefaultKey),
						},
					}),
					withStatus(v1alpha1.SecretObservation{
						LastModifiedAt: &metaTimeNow,
					}),
				),
				result: managed.ExternalObservation{},
			},
		},
		"DescribeFailed": {
			args: args{
				kube: mockKube(),
				styra: mockStyra(
					withMockSecret(func(mcs *mocksecret.MockClientService) {
						mcs.EXPECT().
							GetSecret(&secrets.GetSecretParams{
								SecretID: testSecretID,
								Context:  context.Background(),
							}).
							Return(nil, errBoom)
					}),
				),
				cr: Secret(
					withExternalName(testSecretID),
					withSpec(v1alpha1.SecretParameters{
						Description: testDescription,
						SecretRef: v1alpha1.SecretReference{
							Name:      testSecretRefName,
							Namespace: testSecretRefNameSpace,
							Key:       &testSecretRefKey,
						},
						ChecksumSecretRef: &v1alpha1.SecretReference{
							Name:      generateChecksumSecretName(testSecretID),
							Namespace: testSecretRefNameSpace,
							Key:       styraclient.String(checksumSecretDefaultKey),
						},
					}),
					withStatus(v1alpha1.SecretObservation{
						LastModifiedAt: &metaTimeNow,
					}),
				),
			},
			want: want{
				cr: Secret(
					withExternalName(testSecretID),
					withSpec(v1alpha1.SecretParameters{
						Description: testDescription,
						SecretRef: v1alpha1.SecretReference{
							Name:      testSecretRefName,
							Namespace: testSecretRefNameSpace,
							Key:       &testSecretRefKey,
						},
						ChecksumSecretRef: &v1alpha1.SecretReference{
							Name:      generateChecksumSecretName(testSecretID),
							Namespace: testSecretRefNameSpace,
							Key:       styraclient.String(checksumSecretDefaultKey),
						},
					}),
					withStatus(v1alpha1.SecretObservation{
						LastModifiedAt: &metaTimeNow,
					}),
				),
				result: managed.ExternalObservation{},
				err:    errors.Wrap(errBoom, errDescribeFailed),
			},
		},
		"GetSecretValueFailed": {
			args: args{
				kube: mockKube(
					withMockKubeClient(func(mc *mockkube.MockClient) {
						mc.EXPECT().
							Get(
								context.Background(),
								types.NamespacedName{
									Name:      testSecretRefName,
									Namespace: testSecretRefNameSpace,
								},
								&corev1.Secret{},
							).
							Return(errBoom)
					}),
				),
				styra: mockStyra(
					withMockSecret(func(mcs *mocksecret.MockClientService) {
						mcs.EXPECT().
							GetSecret(&secrets.GetSecretParams{
								SecretID: testSecretID,
								Context:  context.Background(),
							}).
							Return(&secrets.GetSecretOK{
								Payload: &models.SecretsV1SecretsGetResponse{
									Result: &models.SecretsV1Secret{
										Description: &testDescription,
										ID:          &testSecretID,
										Name:        &testSecretID,
										Metadata: &models.MetaV1ObjectMeta{
											LastModifiedAt: strfmt.DateTime(timeNow),
										},
									},
								},
							}, nil)
					}),
				),
				cr: Secret(
					withExternalName(testSecretID),
					withSpec(v1alpha1.SecretParameters{
						Description: testDescription,
						SecretRef: v1alpha1.SecretReference{
							Name:      testSecretRefName,
							Namespace: testSecretRefNameSpace,
							Key:       &testSecretRefKey,
						},
						ChecksumSecretRef: &v1alpha1.SecretReference{
							Name:      generateChecksumSecretName(testSecretID),
							Namespace: testSecretRefNameSpace,
							Key:       styraclient.String(checksumSecretDefaultKey),
						},
					}),
				),
			},
			want: want{
				cr: Secret(
					withExternalName(testSecretID),
					withSpec(v1alpha1.SecretParameters{
						Description: testDescription,
						SecretRef: v1alpha1.SecretReference{
							Name:      testSecretRefName,
							Namespace: testSecretRefNameSpace,
							Key:       &testSecretRefKey,
						},
						ChecksumSecretRef: &v1alpha1.SecretReference{
							Name:      generateChecksumSecretName(testSecretID),
							Namespace: testSecretRefNameSpace,
							Key:       styraclient.String(checksumSecretDefaultKey),
						},
					}),
					withStatus(v1alpha1.SecretObservation{
						LastModifiedAt: &metaTimeNow,
					}),
				),
				err: errors.Wrap(errors.Wrap(errors.Wrap(errBoom, errGetK8sSecretFailed), errGetSecretValueFailed), errIsUpToDateFailed),
			},
		},
		"GetChecksumFailed": {
			args: args{
				kube: mockKube(
					withMockKubeClient(func(mc *mockkube.MockClient) {
						getSecret := mc.EXPECT().
							Get(
								context.Background(),
								types.NamespacedName{
									Name:      testSecretRefName,
									Namespace: testSecretRefNameSpace,
								},
								&corev1.Secret{},
							).
							Do(func(_ context.Context, _ types.NamespacedName, sc *corev1.Secret) {
								sc.Data = map[string][]byte{
									testSecretRefKey: []byte(testSecretValue),
								}
							})
						mc.EXPECT().
							Get(
								context.Background(),
								types.NamespacedName{
									Name:      generateChecksumSecretName(testSecretID),
									Namespace: testSecretRefNameSpace,
								},
								&corev1.Secret{},
							).
							Return(errBoom).
							After(getSecret)
					}),
				),
				styra: mockStyra(
					withMockSecret(func(mcs *mocksecret.MockClientService) {
						mcs.EXPECT().
							GetSecret(&secrets.GetSecretParams{
								SecretID: testSecretID,
								Context:  context.Background(),
							}).
							Return(&secrets.GetSecretOK{
								Payload: &models.SecretsV1SecretsGetResponse{
									Result: &models.SecretsV1Secret{
										Description: &testDescription,
										ID:          &testSecretID,
										Name:        &testSecretID,
										Metadata: &models.MetaV1ObjectMeta{
											LastModifiedAt: strfmt.DateTime(timeNow),
										},
									},
								},
							}, nil)
					}),
				),
				cr: Secret(
					withExternalName(testSecretID),
					withSpec(v1alpha1.SecretParameters{
						Description: testDescription,
						SecretRef: v1alpha1.SecretReference{
							Name:      testSecretRefName,
							Namespace: testSecretRefNameSpace,
							Key:       &testSecretRefKey,
						},
						ChecksumSecretRef: &v1alpha1.SecretReference{
							Name:      generateChecksumSecretName(testSecretID),
							Namespace: testSecretRefNameSpace,
							Key:       styraclient.String(checksumSecretDefaultKey),
						},
					}),
				),
			},
			want: want{
				cr: Secret(
					withExternalName(testSecretID),
					withSpec(v1alpha1.SecretParameters{
						Description: testDescription,
						SecretRef: v1alpha1.SecretReference{
							Name:      testSecretRefName,
							Namespace: testSecretRefNameSpace,
							Key:       &testSecretRefKey,
						},
						ChecksumSecretRef: &v1alpha1.SecretReference{
							Name:      generateChecksumSecretName(testSecretID),
							Namespace: testSecretRefNameSpace,
							Key:       styraclient.String(checksumSecretDefaultKey),
						},
					}),
					withStatus(v1alpha1.SecretObservation{
						LastModifiedAt: &metaTimeNow,
					}),
				),
				err: errors.Wrap(errors.Wrap(errors.Wrap(errBoom, errGetK8sSecretFailed), errGetChecksumFailed), errIsUpToDateFailed),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := &external{client: tc.args.styra(t), kube: tc.args.kube(t)}
			o, err := e.Observe(context.Background(), tc.args.cr)

			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.cr, tc.args.cr, test.EquateConditions()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.result, o); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	type want struct {
		cr     *v1alpha1.Secret
		result managed.ExternalCreation
		err    error
	}

	cases := map[string]struct {
		args
		want
	}{
		"Successful": {
			args: args{
				kube: mockKube(
					withMockKubeClient(func(mc *mockkube.MockClient) {
						mc.EXPECT().
							Get(
								context.Background(),
								types.NamespacedName{
									Name:      testSecretRefName,
									Namespace: testSecretRefNameSpace,
								},
								&corev1.Secret{},
							).
							Do(func(_ context.Context, _ types.NamespacedName, sc *corev1.Secret) {
								sc.Data = map[string][]byte{
									testSecretRefKey: []byte(testSecretValue),
								}
							})
					}),
					withMockApplicator(func(ma *mockresource.MockApplicator) {
						ma.EXPECT().
							Apply(
								context.Background(),
								&corev1.Secret{
									ObjectMeta: metav1.ObjectMeta{
										Name:      generateChecksumSecretName(testSecretID),
										Namespace: testSecretRefNameSpace,
									},
									Data: map[string][]byte{
										checksumSecretDefaultKey: checksum(testSecretValue, timeNow),
									},
								},
							).
							Return(nil)
					}),
				),
				styra: mockStyra(
					withMockSecret(func(mcs *mocksecret.MockClientService) {
						upsertSecret := mcs.EXPECT().
							CreateUpdateSecret(&secrets.CreateUpdateSecretParams{
								Body: &models.SecretsV1SecretsPutRequest{
									Description: &testDescription,
									Name:        &testSecretName,
									Secret:      &testSecretValue,
								},
								SecretID: testSecretID,
								Context:  context.Background(),
							}).
							Return(&secrets.CreateUpdateSecretOK{
								Payload: &models.SecretsV1SecretsPutResponse{},
							}, nil)
						mcs.EXPECT().
							GetSecret(&secrets.GetSecretParams{
								SecretID: testSecretID,
								Context:  context.Background(),
							}).
							Return(&secrets.GetSecretOK{
								Payload: &models.SecretsV1SecretsGetResponse{
									Result: &models.SecretsV1Secret{
										Description: &testDescription,
										ID:          &testSecretID,
										Name:        &testSecretID,
										Metadata: &models.MetaV1ObjectMeta{
											LastModifiedAt: strfmt.DateTime(timeNow),
										},
									},
								},
							}, nil).
							After(upsertSecret)
					}),
				),
				cr: Secret(
					withExternalName(testSecretID),
					withSpec(v1alpha1.SecretParameters{
						Name:        testSecretName,
						Description: testDescription,
						SecretRef: v1alpha1.SecretReference{
							Name:      testSecretRefName,
							Namespace: testSecretRefNameSpace,
							Key:       &testSecretRefKey,
						},
						ChecksumSecretRef: &v1alpha1.SecretReference{
							Name:      generateChecksumSecretName(testSecretID),
							Namespace: testSecretRefNameSpace,
							Key:       styraclient.String(checksumSecretDefaultKey),
						},
					}),
				),
			},
			want: want{
				cr: Secret(
					withExternalName(testSecretID),
					withSpec(v1alpha1.SecretParameters{
						Name:        testSecretName,
						Description: testDescription,
						SecretRef: v1alpha1.SecretReference{
							Name:      testSecretRefName,
							Namespace: testSecretRefNameSpace,
							Key:       &testSecretRefKey,
						},
						ChecksumSecretRef: &v1alpha1.SecretReference{
							Name:      generateChecksumSecretName(testSecretID),
							Namespace: testSecretRefNameSpace,
							Key:       styraclient.String(checksumSecretDefaultKey),
						},
					}),
				),
				result: managed.ExternalCreation{},
			},
		},
		"CreateSecretFailed": {
			args: args{
				kube: mockKube(
					withMockKubeClient(func(mc *mockkube.MockClient) {
						mc.EXPECT().
							Get(
								context.Background(),
								types.NamespacedName{
									Name:      testSecretRefName,
									Namespace: testSecretRefNameSpace,
								},
								&corev1.Secret{},
							).
							Do(func(_ context.Context, _ types.NamespacedName, sc *corev1.Secret) {
								sc.Data = map[string][]byte{
									testSecretRefKey: []byte(testSecretValue),
								}
							})
					}),
				),
				styra: mockStyra(
					withMockSecret(func(mcs *mocksecret.MockClientService) {
						mcs.EXPECT().
							CreateUpdateSecret(&secrets.CreateUpdateSecretParams{
								Body: &models.SecretsV1SecretsPutRequest{
									Description: &testDescription,
									Name:        &testSecretName,
									Secret:      &testSecretValue,
								},
								SecretID: testSecretID,
								Context:  context.Background(),
							}).
							Return(nil, errBoom)
					}),
				),
				cr: Secret(
					withExternalName(testSecretID),
					withSpec(v1alpha1.SecretParameters{
						Name:        testSecretName,
						Description: testDescription,
						SecretRef: v1alpha1.SecretReference{
							Name:      testSecretRefName,
							Namespace: testSecretRefNameSpace,
							Key:       &testSecretRefKey,
						},
						ChecksumSecretRef: &v1alpha1.SecretReference{
							Name:      generateChecksumSecretName(testSecretID),
							Namespace: testSecretRefNameSpace,
							Key:       styraclient.String(checksumSecretDefaultKey),
						},
					}),
				),
			},
			want: want{
				cr: Secret(
					withExternalName(testSecretID),
					withSpec(v1alpha1.SecretParameters{
						Name:        testSecretName,
						Description: testDescription,
						SecretRef: v1alpha1.SecretReference{
							Name:      testSecretRefName,
							Namespace: testSecretRefNameSpace,
							Key:       &testSecretRefKey,
						},
						ChecksumSecretRef: &v1alpha1.SecretReference{
							Name:      generateChecksumSecretName(testSecretID),
							Namespace: testSecretRefNameSpace,
							Key:       styraclient.String(checksumSecretDefaultKey),
						},
					}),
				),
				err: errors.Wrap(errors.Wrap(errBoom, errUpsertSecretFailed), errCreateFailed),
			},
		},
		"GetSecretValueFailed": {
			args: args{
				kube: mockKube(
					withMockKubeClient(func(mc *mockkube.MockClient) {
						mc.EXPECT().
							Get(
								context.Background(),
								types.NamespacedName{
									Name:      testSecretRefName,
									Namespace: testSecretRefNameSpace,
								},
								&corev1.Secret{},
							).
							Return(errBoom)
					}),
				),
				styra: mockStyra(),
				cr: Secret(
					withExternalName(testSecretID),
					withSpec(v1alpha1.SecretParameters{
						Description: testDescription,
						SecretRef: v1alpha1.SecretReference{
							Name:      testSecretRefName,
							Namespace: testSecretRefNameSpace,
							Key:       &testSecretRefKey,
						},
						ChecksumSecretRef: &v1alpha1.SecretReference{
							Name:      generateChecksumSecretName(testSecretID),
							Namespace: testSecretRefNameSpace,
							Key:       styraclient.String(checksumSecretDefaultKey),
						},
					}),
				),
			},
			want: want{
				cr: Secret(
					withExternalName(testSecretID),
					withSpec(v1alpha1.SecretParameters{
						Description: testDescription,
						SecretRef: v1alpha1.SecretReference{
							Name:      testSecretRefName,
							Namespace: testSecretRefNameSpace,
							Key:       &testSecretRefKey,
						},
						ChecksumSecretRef: &v1alpha1.SecretReference{
							Name:      generateChecksumSecretName(testSecretID),
							Namespace: testSecretRefNameSpace,
							Key:       styraclient.String(checksumSecretDefaultKey),
						},
					}),
				),
				result: managed.ExternalCreation{},
				err:    errors.Wrap(errors.Wrap(errors.Wrap(errBoom, errGetK8sSecretFailed), errGetSecretValueFailed), errCreateFailed),
			},
		},
		"UpdateChecksumFailed": {
			args: args{
				kube: mockKube(
					withMockKubeClient(func(mc *mockkube.MockClient) {
						mc.EXPECT().
							Get(
								context.Background(),
								types.NamespacedName{
									Name:      testSecretRefName,
									Namespace: testSecretRefNameSpace,
								},
								&corev1.Secret{},
							).
							Do(func(_ context.Context, _ types.NamespacedName, sc *corev1.Secret) {
								sc.Data = map[string][]byte{
									testSecretRefKey: []byte(testSecretValue),
								}
							})
					}),
					withMockApplicator(func(ma *mockresource.MockApplicator) {
						ma.EXPECT().
							Apply(
								context.Background(),
								&corev1.Secret{
									ObjectMeta: metav1.ObjectMeta{
										Name:      generateChecksumSecretName(testSecretID),
										Namespace: testSecretRefNameSpace,
									},
									Data: map[string][]byte{
										checksumSecretDefaultKey: checksum(testSecretValue, timeNow),
									},
								},
							).
							Return(errBoom)
					}),
				),
				styra: mockStyra(
					withMockSecret(func(mcs *mocksecret.MockClientService) {
						upsertSecret := mcs.EXPECT().
							CreateUpdateSecret(&secrets.CreateUpdateSecretParams{
								Body: &models.SecretsV1SecretsPutRequest{
									Description: &testDescription,
									Name:        &testSecretName,
									Secret:      &testSecretValue,
								},
								SecretID: testSecretID,
								Context:  context.Background(),
							}).
							Return(&secrets.CreateUpdateSecretOK{
								Payload: &models.SecretsV1SecretsPutResponse{},
							}, nil)
						mcs.EXPECT().
							GetSecret(&secrets.GetSecretParams{
								SecretID: testSecretID,
								Context:  context.Background(),
							}).
							Return(&secrets.GetSecretOK{
								Payload: &models.SecretsV1SecretsGetResponse{
									Result: &models.SecretsV1Secret{
										Description: &testDescription,
										ID:          &testSecretID,
										Name:        &testSecretID,
										Metadata: &models.MetaV1ObjectMeta{
											LastModifiedAt: strfmt.DateTime(timeNow),
										},
									},
								},
							}, nil).
							After(upsertSecret)
					}),
				),
				cr: Secret(
					withExternalName(testSecretID),
					withSpec(v1alpha1.SecretParameters{
						Name:        testSecretName,
						Description: testDescription,
						SecretRef: v1alpha1.SecretReference{
							Name:      testSecretRefName,
							Namespace: testSecretRefNameSpace,
							Key:       &testSecretRefKey,
						},
						ChecksumSecretRef: &v1alpha1.SecretReference{
							Name:      generateChecksumSecretName(testSecretID),
							Namespace: testSecretRefNameSpace,
							Key:       styraclient.String(checksumSecretDefaultKey),
						},
					}),
				),
			},
			want: want{
				cr: Secret(
					withExternalName(testSecretID),
					withSpec(v1alpha1.SecretParameters{
						Name:        testSecretName,
						Description: testDescription,
						SecretRef: v1alpha1.SecretReference{
							Name:      testSecretRefName,
							Namespace: testSecretRefNameSpace,
							Key:       &testSecretRefKey,
						},
						ChecksumSecretRef: &v1alpha1.SecretReference{
							Name:      generateChecksumSecretName(testSecretID),
							Namespace: testSecretRefNameSpace,
							Key:       styraclient.String(checksumSecretDefaultKey),
						},
					}),
				),
				result: managed.ExternalCreation{},
				err:    errors.Wrap(errors.Wrap(errBoom, errUpdateChecksumFailed), errCreateFailed),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := &external{client: tc.args.styra(t), kube: tc.args.kube(t)}
			o, err := e.Create(context.Background(), tc.args.cr)

			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.cr, tc.args.cr, test.EquateConditions()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.result, o); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	type want struct {
		cr     *v1alpha1.Secret
		result managed.ExternalUpdate
		err    error
	}

	cases := map[string]struct {
		args
		want
	}{
		"Successful": {
			args: args{
				kube: mockKube(
					withMockKubeClient(func(mc *mockkube.MockClient) {
						mc.EXPECT().
							Get(
								context.Background(),
								types.NamespacedName{
									Name:      testSecretRefName,
									Namespace: testSecretRefNameSpace,
								},
								&corev1.Secret{},
							).
							Do(func(_ context.Context, _ types.NamespacedName, sc *corev1.Secret) {
								sc.Data = map[string][]byte{
									testSecretRefKey: []byte(testSecretValue),
								}
							})
					}),
					withMockApplicator(func(ma *mockresource.MockApplicator) {
						ma.EXPECT().
							Apply(
								context.Background(),
								&corev1.Secret{
									ObjectMeta: metav1.ObjectMeta{
										Name:      generateChecksumSecretName(testSecretID),
										Namespace: testSecretRefNameSpace,
									},
									Data: map[string][]byte{
										checksumSecretDefaultKey: checksum(testSecretValue, timeNow),
									},
								},
							).
							Return(nil)
					}),
				),
				styra: mockStyra(
					withMockSecret(func(mcs *mocksecret.MockClientService) {
						upsertSecret := mcs.EXPECT().
							CreateUpdateSecret(&secrets.CreateUpdateSecretParams{
								Body: &models.SecretsV1SecretsPutRequest{
									Description: &testDescription,
									Name:        &testSecretName,
									Secret:      &testSecretValue,
								},
								SecretID: testSecretID,
								Context:  context.Background(),
							}).
							Return(&secrets.CreateUpdateSecretOK{
								Payload: &models.SecretsV1SecretsPutResponse{},
							}, nil)
						mcs.EXPECT().
							GetSecret(&secrets.GetSecretParams{
								SecretID: testSecretID,
								Context:  context.Background(),
							}).
							Return(&secrets.GetSecretOK{
								Payload: &models.SecretsV1SecretsGetResponse{
									Result: &models.SecretsV1Secret{
										Description: &testDescription,
										ID:          &testSecretID,
										Name:        &testSecretID,
										Metadata: &models.MetaV1ObjectMeta{
											LastModifiedAt: strfmt.DateTime(timeNow),
										},
									},
								},
							}, nil).
							After(upsertSecret)
					}),
				),
				cr: Secret(
					withExternalName(testSecretID),
					withSpec(v1alpha1.SecretParameters{
						Name:        testSecretName,
						Description: testDescription,
						SecretRef: v1alpha1.SecretReference{
							Name:      testSecretRefName,
							Namespace: testSecretRefNameSpace,
							Key:       &testSecretRefKey,
						},
						ChecksumSecretRef: &v1alpha1.SecretReference{
							Name:      generateChecksumSecretName(testSecretID),
							Namespace: testSecretRefNameSpace,
							Key:       styraclient.String(checksumSecretDefaultKey),
						},
					}),
				),
			},
			want: want{
				cr: Secret(
					withExternalName(testSecretID),
					withSpec(v1alpha1.SecretParameters{
						Name:        testSecretName,
						Description: testDescription,
						SecretRef: v1alpha1.SecretReference{
							Name:      testSecretRefName,
							Namespace: testSecretRefNameSpace,
							Key:       &testSecretRefKey,
						},
						ChecksumSecretRef: &v1alpha1.SecretReference{
							Name:      generateChecksumSecretName(testSecretID),
							Namespace: testSecretRefNameSpace,
							Key:       styraclient.String(checksumSecretDefaultKey),
						},
					}),
				),
				result: managed.ExternalUpdate{},
			},
		},
		"CreateSecretFailed": {
			args: args{
				kube: mockKube(
					withMockKubeClient(func(mc *mockkube.MockClient) {
						mc.EXPECT().
							Get(
								context.Background(),
								types.NamespacedName{
									Name:      testSecretRefName,
									Namespace: testSecretRefNameSpace,
								},
								&corev1.Secret{},
							).
							Do(func(_ context.Context, _ types.NamespacedName, sc *corev1.Secret) {
								sc.Data = map[string][]byte{
									testSecretRefKey: []byte(testSecretValue),
								}
							})
					}),
				),
				styra: mockStyra(
					withMockSecret(func(mcs *mocksecret.MockClientService) {
						mcs.EXPECT().
							CreateUpdateSecret(&secrets.CreateUpdateSecretParams{
								Body: &models.SecretsV1SecretsPutRequest{
									Description: &testDescription,
									Name:        &testSecretName,
									Secret:      &testSecretValue,
								},
								SecretID: testSecretID,
								Context:  context.Background(),
							}).
							Return(nil, errBoom)
					}),
				),
				cr: Secret(
					withExternalName(testSecretID),
					withSpec(v1alpha1.SecretParameters{
						Name:        testSecretName,
						Description: testDescription,
						SecretRef: v1alpha1.SecretReference{
							Name:      testSecretRefName,
							Namespace: testSecretRefNameSpace,
							Key:       &testSecretRefKey,
						},
						ChecksumSecretRef: &v1alpha1.SecretReference{
							Name:      generateChecksumSecretName(testSecretID),
							Namespace: testSecretRefNameSpace,
							Key:       styraclient.String(checksumSecretDefaultKey),
						},
					}),
				),
			},
			want: want{
				cr: Secret(
					withExternalName(testSecretID),
					withSpec(v1alpha1.SecretParameters{
						Name:        testSecretName,
						Description: testDescription,
						SecretRef: v1alpha1.SecretReference{
							Name:      testSecretRefName,
							Namespace: testSecretRefNameSpace,
							Key:       &testSecretRefKey,
						},
						ChecksumSecretRef: &v1alpha1.SecretReference{
							Name:      generateChecksumSecretName(testSecretID),
							Namespace: testSecretRefNameSpace,
							Key:       styraclient.String(checksumSecretDefaultKey),
						},
					}),
				),
				err: errors.Wrap(errors.Wrap(errBoom, errUpsertSecretFailed), errUpdateFailed),
			},
		},
		"GetSecretValueFailed": {
			args: args{
				kube: mockKube(
					withMockKubeClient(func(mc *mockkube.MockClient) {
						mc.EXPECT().
							Get(
								context.Background(),
								types.NamespacedName{
									Name:      testSecretRefName,
									Namespace: testSecretRefNameSpace,
								},
								&corev1.Secret{},
							).
							Return(errBoom)
					}),
				),
				styra: mockStyra(),
				cr: Secret(
					withExternalName(testSecretID),
					withSpec(v1alpha1.SecretParameters{
						Name:        testSecretName,
						Description: testDescription,
						SecretRef: v1alpha1.SecretReference{
							Name:      testSecretRefName,
							Namespace: testSecretRefNameSpace,
							Key:       &testSecretRefKey,
						},
						ChecksumSecretRef: &v1alpha1.SecretReference{
							Name:      generateChecksumSecretName(testSecretID),
							Namespace: testSecretRefNameSpace,
							Key:       styraclient.String(checksumSecretDefaultKey),
						},
					}),
				),
			},
			want: want{
				cr: Secret(
					withExternalName(testSecretID),
					withSpec(v1alpha1.SecretParameters{
						Name:        testSecretName,
						Description: testDescription,
						SecretRef: v1alpha1.SecretReference{
							Name:      testSecretRefName,
							Namespace: testSecretRefNameSpace,
							Key:       &testSecretRefKey,
						},
						ChecksumSecretRef: &v1alpha1.SecretReference{
							Name:      generateChecksumSecretName(testSecretID),
							Namespace: testSecretRefNameSpace,
							Key:       styraclient.String(checksumSecretDefaultKey),
						},
					}),
				),
				result: managed.ExternalUpdate{},
				err:    errors.Wrap(errors.Wrap(errors.Wrap(errBoom, errGetK8sSecretFailed), errGetSecretValueFailed), errUpdateFailed),
			},
		},
		"UpdateChecksumFailed": {
			args: args{
				kube: mockKube(
					withMockKubeClient(func(mc *mockkube.MockClient) {
						mc.EXPECT().
							Get(
								context.Background(),
								types.NamespacedName{
									Name:      testSecretRefName,
									Namespace: testSecretRefNameSpace,
								},
								&corev1.Secret{},
							).
							Do(func(_ context.Context, _ types.NamespacedName, sc *corev1.Secret) {
								sc.Data = map[string][]byte{
									testSecretRefKey: []byte(testSecretValue),
								}
							})
					}),
					withMockApplicator(func(ma *mockresource.MockApplicator) {
						ma.EXPECT().
							Apply(
								context.Background(),
								&corev1.Secret{
									ObjectMeta: metav1.ObjectMeta{
										Name:      generateChecksumSecretName(testSecretID),
										Namespace: testSecretRefNameSpace,
									},
									Data: map[string][]byte{
										checksumSecretDefaultKey: checksum(testSecretValue, timeNow),
									},
								},
							).
							Return(errBoom)
					}),
				),
				styra: mockStyra(
					withMockSecret(func(mcs *mocksecret.MockClientService) {
						upsertSecret := mcs.EXPECT().
							CreateUpdateSecret(&secrets.CreateUpdateSecretParams{
								Body: &models.SecretsV1SecretsPutRequest{
									Description: &testDescription,
									Name:        &testSecretName,
									Secret:      &testSecretValue,
								},
								SecretID: testSecretID,
								Context:  context.Background(),
							}).
							Return(&secrets.CreateUpdateSecretOK{
								Payload: &models.SecretsV1SecretsPutResponse{},
							}, nil)
						mcs.EXPECT().
							GetSecret(&secrets.GetSecretParams{
								SecretID: testSecretID,
								Context:  context.Background(),
							}).
							Return(&secrets.GetSecretOK{
								Payload: &models.SecretsV1SecretsGetResponse{
									Result: &models.SecretsV1Secret{
										Description: &testDescription,
										ID:          &testSecretID,
										Name:        &testSecretID,
										Metadata: &models.MetaV1ObjectMeta{
											LastModifiedAt: strfmt.DateTime(timeNow),
										},
									},
								},
							}, nil).
							After(upsertSecret)
					}),
				),
				cr: Secret(
					withExternalName(testSecretID),
					withSpec(v1alpha1.SecretParameters{
						Name:        testSecretName,
						Description: testDescription,
						SecretRef: v1alpha1.SecretReference{
							Name:      testSecretRefName,
							Namespace: testSecretRefNameSpace,
							Key:       &testSecretRefKey,
						},
						ChecksumSecretRef: &v1alpha1.SecretReference{
							Name:      generateChecksumSecretName(testSecretID),
							Namespace: testSecretRefNameSpace,
							Key:       styraclient.String(checksumSecretDefaultKey),
						},
					}),
				),
			},
			want: want{
				cr: Secret(
					withExternalName(testSecretID),
					withSpec(v1alpha1.SecretParameters{
						Name:        testSecretName,
						Description: testDescription,
						SecretRef: v1alpha1.SecretReference{
							Name:      testSecretRefName,
							Namespace: testSecretRefNameSpace,
							Key:       &testSecretRefKey,
						},
						ChecksumSecretRef: &v1alpha1.SecretReference{
							Name:      generateChecksumSecretName(testSecretID),
							Namespace: testSecretRefNameSpace,
							Key:       styraclient.String(checksumSecretDefaultKey),
						},
					}),
				),
				result: managed.ExternalUpdate{},
				err:    errors.Wrap(errors.Wrap(errBoom, errUpdateChecksumFailed), errUpdateFailed),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := &external{client: tc.args.styra(t), kube: tc.args.kube(t)}
			u, err := e.Update(context.Background(), tc.args.cr)

			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.cr, tc.args.cr, test.EquateConditions()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.result, u); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type want struct {
		cr  *v1alpha1.Secret
		err error
	}

	cases := map[string]struct {
		args
		want
	}{
		"Successful": {
			args: args{
				kube: mockKube(
					withMockKubeClient(func(mc *mockkube.MockClient) {
						mc.EXPECT().Delete(
							context.Background(),
							&corev1.Secret{
								ObjectMeta: metav1.ObjectMeta{
									Name:      generateChecksumSecretName(testSecretID),
									Namespace: testSecretRefNameSpace,
								},
							})
					}),
				),
				styra: mockStyra(
					withMockSecret(func(mcs *mocksecret.MockClientService) {
						mcs.EXPECT().
							DeleteSecret(
								&secrets.DeleteSecretParams{
									SecretID: testSecretID,
									Context:  context.Background(),
								},
								gomock.Any(),
							).
							Return(&secrets.DeleteSecretOK{}, nil)
					}),
				),
				cr: Secret(
					withExternalName(testSecretID),
					withSpec(v1alpha1.SecretParameters{
						ChecksumSecretRef: &v1alpha1.SecretReference{
							Name:      generateChecksumSecretName(testSecretID),
							Namespace: testSecretRefNameSpace,
							Key:       styraclient.String(checksumSecretDefaultKey),
						},
					}),
				),
			},
			want: want{
				cr: Secret(
					withExternalName(testSecretID),
					withSpec(v1alpha1.SecretParameters{
						ChecksumSecretRef: &v1alpha1.SecretReference{
							Name:      generateChecksumSecretName(testSecretID),
							Namespace: testSecretRefNameSpace,
							Key:       styraclient.String(checksumSecretDefaultKey),
						},
					}),
				),
			},
		},
		"DeleteSecretFailed": {
			args: args{
				kube: mockKube(
					withMockKubeClient(func(mc *mockkube.MockClient) {
						mc.EXPECT().Delete(
							context.Background(),
							&corev1.Secret{
								ObjectMeta: metav1.ObjectMeta{
									Name:      generateChecksumSecretName(testSecretID),
									Namespace: testSecretRefNameSpace,
								},
							})
					}),
				),
				styra: mockStyra(
					withMockSecret(func(mcs *mocksecret.MockClientService) {
						mcs.EXPECT().
							DeleteSecret(
								&secrets.DeleteSecretParams{
									SecretID: testSecretID,
									Context:  context.Background(),
								},
								gomock.Any(),
							).
							Return(nil, errBoom)
					}),
				),
				cr: Secret(
					withExternalName(testSecretID),
					withSpec(v1alpha1.SecretParameters{
						ChecksumSecretRef: &v1alpha1.SecretReference{
							Name:      generateChecksumSecretName(testSecretID),
							Namespace: testSecretRefNameSpace,
							Key:       styraclient.String(checksumSecretDefaultKey),
						},
					}),
				),
			},
			want: want{
				cr: Secret(
					withExternalName(testSecretID),
					withSpec(v1alpha1.SecretParameters{
						ChecksumSecretRef: &v1alpha1.SecretReference{
							Name:      generateChecksumSecretName(testSecretID),
							Namespace: testSecretRefNameSpace,
							Key:       styraclient.String(checksumSecretDefaultKey),
						},
					}),
				),
				err: errors.Wrap(errBoom, errDeleteFailed),
			},
		},
		"DeleteChecksumFailed": {
			args: args{
				kube: mockKube(
					withMockKubeClient(func(mc *mockkube.MockClient) {
						mc.EXPECT().
							Delete(
								context.Background(),
								&corev1.Secret{
									ObjectMeta: metav1.ObjectMeta{
										Name:      generateChecksumSecretName(testSecretID),
										Namespace: testSecretRefNameSpace,
									},
								}).
							Return(errBoom)
					}),
				),
				styra: mockStyra(),
				cr: Secret(
					withExternalName(testSecretID),
					withSpec(v1alpha1.SecretParameters{
						ChecksumSecretRef: &v1alpha1.SecretReference{
							Name:      generateChecksumSecretName(testSecretID),
							Namespace: testSecretRefNameSpace,
							Key:       styraclient.String(checksumSecretDefaultKey),
						},
					}),
				),
			},
			want: want{
				cr: Secret(
					withExternalName(testSecretID),
					withSpec(v1alpha1.SecretParameters{
						ChecksumSecretRef: &v1alpha1.SecretReference{
							Name:      generateChecksumSecretName(testSecretID),
							Namespace: testSecretRefNameSpace,
							Key:       styraclient.String(checksumSecretDefaultKey),
						},
					}),
				),
				err: errors.Wrap(errors.Wrap(errBoom, errDeleteK8sSecretFailed), errDeleteFailed),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := &external{client: tc.args.styra(t), kube: tc.args.kube(t)}
			err := e.Delete(context.Background(), tc.args.cr)

			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.cr, tc.args.cr, test.EquateConditions()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
		})
	}
}

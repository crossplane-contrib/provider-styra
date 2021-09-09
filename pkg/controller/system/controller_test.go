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
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/test"

	styra "github.com/mistermx/styra-go-client/pkg/client"
	"github.com/mistermx/styra-go-client/pkg/client/policies"
	"github.com/mistermx/styra-go-client/pkg/client/systems"
	"github.com/mistermx/styra-go-client/pkg/models"

	v1alpha1 "github.com/crossplane-contrib/provider-styra/apis/system/v1alpha1"
	styraclient "github.com/crossplane-contrib/provider-styra/pkg/client"
	mockpolicies "github.com/crossplane-contrib/provider-styra/pkg/client/mock/policies"
	mocksystem "github.com/crossplane-contrib/provider-styra/pkg/client/mock/systems"
)

var (
	errBoom = errors.New("boom")

	testSystemID    = "testsystem"
	testSystemName  = "testname"
	testType        = "kubernetes:test"
	testExternalID  = "test-external-ID"
	testDescription = "test-description"
	testAsset       = "test-asset"
	testLabelKey    = "foo"
	testLabelValue  = "bar"
	testLabelsRego  = `
package metadata.testsystem.labels

labels := {
    "foo": "bar",
    "system-type": "kubernetes:test",
}
`
)

type args struct {
	styra styra.StyraAPI
	cr    *v1alpha1.System
}

type mockSystemModifier func(*mocksystem.MockClientService)

func withMockSystem(t *testing.T, mod mockSystemModifier) *mocksystem.MockClientService {
	ctrl := gomock.NewController(t)
	mock := mocksystem.NewMockClientService(ctrl)
	mod(mock)
	return mock
}

type mockPolicyModifier func(*mockpolicies.MockClientService)

func withMockPolicies(t *testing.T, mod mockPolicyModifier) *mockpolicies.MockClientService {
	ctrl := gomock.NewController(t)
	mock := mockpolicies.NewMockClientService(ctrl)
	mod(mock)
	return mock
}

type SystemModifier func(*v1alpha1.System)

func withName(v string) SystemModifier {
	return func(s *v1alpha1.System) {
		s.ObjectMeta.Name = v
	}
}

func withExternalName(v string) SystemModifier {
	return func(s *v1alpha1.System) {
		meta.SetExternalName(s, v)
	}
}

func withConditions(c ...xpv1.Condition) SystemModifier {
	return func(r *v1alpha1.System) { r.Status.ConditionedStatus.Conditions = c }
}

func withSpec(p v1alpha1.SystemParameters) SystemModifier {
	return func(r *v1alpha1.System) { r.Spec.ForProvider = p }
}

func System(m ...SystemModifier) *v1alpha1.System {
	cr := &v1alpha1.System{}
	for _, f := range m {
		f(cr)
	}
	return cr
}

func TestObserve(t *testing.T) {
	type want struct {
		cr     *v1alpha1.System
		result managed.ExternalObservation
		err    error
	}

	cases := map[string]struct {
		args
		want
	}{
		"SuccessfulAvailable": {
			args: args{
				styra: styra.StyraAPI{
					Systems: withMockSystem(t, func(mcs *mocksystem.MockClientService) {
						mcs.EXPECT().
							GetSystem(&systems.GetSystemParams{
								System:  testSystemID,
								Context: context.Background(),
							}).
							Return(&systems.GetSystemOK{
								Payload: &models.V1SystemsGetResponse{
									Result: &models.V1SystemConfig{
										Description:          testDescription,
										DeploymentParameters: &models.V1SystemDeploymentParameters{},
										ReadOnly:             styraclient.Bool(true),
										Type:                 &testType,
										ExternalID:           testExternalID,
									},
								},
							}, nil)
						mcs.EXPECT().
							GetAsset(&systems.GetAssetParams{
								Assettype: helmValuesAssetType,
								System:    testSystemID,
								Context:   context.Background(),
							}, &bytes.Buffer{}, gomock.Any()).
							DoAndReturn(func(params *systems.GetAssetParams, writer io.Writer, _ ...systems.ClientOption) (*systems.GetAssetOK, error) {
								writer.Write([]byte(testAsset))
								return nil, nil
							})
					}),
					Policies: withMockPolicies(t, func(mcs *mockpolicies.MockClientService) {
						mcs.EXPECT().
							GetPolicy(&policies.GetPolicyParams{
								Policy:  fmt.Sprintf("metadata/%s/labels", testSystemID),
								Context: context.Background(),
							}).
							Return(&policies.GetPolicyOK{
								Payload: &models.V1PolicyGetResponse{
									Result: map[string]interface{}{
										"modules": map[string]interface{}{
											"labels.rego": testLabelsRego,
										},
									},
								},
							}, nil)
					}),
				},
				cr: System(
					withExternalName(testSystemID),
					withSpec(v1alpha1.SystemParameters{
						CustomSystemParameters: v1alpha1.CustomSystemParameters{
							Labels: map[string]string{
								testLabelKey: testLabelValue,
							},
						},
						Description: &testDescription,
						DeploymentParameters: &v1alpha1.V1SystemDeploymentParameters{
							HTTPProxy:                styraclient.String(""),
							HTTPSProxy:               styraclient.String(""),
							KubernetesVersion:        styraclient.String(""),
							Namespace:                styraclient.String(""),
							NoProxy:                  styraclient.String(""),
							TimeoutSeconds:           styraclient.Int32(0),
							TrustedContainerRegistry: styraclient.String(""),
						},
						ReadOnly:   styraclient.Bool(true),
						Type:       testType,
						ExternalID: &testExternalID,
					}),
				),
			},
			want: want{
				cr: System(
					withSpec(v1alpha1.SystemParameters{
						CustomSystemParameters: v1alpha1.CustomSystemParameters{
							Labels: map[string]string{
								testLabelKey: testLabelValue,
							},
						},
						Description: &testDescription,
						DeploymentParameters: &v1alpha1.V1SystemDeploymentParameters{
							HTTPProxy:                styraclient.String(""),
							HTTPSProxy:               styraclient.String(""),
							KubernetesVersion:        styraclient.String(""),
							Namespace:                styraclient.String(""),
							NoProxy:                  styraclient.String(""),
							TimeoutSeconds:           styraclient.Int32(0),
							TrustedContainerRegistry: styraclient.String(""),
						},
						ReadOnly:   styraclient.Bool(true),
						Type:       testType,
						ExternalID: &testExternalID,
					}),
					withExternalName(testSystemID),
					withConditions(xpv1.Available()),
				),
				result: managed.ExternalObservation{
					ResourceExists:   true,
					ResourceUpToDate: true,
					ConnectionDetails: managed.ConnectionDetails{
						helmValuesConnectionDetailsKey: []byte(testAsset),
					},
				},
			},
		},
		"SuccessfulLateInitialize": {
			args: args{
				styra: styra.StyraAPI{
					Systems: withMockSystem(t, func(mcs *mocksystem.MockClientService) {
						mcs.EXPECT().
							GetSystem(&systems.GetSystemParams{
								System:  testSystemID,
								Context: context.Background(),
							}).
							Return(&systems.GetSystemOK{
								Payload: &models.V1SystemsGetResponse{
									Result: &models.V1SystemConfig{
										Description:          testDescription,
										DeploymentParameters: &models.V1SystemDeploymentParameters{},
										ReadOnly:             styraclient.Bool(true),
										Type:                 &testType,
										ExternalID:           testExternalID,
									},
								},
							}, nil)
						mcs.EXPECT().
							GetAsset(&systems.GetAssetParams{
								Assettype: helmValuesAssetType,
								System:    testSystemID,
								Context:   context.Background(),
							}, &bytes.Buffer{}, gomock.Any()).
							DoAndReturn(func(params *systems.GetAssetParams, writer io.Writer, _ ...systems.ClientOption) (*systems.GetAssetOK, error) {
								writer.Write([]byte(testAsset))
								return nil, nil
							})
					}),
					Policies: withMockPolicies(t, func(mcs *mockpolicies.MockClientService) {
						mcs.EXPECT().
							GetPolicy(&policies.GetPolicyParams{
								Policy:  fmt.Sprintf("metadata/%s/labels", testSystemID),
								Context: context.Background(),
							}).
							Return(&policies.GetPolicyOK{
								Payload: &models.V1PolicyGetResponse{
									Result: map[string]interface{}{
										"modules": map[string]interface{}{
											"labels.rego": testLabelsRego,
										},
									},
								},
							}, nil)
					}),
				},
				cr: System(
					withExternalName(testSystemID),
					withSpec(v1alpha1.SystemParameters{
						CustomSystemParameters: v1alpha1.CustomSystemParameters{
							Labels: map[string]string{
								testLabelKey: testLabelValue,
							},
						},
						Type: testType,
					}),
				),
			},
			want: want{
				cr: System(
					withSpec(v1alpha1.SystemParameters{
						CustomSystemParameters: v1alpha1.CustomSystemParameters{
							Labels: map[string]string{
								testLabelKey: testLabelValue,
							},
						},
						Description: &testDescription,
						DeploymentParameters: &v1alpha1.V1SystemDeploymentParameters{
							HTTPProxy:                styraclient.String(""),
							HTTPSProxy:               styraclient.String(""),
							KubernetesVersion:        styraclient.String(""),
							Namespace:                styraclient.String(""),
							NoProxy:                  styraclient.String(""),
							TimeoutSeconds:           styraclient.Int32(0),
							TrustedContainerRegistry: styraclient.String(""),
						},
						ReadOnly:   styraclient.Bool(true),
						Type:       testType,
						ExternalID: &testExternalID,
					}),
					withExternalName(testSystemID),
					withConditions(xpv1.Available()),
				),
				result: managed.ExternalObservation{
					ResourceExists:          true,
					ResourceUpToDate:        true,
					ResourceLateInitialized: true,
					ConnectionDetails: managed.ConnectionDetails{
						helmValuesConnectionDetailsKey: []byte(testAsset),
					},
				},
			},
		},
		"LabelsNotUpToDate": {
			args: args{
				styra: styra.StyraAPI{
					Systems: withMockSystem(t, func(mcs *mocksystem.MockClientService) {
						mcs.EXPECT().
							GetSystem(&systems.GetSystemParams{
								System:  testSystemID,
								Context: context.Background(),
							}).
							Return(&systems.GetSystemOK{
								Payload: &models.V1SystemsGetResponse{
									Result: &models.V1SystemConfig{
										Description:          testDescription,
										DeploymentParameters: &models.V1SystemDeploymentParameters{},
										ReadOnly:             styraclient.Bool(true),
										Type:                 &testType,
										ExternalID:           testExternalID,
									},
								},
							}, nil)
						mcs.EXPECT().
							GetAsset(&systems.GetAssetParams{
								Assettype: helmValuesAssetType,
								System:    testSystemID,
								Context:   context.Background(),
							}, &bytes.Buffer{}, gomock.Any()).
							DoAndReturn(func(params *systems.GetAssetParams, writer io.Writer, _ ...systems.ClientOption) (*systems.GetAssetOK, error) {
								writer.Write([]byte(testAsset))
								return nil, nil
							})
					}),
					Policies: withMockPolicies(t, func(mcs *mockpolicies.MockClientService) {
						mcs.EXPECT().
							GetPolicy(&policies.GetPolicyParams{
								Policy:  fmt.Sprintf("metadata/%s/labels", testSystemID),
								Context: context.Background(),
							}).
							Return(&policies.GetPolicyOK{
								Payload: &models.V1PolicyGetResponse{
									Result: map[string]interface{}{
										"modules": map[string]interface{}{
											"labels.rego": testLabelsRego,
										},
									},
								},
							}, nil)
					}),
				},
				cr: System(
					withExternalName(testSystemID),
					withSpec(v1alpha1.SystemParameters{
						CustomSystemParameters: v1alpha1.CustomSystemParameters{
							Labels: map[string]string{},
						},
						Description: &testDescription,
						DeploymentParameters: &v1alpha1.V1SystemDeploymentParameters{
							HTTPProxy:                styraclient.String(""),
							HTTPSProxy:               styraclient.String(""),
							KubernetesVersion:        styraclient.String(""),
							Namespace:                styraclient.String(""),
							NoProxy:                  styraclient.String(""),
							TimeoutSeconds:           styraclient.Int32(0),
							TrustedContainerRegistry: styraclient.String(""),
						},
						ReadOnly:   styraclient.Bool(true),
						Type:       testType,
						ExternalID: &testExternalID,
					}),
				),
			},
			want: want{
				cr: System(
					withSpec(v1alpha1.SystemParameters{
						CustomSystemParameters: v1alpha1.CustomSystemParameters{
							Labels: map[string]string{},
						},
						Description: &testDescription,
						DeploymentParameters: &v1alpha1.V1SystemDeploymentParameters{
							HTTPProxy:                styraclient.String(""),
							HTTPSProxy:               styraclient.String(""),
							KubernetesVersion:        styraclient.String(""),
							Namespace:                styraclient.String(""),
							NoProxy:                  styraclient.String(""),
							TimeoutSeconds:           styraclient.Int32(0),
							TrustedContainerRegistry: styraclient.String(""),
						},
						ReadOnly:   styraclient.Bool(true),
						Type:       testType,
						ExternalID: &testExternalID,
					}),
					withExternalName(testSystemID),
					withConditions(xpv1.Available()),
				),
				result: managed.ExternalObservation{
					ResourceExists:   true,
					ResourceUpToDate: false,
					ConnectionDetails: managed.ConnectionDetails{
						helmValuesConnectionDetailsKey: []byte(testAsset),
					},
				},
			},
		},
		"GetSystemFailed": {
			args: args{
				styra: styra.StyraAPI{
					Systems: withMockSystem(t, func(mcs *mocksystem.MockClientService) {
						mcs.EXPECT().
							GetSystem(&systems.GetSystemParams{
								System:  testSystemID,
								Context: context.Background(),
							}).
							Return(nil, errBoom)
					}),
				},
				cr: System(
					withExternalName(testSystemID),
				),
			},
			want: want{
				cr: System(
					withExternalName(testSystemID),
				),
				err: errors.Wrap(errBoom, errDescribeFailed),
			},
		},
		"GetLabelsFailed": {
			args: args{
				styra: styra.StyraAPI{
					Systems: withMockSystem(t, func(mcs *mocksystem.MockClientService) {
						mcs.EXPECT().
							GetSystem(&systems.GetSystemParams{
								System:  testSystemID,
								Context: context.Background(),
							}).
							Return(&systems.GetSystemOK{
								Payload: &models.V1SystemsGetResponse{
									Result: &models.V1SystemConfig{
										Description:          testDescription,
										DeploymentParameters: &models.V1SystemDeploymentParameters{},
										ReadOnly:             styraclient.Bool(true),
										Type:                 &testType,
										ExternalID:           testExternalID,
									},
								},
							}, nil)
					}),
					Policies: withMockPolicies(t, func(mcs *mockpolicies.MockClientService) {
						mcs.EXPECT().
							GetPolicy(&policies.GetPolicyParams{
								Policy:  fmt.Sprintf("metadata/%s/labels", testSystemID),
								Context: context.Background(),
							}).
							Return(nil, errBoom)
					}),
				},
				cr: System(
					withExternalName(testSystemID),
					withSpec(v1alpha1.SystemParameters{
						Type: testType,
					}),
				),
			},
			want: want{
				cr: System(
					withExternalName(testSystemID),
					withSpec(v1alpha1.SystemParameters{
						Description: &testDescription,
						DeploymentParameters: &v1alpha1.V1SystemDeploymentParameters{
							HTTPProxy:                styraclient.String(""),
							HTTPSProxy:               styraclient.String(""),
							KubernetesVersion:        styraclient.String(""),
							Namespace:                styraclient.String(""),
							NoProxy:                  styraclient.String(""),
							TimeoutSeconds:           styraclient.Int32(0),
							TrustedContainerRegistry: styraclient.String(""),
						},
						ReadOnly:   styraclient.Bool(true),
						Type:       testType,
						ExternalID: &testExternalID,
					}),
				),
				err: errors.Wrap(errors.Wrap(errBoom, errGetLabels), errIsUpToDateFailed),
			},
		},
		"GetAssetFailed": {
			args: args{
				styra: styra.StyraAPI{
					Systems: withMockSystem(t, func(mcs *mocksystem.MockClientService) {
						mcs.EXPECT().
							GetSystem(&systems.GetSystemParams{
								System:  testSystemID,
								Context: context.Background(),
							}).
							Return(&systems.GetSystemOK{
								Payload: &models.V1SystemsGetResponse{
									Result: &models.V1SystemConfig{
										Description:          testDescription,
										DeploymentParameters: &models.V1SystemDeploymentParameters{},
										ReadOnly:             styraclient.Bool(true),
										Type:                 &testType,
										ExternalID:           testExternalID,
									},
								},
							}, nil)
						mcs.EXPECT().
							GetAsset(&systems.GetAssetParams{
								Assettype: helmValuesAssetType,
								System:    testSystemID,
								Context:   context.Background(),
							}, &bytes.Buffer{}, gomock.Any()).
							Return(nil, errBoom)
					}),
					Policies: withMockPolicies(t, func(mcs *mockpolicies.MockClientService) {
						mcs.EXPECT().
							GetPolicy(&policies.GetPolicyParams{
								Policy:  fmt.Sprintf("metadata/%s/labels", testSystemID),
								Context: context.Background(),
							}).
							Return(&policies.GetPolicyOK{
								Payload: &models.V1PolicyGetResponse{
									Result: map[string]interface{}{
										"modules": map[string]interface{}{
											"labels.rego": testLabelsRego,
										},
									},
								},
							}, nil)
					}),
				},
				cr: System(
					withExternalName(testSystemID),
					withSpec(v1alpha1.SystemParameters{
						CustomSystemParameters: v1alpha1.CustomSystemParameters{
							Labels: map[string]string{
								testLabelKey: testLabelValue,
							},
						},
						Description: &testDescription,
						DeploymentParameters: &v1alpha1.V1SystemDeploymentParameters{
							HTTPProxy:                styraclient.String(""),
							HTTPSProxy:               styraclient.String(""),
							KubernetesVersion:        styraclient.String(""),
							Namespace:                styraclient.String(""),
							NoProxy:                  styraclient.String(""),
							TimeoutSeconds:           styraclient.Int32(0),
							TrustedContainerRegistry: styraclient.String(""),
						},
						ReadOnly:   styraclient.Bool(true),
						Type:       testType,
						ExternalID: &testExternalID,
					}),
				),
			},
			want: want{
				cr: System(
					withSpec(v1alpha1.SystemParameters{
						CustomSystemParameters: v1alpha1.CustomSystemParameters{
							Labels: map[string]string{
								testLabelKey: testLabelValue,
							},
						},
						Description: &testDescription,
						DeploymentParameters: &v1alpha1.V1SystemDeploymentParameters{
							HTTPProxy:                styraclient.String(""),
							HTTPSProxy:               styraclient.String(""),
							KubernetesVersion:        styraclient.String(""),
							Namespace:                styraclient.String(""),
							NoProxy:                  styraclient.String(""),
							TimeoutSeconds:           styraclient.Int32(0),
							TrustedContainerRegistry: styraclient.String(""),
						},
						ReadOnly:   styraclient.Bool(true),
						Type:       testType,
						ExternalID: &testExternalID,
					}),
					withExternalName(testSystemID),
					withConditions(xpv1.Available()),
				),
				result: managed.ExternalObservation{},
				err:    errors.Wrap(errors.Wrap(errBoom, errGetHelmValues), errGetConnectionDetails),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := &external{client: &tc.styra}
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
		cr     *v1alpha1.System
		result managed.ExternalCreation
		err    error
	}

	cases := map[string]struct {
		args
		want
	}{
		"Successful": {
			args: args{
				styra: styra.StyraAPI{
					Systems: withMockSystem(t, func(mcs *mocksystem.MockClientService) {
						mcs.EXPECT().
							CreateSystem(&systems.CreateSystemParams{
								Body: &models.V1SystemsPostRequest{
									Description: testDescription,
									DeploymentParameters: &models.V1SystemDeploymentParameters{
										DenyOnOpaFail: styraclient.Bool(true),
									},
									Name:       &testSystemName,
									ReadOnly:   styraclient.Bool(true),
									Type:       &testType,
									ExternalID: testExternalID,
								},
								Context: context.Background(),
							}).
							Return(&systems.CreateSystemOK{
								Payload: &models.V1SystemsPostResponse{
									Result: &models.V1SystemConfig{
										ID: testSystemID,
									},
								},
							}, nil)
					}),
				},
				cr: System(
					withName(testSystemName),
					withSpec(v1alpha1.SystemParameters{
						CustomSystemParameters: v1alpha1.CustomSystemParameters{
							Labels: map[string]string{
								testLabelKey: testLabelValue,
							},
						},
						Description: &testDescription,
						DeploymentParameters: &v1alpha1.V1SystemDeploymentParameters{
							DenyOnOpaFail:            styraclient.Bool(true),
							HTTPProxy:                styraclient.String(""),
							HTTPSProxy:               styraclient.String(""),
							KubernetesVersion:        styraclient.String(""),
							Namespace:                styraclient.String(""),
							NoProxy:                  styraclient.String(""),
							TimeoutSeconds:           styraclient.Int32(0),
							TrustedContainerRegistry: styraclient.String(""),
						},
						ReadOnly:   styraclient.Bool(true),
						Type:       testType,
						ExternalID: &testExternalID,
					}),
				),
			},
			want: want{
				cr: System(
					withName(testSystemName),
					withExternalName(testSystemID),
					withSpec(v1alpha1.SystemParameters{
						CustomSystemParameters: v1alpha1.CustomSystemParameters{
							Labels: map[string]string{
								testLabelKey: testLabelValue,
							},
						},
						Description: &testDescription,
						DeploymentParameters: &v1alpha1.V1SystemDeploymentParameters{
							DenyOnOpaFail:            styraclient.Bool(true),
							HTTPProxy:                styraclient.String(""),
							HTTPSProxy:               styraclient.String(""),
							KubernetesVersion:        styraclient.String(""),
							Namespace:                styraclient.String(""),
							NoProxy:                  styraclient.String(""),
							TimeoutSeconds:           styraclient.Int32(0),
							TrustedContainerRegistry: styraclient.String(""),
						},
						ReadOnly:   styraclient.Bool(true),
						Type:       testType,
						ExternalID: &testExternalID,
					}),
				),
				result: managed.ExternalCreation{
					ExternalNameAssigned: true,
				},
			},
		},
		"CreateSystemFailed": {
			args: args{
				styra: styra.StyraAPI{
					Systems: withMockSystem(t, func(mcs *mocksystem.MockClientService) {
						mcs.EXPECT().
							CreateSystem(&systems.CreateSystemParams{
								Body: &models.V1SystemsPostRequest{
									Name: &testSystemName,
									Type: &testType,
								},
								Context: context.Background(),
							}).
							Return(nil, errBoom)
					}),
				},
				cr: System(
					withName(testSystemName),
					withSpec(v1alpha1.SystemParameters{
						Type: testType,
					}),
				),
			},
			want: want{
				cr: System(
					withName(testSystemName),
					withSpec(v1alpha1.SystemParameters{
						Type: testType,
					}),
				),
				err: errors.Wrap(errBoom, errCreateFailed),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := &external{client: &tc.styra}
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
		cr     *v1alpha1.System
		result managed.ExternalUpdate
		err    error
	}

	cases := map[string]struct {
		args
		want
	}{
		"Successful": {
			args: args{
				styra: styra.StyraAPI{
					Systems: withMockSystem(t, func(mcs *mocksystem.MockClientService) {
						mcs.EXPECT().
							UpdateSystem(&systems.UpdateSystemParams{
								System: testSystemID,
								Body: &models.V1SystemsPutRequest{
									Description: testDescription,
									DeploymentParameters: &models.V1SystemDeploymentParameters{
										DenyOnOpaFail: styraclient.Bool(true),
									},
									Name:       &testSystemName,
									ReadOnly:   styraclient.Bool(true),
									Type:       &testType,
									ExternalID: testExternalID,
								},
								Context: context.Background(),
							}).
							Return(&systems.UpdateSystemOK{
								Payload: &models.V1SystemsPutResponse{
									Result: &models.V1SystemConfig{},
								},
							}, nil)
					}),
					Policies: withMockPolicies(t, func(mcs *mockpolicies.MockClientService) {
						mcs.EXPECT().
							UpdatePolicy(&policies.UpdatePolicyParams{
								Policy: fmt.Sprintf("metadata/%s/labels", testSystemID),
								Body: &models.V1PoliciesPutRequest{
									Modules: map[string]string{
										"labels.rego": testLabelsRego,
									},
								},
								Context: context.Background(),
							}).
							Return(&policies.UpdatePolicyOK{}, nil)
					}),
				},
				cr: System(
					withName(testSystemName),
					withExternalName(testSystemID),
					withSpec(v1alpha1.SystemParameters{
						CustomSystemParameters: v1alpha1.CustomSystemParameters{
							Labels: map[string]string{
								testLabelKey: testLabelValue,
							},
						},
						Description: &testDescription,
						DeploymentParameters: &v1alpha1.V1SystemDeploymentParameters{
							DenyOnOpaFail:            styraclient.Bool(true),
							HTTPProxy:                styraclient.String(""),
							HTTPSProxy:               styraclient.String(""),
							KubernetesVersion:        styraclient.String(""),
							Namespace:                styraclient.String(""),
							NoProxy:                  styraclient.String(""),
							TimeoutSeconds:           styraclient.Int32(0),
							TrustedContainerRegistry: styraclient.String(""),
						},
						ReadOnly:   styraclient.Bool(true),
						Type:       testType,
						ExternalID: &testExternalID,
					}),
				),
			},
			want: want{
				cr: System(
					withName(testSystemName),
					withExternalName(testSystemID),
					withSpec(v1alpha1.SystemParameters{
						CustomSystemParameters: v1alpha1.CustomSystemParameters{
							Labels: map[string]string{
								testLabelKey: testLabelValue,
							},
						},
						Description: &testDescription,
						DeploymentParameters: &v1alpha1.V1SystemDeploymentParameters{
							DenyOnOpaFail:            styraclient.Bool(true),
							HTTPProxy:                styraclient.String(""),
							HTTPSProxy:               styraclient.String(""),
							KubernetesVersion:        styraclient.String(""),
							Namespace:                styraclient.String(""),
							NoProxy:                  styraclient.String(""),
							TimeoutSeconds:           styraclient.Int32(0),
							TrustedContainerRegistry: styraclient.String(""),
						},
						ReadOnly:   styraclient.Bool(true),
						Type:       testType,
						ExternalID: &testExternalID,
					}),
				),
				result: managed.ExternalUpdate{},
			},
		},
		"UpdateSystemFailed": {
			args: args{
				styra: styra.StyraAPI{
					Systems: withMockSystem(t, func(mcs *mocksystem.MockClientService) {
						mcs.EXPECT().
							UpdateSystem(&systems.UpdateSystemParams{
								Body: &models.V1SystemsPutRequest{
									Name: &testSystemName,
									Type: &testType,
								},
								Context: context.Background(),
							}).
							Return(nil, errBoom)
					}),
				},
				cr: System(
					withName(testSystemName),
					withSpec(v1alpha1.SystemParameters{
						Type: testType,
					}),
				),
			},
			want: want{
				cr: System(
					withName(testSystemName),
					withSpec(v1alpha1.SystemParameters{
						Type: testType,
					}),
				),
				err: errors.Wrap(errBoom, errUpdateFailed),
			},
		},
		"UpdatePoliciesFailed": {
			args: args{
				styra: styra.StyraAPI{
					Systems: withMockSystem(t, func(mcs *mocksystem.MockClientService) {
						mcs.EXPECT().
							UpdateSystem(&systems.UpdateSystemParams{
								System: testSystemID,
								Body: &models.V1SystemsPutRequest{
									Description: testDescription,
									DeploymentParameters: &models.V1SystemDeploymentParameters{
										DenyOnOpaFail: styraclient.Bool(true),
									},
									Name:       &testSystemName,
									ReadOnly:   styraclient.Bool(true),
									Type:       &testType,
									ExternalID: testExternalID,
								},
								Context: context.Background(),
							}).
							Return(&systems.UpdateSystemOK{
								Payload: &models.V1SystemsPutResponse{
									Result: &models.V1SystemConfig{},
								},
							}, nil)
					}),
					Policies: withMockPolicies(t, func(mcs *mockpolicies.MockClientService) {
						mcs.EXPECT().
							UpdatePolicy(&policies.UpdatePolicyParams{
								Policy: fmt.Sprintf("metadata/%s/labels", testSystemID),
								Body: &models.V1PoliciesPutRequest{
									Modules: map[string]string{
										"labels.rego": testLabelsRego,
									},
								},
								Context: context.Background(),
							}).
							Return(nil, errBoom)
					}),
				},
				cr: System(
					withName(testSystemName),
					withExternalName(testSystemID),
					withSpec(v1alpha1.SystemParameters{
						CustomSystemParameters: v1alpha1.CustomSystemParameters{
							Labels: map[string]string{
								testLabelKey: testLabelValue,
							},
						},
						Description: &testDescription,
						DeploymentParameters: &v1alpha1.V1SystemDeploymentParameters{
							DenyOnOpaFail:            styraclient.Bool(true),
							HTTPProxy:                styraclient.String(""),
							HTTPSProxy:               styraclient.String(""),
							KubernetesVersion:        styraclient.String(""),
							Namespace:                styraclient.String(""),
							NoProxy:                  styraclient.String(""),
							TimeoutSeconds:           styraclient.Int32(0),
							TrustedContainerRegistry: styraclient.String(""),
						},
						ReadOnly:   styraclient.Bool(true),
						Type:       testType,
						ExternalID: &testExternalID,
					}),
				),
			},
			want: want{
				cr: System(
					withName(testSystemName),
					withExternalName(testSystemID),
					withSpec(v1alpha1.SystemParameters{
						CustomSystemParameters: v1alpha1.CustomSystemParameters{
							Labels: map[string]string{
								testLabelKey: testLabelValue,
							},
						},
						Description: &testDescription,
						DeploymentParameters: &v1alpha1.V1SystemDeploymentParameters{
							DenyOnOpaFail:            styraclient.Bool(true),
							HTTPProxy:                styraclient.String(""),
							HTTPSProxy:               styraclient.String(""),
							KubernetesVersion:        styraclient.String(""),
							Namespace:                styraclient.String(""),
							NoProxy:                  styraclient.String(""),
							TimeoutSeconds:           styraclient.Int32(0),
							TrustedContainerRegistry: styraclient.String(""),
						},
						ReadOnly:   styraclient.Bool(true),
						Type:       testType,
						ExternalID: &testExternalID,
					}),
				),
				result: managed.ExternalUpdate{},
				err:    errors.Wrap(errors.Wrap(errBoom, errUpdateLabels), errUpdateFailed),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := &external{client: &tc.styra}
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
		cr  *v1alpha1.System
		err error
	}

	cases := map[string]struct {
		args
		want
	}{
		"Successful": {
			args: args{
				styra: styra.StyraAPI{
					Systems: withMockSystem(t, func(mcs *mocksystem.MockClientService) {
						mcs.EXPECT().
							DeleteSystem(&systems.DeleteSystemParams{
								System:  testSystemID,
								Context: context.Background(),
							}).
							Return(&systems.DeleteSystemOK{}, nil)
					}),
				},
				cr: System(
					withExternalName(testSystemID),
				),
			},
			want: want{
				cr: System(
					withExternalName(testSystemID),
				),
			},
		},
		"DeleteFailed": {
			args: args{
				styra: styra.StyraAPI{
					Systems: withMockSystem(t, func(mcs *mocksystem.MockClientService) {
						mcs.EXPECT().
							DeleteSystem(&systems.DeleteSystemParams{
								System:  testSystemID,
								Context: context.Background(),
							}).
							Return(nil, errBoom)
					}),
				},
				cr: System(
					withExternalName(testSystemID),
				),
			},
			want: want{
				cr: System(
					withExternalName(testSystemID),
				),
				err: errors.Wrap(errBoom, errDeleteFailed),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := &external{client: &tc.styra}
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

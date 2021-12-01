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
	"github.com/mistermx/styra-go-client/pkg/client/stacks"
	"github.com/mistermx/styra-go-client/pkg/models"

	v1alpha1 "github.com/crossplane-contrib/provider-styra/apis/stack/v1alpha1"
	styraclient "github.com/crossplane-contrib/provider-styra/pkg/client"
	mockpolicies "github.com/crossplane-contrib/provider-styra/pkg/client/mock/policies"
	mockstack "github.com/crossplane-contrib/provider-styra/pkg/client/mock/stacks"
)

var (
	errBoom = errors.New("boom")

	testStackID       = "teststack"
	testStackName     = "testname"
	testType          = "kubernetes:test"
	testDescription   = "test-description"
	testSelectorKey   = "key"
	testSelectorValue = "value"
	testCredentials   = "test-credentials"
	testURL           = "test-url"
	testReference     = "test-reference"
	testPath          = "test-path"
	testSelectorRego  = `
package stacks.teststack.selectors
import data.library.v1.utils.labels.match.v1 as match

systems[system_id] {
  include := {
    "key": {
      "value",
    },
  }

  exclude := {
    "key": {
      "value",
    },
  }

  metadata := data.metadata[system_id]
  match.all(metadata.labels.labels, include, exclude)
}
`
)

type args struct {
	styra styra.StyraAPI
	cr    *v1alpha1.Stack
}

type mockStackModifier func(*mockstack.MockClientService)

func withMockStack(t *testing.T, mod mockStackModifier) *mockstack.MockClientService {
	ctrl := gomock.NewController(t)
	mock := mockstack.NewMockClientService(ctrl)
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

type StackModifier func(*v1alpha1.Stack)

func withName(v string) StackModifier {
	return func(s *v1alpha1.Stack) {
		s.ObjectMeta.Name = v
	}
}

func withExternalName(v string) StackModifier {
	return func(s *v1alpha1.Stack) {
		meta.SetExternalName(s, v)
	}
}

func withConditions(c ...xpv1.Condition) StackModifier {
	return func(r *v1alpha1.Stack) { r.Status.ConditionedStatus.Conditions = c }
}

func withSpec(p v1alpha1.StackParameters) StackModifier {
	return func(r *v1alpha1.Stack) { r.Spec.ForProvider = p }
}

func Stack(m ...StackModifier) *v1alpha1.Stack {
	cr := &v1alpha1.Stack{}
	for _, f := range m {
		f(cr)
	}
	return cr
}

func TestObserve(t *testing.T) {
	type want struct {
		cr     *v1alpha1.Stack
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
					Stacks: withMockStack(t, func(mcs *mockstack.MockClientService) {
						mcs.EXPECT().
							GetStack(&stacks.GetStackParams{
								Stack:   testStackID,
								Context: context.Background(),
							}).
							Return(&stacks.GetStackOK{
								Payload: &models.StacksV1StacksGetResponse{
									Result: &models.StacksV1StackConfig{
										Description: &testDescription,
										ReadOnly:    styraclient.Bool(true),
										Type:        &testType,
										SourceControl: &models.StacksV1SourceControlConfig{
											Origin: &models.GitV1GitRepoConfig{
												Credentials: &testCredentials,
												Path:        &testPath,
												Reference:   &testReference,
												URL:         &testURL,
											},
										},
									},
								},
							}, nil)
					}),
					Policies: withMockPolicies(t, func(mcs *mockpolicies.MockClientService) {
						mcs.EXPECT().
							GetPolicy(&policies.GetPolicyParams{
								Policy:  fmt.Sprintf("stacks/%s/selectors", testStackID),
								Context: context.Background(),
							}).
							Return(&policies.GetPolicyOK{
								Payload: &models.PoliciesV1PolicyGetResponse{
									Result: map[string]interface{}{
										"modules": map[string]interface{}{
											"selector.rego": testSelectorRego,
										},
									},
								},
							}, nil)
					}),
				},
				cr: Stack(
					withExternalName(testStackID),
					withSpec(v1alpha1.StackParameters{
						CustomStackParameters: v1alpha1.CustomStackParameters{
							SelectorInclude: map[string][]string{
								testSelectorKey: {testSelectorValue},
							},
							SelectorExclude: map[string][]string{
								testSelectorKey: {testSelectorValue},
							},
						},
						Description: testDescription,
						ReadOnly:    true,
						Type:        testType,
						SourceControl: &v1alpha1.V1SourceControlConfig{
							Origin: v1alpha1.V1GitRepoConfig{
								Credentials: testCredentials,
								Path:        testPath,
								Reference:   testReference,
								URL:         testURL,
							},
						},
					}),
				),
			},
			want: want{
				cr: Stack(
					withExternalName(testStackID),
					withSpec(v1alpha1.StackParameters{
						CustomStackParameters: v1alpha1.CustomStackParameters{
							SelectorInclude: map[string][]string{
								testSelectorKey: {testSelectorValue},
							},
							SelectorExclude: map[string][]string{
								testSelectorKey: {testSelectorValue},
							},
						},
						Description: testDescription,
						ReadOnly:    true,
						Type:        testType,
						SourceControl: &v1alpha1.V1SourceControlConfig{
							Origin: v1alpha1.V1GitRepoConfig{
								Credentials: testCredentials,
								Path:        testPath,
								Reference:   testReference,
								URL:         testURL,
							},
						},
					}),
					withConditions(xpv1.Available()),
				),
				result: managed.ExternalObservation{
					ResourceExists:   true,
					ResourceUpToDate: true,
				},
			},
		},
		"SuccessfulLateInitialize": {
			args: args{
				styra: styra.StyraAPI{
					Stacks: withMockStack(t, func(mcs *mockstack.MockClientService) {
						mcs.EXPECT().
							GetStack(&stacks.GetStackParams{
								Stack:   testStackID,
								Context: context.Background(),
							}).
							Return(&stacks.GetStackOK{
								Payload: &models.StacksV1StacksGetResponse{
									Result: &models.StacksV1StackConfig{
										Description: &testDescription,
										ReadOnly:    styraclient.Bool(true),
										Type:        &testType,
										SourceControl: &models.StacksV1SourceControlConfig{
											Origin: &models.GitV1GitRepoConfig{
												Credentials: &testCredentials,
												Path:        &testPath,
												Reference:   &testReference,
												URL:         &testURL,
											},
										},
									},
								},
							}, nil)
					}),
					Policies: withMockPolicies(t, func(mcs *mockpolicies.MockClientService) {
						mcs.EXPECT().
							GetPolicy(&policies.GetPolicyParams{
								Policy:  fmt.Sprintf("stacks/%s/selectors", testStackID),
								Context: context.Background(),
							}).
							Return(&policies.GetPolicyOK{
								Payload: &models.PoliciesV1PolicyGetResponse{
									Result: map[string]interface{}{
										"modules": map[string]interface{}{
											"selector.rego": testSelectorRego,
										},
									},
								},
							}, nil)
					}),
				},
				cr: Stack(
					withExternalName(testStackID),
					withSpec(v1alpha1.StackParameters{
						CustomStackParameters: v1alpha1.CustomStackParameters{
							SelectorInclude: map[string][]string{
								testSelectorKey: {testSelectorValue},
							},
							SelectorExclude: map[string][]string{
								testSelectorKey: {testSelectorValue},
							},
						},
						Description: testDescription,
						ReadOnly:    true,
						Type:        testType,
					}),
				),
			},
			want: want{
				cr: Stack(
					withExternalName(testStackID),
					withSpec(v1alpha1.StackParameters{
						CustomStackParameters: v1alpha1.CustomStackParameters{
							SelectorInclude: map[string][]string{
								testSelectorKey: {testSelectorValue},
							},
							SelectorExclude: map[string][]string{
								testSelectorKey: {testSelectorValue},
							},
						},
						Description: testDescription,
						ReadOnly:    true,
						Type:        testType,
						SourceControl: &v1alpha1.V1SourceControlConfig{
							Origin: v1alpha1.V1GitRepoConfig{
								Credentials: testCredentials,
								Path:        testPath,
								Reference:   testReference,
								URL:         testURL,
							},
						},
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
		"SelectorsNotUpToDate": {
			args: args{
				styra: styra.StyraAPI{
					Stacks: withMockStack(t, func(mcs *mockstack.MockClientService) {
						mcs.EXPECT().
							GetStack(&stacks.GetStackParams{
								Stack:   testStackID,
								Context: context.Background(),
							}).
							Return(&stacks.GetStackOK{
								Payload: &models.StacksV1StacksGetResponse{
									Result: &models.StacksV1StackConfig{
										Description: &testDescription,
										ReadOnly:    styraclient.Bool(true),
										Type:        &testType,
									},
								},
							}, nil)
					}),
					Policies: withMockPolicies(t, func(mcs *mockpolicies.MockClientService) {
						mcs.EXPECT().
							GetPolicy(&policies.GetPolicyParams{
								Policy:  fmt.Sprintf("stacks/%s/selectors", testStackID),
								Context: context.Background(),
							}).
							Return(&policies.GetPolicyOK{
								Payload: &models.PoliciesV1PolicyGetResponse{
									Result: map[string]interface{}{
										"modules": map[string]interface{}{
											"selector.rego": testSelectorRego,
										},
									},
								},
							}, nil)
					}),
				},
				cr: Stack(
					withExternalName(testStackID),
					withSpec(v1alpha1.StackParameters{
						CustomStackParameters: v1alpha1.CustomStackParameters{
							SelectorInclude: map[string][]string{
								testSelectorKey: {testSelectorValue, testSelectorValue},
							},
						},
						Description: testDescription,
						ReadOnly:    true,
						Type:        testType,
					}),
				),
			},
			want: want{
				cr: Stack(
					withExternalName(testStackID),
					withSpec(v1alpha1.StackParameters{
						CustomStackParameters: v1alpha1.CustomStackParameters{
							SelectorInclude: map[string][]string{
								testSelectorKey: {testSelectorValue, testSelectorValue},
							},
						},
						Description: testDescription,
						ReadOnly:    true,
						Type:        testType,
					}),
					withConditions(xpv1.Available()),
				),
				result: managed.ExternalObservation{
					ResourceExists:   true,
					ResourceUpToDate: false,
				},
			},
		},
		"GetStackFailed": {
			args: args{
				styra: styra.StyraAPI{
					Stacks: withMockStack(t, func(mcs *mockstack.MockClientService) {
						mcs.EXPECT().
							GetStack(&stacks.GetStackParams{
								Stack:   testStackID,
								Context: context.Background(),
							}).
							Return(nil, errBoom)
					}),
				},
				cr: Stack(
					withExternalName(testStackID),
				),
			},
			want: want{
				cr: Stack(
					withExternalName(testStackID),
				),
				err: errors.Wrap(errBoom, errDescribeFailed),
			},
		},
		"GetSelectorsFailed": {
			args: args{
				styra: styra.StyraAPI{
					Stacks: withMockStack(t, func(mcs *mockstack.MockClientService) {
						mcs.EXPECT().
							GetStack(&stacks.GetStackParams{
								Stack:   testStackID,
								Context: context.Background(),
							}).
							Return(&stacks.GetStackOK{
								Payload: &models.StacksV1StacksGetResponse{
									Result: &models.StacksV1StackConfig{
										Description: &testDescription,
										ReadOnly:    styraclient.Bool(true),
										Type:        &testType,
									},
								},
							}, nil)
					}),
					Policies: withMockPolicies(t, func(mcs *mockpolicies.MockClientService) {
						mcs.EXPECT().
							GetPolicy(&policies.GetPolicyParams{
								Policy:  fmt.Sprintf("stacks/%s/selectors", testStackID),
								Context: context.Background(),
							}).
							Return(nil, errBoom)
					}),
				},
				cr: Stack(
					withExternalName(testStackID),
					withSpec(v1alpha1.StackParameters{
						Type:        testType,
						Description: testDescription,
						ReadOnly:    true,
					}),
				),
			},
			want: want{
				cr: Stack(
					withExternalName(testStackID),
					withSpec(v1alpha1.StackParameters{
						Type:        testType,
						Description: testDescription,
						ReadOnly:    true,
					}),
				),
				err: errors.Wrap(errors.Wrap(errBoom, errGetSelectors), errIsUpToDateFailed),
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
		cr     *v1alpha1.Stack
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
					Stacks: withMockStack(t, func(mcs *mockstack.MockClientService) {
						mcs.EXPECT().
							CreateStack(&stacks.CreateStackParams{
								Body: &models.StacksV1StacksPostRequest{
									Description: &testDescription,
									Name:        &testStackName,
									ReadOnly:    styraclient.Bool(true),
									Type:        &testType,
									SourceControl: &models.StacksV1SourceControlConfig{
										Origin: &models.GitV1GitRepoConfig{
											Credentials: &testCredentials,
											Path:        &testPath,
											Reference:   &testReference,
											URL:         &testURL,
										},
									},
								},
								Context: context.Background(),
							}).
							Return(&stacks.CreateStackOK{
								Payload: &models.StacksV1StacksPostResponse{
									Result: &models.StacksV1StackConfig{
										ID: &testStackID,
									},
								},
							}, nil)
					}),
				},
				cr: Stack(
					withName(testStackName),
					withSpec(v1alpha1.StackParameters{
						CustomStackParameters: v1alpha1.CustomStackParameters{
							SelectorInclude: map[string][]string{
								testSelectorKey: {testSelectorValue},
							},
							SelectorExclude: map[string][]string{
								testSelectorKey: {testSelectorValue},
							},
						},
						Description: testDescription,
						ReadOnly:    true,
						Type:        testType,
						SourceControl: &v1alpha1.V1SourceControlConfig{
							Origin: v1alpha1.V1GitRepoConfig{
								Credentials: testCredentials,
								Path:        testPath,
								Reference:   testReference,
								URL:         testURL,
							},
						},
					}),
				),
			},
			want: want{
				cr: Stack(
					withName(testStackName),
					withExternalName(testStackID),
					withSpec(v1alpha1.StackParameters{
						CustomStackParameters: v1alpha1.CustomStackParameters{
							SelectorInclude: map[string][]string{
								testSelectorKey: {testSelectorValue},
							},
							SelectorExclude: map[string][]string{
								testSelectorKey: {testSelectorValue},
							},
						},
						Description: testDescription,
						ReadOnly:    true,
						Type:        testType,
						SourceControl: &v1alpha1.V1SourceControlConfig{
							Origin: v1alpha1.V1GitRepoConfig{
								Credentials: testCredentials,
								Path:        testPath,
								Reference:   testReference,
								URL:         testURL,
							},
						},
					}),
				),
				result: managed.ExternalCreation{
					ExternalNameAssigned: true,
				},
			},
		},
		"CreateStackFailed": {
			args: args{
				styra: styra.StyraAPI{
					Stacks: withMockStack(t, func(mcs *mockstack.MockClientService) {
						mcs.EXPECT().
							CreateStack(&stacks.CreateStackParams{
								Body: &models.StacksV1StacksPostRequest{
									Name:        &testStackName,
									Type:        &testType,
									Description: &testDescription,
									ReadOnly:    styraclient.Bool(true),
								},
								Context: context.Background(),
							}).
							Return(nil, errBoom)
					}),
				},
				cr: Stack(
					withName(testStackName),
					withSpec(v1alpha1.StackParameters{
						Type:        testType,
						Description: testDescription,
						ReadOnly:    true,
					}),
				),
			},
			want: want{
				cr: Stack(
					withName(testStackName),
					withSpec(v1alpha1.StackParameters{
						Type:        testType,
						Description: testDescription,
						ReadOnly:    true,
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
		cr     *v1alpha1.Stack
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
					Stacks: withMockStack(t, func(mcs *mockstack.MockClientService) {
						mcs.EXPECT().
							UpdateStack(&stacks.UpdateStackParams{
								Stack: testStackID,
								Body: &models.StacksV1StacksPutRequest{
									Name:        &testStackName,
									Description: &testDescription,
									ReadOnly:    styraclient.Bool(true),
									Type:        &testType,
									SourceControl: &models.StacksV1SourceControlConfig{
										Origin: &models.GitV1GitRepoConfig{
											Credentials: &testCredentials,
											Path:        &testPath,
											Reference:   &testReference,
											URL:         &testURL,
										},
									},
								},
								Context: context.Background(),
							}).
							Return(&stacks.UpdateStackOK{
								Payload: &models.StacksV1StacksPutResponse{},
							}, nil)
					}),
					Policies: withMockPolicies(t, func(mcs *mockpolicies.MockClientService) {
						mcs.EXPECT().
							UpdatePolicy(&policies.UpdatePolicyParams{
								Policy: fmt.Sprintf("stacks/%s/selectors", testStackID),
								Body: &models.PoliciesV1PoliciesPutRequest{
									Modules: map[string]string{
										"selector.rego": testSelectorRego,
									},
								},
								Context: context.Background(),
							}).
							Return(&policies.UpdatePolicyOK{}, nil)
					}),
				},
				cr: Stack(
					withName(testStackName),
					withExternalName(testStackID),
					withSpec(v1alpha1.StackParameters{
						CustomStackParameters: v1alpha1.CustomStackParameters{
							SelectorInclude: map[string][]string{
								testSelectorKey: {testSelectorValue},
							},
							SelectorExclude: map[string][]string{
								testSelectorKey: {testSelectorValue},
							},
						},
						Description: testDescription,
						ReadOnly:    true,
						Type:        testType,
						SourceControl: &v1alpha1.V1SourceControlConfig{
							Origin: v1alpha1.V1GitRepoConfig{
								Credentials: testCredentials,
								Path:        testPath,
								Reference:   testReference,
								URL:         testURL,
							},
						},
					}),
				),
			},
			want: want{
				cr: Stack(
					withName(testStackName),
					withExternalName(testStackID),
					withSpec(v1alpha1.StackParameters{
						CustomStackParameters: v1alpha1.CustomStackParameters{
							SelectorInclude: map[string][]string{
								testSelectorKey: {testSelectorValue},
							},
							SelectorExclude: map[string][]string{
								testSelectorKey: {testSelectorValue},
							},
						},
						Description: testDescription,
						ReadOnly:    true,
						Type:        testType,
						SourceControl: &v1alpha1.V1SourceControlConfig{
							Origin: v1alpha1.V1GitRepoConfig{
								Credentials: testCredentials,
								Path:        testPath,
								Reference:   testReference,
								URL:         testURL,
							},
						},
					}),
				),
				result: managed.ExternalUpdate{},
			},
		},
		"UpdateStackFailed": {
			args: args{
				styra: styra.StyraAPI{
					Stacks: withMockStack(t, func(mcs *mockstack.MockClientService) {
						mcs.EXPECT().
							UpdateStack(&stacks.UpdateStackParams{
								Stack: testStackID,
								Body: &models.StacksV1StacksPutRequest{
									Name:        &testStackName,
									Description: &testDescription,
									ReadOnly:    styraclient.Bool(true),
									Type:        &testType,
								},
								Context: context.Background(),
							}).
							Return(nil, errBoom)
					}),
				},
				cr: Stack(
					withName(testStackName),
					withExternalName(testStackID),
					withSpec(v1alpha1.StackParameters{
						Description: testDescription,
						ReadOnly:    true,
						Type:        testType,
					}),
				),
			},
			want: want{
				cr: Stack(
					withName(testStackName),
					withExternalName(testStackID),
					withSpec(v1alpha1.StackParameters{
						Description: testDescription,
						ReadOnly:    true,
						Type:        testType,
					}),
				),
				err: errors.Wrap(errBoom, errUpdateFailed),
			},
		},
		"UpdatePoliciesFailed": {
			args: args{
				styra: styra.StyraAPI{
					Stacks: withMockStack(t, func(mcs *mockstack.MockClientService) {
						mcs.EXPECT().
							UpdateStack(&stacks.UpdateStackParams{
								Stack: testStackID,
								Body: &models.StacksV1StacksPutRequest{
									Name:        &testStackName,
									Description: &testDescription,
									ReadOnly:    styraclient.Bool(true),
									Type:        &testType,
								},
								Context: context.Background(),
							}).
							Return(&stacks.UpdateStackOK{
								Payload: &models.StacksV1StacksPutResponse{},
							}, nil)
					}),
					Policies: withMockPolicies(t, func(mcs *mockpolicies.MockClientService) {
						mcs.EXPECT().
							UpdatePolicy(&policies.UpdatePolicyParams{
								Policy: fmt.Sprintf("stacks/%s/selectors", testStackID),
								Body: &models.PoliciesV1PoliciesPutRequest{
									Modules: map[string]string{
										"selector.rego": testSelectorRego,
									},
								},
								Context: context.Background(),
							}).
							Return(nil, errBoom)
					}),
				},
				cr: Stack(
					withName(testStackName),
					withExternalName(testStackID),
					withSpec(v1alpha1.StackParameters{
						CustomStackParameters: v1alpha1.CustomStackParameters{
							SelectorInclude: map[string][]string{
								testSelectorKey: {testSelectorValue},
							},
							SelectorExclude: map[string][]string{
								testSelectorKey: {testSelectorValue},
							},
						},
						Description: testDescription,
						ReadOnly:    true,
						Type:        testType,
					}),
				),
			},
			want: want{
				cr: Stack(
					withName(testStackName),
					withExternalName(testStackID),
					withSpec(v1alpha1.StackParameters{
						CustomStackParameters: v1alpha1.CustomStackParameters{
							SelectorInclude: map[string][]string{
								testSelectorKey: {testSelectorValue},
							},
							SelectorExclude: map[string][]string{
								testSelectorKey: {testSelectorValue},
							},
						},
						Description: testDescription,
						ReadOnly:    true,
						Type:        testType,
					}),
				),
				result: managed.ExternalUpdate{},
				err:    errors.Wrap(errors.Wrap(errBoom, errUpdateSelectors), errUpdateFailed),
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
		cr  *v1alpha1.Stack
		err error
	}

	cases := map[string]struct {
		args
		want
	}{
		"Successful": {
			args: args{
				styra: styra.StyraAPI{
					Stacks: withMockStack(t, func(mcs *mockstack.MockClientService) {
						mcs.EXPECT().
							DeleteStack(&stacks.DeleteStackParams{
								Stack:   testStackID,
								Context: context.Background(),
							}).
							Return(&stacks.DeleteStackOK{}, nil)
					}),
				},
				cr: Stack(
					withExternalName(testStackID),
				),
			},
			want: want{
				cr: Stack(
					withExternalName(testStackID),
				),
			},
		},
		"DeleteFailed": {
			args: args{
				styra: styra.StyraAPI{
					Stacks: withMockStack(t, func(mcs *mockstack.MockClientService) {
						mcs.EXPECT().
							DeleteStack(&stacks.DeleteStackParams{
								Stack:   testStackID,
								Context: context.Background(),
							}).
							Return(nil, errBoom)
					}),
				},
				cr: Stack(
					withExternalName(testStackID),
				),
			},
			want: want{
				cr: Stack(
					withExternalName(testStackID),
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

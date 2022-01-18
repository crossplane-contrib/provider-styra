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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/resource"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/test"

	styra "github.com/mistermx/styra-go-client/pkg/client"
	"github.com/mistermx/styra-go-client/pkg/client/datasources"
	"github.com/mistermx/styra-go-client/pkg/models"

	v1alpha1 "github.com/crossplane-contrib/provider-styra/apis/datasource/v1alpha1"
	styraclient "github.com/crossplane-contrib/provider-styra/pkg/client"
	mockdatasource "github.com/crossplane-contrib/provider-styra/pkg/client/mock/datasources"
)

var (
	errBoom = errors.New("boom")

	testDataSourceID             = "test-datasource"
	testCategory                 = "aws/ecr"
	testType                     = "pull"
	testRegistryID               = "test-registry"
	testCredentials              = "test-creds"
	testDescription              = "test-desc"
	testPollingInterval    int64 = 10
	testPollingIntervalStr       = "10s"
	testRateLimitNumber          = 3.0
	testRegion                   = "test-region"

	empty = ""

	testRateLimitQuantity       = resource.NewQuantity(3, resource.DecimalSI)
	testPollingIntervalDuration = generateDurationFromSeconds(testPollingInterval)
)

type args struct {
	styra styra.StyraAPI
	cr    *v1alpha1.DataSource
}

type mockDataSourceModifier func(*mockdatasource.MockClientService)

func withMockDataSource(t *testing.T, mod mockDataSourceModifier) *mockdatasource.MockClientService {
	ctrl := gomock.NewController(t)
	mock := mockdatasource.NewMockClientService(ctrl)
	mod(mock)
	return mock
}

type DataSourceModifier func(*v1alpha1.DataSource)

func withExternalName(v string) DataSourceModifier {
	return func(s *v1alpha1.DataSource) {
		meta.SetExternalName(s, v)
	}
}

func withConditions(c ...xpv1.Condition) DataSourceModifier {
	return func(r *v1alpha1.DataSource) { r.Status.ConditionedStatus.Conditions = c }
}

func withSpec(p v1alpha1.DataSourceParameters) DataSourceModifier {
	return func(r *v1alpha1.DataSource) { r.Spec.ForProvider = p }
}

func withStatus(s v1alpha1.DataSourceObservation) DataSourceModifier {
	return func(r *v1alpha1.DataSource) { r.Status.AtProvider = s }
}

func DataSource(m ...DataSourceModifier) *v1alpha1.DataSource {
	cr := &v1alpha1.DataSource{}
	for _, f := range m {
		f(cr)
	}
	return cr
}

func TestObserve(t *testing.T) {
	type want struct {
		cr     *v1alpha1.DataSource
		result managed.ExternalObservation
		err    error
	}

	cases := map[string]struct {
		args
		want
	}{
		"AWS_ECR_SuccessfulAvailable": {
			args: args{
				styra: styra.StyraAPI{
					Datasources: withMockDataSource(t, func(mcs *mockdatasource.MockClientService) {
						mcs.EXPECT().
							GetDatasource(&datasources.GetDatasourceParams{
								Datasource: testDataSourceID,
								Context:    context.Background(),
							}).
							Return(&datasources.GetDatasourceOK{
								Payload: &models.DatasourcesV1DatasourcesGetResponse{
									Result: &models.DatasourcesV1DatasourcesGetResponseResult{
										Category:        testCategory,
										Type:            testType,
										Enabled:         false,
										OnPremises:      true,
										Description:     testDescription,
										RateLimit:       &testRateLimitNumber,
										PollingInterval: testPollingInterval,
										Credentials:     testCredentials,
										Region:          testRegion,
										RegistryID:      testRegistryID,
									},
								},
							}, nil)
					}),
				},
				cr: DataSource(
					withExternalName(testDataSourceID),
					withSpec(v1alpha1.DataSourceParameters{
						DatasourcesV1Common: v1alpha1.DatasourcesV1Common{
							Category:    testCategory,
							Type:        testType,
							Enabled:     styraclient.Bool(false),
							OnPremises:  true,
							Description: &testDescription,
						},
						AWSECR: &v1alpha1.DatasourcesV1AWSECR{
							DatasourcesV1RateLimiter: v1alpha1.DatasourcesV1RateLimiter{
								RateLimit: testRateLimitQuantity,
							},
							DatasourcesV1Poller: v1alpha1.DatasourcesV1Poller{
								PollingInterval: testPollingIntervalDuration,
							},
							DatasourcesV1AWSCommon: v1alpha1.DatasourcesV1AWSCommon{
								Credentials: testCredentials,
								Region:      testRegion,
							},
							RegistryID: &testRegistryID,
						},
					}),
				),
			},
			want: want{
				cr: DataSource(
					withExternalName(testDataSourceID),
					withSpec(v1alpha1.DataSourceParameters{
						DatasourcesV1Common: v1alpha1.DatasourcesV1Common{
							Category:    testCategory,
							Type:        testType,
							Enabled:     styraclient.Bool(false),
							OnPremises:  true,
							Description: &testDescription,
						},
						AWSECR: &v1alpha1.DatasourcesV1AWSECR{
							DatasourcesV1RateLimiter: v1alpha1.DatasourcesV1RateLimiter{
								RateLimit: testRateLimitQuantity,
							},
							DatasourcesV1Poller: v1alpha1.DatasourcesV1Poller{
								PollingInterval: testPollingIntervalDuration,
							},
							DatasourcesV1AWSCommon: v1alpha1.DatasourcesV1AWSCommon{
								Credentials: testCredentials,
								Region:      testRegion,
							},
							RegistryID: &testRegistryID,
						},
					}),
					withConditions(xpv1.Available()),
					withStatus(v1alpha1.DataSourceObservation{
						Executed: &empty,
					}),
				),
				result: managed.ExternalObservation{
					ResourceExists:   true,
					ResourceUpToDate: true,
				},
			},
		},
		"AWS_ECR_SuccessfulLateInitialize": {
			args: args{
				styra: styra.StyraAPI{
					Datasources: withMockDataSource(t, func(mcs *mockdatasource.MockClientService) {
						mcs.EXPECT().
							GetDatasource(&datasources.GetDatasourceParams{
								Datasource: testDataSourceID,
								Context:    context.Background(),
							}).
							Return(&datasources.GetDatasourceOK{
								Payload: &models.DatasourcesV1DatasourcesGetResponse{
									Result: &models.DatasourcesV1DatasourcesGetResponseResult{
										Category:        testCategory,
										Type:            testType,
										Enabled:         false,
										OnPremises:      true,
										Description:     testDescription,
										RateLimit:       &testRateLimitNumber,
										PollingInterval: testPollingInterval,
										Credentials:     testCredentials,
										Region:          testRegion,
										RegistryID:      testRegistryID,
									},
								},
							}, nil)
					}),
				},
				cr: DataSource(
					withExternalName(testDataSourceID),
					withSpec(v1alpha1.DataSourceParameters{
						DatasourcesV1Common: v1alpha1.DatasourcesV1Common{
							Category:    testCategory,
							Type:        testType,
							Enabled:     styraclient.Bool(false),
							OnPremises:  true,
							Description: &testDescription,
						},
					}),
					withConditions(xpv1.Available()),
					withStatus(v1alpha1.DataSourceObservation{
						Executed: &empty,
					}),
				),
			},
			want: want{
				cr: DataSource(
					withExternalName(testDataSourceID),
					withSpec(v1alpha1.DataSourceParameters{
						DatasourcesV1Common: v1alpha1.DatasourcesV1Common{
							Category:    testCategory,
							Type:        testType,
							Enabled:     styraclient.Bool(false),
							OnPremises:  true,
							Description: &testDescription,
						},
						AWSECR: &v1alpha1.DatasourcesV1AWSECR{
							DatasourcesV1RateLimiter: v1alpha1.DatasourcesV1RateLimiter{
								RateLimit: testRateLimitQuantity,
							},
							DatasourcesV1Poller: v1alpha1.DatasourcesV1Poller{
								PollingInterval: testPollingIntervalDuration,
							},
							DatasourcesV1AWSCommon: v1alpha1.DatasourcesV1AWSCommon{
								Credentials: testCredentials,
								Region:      testRegion,
							},
							RegistryID: &testRegistryID,
						},
					}),
					withConditions(xpv1.Available()),
					withStatus(v1alpha1.DataSourceObservation{
						Executed: &empty,
					}),
				),
				result: managed.ExternalObservation{
					ResourceExists:          true,
					ResourceUpToDate:        true,
					ResourceLateInitialized: true,
				},
			},
		},
		"AWS_ECR_GetDataSourceFailed": {
			args: args{
				styra: styra.StyraAPI{
					Datasources: withMockDataSource(t, func(mcs *mockdatasource.MockClientService) {
						mcs.EXPECT().
							GetDatasource(&datasources.GetDatasourceParams{
								Datasource: testDataSourceID,
								Context:    context.Background(),
							}).
							Return(nil, errBoom)
					}),
				},
				cr: DataSource(
					withExternalName(testDataSourceID),
				),
			},
			want: want{
				cr: DataSource(
					withExternalName(testDataSourceID),
				),
				err: errors.Wrap(errBoom, errDescribeFailed),
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
		cr     *v1alpha1.DataSource
		result managed.ExternalCreation
		err    error
	}

	cases := map[string]struct {
		args
		want
	}{
		"AWS_ECR_Successful": {
			args: args{
				styra: styra.StyraAPI{
					Datasources: withMockDataSource(t, func(mcs *mockdatasource.MockClientService) {
						mcs.EXPECT().
							UpsertDatasource(&datasources.UpsertDatasourceParams{
								Body: &models.DatasourcesV1AWSECR{
									DatasourcesV1Common: models.DatasourcesV1Common{
										Category:    &testCategory,
										Type:        &testType,
										Enabled:     styraclient.Bool(false),
										OnPremises:  styraclient.Bool(true),
										Description: testDescription,
									},
									DatasourcesV1RateLimiter: models.DatasourcesV1RateLimiter{

										RateLimit: &testRateLimitNumber,
									},
									DatasourcesV1Poller: models.DatasourcesV1Poller{

										PollingInterval: &testPollingIntervalStr,
									},
									DatasourcesV1AWSCommon: models.DatasourcesV1AWSCommon{
										Credentials: &testCredentials,
										Region:      &testRegion,
									},
									RegistryID: testRegistryID,
								},
								Datasource: testDataSourceID,
								Context:    context.Background(),
							}).
							Return(&datasources.UpsertDatasourceOK{
								Payload: &models.DatasourcesV1DatasourcesPutResponse{},
							}, nil)
					}),
				},
				cr: DataSource(
					withExternalName(testDataSourceID),
					withSpec(v1alpha1.DataSourceParameters{
						DatasourcesV1Common: v1alpha1.DatasourcesV1Common{
							Category:    testCategory,
							Type:        testType,
							Enabled:     styraclient.Bool(false),
							OnPremises:  true,
							Description: &testDescription,
						},
						AWSECR: &v1alpha1.DatasourcesV1AWSECR{
							DatasourcesV1RateLimiter: v1alpha1.DatasourcesV1RateLimiter{
								RateLimit: testRateLimitQuantity,
							},
							DatasourcesV1Poller: v1alpha1.DatasourcesV1Poller{
								PollingInterval: testPollingIntervalDuration,
							},
							DatasourcesV1AWSCommon: v1alpha1.DatasourcesV1AWSCommon{
								Credentials: testCredentials,
								Region:      testRegion,
							},
							RegistryID: &testRegistryID,
						},
					}),
				),
			},
			want: want{
				cr: DataSource(
					withExternalName(testDataSourceID),
					withSpec(v1alpha1.DataSourceParameters{
						DatasourcesV1Common: v1alpha1.DatasourcesV1Common{
							Category:    testCategory,
							Type:        testType,
							Enabled:     styraclient.Bool(false),
							OnPremises:  true,
							Description: &testDescription,
						},
						AWSECR: &v1alpha1.DatasourcesV1AWSECR{
							DatasourcesV1RateLimiter: v1alpha1.DatasourcesV1RateLimiter{
								RateLimit: testRateLimitQuantity,
							},
							DatasourcesV1Poller: v1alpha1.DatasourcesV1Poller{
								PollingInterval: testPollingIntervalDuration,
							},
							DatasourcesV1AWSCommon: v1alpha1.DatasourcesV1AWSCommon{
								Credentials: testCredentials,
								Region:      testRegion,
							},
							RegistryID: &testRegistryID,
						},
					}),
				),
				result: managed.ExternalCreation{},
			},
		},
		"AWS_ECR_CreateDataSourceFailed": {
			args: args{
				styra: styra.StyraAPI{
					Datasources: withMockDataSource(t, func(mcs *mockdatasource.MockClientService) {
						mcs.EXPECT().
							UpsertDatasource(&datasources.UpsertDatasourceParams{
								Body: &models.DatasourcesV1AWSECR{
									DatasourcesV1Common: models.DatasourcesV1Common{
										Category:    &testCategory,
										Type:        &testType,
										Enabled:     styraclient.Bool(false),
										OnPremises:  styraclient.Bool(true),
										Description: testDescription,
									},
									DatasourcesV1RateLimiter: models.DatasourcesV1RateLimiter{

										RateLimit: &testRateLimitNumber,
									},
									DatasourcesV1Poller: models.DatasourcesV1Poller{

										PollingInterval: &testPollingIntervalStr,
									},
									DatasourcesV1AWSCommon: models.DatasourcesV1AWSCommon{
										Credentials: &testCredentials,
										Region:      &testRegion,
									},
									RegistryID: testRegistryID,
								},
								Datasource: testDataSourceID,
								Context:    context.Background(),
							}).
							Return(nil, errBoom)
					}),
				},
				cr: DataSource(
					withExternalName(testDataSourceID),
					withSpec(v1alpha1.DataSourceParameters{
						DatasourcesV1Common: v1alpha1.DatasourcesV1Common{
							Category:    testCategory,
							Type:        testType,
							Enabled:     styraclient.Bool(false),
							OnPremises:  true,
							Description: &testDescription,
						},
						AWSECR: &v1alpha1.DatasourcesV1AWSECR{
							DatasourcesV1RateLimiter: v1alpha1.DatasourcesV1RateLimiter{
								RateLimit: testRateLimitQuantity,
							},
							DatasourcesV1Poller: v1alpha1.DatasourcesV1Poller{
								PollingInterval: testPollingIntervalDuration,
							},
							DatasourcesV1AWSCommon: v1alpha1.DatasourcesV1AWSCommon{
								Credentials: testCredentials,
								Region:      testRegion,
							},
							RegistryID: &testRegistryID,
						},
					}),
				),
			},
			want: want{
				cr: DataSource(
					withExternalName(testDataSourceID),
					withSpec(v1alpha1.DataSourceParameters{
						DatasourcesV1Common: v1alpha1.DatasourcesV1Common{
							Category:    testCategory,
							Type:        testType,
							Enabled:     styraclient.Bool(false),
							OnPremises:  true,
							Description: &testDescription,
						},
						AWSECR: &v1alpha1.DatasourcesV1AWSECR{
							DatasourcesV1RateLimiter: v1alpha1.DatasourcesV1RateLimiter{
								RateLimit: testRateLimitQuantity,
							},
							DatasourcesV1Poller: v1alpha1.DatasourcesV1Poller{
								PollingInterval: testPollingIntervalDuration,
							},
							DatasourcesV1AWSCommon: v1alpha1.DatasourcesV1AWSCommon{
								Credentials: testCredentials,
								Region:      testRegion,
							},
							RegistryID: &testRegistryID,
						},
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
		cr     *v1alpha1.DataSource
		result managed.ExternalUpdate
		err    error
	}

	cases := map[string]struct {
		args
		want
	}{
		"AWS_ECR_Successful": {
			args: args{
				styra: styra.StyraAPI{
					Datasources: withMockDataSource(t, func(mcs *mockdatasource.MockClientService) {
						mcs.EXPECT().
							UpsertDatasource(&datasources.UpsertDatasourceParams{
								Body: &models.DatasourcesV1AWSECR{
									DatasourcesV1Common: models.DatasourcesV1Common{
										Category:    &testCategory,
										Type:        &testType,
										Enabled:     styraclient.Bool(false),
										OnPremises:  styraclient.Bool(true),
										Description: testDescription,
									},
									DatasourcesV1RateLimiter: models.DatasourcesV1RateLimiter{

										RateLimit: &testRateLimitNumber,
									},
									DatasourcesV1Poller: models.DatasourcesV1Poller{

										PollingInterval: &testPollingIntervalStr,
									},
									DatasourcesV1AWSCommon: models.DatasourcesV1AWSCommon{
										Credentials: &testCredentials,
										Region:      &testRegion,
									},
									RegistryID: testRegistryID,
								},
								Datasource: testDataSourceID,
								Context:    context.Background(),
							}).
							Return(&datasources.UpsertDatasourceOK{
								Payload: &models.DatasourcesV1DatasourcesPutResponse{},
							}, nil)
					}),
				},
				cr: DataSource(
					withExternalName(testDataSourceID),
					withSpec(v1alpha1.DataSourceParameters{
						DatasourcesV1Common: v1alpha1.DatasourcesV1Common{
							Category:    testCategory,
							Type:        testType,
							Enabled:     styraclient.Bool(false),
							OnPremises:  true,
							Description: &testDescription,
						},
						AWSECR: &v1alpha1.DatasourcesV1AWSECR{
							DatasourcesV1RateLimiter: v1alpha1.DatasourcesV1RateLimiter{
								RateLimit: testRateLimitQuantity,
							},
							DatasourcesV1Poller: v1alpha1.DatasourcesV1Poller{
								PollingInterval: testPollingIntervalDuration,
							},
							DatasourcesV1AWSCommon: v1alpha1.DatasourcesV1AWSCommon{
								Credentials: testCredentials,
								Region:      testRegion,
							},
							RegistryID: &testRegistryID,
						},
					}),
				),
			},
			want: want{
				cr: DataSource(
					withExternalName(testDataSourceID),
					withSpec(v1alpha1.DataSourceParameters{
						DatasourcesV1Common: v1alpha1.DatasourcesV1Common{
							Category:    testCategory,
							Type:        testType,
							Enabled:     styraclient.Bool(false),
							OnPremises:  true,
							Description: &testDescription,
						},
						AWSECR: &v1alpha1.DatasourcesV1AWSECR{
							DatasourcesV1RateLimiter: v1alpha1.DatasourcesV1RateLimiter{
								RateLimit: testRateLimitQuantity,
							},
							DatasourcesV1Poller: v1alpha1.DatasourcesV1Poller{
								PollingInterval: testPollingIntervalDuration,
							},
							DatasourcesV1AWSCommon: v1alpha1.DatasourcesV1AWSCommon{
								Credentials: testCredentials,
								Region:      testRegion,
							},
							RegistryID: &testRegistryID,
						},
					}),
				),
				result: managed.ExternalUpdate{},
			},
		},
		"AWS_ECR_UpdateDataSourceFailed": {
			args: args{
				styra: styra.StyraAPI{
					Datasources: withMockDataSource(t, func(mcs *mockdatasource.MockClientService) {
						mcs.EXPECT().
							UpsertDatasource(&datasources.UpsertDatasourceParams{
								Body: &models.DatasourcesV1AWSECR{
									DatasourcesV1Common: models.DatasourcesV1Common{
										Category:    &testCategory,
										Type:        &testType,
										Enabled:     styraclient.Bool(false),
										OnPremises:  styraclient.Bool(true),
										Description: testDescription,
									},
									DatasourcesV1RateLimiter: models.DatasourcesV1RateLimiter{

										RateLimit: &testRateLimitNumber,
									},
									DatasourcesV1Poller: models.DatasourcesV1Poller{

										PollingInterval: &testPollingIntervalStr,
									},
									DatasourcesV1AWSCommon: models.DatasourcesV1AWSCommon{
										Credentials: &testCredentials,
										Region:      &testRegion,
									},
									RegistryID: testRegistryID,
								},
								Datasource: testDataSourceID,
								Context:    context.Background(),
							}).
							Return(nil, errBoom)
					}),
				},
				cr: DataSource(
					withExternalName(testDataSourceID),
					withSpec(v1alpha1.DataSourceParameters{
						DatasourcesV1Common: v1alpha1.DatasourcesV1Common{
							Category:    testCategory,
							Type:        testType,
							Enabled:     styraclient.Bool(false),
							OnPremises:  true,
							Description: &testDescription,
						},
						AWSECR: &v1alpha1.DatasourcesV1AWSECR{
							DatasourcesV1RateLimiter: v1alpha1.DatasourcesV1RateLimiter{
								RateLimit: testRateLimitQuantity,
							},
							DatasourcesV1Poller: v1alpha1.DatasourcesV1Poller{
								PollingInterval: testPollingIntervalDuration,
							},
							DatasourcesV1AWSCommon: v1alpha1.DatasourcesV1AWSCommon{
								Credentials: testCredentials,
								Region:      testRegion,
							},
							RegistryID: &testRegistryID,
						},
					}),
				),
			},
			want: want{
				cr: DataSource(
					withExternalName(testDataSourceID),
					withSpec(v1alpha1.DataSourceParameters{
						DatasourcesV1Common: v1alpha1.DatasourcesV1Common{
							Category:    testCategory,
							Type:        testType,
							Enabled:     styraclient.Bool(false),
							OnPremises:  true,
							Description: &testDescription,
						},
						AWSECR: &v1alpha1.DatasourcesV1AWSECR{
							DatasourcesV1RateLimiter: v1alpha1.DatasourcesV1RateLimiter{
								RateLimit: testRateLimitQuantity,
							},
							DatasourcesV1Poller: v1alpha1.DatasourcesV1Poller{
								PollingInterval: testPollingIntervalDuration,
							},
							DatasourcesV1AWSCommon: v1alpha1.DatasourcesV1AWSCommon{
								Credentials: testCredentials,
								Region:      testRegion,
							},
							RegistryID: &testRegistryID,
						},
					}),
				),
				err: errors.Wrap(errBoom, errUpdateFailed),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := &external{client: &tc.styra}
			o, err := e.Update(context.Background(), tc.args.cr)

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

func TestDelete(t *testing.T) {
	type want struct {
		cr  *v1alpha1.DataSource
		err error
	}

	cases := map[string]struct {
		args
		want
	}{
		"Successful": {
			args: args{
				styra: styra.StyraAPI{
					Datasources: withMockDataSource(t, func(mcs *mockdatasource.MockClientService) {
						mcs.EXPECT().
							DeleteDatasource(&datasources.DeleteDatasourceParams{
								Datasource: testDataSourceID,
								Context:    context.Background(),
							}).
							Return(&datasources.DeleteDatasourceOK{}, nil)
					}),
				},
				cr: DataSource(
					withExternalName(testDataSourceID),
				),
			},
			want: want{
				cr: DataSource(
					withExternalName(testDataSourceID),
				),
			},
		},
		"DeleteFailed": {
			args: args{
				styra: styra.StyraAPI{
					Datasources: withMockDataSource(t, func(mcs *mockdatasource.MockClientService) {
						mcs.EXPECT().
							DeleteDatasource(&datasources.DeleteDatasourceParams{
								Datasource: testDataSourceID,
								Context:    context.Background(),
							}).
							Return(nil, errBoom)
					}),
				},
				cr: DataSource(
					withExternalName(testDataSourceID),
				),
			},
			want: want{
				cr: DataSource(
					withExternalName(testDataSourceID),
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

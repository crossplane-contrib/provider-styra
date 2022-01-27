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

package v1alpha1

import (
	"testing"

	"github.com/crossplane/crossplane-runtime/pkg/test"
	"github.com/google/go-cmp/cmp"
)

type SystemModifier func(*System)

func withSpec(p SystemParameters) SystemModifier {
	return func(r *System) { r.Spec.ForProvider = p }
}

func getSystem(m ...SystemModifier) *System {
	cr := &System{}
	for _, f := range m {
		f(cr)
	}
	return cr
}

func TestHasAssets(t *testing.T) {
	type want struct {
		hasAssets bool
	}

	type args struct {
		cr *System
	}

	cases := map[string]struct {
		args
		want
	}{
		"UnsupportedSystem": {
			args: args{
				cr: getSystem(
					withSpec(SystemParameters{
						Type: "fooType",
					}),
				),
			},
			want: want{
				false,
			},
		},
		"KubernetesSystemV123": {
			args: args{
				cr: getSystem(
					withSpec(SystemParameters{
						Type: "kubernetes:v123",
					}),
				),
			},
			want: want{
				true,
			},
		},
		"OPASystem": {
			args: args{
				cr: getSystem(
					withSpec(SystemParameters{
						Type: "custom",
					}),
				),
			},
			want: want{
				true,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			actlHasAssets := tc.args.cr.Spec.ForProvider.HasAssets()

			if diff := cmp.Diff(tc.want.hasAssets, actlHasAssets, test.EquateConditions()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
		})
	}
}

func TestGetAssetTypes(t *testing.T) {
	type want struct {
		assetTypes []string
	}

	type args struct {
		cr *System
	}

	cases := map[string]struct {
		args
		want
	}{
		"UnsupportedSystem": {
			args: args{
				cr: getSystem(
					withSpec(SystemParameters{
						Type: "fooType",
					}),
				),
			},
			want: want{
				[]string{},
			},
		},
		"KubernetesSystem": {
			args: args{
				cr: getSystem(
					withSpec(SystemParameters{
						Type: "kubernetes:v123",
					}),
				),
			},
			want: want{
				[]string{"helm-values"},
			},
		},
		"OPASystem": {
			args: args{
				cr: getSystem(
					withSpec(SystemParameters{
						Type: "custom",
					}),
				),
			},
			want: want{
				[]string{"opa-config"},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			actlDetails := tc.args.cr.Spec.ForProvider.GetAssetTypes()

			if diff := cmp.Diff(tc.want.assetTypes, actlDetails, test.EquateConditions()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
		})
	}
}

func TestHasLabels(t *testing.T) {
	type want struct {
		hasLabels bool
	}

	type args struct {
		cr *System
	}

	cases := map[string]struct {
		args
		want
	}{
		"UnsupportedSystem": {
			args: args{
				cr: getSystem(
					withSpec(SystemParameters{
						Type: "fooType",
					}),
				),
			},
			want: want{
				false,
			},
		},
		"KubernetesSystemV2": {
			args: args{
				cr: getSystem(
					withSpec(SystemParameters{
						Type: "kubernetes:v2",
					}),
				),
			},
			want: want{
				true,
			},
		},
		"KubernetesSystemV123": {
			args: args{
				cr: getSystem(
					withSpec(SystemParameters{
						Type: "kubernetes:v123",
					}),
				),
			},
			want: want{
				false,
			},
		},
		"OPASystem": {
			args: args{
				cr: getSystem(
					withSpec(SystemParameters{
						Type: "custom",
					}),
				),
			},
			want: want{
				false,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			actlHasLabels := tc.args.cr.Spec.ForProvider.HasLabels()

			if diff := cmp.Diff(tc.want.hasLabels, actlHasLabels, test.EquateConditions()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
		})
	}
}

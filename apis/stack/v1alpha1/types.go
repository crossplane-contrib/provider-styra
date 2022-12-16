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

// V1SourceControlConfig v1 source control config
type V1SourceControlConfig struct {

	// origin
	// +kubebuilder:validation:Required
	Origin V1GitRepoConfig `json:"origin"`
}

// V1GitRepoConfig v1 git repo config
type V1GitRepoConfig struct {

	// Credentials are looked under the key <name>/<creds>
	// +kubebuilder:validation:Required
	Credentials string `json:"credentials"`

	// Path to limit the import to
	// +kubebuilder:validation:Required
	Path string `json:"path"`

	// Remote reference, defaults to refs/heads/master
	// +kubebuilder:validation:Required
	Reference string `json:"reference"`

	// Repository URL
	// +kubebuilder:validation:Required
	URL string `json:"url"`
}

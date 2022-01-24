//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"github.com/crossplane/crossplane-runtime/apis/common/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DataSource) DeepCopyInto(out *DataSource) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DataSource.
func (in *DataSource) DeepCopy() *DataSource {
	if in == nil {
		return nil
	}
	out := new(DataSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DataSource) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DataSourceExternalStatus) DeepCopyInto(out *DataSourceExternalStatus) {
	*out = *in
	if in.Code != nil {
		in, out := &in.Code, &out.Code
		*out = new(string)
		**out = **in
	}
	if in.Message != nil {
		in, out := &in.Message, &out.Message
		*out = new(string)
		**out = **in
	}
	if in.Timestamp != nil {
		in, out := &in.Timestamp, &out.Timestamp
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DataSourceExternalStatus.
func (in *DataSourceExternalStatus) DeepCopy() *DataSourceExternalStatus {
	if in == nil {
		return nil
	}
	out := new(DataSourceExternalStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DataSourceList) DeepCopyInto(out *DataSourceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DataSource, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DataSourceList.
func (in *DataSourceList) DeepCopy() *DataSourceList {
	if in == nil {
		return nil
	}
	out := new(DataSourceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DataSourceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DataSourceObservation) DeepCopyInto(out *DataSourceObservation) {
	*out = *in
	if in.Executed != nil {
		in, out := &in.Executed, &out.Executed
		*out = new(string)
		**out = **in
	}
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(DataSourceExternalStatus)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DataSourceObservation.
func (in *DataSourceObservation) DeepCopy() *DataSourceObservation {
	if in == nil {
		return nil
	}
	out := new(DataSourceObservation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DataSourceParameters) DeepCopyInto(out *DataSourceParameters) {
	*out = *in
	in.DatasourcesV1Common.DeepCopyInto(&out.DatasourcesV1Common)
	if in.AWSECR != nil {
		in, out := &in.AWSECR, &out.AWSECR
		*out = new(DatasourcesV1AWSECR)
		(*in).DeepCopyInto(*out)
	}
	if in.BundleS3 != nil {
		in, out := &in.BundleS3, &out.BundleS3
		*out = new(DatasourcesV1BundleS3)
		(*in).DeepCopyInto(*out)
	}
	if in.GitBlame != nil {
		in, out := &in.GitBlame, &out.GitBlame
		*out = new(DatasourcesV1GitBlame)
		(*in).DeepCopyInto(*out)
	}
	if in.GitContent != nil {
		in, out := &in.GitContent, &out.GitContent
		*out = new(DatasourcesV1GitContent)
		(*in).DeepCopyInto(*out)
	}
	if in.GitRego != nil {
		in, out := &in.GitRego, &out.GitRego
		*out = new(DatasourcesV1GitRego)
		(*in).DeepCopyInto(*out)
	}
	if in.HTTP != nil {
		in, out := &in.HTTP, &out.HTTP
		*out = new(DatasourcesV1HTTP)
		(*in).DeepCopyInto(*out)
	}
	if in.KubernetesResources != nil {
		in, out := &in.KubernetesResources, &out.KubernetesResources
		*out = new(DatasourcesV1KubernetesResources)
		(*in).DeepCopyInto(*out)
	}
	if in.LDAP != nil {
		in, out := &in.LDAP, &out.LDAP
		*out = new(DatasourcesV1LDAP)
		(*in).DeepCopyInto(*out)
	}
	if in.PolicyLibrary != nil {
		in, out := &in.PolicyLibrary, &out.PolicyLibrary
		*out = new(DatasourcesV1PolicyLibrary)
		(*in).DeepCopyInto(*out)
	}
	if in.Rest != nil {
		in, out := &in.Rest, &out.Rest
		*out = new(DatasourcesV1Rest)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DataSourceParameters.
func (in *DataSourceParameters) DeepCopy() *DataSourceParameters {
	if in == nil {
		return nil
	}
	out := new(DataSourceParameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DataSourceSpec) DeepCopyInto(out *DataSourceSpec) {
	*out = *in
	in.ResourceSpec.DeepCopyInto(&out.ResourceSpec)
	in.ForProvider.DeepCopyInto(&out.ForProvider)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DataSourceSpec.
func (in *DataSourceSpec) DeepCopy() *DataSourceSpec {
	if in == nil {
		return nil
	}
	out := new(DataSourceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DataSourceStatus) DeepCopyInto(out *DataSourceStatus) {
	*out = *in
	in.ResourceStatus.DeepCopyInto(&out.ResourceStatus)
	in.AtProvider.DeepCopyInto(&out.AtProvider)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DataSourceStatus.
func (in *DataSourceStatus) DeepCopy() *DataSourceStatus {
	if in == nil {
		return nil
	}
	out := new(DataSourceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatasourcesV1AWSCommon) DeepCopyInto(out *DatasourcesV1AWSCommon) {
	*out = *in
	if in.CredentialsRef != nil {
		in, out := &in.CredentialsRef, &out.CredentialsRef
		*out = new(v1.Reference)
		**out = **in
	}
	if in.CredentialsSelector != nil {
		in, out := &in.CredentialsSelector, &out.CredentialsSelector
		*out = new(v1.Selector)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatasourcesV1AWSCommon.
func (in *DatasourcesV1AWSCommon) DeepCopy() *DatasourcesV1AWSCommon {
	if in == nil {
		return nil
	}
	out := new(DatasourcesV1AWSCommon)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatasourcesV1AWSECR) DeepCopyInto(out *DatasourcesV1AWSECR) {
	*out = *in
	in.DatasourcesV1RateLimiter.DeepCopyInto(&out.DatasourcesV1RateLimiter)
	in.DatasourcesV1Poller.DeepCopyInto(&out.DatasourcesV1Poller)
	in.DatasourcesV1AWSCommon.DeepCopyInto(&out.DatasourcesV1AWSCommon)
	if in.RegistryID != nil {
		in, out := &in.RegistryID, &out.RegistryID
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatasourcesV1AWSECR.
func (in *DatasourcesV1AWSECR) DeepCopy() *DatasourcesV1AWSECR {
	if in == nil {
		return nil
	}
	out := new(DatasourcesV1AWSECR)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatasourcesV1BundleS3) DeepCopyInto(out *DatasourcesV1BundleS3) {
	*out = *in
	in.DatasourcesV1Poller.DeepCopyInto(&out.DatasourcesV1Poller)
	in.DatasourcesV1AWSCommon.DeepCopyInto(&out.DatasourcesV1AWSCommon)
	if in.Endpoint != nil {
		in, out := &in.Endpoint, &out.Endpoint
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatasourcesV1BundleS3.
func (in *DatasourcesV1BundleS3) DeepCopy() *DatasourcesV1BundleS3 {
	if in == nil {
		return nil
	}
	out := new(DatasourcesV1BundleS3)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatasourcesV1Common) DeepCopyInto(out *DatasourcesV1Common) {
	*out = *in
	if in.Description != nil {
		in, out := &in.Description, &out.Description
		*out = new(string)
		**out = **in
	}
	if in.Enabled != nil {
		in, out := &in.Enabled, &out.Enabled
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatasourcesV1Common.
func (in *DatasourcesV1Common) DeepCopy() *DatasourcesV1Common {
	if in == nil {
		return nil
	}
	out := new(DatasourcesV1Common)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatasourcesV1GitBlame) DeepCopyInto(out *DatasourcesV1GitBlame) {
	*out = *in
	in.DatasourcesV1GitCommon.DeepCopyInto(&out.DatasourcesV1GitCommon)
	if in.PathRegexp != nil {
		in, out := &in.PathRegexp, &out.PathRegexp
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatasourcesV1GitBlame.
func (in *DatasourcesV1GitBlame) DeepCopy() *DatasourcesV1GitBlame {
	if in == nil {
		return nil
	}
	out := new(DatasourcesV1GitBlame)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatasourcesV1GitCommon) DeepCopyInto(out *DatasourcesV1GitCommon) {
	*out = *in
	in.DatasourcesV1Poller.DeepCopyInto(&out.DatasourcesV1Poller)
	in.DatasourcesV1RateLimiter.DeepCopyInto(&out.DatasourcesV1RateLimiter)
	if in.Credentials != nil {
		in, out := &in.Credentials, &out.Credentials
		*out = new(string)
		**out = **in
	}
	if in.CredentialsRef != nil {
		in, out := &in.CredentialsRef, &out.CredentialsRef
		*out = new(v1.Reference)
		**out = **in
	}
	if in.CredentialsSelector != nil {
		in, out := &in.CredentialsSelector, &out.CredentialsSelector
		*out = new(v1.Selector)
		(*in).DeepCopyInto(*out)
	}
	if in.Reference != nil {
		in, out := &in.Reference, &out.Reference
		*out = new(string)
		**out = **in
	}
	if in.SSHCredentials != nil {
		in, out := &in.SSHCredentials, &out.SSHCredentials
		*out = new(DatasourcesV1GitCommonAO3SSHCredentials)
		(*in).DeepCopyInto(*out)
	}
	if in.Timeout != nil {
		in, out := &in.Timeout, &out.Timeout
		*out = new(metav1.Duration)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatasourcesV1GitCommon.
func (in *DatasourcesV1GitCommon) DeepCopy() *DatasourcesV1GitCommon {
	if in == nil {
		return nil
	}
	out := new(DatasourcesV1GitCommon)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatasourcesV1GitCommonAO3SSHCredentials) DeepCopyInto(out *DatasourcesV1GitCommonAO3SSHCredentials) {
	*out = *in
	if in.Passphrase != nil {
		in, out := &in.Passphrase, &out.Passphrase
		*out = new(string)
		**out = **in
	}
	if in.PassphraseRef != nil {
		in, out := &in.PassphraseRef, &out.PassphraseRef
		*out = new(v1.Reference)
		**out = **in
	}
	if in.PassphraseSelector != nil {
		in, out := &in.PassphraseSelector, &out.PassphraseSelector
		*out = new(v1.Selector)
		(*in).DeepCopyInto(*out)
	}
	if in.PrivateKeyRef != nil {
		in, out := &in.PrivateKeyRef, &out.PrivateKeyRef
		*out = new(v1.Reference)
		**out = **in
	}
	if in.PrivateKeySelector != nil {
		in, out := &in.PrivateKeySelector, &out.PrivateKeySelector
		*out = new(v1.Selector)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatasourcesV1GitCommonAO3SSHCredentials.
func (in *DatasourcesV1GitCommonAO3SSHCredentials) DeepCopy() *DatasourcesV1GitCommonAO3SSHCredentials {
	if in == nil {
		return nil
	}
	out := new(DatasourcesV1GitCommonAO3SSHCredentials)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatasourcesV1GitContent) DeepCopyInto(out *DatasourcesV1GitContent) {
	*out = *in
	in.DatasourcesV1GitCommon.DeepCopyInto(&out.DatasourcesV1GitCommon)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatasourcesV1GitContent.
func (in *DatasourcesV1GitContent) DeepCopy() *DatasourcesV1GitContent {
	if in == nil {
		return nil
	}
	out := new(DatasourcesV1GitContent)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatasourcesV1GitRego) DeepCopyInto(out *DatasourcesV1GitRego) {
	*out = *in
	in.DatasourcesV1GitCommon.DeepCopyInto(&out.DatasourcesV1GitCommon)
	if in.Path != nil {
		in, out := &in.Path, &out.Path
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatasourcesV1GitRego.
func (in *DatasourcesV1GitRego) DeepCopy() *DatasourcesV1GitRego {
	if in == nil {
		return nil
	}
	out := new(DatasourcesV1GitRego)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatasourcesV1HTTP) DeepCopyInto(out *DatasourcesV1HTTP) {
	*out = *in
	in.DatasourcesV1Poller.DeepCopyInto(&out.DatasourcesV1Poller)
	in.DatasourcesV1RegoFiltering.DeepCopyInto(&out.DatasourcesV1RegoFiltering)
	in.DatasourcesV1TLSSettings.DeepCopyInto(&out.DatasourcesV1TLSSettings)
	if in.Headers != nil {
		in, out := &in.Headers, &out.Headers
		*out = make([]DatasourcesV1HTTPHeader, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatasourcesV1HTTP.
func (in *DatasourcesV1HTTP) DeepCopy() *DatasourcesV1HTTP {
	if in == nil {
		return nil
	}
	out := new(DatasourcesV1HTTP)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatasourcesV1HTTPHeader) DeepCopyInto(out *DatasourcesV1HTTPHeader) {
	*out = *in
	if in.SecretID != nil {
		in, out := &in.SecretID, &out.SecretID
		*out = new(string)
		**out = **in
	}
	if in.Value != nil {
		in, out := &in.Value, &out.Value
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatasourcesV1HTTPHeader.
func (in *DatasourcesV1HTTPHeader) DeepCopy() *DatasourcesV1HTTPHeader {
	if in == nil {
		return nil
	}
	out := new(DatasourcesV1HTTPHeader)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatasourcesV1KubernetesResources) DeepCopyInto(out *DatasourcesV1KubernetesResources) {
	*out = *in
	in.DatasourcesV1RateLimiter.DeepCopyInto(&out.DatasourcesV1RateLimiter)
	in.DatasourcesV1Poller.DeepCopyInto(&out.DatasourcesV1Poller)
	if in.Masks != nil {
		in, out := &in.Masks, &out.Masks
		*out = make(map[string][]string, len(*in))
		for key, val := range *in {
			var outVal []string
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = make([]string, len(*in))
				copy(*out, *in)
			}
			(*out)[key] = outVal
		}
	}
	if in.Namespaces != nil {
		in, out := &in.Namespaces, &out.Namespaces
		*out = make(map[string]bool, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Selectors != nil {
		in, out := &in.Selectors, &out.Selectors
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatasourcesV1KubernetesResources.
func (in *DatasourcesV1KubernetesResources) DeepCopy() *DatasourcesV1KubernetesResources {
	if in == nil {
		return nil
	}
	out := new(DatasourcesV1KubernetesResources)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatasourcesV1LDAP) DeepCopyInto(out *DatasourcesV1LDAP) {
	*out = *in
	in.DatasourcesV1RateLimiter.DeepCopyInto(&out.DatasourcesV1RateLimiter)
	in.DatasourcesV1Poller.DeepCopyInto(&out.DatasourcesV1Poller)
	in.DatasourcesV1RegoFiltering.DeepCopyInto(&out.DatasourcesV1RegoFiltering)
	in.DatasourcesV1TLSSettings.DeepCopyInto(&out.DatasourcesV1TLSSettings)
	if in.Credentials != nil {
		in, out := &in.Credentials, &out.Credentials
		*out = new(string)
		**out = **in
	}
	if in.CredentialsRef != nil {
		in, out := &in.CredentialsRef, &out.CredentialsRef
		*out = new(v1.Reference)
		**out = **in
	}
	if in.CredentialsSelector != nil {
		in, out := &in.CredentialsSelector, &out.CredentialsSelector
		*out = new(v1.Selector)
		(*in).DeepCopyInto(*out)
	}
	if in.Search != nil {
		in, out := &in.Search, &out.Search
		*out = new(DatasourcesV1LDAPAO5Search)
		(*in).DeepCopyInto(*out)
	}
	if in.Urls != nil {
		in, out := &in.Urls, &out.Urls
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatasourcesV1LDAP.
func (in *DatasourcesV1LDAP) DeepCopy() *DatasourcesV1LDAP {
	if in == nil {
		return nil
	}
	out := new(DatasourcesV1LDAP)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatasourcesV1LDAPAO5Search) DeepCopyInto(out *DatasourcesV1LDAPAO5Search) {
	*out = *in
	if in.Attributes != nil {
		in, out := &in.Attributes, &out.Attributes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Deref != nil {
		in, out := &in.Deref, &out.Deref
		*out = new(string)
		**out = **in
	}
	if in.PageSize != nil {
		in, out := &in.PageSize, &out.PageSize
		*out = new(int64)
		**out = **in
	}
	if in.Scope != nil {
		in, out := &in.Scope, &out.Scope
		*out = new(string)
		**out = **in
	}
	if in.SizeLimit != nil {
		in, out := &in.SizeLimit, &out.SizeLimit
		*out = new(int64)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatasourcesV1LDAPAO5Search.
func (in *DatasourcesV1LDAPAO5Search) DeepCopy() *DatasourcesV1LDAPAO5Search {
	if in == nil {
		return nil
	}
	out := new(DatasourcesV1LDAPAO5Search)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatasourcesV1PolicyLibrary) DeepCopyInto(out *DatasourcesV1PolicyLibrary) {
	*out = *in
	in.DatasourcesV1Poller.DeepCopyInto(&out.DatasourcesV1Poller)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatasourcesV1PolicyLibrary.
func (in *DatasourcesV1PolicyLibrary) DeepCopy() *DatasourcesV1PolicyLibrary {
	if in == nil {
		return nil
	}
	out := new(DatasourcesV1PolicyLibrary)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatasourcesV1Poller) DeepCopyInto(out *DatasourcesV1Poller) {
	*out = *in
	if in.PollingInterval != nil {
		in, out := &in.PollingInterval, &out.PollingInterval
		*out = new(metav1.Duration)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatasourcesV1Poller.
func (in *DatasourcesV1Poller) DeepCopy() *DatasourcesV1Poller {
	if in == nil {
		return nil
	}
	out := new(DatasourcesV1Poller)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatasourcesV1RateLimiter) DeepCopyInto(out *DatasourcesV1RateLimiter) {
	*out = *in
	if in.RateLimit != nil {
		in, out := &in.RateLimit, &out.RateLimit
		x := (*in).DeepCopy()
		*out = &x
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatasourcesV1RateLimiter.
func (in *DatasourcesV1RateLimiter) DeepCopy() *DatasourcesV1RateLimiter {
	if in == nil {
		return nil
	}
	out := new(DatasourcesV1RateLimiter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatasourcesV1RegoFiltering) DeepCopyInto(out *DatasourcesV1RegoFiltering) {
	*out = *in
	if in.PolicyFilter != nil {
		in, out := &in.PolicyFilter, &out.PolicyFilter
		*out = new(string)
		**out = **in
	}
	if in.PolicyQuery != nil {
		in, out := &in.PolicyQuery, &out.PolicyQuery
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatasourcesV1RegoFiltering.
func (in *DatasourcesV1RegoFiltering) DeepCopy() *DatasourcesV1RegoFiltering {
	if in == nil {
		return nil
	}
	out := new(DatasourcesV1RegoFiltering)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatasourcesV1Rest) DeepCopyInto(out *DatasourcesV1Rest) {
	*out = *in
	if in.ContentType != nil {
		in, out := &in.ContentType, &out.ContentType
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatasourcesV1Rest.
func (in *DatasourcesV1Rest) DeepCopy() *DatasourcesV1Rest {
	if in == nil {
		return nil
	}
	out := new(DatasourcesV1Rest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatasourcesV1TLSSettings) DeepCopyInto(out *DatasourcesV1TLSSettings) {
	*out = *in
	if in.CaCertificate != nil {
		in, out := &in.CaCertificate, &out.CaCertificate
		*out = new(string)
		**out = **in
	}
	if in.SkipTLSVerification != nil {
		in, out := &in.SkipTLSVerification, &out.SkipTLSVerification
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatasourcesV1TLSSettings.
func (in *DatasourcesV1TLSSettings) DeepCopy() *DatasourcesV1TLSSettings {
	if in == nil {
		return nil
	}
	out := new(DatasourcesV1TLSSettings)
	in.DeepCopyInto(out)
	return out
}

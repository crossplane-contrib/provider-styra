// +build generate

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

//go:generate go run -tags generate github.com/golang/mock/mockgen -package kube -destination ./kube/mock.go sigs.k8s.io/controller-runtime/pkg/client Client

//go:generate go run -tags generate github.com/golang/mock/mockgen -package resource -destination ./resource/mock.go github.com/crossplane/crossplane-runtime/pkg/resource Applicator

//go:generate go run -tags generate github.com/golang/mock/mockgen -package policies -destination ./policies/mock.go github.com/mistermx/styra-go-client/pkg/client/policies ClientService

//go:generate go run -tags generate github.com/golang/mock/mockgen -package secrets -destination ./secrets/mock.go github.com/mistermx/styra-go-client/pkg/client/secrets ClientService

//go:generate go run -tags generate github.com/golang/mock/mockgen -package stacks -destination ./stacks/mock.go github.com/mistermx/styra-go-client/pkg/client/stacks ClientService

//go:generate go run -tags generate github.com/golang/mock/mockgen -package systems -destination ./systems/mock.go github.com/mistermx/styra-go-client/pkg/client/systems ClientService

package mock

import (
	// Workaround to vendor mockgen (https://github.com/golang/mock/issues/415#issuecomment-602547154)
	_ "github.com/golang/mock/mockgen" //nolint:typecheck
)

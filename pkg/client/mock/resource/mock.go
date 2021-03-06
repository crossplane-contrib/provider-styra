// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/crossplane/crossplane-runtime/pkg/resource (interfaces: Applicator)

// Package resource is a generated GoMock package.
package resource

import (
	context "context"
	reflect "reflect"

	resource "github.com/crossplane/crossplane-runtime/pkg/resource"
	gomock "github.com/golang/mock/gomock"
	client "sigs.k8s.io/controller-runtime/pkg/client"
)

// MockApplicator is a mock of Applicator interface.
type MockApplicator struct {
	ctrl     *gomock.Controller
	recorder *MockApplicatorMockRecorder
}

// MockApplicatorMockRecorder is the mock recorder for MockApplicator.
type MockApplicatorMockRecorder struct {
	mock *MockApplicator
}

// NewMockApplicator creates a new mock instance.
func NewMockApplicator(ctrl *gomock.Controller) *MockApplicator {
	mock := &MockApplicator{ctrl: ctrl}
	mock.recorder = &MockApplicatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApplicator) EXPECT() *MockApplicatorMockRecorder {
	return m.recorder
}

// Apply mocks base method.
func (m *MockApplicator) Apply(arg0 context.Context, arg1 client.Object, arg2 ...resource.ApplyOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Apply", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Apply indicates an expected call of Apply.
func (mr *MockApplicatorMockRecorder) Apply(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Apply", reflect.TypeOf((*MockApplicator)(nil).Apply), varargs...)
}

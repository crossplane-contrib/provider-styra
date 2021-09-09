// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mistermx/styra-go-client/pkg/client/systems (interfaces: ClientService)

// Package systems is a generated GoMock package.
package systems

import (
	io "io"
	reflect "reflect"

	runtime "github.com/go-openapi/runtime"
	gomock "github.com/golang/mock/gomock"
	systems "github.com/mistermx/styra-go-client/pkg/client/systems"
)

// MockClientService is a mock of ClientService interface.
type MockClientService struct {
	ctrl     *gomock.Controller
	recorder *MockClientServiceMockRecorder
}

// MockClientServiceMockRecorder is the mock recorder for MockClientService.
type MockClientServiceMockRecorder struct {
	mock *MockClientService
}

// NewMockClientService creates a new mock instance.
func NewMockClientService(ctrl *gomock.Controller) *MockClientService {
	mock := &MockClientService{ctrl: ctrl}
	mock.recorder = &MockClientServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClientService) EXPECT() *MockClientServiceMockRecorder {
	return m.recorder
}

// CommitFilesToSourceControlSystem mocks base method.
func (m *MockClientService) CommitFilesToSourceControlSystem(arg0 *systems.CommitFilesToSourceControlSystemParams, arg1 ...systems.ClientOption) (*systems.CommitFilesToSourceControlSystemOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CommitFilesToSourceControlSystem", varargs...)
	ret0, _ := ret[0].(*systems.CommitFilesToSourceControlSystemOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CommitFilesToSourceControlSystem indicates an expected call of CommitFilesToSourceControlSystem.
func (mr *MockClientServiceMockRecorder) CommitFilesToSourceControlSystem(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CommitFilesToSourceControlSystem", reflect.TypeOf((*MockClientService)(nil).CommitFilesToSourceControlSystem), varargs...)
}

// CreateSystem mocks base method.
func (m *MockClientService) CreateSystem(arg0 *systems.CreateSystemParams, arg1 ...systems.ClientOption) (*systems.CreateSystemOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateSystem", varargs...)
	ret0, _ := ret[0].(*systems.CreateSystemOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSystem indicates an expected call of CreateSystem.
func (mr *MockClientServiceMockRecorder) CreateSystem(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSystem", reflect.TypeOf((*MockClientService)(nil).CreateSystem), varargs...)
}

// DeleteSystem mocks base method.
func (m *MockClientService) DeleteSystem(arg0 *systems.DeleteSystemParams, arg1 ...systems.ClientOption) (*systems.DeleteSystemOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteSystem", varargs...)
	ret0, _ := ret[0].(*systems.DeleteSystemOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteSystem indicates an expected call of DeleteSystem.
func (mr *MockClientServiceMockRecorder) DeleteSystem(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSystem", reflect.TypeOf((*MockClientService)(nil).DeleteSystem), varargs...)
}

// DeleteUserBranchSystem mocks base method.
func (m *MockClientService) DeleteUserBranchSystem(arg0 *systems.DeleteUserBranchSystemParams, arg1 ...systems.ClientOption) (*systems.DeleteUserBranchSystemOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteUserBranchSystem", varargs...)
	ret0, _ := ret[0].(*systems.DeleteUserBranchSystemOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteUserBranchSystem indicates an expected call of DeleteUserBranchSystem.
func (mr *MockClientServiceMockRecorder) DeleteUserBranchSystem(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUserBranchSystem", reflect.TypeOf((*MockClientService)(nil).DeleteUserBranchSystem), varargs...)
}

// GetAsset mocks base method.
func (m *MockClientService) GetAsset(arg0 *systems.GetAssetParams, arg1 io.Writer, arg2 ...systems.ClientOption) (*systems.GetAssetOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetAsset", varargs...)
	ret0, _ := ret[0].(*systems.GetAssetOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAsset indicates an expected call of GetAsset.
func (mr *MockClientServiceMockRecorder) GetAsset(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAsset", reflect.TypeOf((*MockClientService)(nil).GetAsset), varargs...)
}

// GetInstructions mocks base method.
func (m *MockClientService) GetInstructions(arg0 *systems.GetInstructionsParams, arg1 ...systems.ClientOption) (*systems.GetInstructionsOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetInstructions", varargs...)
	ret0, _ := ret[0].(*systems.GetInstructionsOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInstructions indicates an expected call of GetInstructions.
func (mr *MockClientServiceMockRecorder) GetInstructions(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInstructions", reflect.TypeOf((*MockClientService)(nil).GetInstructions), varargs...)
}

// GetOPADiscoveryConfig mocks base method.
func (m *MockClientService) GetOPADiscoveryConfig(arg0 *systems.GetOPADiscoveryConfigParams, arg1 ...systems.ClientOption) (*systems.GetOPADiscoveryConfigOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetOPADiscoveryConfig", varargs...)
	ret0, _ := ret[0].(*systems.GetOPADiscoveryConfigOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOPADiscoveryConfig indicates an expected call of GetOPADiscoveryConfig.
func (mr *MockClientServiceMockRecorder) GetOPADiscoveryConfig(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOPADiscoveryConfig", reflect.TypeOf((*MockClientService)(nil).GetOPADiscoveryConfig), varargs...)
}

// GetSourceControlFilesBranchSystem mocks base method.
func (m *MockClientService) GetSourceControlFilesBranchSystem(arg0 *systems.GetSourceControlFilesBranchSystemParams, arg1 ...systems.ClientOption) (*systems.GetSourceControlFilesBranchSystemOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetSourceControlFilesBranchSystem", varargs...)
	ret0, _ := ret[0].(*systems.GetSourceControlFilesBranchSystemOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSourceControlFilesBranchSystem indicates an expected call of GetSourceControlFilesBranchSystem.
func (mr *MockClientServiceMockRecorder) GetSourceControlFilesBranchSystem(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSourceControlFilesBranchSystem", reflect.TypeOf((*MockClientService)(nil).GetSourceControlFilesBranchSystem), varargs...)
}

// GetSourceControlFilesMasterSystem mocks base method.
func (m *MockClientService) GetSourceControlFilesMasterSystem(arg0 *systems.GetSourceControlFilesMasterSystemParams, arg1 ...systems.ClientOption) (*systems.GetSourceControlFilesMasterSystemOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetSourceControlFilesMasterSystem", varargs...)
	ret0, _ := ret[0].(*systems.GetSourceControlFilesMasterSystemOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSourceControlFilesMasterSystem indicates an expected call of GetSourceControlFilesMasterSystem.
func (mr *MockClientServiceMockRecorder) GetSourceControlFilesMasterSystem(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSourceControlFilesMasterSystem", reflect.TypeOf((*MockClientService)(nil).GetSourceControlFilesMasterSystem), varargs...)
}

// GetSystem mocks base method.
func (m *MockClientService) GetSystem(arg0 *systems.GetSystemParams, arg1 ...systems.ClientOption) (*systems.GetSystemOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetSystem", varargs...)
	ret0, _ := ret[0].(*systems.GetSystemOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSystem indicates an expected call of GetSystem.
func (mr *MockClientServiceMockRecorder) GetSystem(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSystem", reflect.TypeOf((*MockClientService)(nil).GetSystem), varargs...)
}

// GetSystemAgents mocks base method.
func (m *MockClientService) GetSystemAgents(arg0 *systems.GetSystemAgentsParams, arg1 ...systems.ClientOption) (*systems.GetSystemAgentsOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetSystemAgents", varargs...)
	ret0, _ := ret[0].(*systems.GetSystemAgentsOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSystemAgents indicates an expected call of GetSystemAgents.
func (mr *MockClientServiceMockRecorder) GetSystemAgents(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSystemAgents", reflect.TypeOf((*MockClientService)(nil).GetSystemAgents), varargs...)
}

// GetSystemBundle mocks base method.
func (m *MockClientService) GetSystemBundle(arg0 *systems.GetSystemBundleParams, arg1 ...systems.ClientOption) (*systems.GetSystemBundleOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetSystemBundle", varargs...)
	ret0, _ := ret[0].(*systems.GetSystemBundleOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSystemBundle indicates an expected call of GetSystemBundle.
func (mr *MockClientServiceMockRecorder) GetSystemBundle(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSystemBundle", reflect.TypeOf((*MockClientService)(nil).GetSystemBundle), varargs...)
}

// GetSystemBundleDeploy mocks base method.
func (m *MockClientService) GetSystemBundleDeploy(arg0 *systems.GetSystemBundleDeployParams, arg1 ...systems.ClientOption) (*systems.GetSystemBundleDeployOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetSystemBundleDeploy", varargs...)
	ret0, _ := ret[0].(*systems.GetSystemBundleDeployOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSystemBundleDeploy indicates an expected call of GetSystemBundleDeploy.
func (mr *MockClientServiceMockRecorder) GetSystemBundleDeploy(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSystemBundleDeploy", reflect.TypeOf((*MockClientService)(nil).GetSystemBundleDeploy), varargs...)
}

// GetSystemBundleDetails mocks base method.
func (m *MockClientService) GetSystemBundleDetails(arg0 *systems.GetSystemBundleDetailsParams, arg1 ...systems.ClientOption) (*systems.GetSystemBundleDetailsOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetSystemBundleDetails", varargs...)
	ret0, _ := ret[0].(*systems.GetSystemBundleDetailsOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSystemBundleDetails indicates an expected call of GetSystemBundleDetails.
func (mr *MockClientServiceMockRecorder) GetSystemBundleDetails(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSystemBundleDetails", reflect.TypeOf((*MockClientService)(nil).GetSystemBundleDetails), varargs...)
}

// GetSystemBundles mocks base method.
func (m *MockClientService) GetSystemBundles(arg0 *systems.GetSystemBundlesParams, arg1 ...systems.ClientOption) (*systems.GetSystemBundlesOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetSystemBundles", varargs...)
	ret0, _ := ret[0].(*systems.GetSystemBundlesOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSystemBundles indicates an expected call of GetSystemBundles.
func (mr *MockClientServiceMockRecorder) GetSystemBundles(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSystemBundles", reflect.TypeOf((*MockClientService)(nil).GetSystemBundles), varargs...)
}

// HandleSystemMetrics mocks base method.
func (m *MockClientService) HandleSystemMetrics(arg0 *systems.HandleSystemMetricsParams, arg1 ...systems.ClientOption) (*systems.HandleSystemMetricsOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "HandleSystemMetrics", varargs...)
	ret0, _ := ret[0].(*systems.HandleSystemMetricsOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HandleSystemMetrics indicates an expected call of HandleSystemMetrics.
func (mr *MockClientServiceMockRecorder) HandleSystemMetrics(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleSystemMetrics", reflect.TypeOf((*MockClientService)(nil).HandleSystemMetrics), varargs...)
}

// ListSystems mocks base method.
func (m *MockClientService) ListSystems(arg0 *systems.ListSystemsParams, arg1 ...systems.ClientOption) (*systems.ListSystemsOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListSystems", varargs...)
	ret0, _ := ret[0].(*systems.ListSystemsOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSystems indicates an expected call of ListSystems.
func (mr *MockClientServiceMockRecorder) ListSystems(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSystems", reflect.TypeOf((*MockClientService)(nil).ListSystems), varargs...)
}

// RuleSuggestions mocks base method.
func (m *MockClientService) RuleSuggestions(arg0 *systems.RuleSuggestionsParams, arg1 ...systems.ClientOption) (*systems.RuleSuggestionsOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RuleSuggestions", varargs...)
	ret0, _ := ret[0].(*systems.RuleSuggestionsOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RuleSuggestions indicates an expected call of RuleSuggestions.
func (mr *MockClientServiceMockRecorder) RuleSuggestions(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RuleSuggestions", reflect.TypeOf((*MockClientService)(nil).RuleSuggestions), varargs...)
}

// SetTransport mocks base method.
func (m *MockClientService) SetTransport(arg0 runtime.ClientTransport) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTransport", arg0)
}

// SetTransport indicates an expected call of SetTransport.
func (mr *MockClientServiceMockRecorder) SetTransport(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTransport", reflect.TypeOf((*MockClientService)(nil).SetTransport), arg0)
}

// TranslateExternalIds mocks base method.
func (m *MockClientService) TranslateExternalIds(arg0 *systems.TranslateExternalIdsParams, arg1 ...systems.ClientOption) (*systems.TranslateExternalIdsOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "TranslateExternalIds", varargs...)
	ret0, _ := ret[0].(*systems.TranslateExternalIdsOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TranslateExternalIds indicates an expected call of TranslateExternalIds.
func (mr *MockClientServiceMockRecorder) TranslateExternalIds(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TranslateExternalIds", reflect.TypeOf((*MockClientService)(nil).TranslateExternalIds), varargs...)
}

// UpdateSystem mocks base method.
func (m *MockClientService) UpdateSystem(arg0 *systems.UpdateSystemParams, arg1 ...systems.ClientOption) (*systems.UpdateSystemOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateSystem", varargs...)
	ret0, _ := ret[0].(*systems.UpdateSystemOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateSystem indicates an expected call of UpdateSystem.
func (mr *MockClientServiceMockRecorder) UpdateSystem(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSystem", reflect.TypeOf((*MockClientService)(nil).UpdateSystem), varargs...)
}

// UpdateSystemBundleDeploy mocks base method.
func (m *MockClientService) UpdateSystemBundleDeploy(arg0 *systems.UpdateSystemBundleDeployParams, arg1 ...systems.ClientOption) (*systems.UpdateSystemBundleDeployOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateSystemBundleDeploy", varargs...)
	ret0, _ := ret[0].(*systems.UpdateSystemBundleDeployOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateSystemBundleDeploy indicates an expected call of UpdateSystemBundleDeploy.
func (mr *MockClientServiceMockRecorder) UpdateSystemBundleDeploy(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSystemBundleDeploy", reflect.TypeOf((*MockClientService)(nil).UpdateSystemBundleDeploy), varargs...)
}

// ValidateSystemCompliance mocks base method.
func (m *MockClientService) ValidateSystemCompliance(arg0 *systems.ValidateSystemComplianceParams, arg1 ...systems.ClientOption) (*systems.ValidateSystemComplianceOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ValidateSystemCompliance", varargs...)
	ret0, _ := ret[0].(*systems.ValidateSystemComplianceOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateSystemCompliance indicates an expected call of ValidateSystemCompliance.
func (mr *MockClientServiceMockRecorder) ValidateSystemCompliance(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateSystemCompliance", reflect.TypeOf((*MockClientService)(nil).ValidateSystemCompliance), varargs...)
}

// ValidateSystemTests mocks base method.
func (m *MockClientService) ValidateSystemTests(arg0 *systems.ValidateSystemTestsParams, arg1 ...systems.ClientOption) (*systems.ValidateSystemTestsOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ValidateSystemTests", varargs...)
	ret0, _ := ret[0].(*systems.ValidateSystemTestsOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateSystemTests indicates an expected call of ValidateSystemTests.
func (mr *MockClientServiceMockRecorder) ValidateSystemTests(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateSystemTests", reflect.TypeOf((*MockClientService)(nil).ValidateSystemTests), varargs...)
}

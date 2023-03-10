// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/chhabriv/search-results-aggregator/providerconnect (interfaces: ProviderService)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	providerconnect "github.com/chhabriv/search-results-aggregator/providerconnect"
	gomock "github.com/golang/mock/gomock"
)

// MockProviderService is a mock of ProviderService interface.
type MockProviderService struct {
	ctrl     *gomock.Controller
	recorder *MockProviderServiceMockRecorder
}

// MockProviderServiceMockRecorder is the mock recorder for MockProviderService.
type MockProviderServiceMockRecorder struct {
	mock *MockProviderService
}

// NewMockProviderService creates a new mock instance.
func NewMockProviderService(ctrl *gomock.Controller) *MockProviderService {
	mock := &MockProviderService{ctrl: ctrl}
	mock.recorder = &MockProviderServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProviderService) EXPECT() *MockProviderServiceMockRecorder {
	return m.recorder
}

// QueryProviders mocks base method.
func (m *MockProviderService) QueryProviders(arg0 context.Context) providerconnect.QueryProvidersResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryProviders", arg0)
	ret0, _ := ret[0].(providerconnect.QueryProvidersResponse)
	return ret0
}

// QueryProviders indicates an expected call of QueryProviders.
func (mr *MockProviderServiceMockRecorder) QueryProviders(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryProviders", reflect.TypeOf((*MockProviderService)(nil).QueryProviders), arg0)
}

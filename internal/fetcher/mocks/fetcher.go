// Code generated by MockGen. DO NOT EDIT.
// Source: fetcher.go
//
// Generated by this command:
//
//	mockgen -source=fetcher.go -destination=mocks/fetcher.go -package=mock_fetcher Fetcher
//

// Package mock_fetcher is a generated GoMock package.
package mock_fetcher

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockFetcher is a mock of Fetcher interface.
type MockFetcher struct {
	ctrl     *gomock.Controller
	recorder *MockFetcherMockRecorder
}

// MockFetcherMockRecorder is the mock recorder for MockFetcher.
type MockFetcherMockRecorder struct {
	mock *MockFetcher
}

// NewMockFetcher creates a new mock instance.
func NewMockFetcher(ctrl *gomock.Controller) *MockFetcher {
	mock := &MockFetcher{ctrl: ctrl}
	mock.recorder = &MockFetcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFetcher) EXPECT() *MockFetcherMockRecorder {
	return m.recorder
}

// FetchMasterChainBlockNumber mocks base method.
func (m *MockFetcher) FetchMasterChainBlockNumber() (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchMasterChainBlockNumber")
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchMasterChainBlockNumber indicates an expected call of FetchMasterChainBlockNumber.
func (mr *MockFetcherMockRecorder) FetchMasterChainBlockNumber() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchMasterChainBlockNumber", reflect.TypeOf((*MockFetcher)(nil).FetchMasterChainBlockNumber))
}

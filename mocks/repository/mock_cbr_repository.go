// Code generated by MockGen. DO NOT EDIT.
// Source: repository/cbr_repository.go

// Package mock_repository is a generated GoMock package.
package repository

import (
	context "context"
	entity "my_go/entity"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCBR is a mock of CBR interface.
type MockCBR struct {
	ctrl     *gomock.Controller
	recorder *MockCBRMockRecorder
}

// MockCBRMockRecorder is the mock recorder for MockCBR.
type MockCBRMockRecorder struct {
	mock *MockCBR
}

// NewMockCBR creates a new mock instance.
func NewMockCBR(ctrl *gomock.Controller) *MockCBR {
	mock := &MockCBR{ctrl: ctrl}
	mock.recorder = &MockCBRMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCBR) EXPECT() *MockCBRMockRecorder {
	return m.recorder
}

// GetCBRates mocks base method.
func (m *MockCBR) GetCBRates(ctx context.Context, req *entity.GetCBRatesRequest) (*entity.GetCBRatesResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCBRates", ctx, req)
	ret0, _ := ret[0].(*entity.GetCBRatesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCBRates indicates an expected call of GetCBRates.
func (mr *MockCBRMockRecorder) GetCBRates(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCBRates", reflect.TypeOf((*MockCBR)(nil).GetCBRates), ctx, req)
}

// GetExchangeRate mocks base method.
func (m *MockCBR) GetExchangeRate(ctx context.Context, req *entity.GetExchangeRateRequest) (*entity.GetExchangeRateResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetExchangeRate", ctx, req)
	ret0, _ := ret[0].(*entity.GetExchangeRateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetExchangeRate indicates an expected call of GetExchangeRate.
func (mr *MockCBRMockRecorder) GetExchangeRate(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetExchangeRate", reflect.TypeOf((*MockCBR)(nil).GetExchangeRate), ctx, req)
}
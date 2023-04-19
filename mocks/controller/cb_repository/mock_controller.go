// Code generated by MockGen. DO NOT EDIT.
// Source: controller/cb_repository/controller.go

// Package mock_cb_repository is a generated GoMock package.
package cb_repository

import (
	context "context"
	entity "my_go/entity"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockController is a mock of Controller interface.
type MockController struct {
	ctrl     *gomock.Controller
	recorder *MockControllerMockRecorder
}

// MockControllerMockRecorder is the mock recorder for MockController.
type MockControllerMockRecorder struct {
	mock *MockController
}

// NewMockController creates a new mock instance.
func NewMockController(ctrl *gomock.Controller) *MockController {
	mock := &MockController{ctrl: ctrl}
	mock.recorder = &MockControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockController) EXPECT() *MockControllerMockRecorder {
	return m.recorder
}

// GetCBRates mocks base method.
func (m *MockController) GetCBRates(ctx context.Context, req *entity.GetExchangeRatesRequest) (*entity.GetExchangeRatesResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCBRates", ctx, req)
	ret0, _ := ret[0].(*entity.GetExchangeRatesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCBRates indicates an expected call of GetCBRates.
func (mr *MockControllerMockRecorder) GetCBRates(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCBRates", reflect.TypeOf((*MockController)(nil).GetCBRates), ctx, req)
}

// GetExchangeRate mocks base method.
func (m *MockController) GetExchangeRate(ctx context.Context, req *entity.GetExchangeRateRequest) (*entity.GetExchangeRateResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetExchangeRate", ctx, req)
	ret0, _ := ret[0].(*entity.GetExchangeRateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetExchangeRate indicates an expected call of GetExchangeRate.
func (mr *MockControllerMockRecorder) GetExchangeRate(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetExchangeRate", reflect.TypeOf((*MockController)(nil).GetExchangeRate), ctx, req)
}

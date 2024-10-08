// Code generated by MockGen. DO NOT EDIT.
// Source: /Users/qlong/indeed/temporoio/sdk-go/converter/data_converter.go

// Package temporal is a generated GoMock package.
package temporal

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	common "go.temporal.io/api/common/v1"
)

// MockDataConverter is a mock of DataConverter interface.
type MockDataConverter struct {
	ctrl     *gomock.Controller
	recorder *MockDataConverterMockRecorder
}

// MockDataConverterMockRecorder is the mock recorder for MockDataConverter.
type MockDataConverterMockRecorder struct {
	mock *MockDataConverter
}

// NewMockDataConverter creates a new mock instance.
func NewMockDataConverter(ctrl *gomock.Controller) *MockDataConverter {
	mock := &MockDataConverter{ctrl: ctrl}
	mock.recorder = &MockDataConverterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDataConverter) EXPECT() *MockDataConverterMockRecorder {
	return m.recorder
}

// FromPayload mocks base method.
func (m *MockDataConverter) FromPayload(payload *common.Payload, valuePtr interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FromPayload", payload, valuePtr)
	ret0, _ := ret[0].(error)
	return ret0
}

// FromPayload indicates an expected call of FromPayload.
func (mr *MockDataConverterMockRecorder) FromPayload(payload, valuePtr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FromPayload", reflect.TypeOf((*MockDataConverter)(nil).FromPayload), payload, valuePtr)
}

// FromPayloads mocks base method.
func (m *MockDataConverter) FromPayloads(payloads *common.Payloads, valuePtrs ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{payloads}
	for _, a := range valuePtrs {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FromPayloads", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// FromPayloads indicates an expected call of FromPayloads.
func (mr *MockDataConverterMockRecorder) FromPayloads(payloads interface{}, valuePtrs ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{payloads}, valuePtrs...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FromPayloads", reflect.TypeOf((*MockDataConverter)(nil).FromPayloads), varargs...)
}

// ToPayload mocks base method.
func (m *MockDataConverter) ToPayload(value interface{}) (*common.Payload, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToPayload", value)
	ret0, _ := ret[0].(*common.Payload)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ToPayload indicates an expected call of ToPayload.
func (mr *MockDataConverterMockRecorder) ToPayload(value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToPayload", reflect.TypeOf((*MockDataConverter)(nil).ToPayload), value)
}

// ToPayloads mocks base method.
func (m *MockDataConverter) ToPayloads(value ...interface{}) (*common.Payloads, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range value {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ToPayloads", varargs...)
	ret0, _ := ret[0].(*common.Payloads)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ToPayloads indicates an expected call of ToPayloads.
func (mr *MockDataConverterMockRecorder) ToPayloads(value ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToPayloads", reflect.TypeOf((*MockDataConverter)(nil).ToPayloads), value...)
}

// ToString mocks base method.
func (m *MockDataConverter) ToString(input *common.Payload) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToString", input)
	ret0, _ := ret[0].(string)
	return ret0
}

// ToString indicates an expected call of ToString.
func (mr *MockDataConverterMockRecorder) ToString(input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToString", reflect.TypeOf((*MockDataConverter)(nil).ToString), input)
}

// ToStrings mocks base method.
func (m *MockDataConverter) ToStrings(input *common.Payloads) []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToStrings", input)
	ret0, _ := ret[0].([]string)
	return ret0
}

// ToStrings indicates an expected call of ToStrings.
func (mr *MockDataConverterMockRecorder) ToStrings(input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToStrings", reflect.TypeOf((*MockDataConverter)(nil).ToStrings), input)
}

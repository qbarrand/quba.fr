// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/qbarrand/quba.fr/internal/imgpro (interfaces: Handler,Processor)

// Package mock_imgpro is a generated GoMock package.
package mock_imgpro

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	imgpro "github.com/qbarrand/quba.fr/internal/imgpro"
)

// MockHandler is a mock of Handler interface.
type MockHandler struct {
	ctrl     *gomock.Controller
	recorder *MockHandlerMockRecorder
}

// MockHandlerMockRecorder is the mock recorder for MockHandler.
type MockHandlerMockRecorder struct {
	mock *MockHandler
}

// NewMockHandler creates a new mock instance.
func NewMockHandler(ctrl *gomock.Controller) *MockHandler {
	mock := &MockHandler{ctrl: ctrl}
	mock.recorder = &MockHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHandler) EXPECT() *MockHandlerMockRecorder {
	return m.recorder
}

// Bytes mocks base method.
func (m *MockHandler) Bytes() ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Bytes")
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Bytes indicates an expected call of Bytes.
func (mr *MockHandlerMockRecorder) Bytes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Bytes", reflect.TypeOf((*MockHandler)(nil).Bytes))
}

// Destroy mocks base method.
func (m *MockHandler) Destroy() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Destroy")
	ret0, _ := ret[0].(error)
	return ret0
}

// Destroy indicates an expected call of Destroy.
func (mr *MockHandlerMockRecorder) Destroy() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Destroy", reflect.TypeOf((*MockHandler)(nil).Destroy))
}

// Resize mocks base method.
func (m *MockHandler) Resize(arg0 context.Context, arg1, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Resize", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Resize indicates an expected call of Resize.
func (mr *MockHandlerMockRecorder) Resize(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Resize", reflect.TypeOf((*MockHandler)(nil).Resize), arg0, arg1, arg2)
}

// SetFormat mocks base method.
func (m *MockHandler) SetFormat(arg0 imgpro.Format) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetFormat", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetFormat indicates an expected call of SetFormat.
func (mr *MockHandlerMockRecorder) SetFormat(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetFormat", reflect.TypeOf((*MockHandler)(nil).SetFormat), arg0)
}

// StripMetadata mocks base method.
func (m *MockHandler) StripMetadata() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StripMetadata")
	ret0, _ := ret[0].(error)
	return ret0
}

// StripMetadata indicates an expected call of StripMetadata.
func (mr *MockHandlerMockRecorder) StripMetadata() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StripMetadata", reflect.TypeOf((*MockHandler)(nil).StripMetadata))
}

// MockProcessor is a mock of Processor interface.
type MockProcessor struct {
	ctrl     *gomock.Controller
	recorder *MockProcessorMockRecorder
}

// MockProcessorMockRecorder is the mock recorder for MockProcessor.
type MockProcessorMockRecorder struct {
	mock *MockProcessor
}

// NewMockProcessor creates a new mock instance.
func NewMockProcessor(ctrl *gomock.Controller) *MockProcessor {
	mock := &MockProcessor{ctrl: ctrl}
	mock.recorder = &MockProcessorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProcessor) EXPECT() *MockProcessorMockRecorder {
	return m.recorder
}

// BestFormats mocks base method.
func (m *MockProcessor) BestFormats() []imgpro.Format {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BestFormats")
	ret0, _ := ret[0].([]imgpro.Format)
	return ret0
}

// BestFormats indicates an expected call of BestFormats.
func (mr *MockProcessorMockRecorder) BestFormats() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BestFormats", reflect.TypeOf((*MockProcessor)(nil).BestFormats))
}

// Destroy mocks base method.
func (m *MockProcessor) Destroy() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Destroy")
	ret0, _ := ret[0].(error)
	return ret0
}

// Destroy indicates an expected call of Destroy.
func (mr *MockProcessorMockRecorder) Destroy() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Destroy", reflect.TypeOf((*MockProcessor)(nil).Destroy))
}

// HandlerFromBytes mocks base method.
func (m *MockProcessor) HandlerFromBytes(arg0 []byte) (imgpro.Handler, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandlerFromBytes", arg0)
	ret0, _ := ret[0].(imgpro.Handler)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HandlerFromBytes indicates an expected call of HandlerFromBytes.
func (mr *MockProcessorMockRecorder) HandlerFromBytes(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandlerFromBytes", reflect.TypeOf((*MockProcessor)(nil).HandlerFromBytes), arg0)
}

// Init mocks base method.
func (m *MockProcessor) Init() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Init")
	ret0, _ := ret[0].(error)
	return ret0
}

// Init indicates an expected call of Init.
func (mr *MockProcessorMockRecorder) Init() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockProcessor)(nil).Init))
}

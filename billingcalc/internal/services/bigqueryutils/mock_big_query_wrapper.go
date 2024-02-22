package bigqueryutils

import (
	"context"
	"reflect"

	"github.com/golang/mock/gomock"
)

// MockBigQueryWrapper is a mock of BigQueryWrapper interface.
type MockBigQueryWrapper struct {
	ctrl     *gomock.Controller
	recorder *MockBigQueryWrapperMockRecorder
}

// MockBigQueryWrapperMockRecorder is the mock recorder for MockBigQueryWrapper.
type MockBigQueryWrapperMockRecorder struct {
	mock *MockBigQueryWrapper
}

// NewMockBigQueryWrapper creates a new mock instance.
func NewMockBigQueryWrapper(ctrl *gomock.Controller) *MockBigQueryWrapper {
	mock := &MockBigQueryWrapper{ctrl: ctrl}
	mock.recorder = &MockBigQueryWrapperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBigQueryWrapper) EXPECT() *MockBigQueryWrapperMockRecorder {
	return m.recorder
}

// ExecuteQuery mocks base method.
func (m *MockBigQueryWrapper) ExecuteQuery(ctx context.Context, query string) (ResultIterator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExecuteQuery", ctx, query)
	ret0, _ := ret[0].(ResultIterator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MockResultIterator is a mock of ResultIterator interface.
type MockResultIterator struct {
	ctrl     *gomock.Controller
	recorder *MockResultIteratorMockRecorder
}

// MockResultIteratorMockRecorder is the mock recorder for MockResultIterator.
type MockResultIteratorMockRecorder struct {
	mock *MockResultIterator
}

// NewMockResultIterator creates a new mock instance.
func NewMockResultIterator(ctrl *gomock.Controller) *MockResultIterator {
	mock := &MockResultIterator{ctrl: ctrl}
	mock.recorder = &MockResultIteratorMockRecorder{mock}
	return mock
}

// Next mocks base method.
func (m *MockResultIterator) Next(dst interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Next", dst)
	ret0, _ := ret[0].(error)
	return ret0
}

// Next indicates an expected call of Next.
func (mr *MockResultIteratorMockRecorder) Next(dst interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Next", reflect.TypeOf((*MockResultIterator)(nil).Next), dst)
}

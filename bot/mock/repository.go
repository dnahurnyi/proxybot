// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/dnahurnyi/proxybot/bot (interfaces: Repository)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	bot "github.com/dnahurnyi/proxybot/bot"
	gomock "github.com/golang/mock/gomock"
	go_uuid "github.com/satori/go.uuid"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// GetSubscription mocks base method.
func (m *MockRepository) GetSubscription(arg0 int64) (*bot.Subscription, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubscription", arg0)
	ret0, _ := ret[0].(*bot.Subscription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubscription indicates an expected call of GetSubscription.
func (mr *MockRepositoryMockRecorder) GetSubscription(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubscription", reflect.TypeOf((*MockRepository)(nil).GetSubscription), arg0)
}

// GetTagByName mocks base method.
func (m *MockRepository) GetTagByName(arg0 string) (*bot.Tag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTagByName", arg0)
	ret0, _ := ret[0].(*bot.Tag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTagByName indicates an expected call of GetTagByName.
func (mr *MockRepositoryMockRecorder) GetTagByName(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTagByName", reflect.TypeOf((*MockRepository)(nil).GetTagByName), arg0)
}

// ListSubscriptions mocks base method.
func (m *MockRepository) ListSubscriptions() ([]bot.Subscription, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSubscriptions")
	ret0, _ := ret[0].([]bot.Subscription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSubscriptions indicates an expected call of ListSubscriptions.
func (mr *MockRepositoryMockRecorder) ListSubscriptions() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSubscriptions", reflect.TypeOf((*MockRepository)(nil).ListSubscriptions))
}

// ListTags mocks base method.
func (m *MockRepository) ListTags() ([]bot.Tag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTags")
	ret0, _ := ret[0].([]bot.Tag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTags indicates an expected call of ListTags.
func (mr *MockRepositoryMockRecorder) ListTags() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTags", reflect.TypeOf((*MockRepository)(nil).ListTags))
}

// SaveSubscription mocks base method.
func (m *MockRepository) SaveSubscription(arg0 *bot.Subscription) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveSubscription", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveSubscription indicates an expected call of SaveSubscription.
func (mr *MockRepositoryMockRecorder) SaveSubscription(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveSubscription", reflect.TypeOf((*MockRepository)(nil).SaveSubscription), arg0)
}

// SaveTag mocks base method.
func (m *MockRepository) SaveTag(arg0 *bot.Tag) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveTag", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveTag indicates an expected call of SaveTag.
func (mr *MockRepositoryMockRecorder) SaveTag(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveTag", reflect.TypeOf((*MockRepository)(nil).SaveTag), arg0)
}

// TagSubscription mocks base method.
func (m *MockRepository) TagSubscription(arg0 go_uuid.UUID, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TagSubscription", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// TagSubscription indicates an expected call of TagSubscription.
func (mr *MockRepositoryMockRecorder) TagSubscription(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TagSubscription", reflect.TypeOf((*MockRepository)(nil).TagSubscription), arg0, arg1)
}

// Transaction mocks base method.
func (m *MockRepository) Transaction(arg0 func(bot.Repository) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Transaction", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Transaction indicates an expected call of Transaction.
func (mr *MockRepositoryMockRecorder) Transaction(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Transaction", reflect.TypeOf((*MockRepository)(nil).Transaction), arg0)
}

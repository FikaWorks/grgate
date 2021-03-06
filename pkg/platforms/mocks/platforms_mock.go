// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/fikaworks/grgate/pkg/platforms (interfaces: Platform)

// Package mock_platforms is a generated GoMock package.
package mock_platforms

import (
	io "io"
	reflect "reflect"

	platforms "github.com/fikaworks/grgate/pkg/platforms"
	gomock "github.com/golang/mock/gomock"
)

// MockPlatform is a mock of Platform interface.
type MockPlatform struct {
	ctrl     *gomock.Controller
	recorder *MockPlatformMockRecorder
}

// MockPlatformMockRecorder is the mock recorder for MockPlatform.
type MockPlatformMockRecorder struct {
	mock *MockPlatform
}

// NewMockPlatform creates a new mock instance.
func NewMockPlatform(ctrl *gomock.Controller) *MockPlatform {
	mock := &MockPlatform{ctrl: ctrl}
	mock.recorder = &MockPlatformMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPlatform) EXPECT() *MockPlatformMockRecorder {
	return m.recorder
}

// CheckAllStatusSucceeded mocks base method.
func (m *MockPlatform) CheckAllStatusSucceeded(arg0, arg1, arg2 string, arg3 []string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckAllStatusSucceeded", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckAllStatusSucceeded indicates an expected call of CheckAllStatusSucceeded.
func (mr *MockPlatformMockRecorder) CheckAllStatusSucceeded(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckAllStatusSucceeded", reflect.TypeOf((*MockPlatform)(nil).CheckAllStatusSucceeded), arg0, arg1, arg2, arg3)
}

// CreateFile mocks base method.
func (m *MockPlatform) CreateFile(arg0, arg1, arg2, arg3, arg4, arg5 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFile", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateFile indicates an expected call of CreateFile.
func (mr *MockPlatformMockRecorder) CreateFile(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFile", reflect.TypeOf((*MockPlatform)(nil).CreateFile), arg0, arg1, arg2, arg3, arg4, arg5)
}

// CreateIssue mocks base method.
func (m *MockPlatform) CreateIssue(arg0, arg1 string, arg2 *platforms.Issue) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateIssue", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateIssue indicates an expected call of CreateIssue.
func (mr *MockPlatformMockRecorder) CreateIssue(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateIssue", reflect.TypeOf((*MockPlatform)(nil).CreateIssue), arg0, arg1, arg2)
}

// CreateRelease mocks base method.
func (m *MockPlatform) CreateRelease(arg0, arg1 string, arg2 *platforms.Release) (*platforms.Release, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRelease", arg0, arg1, arg2)
	ret0, _ := ret[0].(*platforms.Release)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRelease indicates an expected call of CreateRelease.
func (mr *MockPlatformMockRecorder) CreateRelease(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRelease", reflect.TypeOf((*MockPlatform)(nil).CreateRelease), arg0, arg1, arg2)
}

// CreateRepository mocks base method.
func (m *MockPlatform) CreateRepository(arg0, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRepository", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateRepository indicates an expected call of CreateRepository.
func (mr *MockPlatformMockRecorder) CreateRepository(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRepository", reflect.TypeOf((*MockPlatform)(nil).CreateRepository), arg0, arg1, arg2)
}

// CreateStatus mocks base method.
func (m *MockPlatform) CreateStatus(arg0, arg1 string, arg2 *platforms.Status) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateStatus", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateStatus indicates an expected call of CreateStatus.
func (mr *MockPlatformMockRecorder) CreateStatus(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateStatus", reflect.TypeOf((*MockPlatform)(nil).CreateStatus), arg0, arg1, arg2)
}

// DeleteRepository mocks base method.
func (m *MockPlatform) DeleteRepository(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRepository", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteRepository indicates an expected call of DeleteRepository.
func (mr *MockPlatformMockRecorder) DeleteRepository(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRepository", reflect.TypeOf((*MockPlatform)(nil).DeleteRepository), arg0, arg1)
}

// GetStatus mocks base method.
func (m *MockPlatform) GetStatus(arg0, arg1, arg2, arg3 string) (*platforms.Status, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStatus", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*platforms.Status)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStatus indicates an expected call of GetStatus.
func (mr *MockPlatformMockRecorder) GetStatus(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStatus", reflect.TypeOf((*MockPlatform)(nil).GetStatus), arg0, arg1, arg2, arg3)
}

// ListDraftReleases mocks base method.
func (m *MockPlatform) ListDraftReleases(arg0, arg1 string) ([]*platforms.Release, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListDraftReleases", arg0, arg1)
	ret0, _ := ret[0].([]*platforms.Release)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListDraftReleases indicates an expected call of ListDraftReleases.
func (mr *MockPlatformMockRecorder) ListDraftReleases(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListDraftReleases", reflect.TypeOf((*MockPlatform)(nil).ListDraftReleases), arg0, arg1)
}

// ListIssuesByAuthor mocks base method.
func (m *MockPlatform) ListIssuesByAuthor(arg0, arg1 string, arg2 interface{}) ([]*platforms.Issue, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListIssuesByAuthor", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*platforms.Issue)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListIssuesByAuthor indicates an expected call of ListIssuesByAuthor.
func (mr *MockPlatformMockRecorder) ListIssuesByAuthor(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListIssuesByAuthor", reflect.TypeOf((*MockPlatform)(nil).ListIssuesByAuthor), arg0, arg1, arg2)
}

// ListReleases mocks base method.
func (m *MockPlatform) ListReleases(arg0, arg1 string) ([]*platforms.Release, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListReleases", arg0, arg1)
	ret0, _ := ret[0].([]*platforms.Release)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListReleases indicates an expected call of ListReleases.
func (mr *MockPlatformMockRecorder) ListReleases(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListReleases", reflect.TypeOf((*MockPlatform)(nil).ListReleases), arg0, arg1)
}

// ListStatuses mocks base method.
func (m *MockPlatform) ListStatuses(arg0, arg1, arg2 string) ([]*platforms.Status, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListStatuses", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*platforms.Status)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListStatuses indicates an expected call of ListStatuses.
func (mr *MockPlatformMockRecorder) ListStatuses(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListStatuses", reflect.TypeOf((*MockPlatform)(nil).ListStatuses), arg0, arg1, arg2)
}

// PublishRelease mocks base method.
func (m *MockPlatform) PublishRelease(arg0, arg1 string, arg2 *platforms.Release) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublishRelease", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PublishRelease indicates an expected call of PublishRelease.
func (mr *MockPlatformMockRecorder) PublishRelease(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishRelease", reflect.TypeOf((*MockPlatform)(nil).PublishRelease), arg0, arg1, arg2)
}

// ReadFile mocks base method.
func (m *MockPlatform) ReadFile(arg0, arg1, arg2 string) (io.Reader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadFile", arg0, arg1, arg2)
	ret0, _ := ret[0].(io.Reader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadFile indicates an expected call of ReadFile.
func (mr *MockPlatformMockRecorder) ReadFile(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadFile", reflect.TypeOf((*MockPlatform)(nil).ReadFile), arg0, arg1, arg2)
}

// UpdateFile mocks base method.
func (m *MockPlatform) UpdateFile(arg0, arg1, arg2, arg3, arg4, arg5 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateFile", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateFile indicates an expected call of UpdateFile.
func (mr *MockPlatformMockRecorder) UpdateFile(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFile", reflect.TypeOf((*MockPlatform)(nil).UpdateFile), arg0, arg1, arg2, arg3, arg4, arg5)
}

// UpdateIssue mocks base method.
func (m *MockPlatform) UpdateIssue(arg0, arg1 string, arg2 *platforms.Issue) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateIssue", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateIssue indicates an expected call of UpdateIssue.
func (mr *MockPlatformMockRecorder) UpdateIssue(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateIssue", reflect.TypeOf((*MockPlatform)(nil).UpdateIssue), arg0, arg1, arg2)
}

// UpdateRelease mocks base method.
func (m *MockPlatform) UpdateRelease(arg0, arg1 string, arg2 *platforms.Release) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRelease", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateRelease indicates an expected call of UpdateRelease.
func (mr *MockPlatformMockRecorder) UpdateRelease(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRelease", reflect.TypeOf((*MockPlatform)(nil).UpdateRelease), arg0, arg1, arg2)
}

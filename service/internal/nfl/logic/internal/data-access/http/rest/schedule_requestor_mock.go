// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/rmarken5/mini-score/service/internal/nfl/logic/internal/data-access/http/rest (interfaces: Requester)
//
// Generated by this command:
//
//	mockgen -destination ./schedule_requestor_mock.go -package rest . Requester
//
// Package rest is a generated GoMock package.
package rest

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockRequester is a mock of Requester interface.
type MockRequester struct {
	ctrl     *gomock.Controller
	recorder *MockRequesterMockRecorder
}

// MockRequesterMockRecorder is the mock recorder for MockRequester.
type MockRequesterMockRecorder struct {
	mock *MockRequester
}

// NewMockRequester creates a new mock instance.
func NewMockRequester(ctrl *gomock.Controller) *MockRequester {
	mock := &MockRequester{ctrl: ctrl}
	mock.recorder = &MockRequesterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRequester) EXPECT() *MockRequesterMockRecorder {
	return m.recorder
}

// GetScoreboard mocks base method.
func (m *MockRequester) GetScoreboard() (ScoreboardResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetScoreboard")
	ret0, _ := ret[0].(ScoreboardResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetScoreboard indicates an expected call of GetScoreboard.
func (mr *MockRequesterMockRecorder) GetScoreboard() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetScoreboard", reflect.TypeOf((*MockRequester)(nil).GetScoreboard))
}

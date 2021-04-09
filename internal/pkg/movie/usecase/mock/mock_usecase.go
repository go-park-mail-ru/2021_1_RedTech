// Code generated by MockGen. DO NOT EDIT.
// Source: Redioteka/internal/pkg/domain (interfaces: MovieUsecase)

// Package mock is a generated GoMock package.
package mock

import (
	domain "Redioteka/internal/pkg/domain"
	session "Redioteka/internal/pkg/utils/session"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockMovieUsecase is a mock of MovieUsecase interface.
type MockMovieUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockMovieUsecaseMockRecorder
}

// MockMovieUsecaseMockRecorder is the mock recorder for MockMovieUsecase.
type MockMovieUsecaseMockRecorder struct {
	mock *MockMovieUsecase
}

// NewMockMovieUsecase creates a new mock instance.
func NewMockMovieUsecase(ctrl *gomock.Controller) *MockMovieUsecase {
	mock := &MockMovieUsecase{ctrl: ctrl}
	mock.recorder = &MockMovieUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMovieUsecase) EXPECT() *MockMovieUsecaseMockRecorder {
	return m.recorder
}

// AddFavourite mocks base method.
func (m *MockMovieUsecase) AddFavourite(arg0 uint, arg1 *session.Session) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddFavourite", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddFavourite indicates an expected call of AddFavourite.
func (mr *MockMovieUsecaseMockRecorder) AddFavourite(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddFavourite", reflect.TypeOf((*MockMovieUsecase)(nil).AddFavourite), arg0, arg1)
}

// GetById mocks base method.
func (m *MockMovieUsecase) GetById(arg0 uint) (domain.Movie, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", arg0)
	ret0, _ := ret[0].(domain.Movie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockMovieUsecaseMockRecorder) GetById(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockMovieUsecase)(nil).GetById), arg0)
}

// RemoveFavourite mocks base method.
func (m *MockMovieUsecase) RemoveFavourite(arg0 uint, arg1 *session.Session) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveFavourite", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveFavourite indicates an expected call of RemoveFavourite.
func (mr *MockMovieUsecaseMockRecorder) RemoveFavourite(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveFavourite", reflect.TypeOf((*MockMovieUsecase)(nil).RemoveFavourite), arg0, arg1)
}

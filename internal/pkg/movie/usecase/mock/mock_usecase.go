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

// GetByFilter mocks base method.
func (m *MockMovieUsecase) GetByFilter(arg0 domain.MovieFilter) ([]domain.Movie, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByFilter", arg0)
	ret0, _ := ret[0].([]domain.Movie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByFilter indicates an expected call of GetByFilter.
func (mr *MockMovieUsecaseMockRecorder) GetByFilter(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByFilter", reflect.TypeOf((*MockMovieUsecase)(nil).GetByFilter), arg0)
}

// GetByID mocks base method.
func (m *MockMovieUsecase) GetByID(arg0 uint, arg1 *session.Session) (domain.Movie, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", arg0, arg1)
	ret0, _ := ret[0].(domain.Movie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockMovieUsecaseMockRecorder) GetByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockMovieUsecase)(nil).GetByID), arg0, arg1)
}

// GetGenres mocks base method.
func (m *MockMovieUsecase) GetGenres() ([]domain.Genre, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGenres")
	ret0, _ := ret[0].([]domain.Genre)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGenres indicates an expected call of GetGenres.
func (mr *MockMovieUsecaseMockRecorder) GetGenres() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGenres", reflect.TypeOf((*MockMovieUsecase)(nil).GetGenres))
}

// GetStream mocks base method.
func (m *MockMovieUsecase) GetStream(arg0 uint) (domain.Stream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStream", arg0)
	ret0, _ := ret[0].(domain.Stream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStream indicates an expected call of GetStream.
func (mr *MockMovieUsecaseMockRecorder) GetStream(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStream", reflect.TypeOf((*MockMovieUsecase)(nil).GetStream), arg0)
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

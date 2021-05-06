// Code generated by MockGen. DO NOT EDIT.
// Source: Redioteka/internal/pkg/domain (interfaces: MovieRepository)

// Package mock is a generated GoMock package.
package mock

import (
	domain "Redioteka/internal/pkg/domain"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockMovieRepository is a mock of MovieRepository interface.
type MockMovieRepository struct {
	ctrl     *gomock.Controller
	recorder *MockMovieRepositoryMockRecorder
}

// MockMovieRepositoryMockRecorder is the mock recorder for MockMovieRepository.
type MockMovieRepositoryMockRecorder struct {
	mock *MockMovieRepository
}

// NewMockMovieRepository creates a new mock instance.
func NewMockMovieRepository(ctrl *gomock.Controller) *MockMovieRepository {
	mock := &MockMovieRepository{ctrl: ctrl}
	mock.recorder = &MockMovieRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMovieRepository) EXPECT() *MockMovieRepositoryMockRecorder {
	return m.recorder
}

// AddFavouriteByID mocks base method.
func (m *MockMovieRepository) AddFavouriteByID(arg0, arg1 uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddFavouriteByID", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddFavouriteByID indicates an expected call of AddFavouriteByID.
func (mr *MockMovieRepositoryMockRecorder) AddFavouriteByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddFavouriteByID", reflect.TypeOf((*MockMovieRepository)(nil).AddFavouriteByID), arg0, arg1)
}

// CheckFavouriteByID mocks base method.
func (m *MockMovieRepository) CheckFavouriteByID(arg0, arg1 uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckFavouriteByID", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckFavouriteByID indicates an expected call of CheckFavouriteByID.
func (mr *MockMovieRepositoryMockRecorder) CheckFavouriteByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckFavouriteByID", reflect.TypeOf((*MockMovieRepository)(nil).CheckFavouriteByID), arg0, arg1)
}

// CheckVoteByID mocks base method.
func (m *MockMovieRepository) CheckVoteByID(arg0, arg1 uint) int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckVoteByID", arg0, arg1)
	ret0, _ := ret[0].(int)
	return ret0
}

// CheckVoteByID indicates an expected call of CheckVoteByID.
func (mr *MockMovieRepositoryMockRecorder) CheckVoteByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckVoteByID", reflect.TypeOf((*MockMovieRepository)(nil).CheckVoteByID), arg0, arg1)
}

// Dislike mocks base method.
func (m *MockMovieRepository) Dislike(arg0, arg1 uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Dislike", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Dislike indicates an expected call of Dislike.
func (mr *MockMovieRepositoryMockRecorder) Dislike(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dislike", reflect.TypeOf((*MockMovieRepository)(nil).Dislike), arg0, arg1)
}

// GetByFilter mocks base method.
func (m *MockMovieRepository) GetByFilter(arg0 domain.MovieFilter) ([]domain.Movie, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByFilter", arg0)
	ret0, _ := ret[0].([]domain.Movie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByFilter indicates an expected call of GetByFilter.
func (mr *MockMovieRepositoryMockRecorder) GetByFilter(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByFilter", reflect.TypeOf((*MockMovieRepository)(nil).GetByFilter), arg0)
}

// GetById mocks base method.
func (m *MockMovieRepository) GetById(arg0 uint) (domain.Movie, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", arg0)
	ret0, _ := ret[0].(domain.Movie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockMovieRepositoryMockRecorder) GetById(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockMovieRepository)(nil).GetById), arg0)
}

// GetGenres mocks base method.
func (m *MockMovieRepository) GetGenres() ([]domain.Genre, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGenres")
	ret0, _ := ret[0].([]domain.Genre)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGenres indicates an expected call of GetGenres.
func (mr *MockMovieRepositoryMockRecorder) GetGenres() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGenres", reflect.TypeOf((*MockMovieRepository)(nil).GetGenres))
}

// GetSeriesList mocks base method.
func (m *MockMovieRepository) GetSeriesList(arg0 uint) ([]uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSeriesList", arg0)
	ret0, _ := ret[0].([]uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSeriesList indicates an expected call of GetSeriesList.
func (mr *MockMovieRepositoryMockRecorder) GetSeriesList(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSeriesList", reflect.TypeOf((*MockMovieRepository)(nil).GetSeriesList), arg0)
}

// GetStream mocks base method.
func (m *MockMovieRepository) GetStream(arg0 uint) ([]domain.Stream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStream", arg0)
	ret0, _ := ret[0].([]domain.Stream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStream indicates an expected call of GetStream.
func (mr *MockMovieRepositoryMockRecorder) GetStream(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStream", reflect.TypeOf((*MockMovieRepository)(nil).GetStream), arg0)
}

// Like mocks base method.
func (m *MockMovieRepository) Like(arg0, arg1 uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Like", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Like indicates an expected call of Like.
func (mr *MockMovieRepositoryMockRecorder) Like(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Like", reflect.TypeOf((*MockMovieRepository)(nil).Like), arg0, arg1)
}

// RemoveFavouriteByID mocks base method.
func (m *MockMovieRepository) RemoveFavouriteByID(arg0, arg1 uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveFavouriteByID", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveFavouriteByID indicates an expected call of RemoveFavouriteByID.
func (mr *MockMovieRepositoryMockRecorder) RemoveFavouriteByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveFavouriteByID", reflect.TypeOf((*MockMovieRepository)(nil).RemoveFavouriteByID), arg0, arg1)
}

// Search mocks base method.
func (m *MockMovieRepository) Search(arg0 string) ([]domain.Movie, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", arg0)
	ret0, _ := ret[0].([]domain.Movie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search.
func (mr *MockMovieRepositoryMockRecorder) Search(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockMovieRepository)(nil).Search), arg0)
}

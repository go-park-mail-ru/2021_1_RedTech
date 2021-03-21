package repository

import (
	"Redioteka/internal/app/domain"
	"errors"
	"sync"
)

type mapMovieRepository struct {
	sync.Mutex
	movies map[uint]*domain.Movie
}

func NewMapMovieRepository() domain.MovieRepository {
	return &mapMovieRepository{
		movies: make(map[uint]*domain.Movie),
	}
}

func (m *mapMovieRepository) GetById(id uint) (domain.Movie, error) {
	m.Lock()
	movie, exists := m.movies[id]
	m.Unlock()
	if !exists {
		return domain.Movie{}, errors.New("movie not found")
	}
	return *movie, nil
}

func (m *mapMovieRepository) Delete(id uint) error {
	m.Lock()
	defer m.Unlock()
	_, inMap := m.movies[id]
	if !inMap {
		return errors.New("user not in map")
	}
	delete(m.movies, id)
	return nil
}

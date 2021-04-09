package repository

import (
	"Redioteka/internal/pkg/domain"
	"errors"
	"sync"
)

type mapMovieRepository struct {
	sync.Mutex
	movies map[uint]*domain.Movie
}

func NewMapMovieRepository() domain.MovieRepository {
	newMap := &mapMovieRepository{
		movies: make(map[uint]*domain.Movie),
	}
	newMap.fillMap()
	return newMap
}

func (m *mapMovieRepository) fillMap() {
	m.movies[1] = &domain.Movie{
		ID:          1,
		Title:       "Film",
		Description: "Test data",
		Rating:      9,
		Countries:   []string{"Japan", "South Korea"},
		IsFree:      false,
		Genres:      []string{"Comedy"},
		Actors:      []string{"Sana", "Momo", "Mina"},
		Avatar:      "/static/movies/default.jpg",
		Type:        domain.MovieT,
		Year:        "2012",
		Director:    []string{"James Cameron"},
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

func (m *mapMovieRepository) AddFavouriteByID(movieID, userID uint) error {
	return errors.New("TO DO")
}

func (m *mapMovieRepository) RemoveFavouriteByID(movieID, userID uint) error {
	return errors.New("TO DO")
}

func (m *mapMovieRepository) CheckFavouriteByID(movieID, userID uint) error {
	return errors.New("TO DO")
}

package repository

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/movie"
	"Redioteka/internal/pkg/utils/baseutils"
	"Redioteka/internal/pkg/utils/movie_generator"
	"errors"
	"sync"
)

type mapMovieRepository struct {
	sync.Mutex
	movies map[uint]domain.Movie
}

func NewMapMovieRepository() domain.MovieRepository {
	newMap := &mapMovieRepository{
		movies: make(map[uint]domain.Movie),
	}
	newMap.fillMap()
	return newMap
}

// можно было бы заюзать мапку, но тут у нас будут массивы длины <= 20
func hasIntersection(source []string, filter []string) bool {
	if filter == nil {
		return true
	}
	for _, filterValue := range filter {
		for _, sourceValue := range source {
			if filterValue == sourceValue {
				return true
			}
		}
	}
	return false
}

func passFree(isFree bool, priceFilter int) bool {
	if priceFilter == domain.FilterFree && isFree {
		return true
	}
	if priceFilter == domain.FilterSubscription && !isFree {
		return true
	}
	return priceFilter == domain.FilterBoth
}

func passFilter(m domain.Movie, filter domain.MovieFilter) bool {
	return m.Rating >= filter.MinRating &&
		hasIntersection(m.Countries, filter.Countries) &&
		passFree(m.IsFree, filter.IsFree) &&
		hasIntersection(m.Genres, filter.Genres) &&
		hasIntersection(m.Actors, filter.Actors) &&
		(filter.Type == "" || m.Type == filter.Type) &&
		hasIntersection(m.Director, filter.Director)
}

func (m *mapMovieRepository) GetByFilter(filter domain.MovieFilter) ([]domain.Movie, error) {
	m.Lock()
	defer m.Unlock()
	var res []domain.Movie
	for _, movieValue := range m.movies {
		if passFilter(movieValue, filter) {
			res = append(res, movieValue)
		}
	}
	if len(res) == 0 {
		return nil, movie.NotFoundError
	}
	left, right := baseutils.SafePage(len(res), filter.Offset, filter.Limit)
	return res[left:right], nil
}

func (m *mapMovieRepository) fillMap() {
	count := uint(100)
	for i := uint(1); i < count; i++ {
		mov := movie_generator.RandomMovie(i)
		mov.Avatar = "/static/movies/default.jpg"
		mov.Video = "/static/movies/default.mpeg"
		m.movies[i] = mov
	}
}

func (m *mapMovieRepository) GetById(id uint) (domain.Movie, error) {
	m.Lock()
	mov, exists := m.movies[id]
	m.Unlock()
	if !exists {
		return domain.Movie{}, errors.New("movie not found")
	}
	return mov, nil
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

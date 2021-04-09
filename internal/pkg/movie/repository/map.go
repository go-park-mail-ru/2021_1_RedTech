package repository

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/movie"
	"Redioteka/internal/pkg/utils/baseutils"
	"Redioteka/internal/pkg/utils/moviegen"
	"errors"
	"sort"
	"sync"
)

type mapMovieRepository struct {
	sync.Mutex
	movies map[uint]domain.Movie
}

func (m *mapMovieRepository) GetGenres() ([]string, error) {
	panic("delete maap")
}

func (m *mapMovieRepository) GetStream(id uint) (domain.Stream, error) {
	panic("delete this whole file")
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

func isFilterValid(filter domain.MovieFilter) bool {
	return filter.Offset >= 0 && filter.Limit >= 0 &&
		(filter.IsFree == domain.FilterBoth ||
			filter.IsFree == domain.FilterFree ||
			filter.IsFree == domain.FilterSubscription) &&
		(filter.Type == "" ||
			filter.Type == domain.MovieT ||
			filter.Type == domain.SeriesT)
}

func (m *mapMovieRepository) GetByFilter(filter domain.MovieFilter) ([]domain.Movie, error) {
	if !isFilterValid(filter) {
		return nil, movie.InvalidFilterError
	}
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
	sort.Slice(res, func(i, j int) bool {
		return res[i].ID < res[j].ID
	})
	left, right := baseutils.SafePage(len(res), filter.Offset, filter.Limit)
	return res[left:right], nil
}

func (m *mapMovieRepository) fillMap() {
	count := uint(100)
	for i := uint(1); i < count; i++ {
		mov := moviegen.RandomMovie(i)
		mov.Avatar = "/static/movies/default.jpg"
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

func (m *mapMovieRepository) AddFavouriteByID(movieID, userID uint) error {
	return errors.New("TO DO")
}

func (m *mapMovieRepository) RemoveFavouriteByID(movieID, userID uint) error {
	return errors.New("TO DO")
}

func (m *mapMovieRepository) CheckFavouriteByID(movieID, userID uint) error {
	return errors.New("TO DO")
}

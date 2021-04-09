package usecase

import "Redioteka/internal/pkg/domain"

type movieUsecase struct {
	movieRepo domain.MovieRepository
}
func NewMovieUsecase(m domain.MovieRepository) domain.MovieUsecase {
	return &movieUsecase{
		movieRepo: m,
	}
}

func (m *movieUsecase) GetById(id uint) (domain.Movie, error) {
	return m.movieRepo.GetById(id)
}

func (m *movieUsecase) GetByFilter(filter domain.MovieFilter) ([]domain.Movie, error) {
	return m.movieRepo.GetByFilter(filter)
}

func (m *movieUsecase) GetStream(id uint) (domain.Stream, error) {
	return m.movieRepo.GetStream(id)
}


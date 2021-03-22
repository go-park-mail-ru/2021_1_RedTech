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

func (m *movieUsecase) Delete(id uint) error {
	return m.movieRepo.Delete(id)
}

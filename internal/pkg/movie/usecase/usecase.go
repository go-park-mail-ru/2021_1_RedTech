package usecase

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/user"
	"Redioteka/internal/pkg/utils/session"
)

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

func (m *movieUsecase) AddFavourite(id uint, sess *session.Session) error {
	err := session.Manager.Check(sess)
	if err != nil {
		return user.UnauthorizedError
	}

	err = m.movieRepo.CheckFavouriteByID(id, sess.UserID)
	if err != nil {
		return nil
	}

	return m.movieRepo.AddFavouriteByID(id, sess.UserID)
}

func (m *movieUsecase) RemoveFavourite(id uint, sess *session.Session) error {
	err := session.Manager.Check(sess)
	if err != nil {
		return user.UnauthorizedError
	}

	return m.movieRepo.RemoveFavouriteByID(id, sess.UserID)
}

func (m *movieUsecase) GetByFilter(filter domain.MovieFilter) ([]domain.Movie, error) {
	return m.movieRepo.GetByFilter(filter)
}

func (m *movieUsecase) GetGenres() ([]string, error) {
	return m.movieRepo.GetGenres()
}

func (m *movieUsecase) GetStream(id uint) (domain.Stream, error) {
	return m.movieRepo.GetStream(id)
}

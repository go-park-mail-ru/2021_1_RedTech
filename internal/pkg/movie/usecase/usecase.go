package usecase

import (
	"Redioteka/internal/pkg/domain"
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
		return err
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
		return err
	}

	return m.movieRepo.RemoveFavouriteByID(id, sess.UserID)
}

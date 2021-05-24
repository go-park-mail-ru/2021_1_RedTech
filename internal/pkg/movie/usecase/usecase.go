package usecase

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/movie"
	"Redioteka/internal/pkg/user"
	"Redioteka/internal/pkg/utils/session"
	"time"
)

type movieUsecase struct {
	movieRepo domain.MovieRepository
	userRepo  domain.UserRepository
	actorRepo domain.ActorRepository
}

func NewMovieUsecase(m domain.MovieRepository, u domain.UserRepository, a domain.ActorRepository) domain.MovieUsecase {
	return &movieUsecase{
		movieRepo: m,
		userRepo:  u,
		actorRepo: a,
	}
}

func (m *movieUsecase) GetByID(id uint, sess *session.Session) (domain.Movie, error) {
	foundMovie, err := m.movieRepo.GetById(id)
	if err != nil {
		return domain.Movie{}, err
	}
	foundMovie.Actors, err = m.actorRepo.GetByMovie(id)
	if err != nil {
		return domain.Movie{}, err
	}

	if foundMovie.Type == domain.SeriesT {
		foundMovie.Series, err = m.movieRepo.GetSeriesList(id)
		if err != nil {
			return domain.Movie{}, err
		}
	}

	err = session.Manager.Check(sess)
	if err == nil {
		err = m.movieRepo.CheckFavouriteByID(id, sess.UserID)
		if err == movie.AlreadyExists {
			foundMovie.Favourite = 1
		}
		foundMovie.Vote = m.movieRepo.CheckVoteByID(id, sess.UserID)

		subExpire := m.userRepo.CheckSub(sess.UserID)
		if subExpire.Sub(time.Now()) > 0 {
			foundMovie.Availability = 1
		}
	} else {
		foundMovie.Availability = -1
	}
	return foundMovie, nil
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

func (m *movieUsecase) GetGenres() ([]domain.Genre, error) {
	return m.movieRepo.GetGenres()
}

func (m *movieUsecase) GetStream(id uint) ([]domain.Stream, error) {
	return m.movieRepo.GetStream(id)
}

func (m *movieUsecase) Like(userId, movieId uint) error {
	return m.movieRepo.Like(userId, movieId)
}

func (m *movieUsecase) Dislike(userId, movieId uint) error {
	return m.movieRepo.Dislike(userId, movieId)
}

func (m *movieUsecase) Search(query string) ([]domain.Movie, error) {
	return m.movieRepo.Search(query)
}

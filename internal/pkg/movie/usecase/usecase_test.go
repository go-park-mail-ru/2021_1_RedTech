package usecase

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/movie"
	"Redioteka/internal/pkg/movie/repository/mock"
	"Redioteka/internal/pkg/user"
	"Redioteka/internal/pkg/utils/session"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type getByIdTestCase struct {
	ID       uint
	outMovie domain.Movie
	outErr   error
}

var getByIdTests = []getByIdTestCase{
	{
		ID: 1,
		outMovie: domain.Movie{
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
		},
		outErr: nil,
	},
	{
		ID:       2,
		outMovie: domain.Movie{},
		outErr:   movie.NotFoundError,
	},
}

func TestMovieUsecase_GetById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mock.NewMockMovieRepository(ctrl)
	um := NewMovieUsecase(repoMock)

	for _, test := range getByIdTests {
		t.Run(fmt.Sprintf("ID: %v", test.ID),
			func(t *testing.T) {
				repoMock.EXPECT().GetById(test.ID).Times(1).Return(test.outMovie, test.outErr)
				currentMovie, currentErr := um.GetById(test.ID)
				if currentErr != nil {
					require.Equal(t, currentErr, test.outErr)
				} else {
					require.Equal(t, currentMovie, test.outMovie)
				}
			})
	}
}

type addFavouriteTestCase struct {
	sess     *session.Session
	movieID  uint
	checkErr error
	outErr   error
}

var addFavouriteTests = []addFavouriteTestCase{
	{
		&session.Session{},
		1,
		movie.AlreadyExists,
		user.UnauthorizedError,
	},
	{
		&session.Session{UserID: 1},
		2,
		movie.AlreadyExists,
		nil,
	},
	{
		&session.Session{UserID: 3},
		4,
		nil,
		nil,
	},
	{
		&session.Session{UserID: 5},
		6,
		nil,
		movie.NotFoundError,
	},
}

func TestMovieUsecase_AddFavourite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mock.NewMockMovieRepository(ctrl)
	um := NewMovieUsecase(repoMock)

	for _, test := range addFavouriteTests {
		t.Run(fmt.Sprintf("userID: %v movieID: %v err: %v", test.sess.UserID, test.movieID, test.outErr),
			func(t *testing.T) {
				if test.outErr != user.UnauthorizedError {
					err := session.Manager.Create(test.sess)
					require.NoError(t, err)
					repoMock.EXPECT().CheckFavouriteByID(test.movieID, test.sess.UserID).Times(1).Return(test.checkErr)
				}
				if test.checkErr == nil {
					repoMock.EXPECT().AddFavouriteByID(test.movieID, test.sess.UserID).Times(1).Return(test.outErr)
				}
				currentErr := um.AddFavourite(test.movieID, test.sess)
				require.Equal(t, test.outErr, currentErr)
			})
	}
}

type removeFavouriteTestCase struct {
	sess    *session.Session
	movieID uint
	outErr  error
}

var removeFavouriteTests = []removeFavouriteTestCase{
	{
		&session.Session{},
		1,
		user.UnauthorizedError,
	},
	{
		&session.Session{UserID: 1},
		2,
		nil,
	},
	{
		&session.Session{UserID: 3},
		4,
		movie.NotFoundError,
	},
}

func TestMovieUsecase_RemoveFavourite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mock.NewMockMovieRepository(ctrl)
	um := NewMovieUsecase(repoMock)

	for _, test := range addFavouriteTests {
		t.Run(fmt.Sprintf("userID: %v movieID: %v err: %v", test.sess.UserID, test.movieID, test.outErr),
			func(t *testing.T) {
				if test.outErr != user.UnauthorizedError {
					err := session.Manager.Create(test.sess)
					require.NoError(t, err)
					repoMock.EXPECT().RemoveFavouriteByID(test.movieID, test.sess.UserID).Times(1).Return(test.outErr)
				}
				currentErr := um.RemoveFavourite(test.movieID, test.sess)
				require.Equal(t, test.outErr, currentErr)
			})
	}
}

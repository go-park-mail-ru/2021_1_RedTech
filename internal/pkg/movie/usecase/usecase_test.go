package usecase

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/movie"
	"Redioteka/internal/pkg/movie/repository/mock"
	"Redioteka/internal/pkg/user"
	userMock "Redioteka/internal/pkg/user/repository/mock"
	"Redioteka/internal/pkg/utils/session"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type getByIdTestCase struct {
	id       uint
	sess     *session.Session
	outMovie domain.Movie
	outError error
}

var getByIdTest = []getByIdTestCase{
	{
		id:   1,
		sess: &session.Session{},
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
		outError: nil,
	},
	{
		id:   2,
		sess: &session.Session{UserID: 2},
		outMovie: domain.Movie{
			ID:          2,
			Title:       "Another film",
			Description: "Test 2",
			Rating:      0,
			Countries:   []string{"Germany"},
			IsFree:      true,
			Genres:      []string{"Cartoon"},
			Actors:      []string{"Anna"},
			Avatar:      "/static/movies/default.jpg",
			Type:        "Сериал",
			Year:        "2006",
			Director:    []string{"Florian Henckel von Donnersmarck"},
			Favourite:   1,
			Vote:        domain.Like,
			Series:      []uint{7, 3},
		},
	},
	{
		id:       1000,
		sess:     &session.Session{},
		outError: movie.NotFoundError,
	},
}

func TestMovieUsecase_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	movieRepoMock := mock.NewMockMovieRepository(ctrl)
	userRepoMock := userMock.NewMockUserRepository(ctrl)
	uc := NewMovieUsecase(movieRepoMock, userRepoMock)

	for testID, test := range getByIdTest {
		t.Run(fmt.Sprintln(testID, test.outError), func(t *testing.T) {
			movieRepoMock.EXPECT().GetById(test.id).Times(1).Return(test.outMovie, test.outError)
			if test.outMovie.Series != nil {
				movieRepoMock.EXPECT().GetSeriesList(test.id).Times(1).Return(test.outMovie.Series, test.outError)
			}
			if test.sess.UserID != 0 {
				err := session.Manager.Create(test.sess)
				require.NoError(t, err)
				movieRepoMock.EXPECT().CheckFavouriteByID(test.id, test.sess.UserID).Times(1).Return(movie.AlreadyExists)
				movieRepoMock.EXPECT().CheckVoteByID(test.id, test.sess.UserID).Times(1).Return(domain.Like)
			}
			currentMovie, currentError := uc.GetByID(test.id, test.sess)
			require.Equal(t, currentError, test.outError)
			require.Equal(t, currentMovie, test.outMovie)
		})
	}
}

type getByFilterTestCase struct {
	filter   domain.MovieFilter
	outMovie domain.Movie
	outError error
}

var getByFilterTests = []getByFilterTestCase{}

func TestMovieUsecase_GetByFilter(t *testing.T) {
}

type addFavouriteTestCase struct {
	sess     *session.Session
	movieID  uint
	checkErr error
	outErr   error
}

var addFavouriteTests = []addFavouriteTestCase{
	{
		sess:     &session.Session{},
		movieID:  1,
		checkErr: movie.AlreadyExists,
		outErr:   user.UnauthorizedError,
	},
	{
		sess:     &session.Session{UserID: 1},
		movieID:  2,
		checkErr: movie.AlreadyExists,
		outErr:   nil,
	},
	{
		sess:     &session.Session{UserID: 3},
		movieID:  4,
		checkErr: nil,
		outErr:   nil,
	},
	{
		sess:     &session.Session{UserID: 5},
		movieID:  6,
		checkErr: nil,
		outErr:   movie.NotFoundError,
	},
}

func TestMovieUsecase_AddFavourite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mock.NewMockMovieRepository(ctrl)
	userRepoMock := userMock.NewMockUserRepository(ctrl)
	um := NewMovieUsecase(repoMock, userRepoMock)

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
		sess:    &session.Session{},
		movieID: 1,
		outErr:  user.UnauthorizedError,
	},
	{
		sess:    &session.Session{UserID: 1},
		movieID: 2,
		outErr:  nil,
	},
	{
		sess:    &session.Session{UserID: 3},
		movieID: 4,
		outErr:  movie.NotFoundError,
	},
}

func TestMovieUsecase_RemoveFavourite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mock.NewMockMovieRepository(ctrl)
	userRepoMock := userMock.NewMockUserRepository(ctrl)
	um := NewMovieUsecase(repoMock, userRepoMock)

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

func TestMovieUsecase_Like(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mock.NewMockMovieRepository(ctrl)
	userRepoMock := userMock.NewMockUserRepository(ctrl)
	um := NewMovieUsecase(repoMock, userRepoMock)
	repoMock.EXPECT().Like(uint(1), uint(1)).Times(1).Return(nil)
	require.NoError(t, um.Like(uint(1), uint(1)))
}

func TestMovieUsecase_Dislike(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mock.NewMockMovieRepository(ctrl)
	userRepoMock := userMock.NewMockUserRepository(ctrl)
	um := NewMovieUsecase(repoMock, userRepoMock)
	repoMock.EXPECT().Dislike(uint(1), uint(1)).Times(1).Return(nil)
	require.NoError(t, um.Dislike(uint(1), uint(1)))
}

func TestMovieUsecase_Search(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mock.NewMockMovieRepository(ctrl)
	userRepoMock := userMock.NewMockUserRepository(ctrl)
	um := NewMovieUsecase(repoMock, userRepoMock)
	repoMock.EXPECT().Search("Film").Times(1).Return(nil, nil)
	res, err := um.Search("Film")
	require.NoError(t, err)
	require.Equal(t, []domain.Movie(nil), res)
}

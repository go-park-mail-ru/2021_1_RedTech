package usecase

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/movie"
	"Redioteka/internal/pkg/movie/repository/mock"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

type getByIdTestCase struct {
	id       uint
	outMovie domain.Movie
	outError error
}

var getByIdTest = []getByIdTestCase{
	{
		id: 1,
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
		id:       1000,
		outError: movie.NotFoundError,
	},
}

func TestMovieUsecase_GetById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	movieRepoMock := mock.NewMockMovieRepository(ctrl)
	uc := NewMovieUsecase(movieRepoMock)

	for testId, test := range getByIdTest {
		t.Run(fmt.Sprintln(testId, test.outError), func(t *testing.T) {
			movieRepoMock.EXPECT().GetById(test.id).Times(1).Return(test.outMovie, test.outError)
		})
		currentMovie, currentError := uc.GetById(test.id)
		require.Equal(t, currentError, test.outError)
		require.Equal(t, currentMovie, test.outMovie)
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

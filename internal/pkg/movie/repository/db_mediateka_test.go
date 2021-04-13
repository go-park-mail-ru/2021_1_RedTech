package repository

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/movie"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

type filterTestStruct struct {
	TestName string
	Filter   domain.MovieFilter
	Error    error
	Movies   []domain.Movie
}

var getByFilterTests = []filterTestStruct{
	{
		TestName: "invalid filter",
		Filter: domain.MovieFilter{
			Limit:  100,
			Offset: -100,
		},
		Error: movie.InvalidFilterError,
		Movies: nil,
	},
	{
		TestName: "invalid filter",
		Filter: domain.MovieFilter{
			Limit:  100,
			Offset: -100,
		},
		Error: movie.InvalidFilterError,
		Movies: nil,
	},
}

func TestDbMovieRepository_GetByFilter(t *testing.T) {
	db, mock := NewMock()
	repo := NewMovieRepository(db)
	defer mock.Close()

	for testId, test := range getByFilterTests {
		t.Run(fmt.Sprintln(testId, ":", test.TestName), func(t *testing.T) {
			movies, err := repo.GetByFilter(test.Filter)

			require.Equal(t, test.Error, err)
			require.Equal(t, test.Movies, movies)
		})
	}
	//mock.ExpectQuery(regexp.QuoteMeta(querySelectFav)).WithArgs(userID, movieID).WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectCommit()
	//	movies := repo.GetByFilter(testFilter)
	//	require.Nil(t, err)
	//	require.NoError(t, mock.ExpectationsWereMet())
}

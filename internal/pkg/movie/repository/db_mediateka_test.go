package repository

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/movie"
	"Redioteka/internal/pkg/utils/cast"
	"errors"
	"fmt"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

type filterTestStruct struct {
	TestName string
	Filter   domain.MovieFilter
	Error    error
	Movies   []domain.Movie
	Query    string
}

var getByFilterTests = []filterTestStruct{
	{
		TestName: "all ok ",
		Filter: domain.MovieFilter{
			MinRating: 10,
			Genres:    []string{"horror"},
			IsFree:    domain.FilterFree,
			Type:      domain.MovieT,
			Order:     domain.DateOrder,
			Limit:     100,
			Offset:    10,
		},
		Movies: []domain.Movie{
			{
				ID:          1,
				Title:       "Movie1",
				Description: "description1",
				Avatar:      "path/to/avka",
				IsFree:      true,
			},
		},
		Error: nil,
	},
	{
		TestName: "db error",
		Filter: domain.MovieFilter{
			MinRating: 10,
			Genres:    []string{"horror"},
			IsFree:    domain.FilterFree,
			Type:      domain.MovieT,
			Order:     domain.DateOrder,
			Limit:     100,
			Offset:    10,
		},
		Error: movie.NotFoundError,
	},
	{
		TestName: "db errror",
		Filter: domain.MovieFilter{
			MinRating: 10,
			Genres:    []string{"horror"},
			IsFree:    domain.FilterFree,
			Type:      domain.MovieT,
			Order:     domain.DateOrder,
			Limit:     -100,
			Offset:    -10,
		},
		Error: movie.InvalidFilterError,
	},
}

func TestDbMovieRepository_GetByFilter(t *testing.T) {
	db, mock := NewMock()
	repo := NewMovieRepository(db)
	defer mock.Close()

	for testId, test := range getByFilterTests {
		t.Run(fmt.Sprintln(testId, ":", test.TestName), func(t *testing.T) {
			query, args, _ := buildFilterQuery(test.Filter)
			mock.ExpectBegin()
			if test.Error == nil {
				rows := pgxmock.NewRows([]string{"movies.id", "movies.title", "movies.description", "movies.avatar", "is_free"}).
					AddRow(cast.UintToBytes(uint(1)), cast.StrToBytes("Movie1"), cast.StrToBytes("description1"), cast.StrToBytes("path/to/avka"), cast.BoolToBytes(true))
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(args...).WillReturnRows(rows)
			}
			if test.Error == movie.NotFoundError {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(errors.New("some sql error"))
			}
			mock.ExpectCommit()
			fmt.Println(query)
			fmt.Println(args)
			movies, err := repo.GetByFilter(test.Filter)

			require.Equal(t, test.Error, err)
			if test.Error == nil {
				require.Equal(t, test.Movies, movies)
			}
		})
	}
}

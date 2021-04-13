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
	Query    string
	Error    error
	Movies   []domain.Movie
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
		Error:  movie.InvalidFilterError,
		Movies: nil,
		Query:  `SELECT movies.id, movies.title, movies.description, movies.avatar, is_free FROM movies JOIN movie_actors ma on movies.id = ma.movie_id JOIN movie_genres mg on movies.id = mg.movie_id JOIN movie_types mt on movies.type = mt.id JOIN genres g on mg.genre_id = g.id JOIN (select a.id, a.firstname || ' ' || a.lastname as full_actor_name from actors as a) full_acts on full_acts.id = ma.actor_id WHERE rating >= $1 AND lower(g.name) IN ($2) AND is_free = $3 AND mt.type = $4 GROUP BY movies.id ORDER BY movies.add_date desc LIMIT 100 OFFSET 10`,
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
			mock.ExpectQuery(query).WithArgs(args).WillReturnRows()
			mock.ExpectCommit()
			movies, err := repo.GetByFilter(test.Filter)

			require.Equal(t, test.Error, err)
			require.Equal(t, test.Movies, movies)
		})
	}
	//mock.ExpectQuery(regexp.QuoteMeta(querySelectFav)).WithArgs(userID, movieID).WillReturnRows(rows)

	//	movies := repo.GetByFilter(testFilter)
	//	require.Nil(t, err)
	//	require.NoError(t, mock.ExpectationsWereMet())
}

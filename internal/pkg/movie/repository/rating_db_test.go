package repository

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/movie"
	"Redioteka/internal/pkg/utils/cast"
	"fmt"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/require"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

type voteTestCase struct {
	name    string
	userId  uint
	movieId uint
	err     error
}

var likeTests = []voteTestCase{
	{
		name:    "all ok",
		userId:  1,
		movieId: 1,
		err:     nil,
	},
	{
		name:    "vote query error",
		userId:  1,
		movieId: 1,
		err:     movie.InvalidVoteError,
	},
	{
		name:    "vote query error",
		userId:  1,
		movieId: 1,
	},
}

func TestDbMovieRepository_Like(t *testing.T) {
	mockDb, dbMock := NewMock()
	movieRepo := NewMovieRepository(mockDb)
	defer dbMock.Close()

	for testId, test := range likeTests {
		t.Run(fmt.Sprintln(testId, test.name), func(t *testing.T) {
			dbMock.ExpectBegin()
			dbMock.ExpectQuery(regexp.QuoteMeta(queryCheckView)).WithArgs(test.userId, test.movieId)
			dbMock.ExpectCommit()

			movieRepo.Like(test.userId, test.movieId)
		})
	}

}

func TestExample(t *testing.T) {
	db, mock := NewMock()
	repo := NewMovieRepository(db)
	defer mock.Close()

	var m = domain.Movie{
		ID:          1,
		Title:       "Film",
		Description: "Test data",
		Rating:      9,
		Countries:   []string{"Japan", "South Korea"},
		IsFree:      true,
		Genres:      []string{"Comedy"},
		Actors:      []string{"Sana", "Momo", "Mina"},
		Avatar:      "/static/movies/default.jpg",
		Type:        domain.MovieT,
		Year:        "0",
		Director:    []string{"James Cameron"},
	}
	year, _ := strconv.Atoi(m.Year)
	rows := pgxmock.NewRows([]string{"m.id", "m.title", "m.description", "m.avatar", "m.rating", "m.countries",
		"m.directors", "m.release_year", "m.is_free", "mt.type", "acts", "gns"}).
		AddRow(cast.UintToBytes(m.ID), cast.StrToBytes(m.Title), cast.StrToBytes(m.Description), cast.StrToBytes(m.Avatar),
			cast.FloatToBytes(m.Rating), cast.StrToBytes(strings.Join(m.Countries, ", ")), cast.StrToBytes(strings.Join(m.Director, ", ")),
			cast.SmallIntToBytes(year), cast.BoolToBytes(m.IsFree), cast.StrToBytes(string(m.Type)), cast.StrToBytes(strings.Join(m.Actors, ";")),
			cast.StrToBytes(strings.Join(m.Genres, ";")))

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(querySelectID)).WithArgs(m.ID).WillReturnRows(rows)
	mock.ExpectCommit()

	actual, err := repo.GetById(m.ID)
	require.NoError(t, err)
	require.Equal(t, m, actual)
	require.NoError(t, mock.ExpectationsWereMet())
}

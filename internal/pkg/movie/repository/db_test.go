package repository

import (
	"Redioteka/internal/pkg/database"
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/utils/cast"
	"Redioteka/internal/pkg/utils/log"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/require"
)

func NewMock() (*database.DBManager, pgxmock.PgxPoolIface) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		log.Log.Error(err)
	}
	return &database.DBManager{Pool: mock}, mock
}

func TestGetByIDSuccess(t *testing.T) {
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
		Year:        "2012",
		Director:    []string{"James Cameron"},
	}
	year, _ := strconv.Atoi(m.Year)
	bytes := cast.SmallIntToBytes(year)
	ret := cast.ToSmallInt(bytes)
	fmt.Println(year, bytes, ret)
	rows := pgxmock.NewRows([]string{"m.id", "m.title", "m.description", "m.avatar", "m.rating", "m.countries",
		"m.directors", "m.release_year", "m.price", "mt.type", "acts", "gns"}).
		AddRow(cast.UintToBytes(m.ID), cast.StrToBytes(m.Title), cast.StrToBytes(m.Description), cast.StrToBytes(m.Avatar),
			cast.FloatToBytes(m.Rating), cast.StrToBytes(strings.Join(m.Countries, ", ")), cast.StrToBytes(strings.Join(m.Director, ", ")),
			cast.SmallIntToBytes(year), cast.FloatToBytes(0), cast.StrToBytes(string(m.Type)), cast.StrToBytes(strings.Join(m.Actors, ";")),
			cast.StrToBytes(strings.Join(m.Genres, ";")))

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(querySelectID)).WithArgs(m.ID).WillReturnRows(rows)
	mock.ExpectCommit()

	actual, err := repo.GetById(m.ID)
	require.NoError(t, err)
	require.Equal(t, m, actual)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByIDFailure(t *testing.T) {
	db, mock := NewMock()
	repo := NewMovieRepository(db)
	defer mock.Close()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(querySelectID)).WithArgs(uint(0))
	mock.ExpectRollback()

	actual, err := repo.GetById(0)
	require.NotNil(t, err)
	require.Equal(t, domain.Movie{}, actual)
	require.NoError(t, mock.ExpectationsWereMet())
}

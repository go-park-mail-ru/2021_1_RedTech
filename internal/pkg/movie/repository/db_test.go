package repository

import (
	"Redioteka/internal/pkg/database"
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/movie"
	"Redioteka/internal/pkg/utils/baseutils"
	"Redioteka/internal/pkg/utils/cast"
	"Redioteka/internal/pkg/utils/log"
	"errors"
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
		Avatar:      "/static/movies/default.jpg",
		Type:        domain.MovieT,
		Year:        "0",
		Director:    []string{"James Cameron"},
	}
	year, _ := strconv.Atoi(m.Year)
	rows := pgxmock.NewRows([]string{"m.id", "m.title", "m.description", "m.avatar", "m.rating", "m.countries",
		"m.directors", "m.release_year", "m.is_free", "mt.type", "gns"}).
		AddRow(cast.UintToBytes(m.ID), cast.StrToBytes(m.Title), cast.StrToBytes(m.Description), cast.StrToBytes(m.Avatar),
			cast.FloatToBytes(m.Rating), cast.StrToBytes(strings.Join(m.Countries, ", ")), cast.StrToBytes(strings.Join(m.Director, ", ")),
			cast.SmallIntToBytes(year), cast.BoolToBytes(m.IsFree), cast.StrToBytes(string(m.Type)),
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
	mock.ExpectQuery(regexp.QuoteMeta(querySelectID)).WithArgs(uint(0)).WillReturnError(errors.New(""))
	mock.ExpectRollback()

	actual, err := repo.GetById(0)
	require.NotNil(t, err)
	require.Equal(t, domain.Movie{}, actual)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGestSeriesListSuccess(t *testing.T) {
	db, mock := NewMock()
	repo := NewMovieRepository(db)
	defer mock.Close()

	var id uint = 1
	expected := []uint{7, 3}
	rows := pgxmock.NewRows([]string{"count"}).AddRow(cast.Uint64ToBytes(uint64(expected[0]))).
		AddRow(cast.Uint64ToBytes(uint64(expected[1])))

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(querySelectSeries)).WithArgs(id).WillReturnRows(rows)
	mock.ExpectCommit()

	actual, err := repo.GetSeriesList(id)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGestSeriesListFailure(t *testing.T) {
	db, mock := NewMock()
	repo := NewMovieRepository(db)
	defer mock.Close()

	var id uint = 1

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(querySelectSeries)).WithArgs(id).WillReturnError(errors.New(""))
	mock.ExpectRollback()

	actual, err := repo.GetSeriesList(id)
	require.NotNil(t, err)
	require.Nil(t, actual)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAddFavouriteByIDSuccess(t *testing.T) {
	db, mock := NewMock()
	repo := NewMovieRepository(db)
	defer mock.Close()

	var movieID, userID uint = 1, 1

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(queryInsertFav)).WithArgs(userID, movieID).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectCommit()

	err := repo.AddFavouriteByID(movieID, userID)
	require.Nil(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAddFavouriteByIDFailure(t *testing.T) {
	db, mock := NewMock()
	repo := NewMovieRepository(db)
	defer mock.Close()

	var movieID, userID uint = 0, 5

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(queryInsertFav)).WithArgs(userID, movieID).
		WillReturnError(errors.New(""))
	mock.ExpectRollback()

	err := repo.AddFavouriteByID(movieID, userID)
	require.NotNil(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRemoveFavouriteByIDSuccess(t *testing.T) {
	db, mock := NewMock()
	repo := NewMovieRepository(db)
	defer mock.Close()

	var movieID, userID uint = 1, 1

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(queryDeleteFav)).WithArgs(userID, movieID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()

	err := repo.RemoveFavouriteByID(movieID, userID)
	require.Nil(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRemoveFavouriteByIDFailure(t *testing.T) {
	db, mock := NewMock()
	repo := NewMovieRepository(db)
	defer mock.Close()

	var movieID, userID uint = 0, 5

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(queryDeleteFav)).WithArgs(userID, movieID).
		WillReturnError(errors.New(""))
	mock.ExpectRollback()

	err := repo.RemoveFavouriteByID(movieID, userID)
	require.NotNil(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckFavouriteByIDSuccess(t *testing.T) {
	db, mock := NewMock()
	repo := NewMovieRepository(db)
	defer mock.Close()

	var movieID, userID uint = 1, 1
	rows := pgxmock.NewRows([]string{"id"})

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(querySelectFav)).WithArgs(userID, movieID).WillReturnRows(rows)
	mock.ExpectCommit()

	err := repo.CheckFavouriteByID(movieID, userID)
	require.Nil(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckFavouriteByIDFailure(t *testing.T) {
	db, mock := NewMock()
	repo := NewMovieRepository(db)
	defer mock.Close()

	var movieID, userID uint = 0, 5
	var expectedID uint = 1
	rows := pgxmock.NewRows([]string{"id"}).AddRow(cast.UintToBytes(expectedID))

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(querySelectFav)).WithArgs(userID, movieID).WillReturnRows(rows)
	mock.ExpectCommit()

	err := repo.CheckFavouriteByID(movieID, userID)
	require.NotNil(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckVoteByIDSuccess(t *testing.T) {
	db, mock := NewMock()
	repo := NewMovieRepository(db)
	defer mock.Close()

	var movieID, userID uint = 1, 1
	var expectedVote int = domain.Dislike
	rows := pgxmock.NewRows([]string{"value"}).AddRow(cast.SmallIntToBytes(expectedVote))

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(querySelectVote)).WithArgs(userID, movieID).WillReturnRows(rows)
	mock.ExpectCommit()

	actualVote := repo.CheckVoteByID(movieID, userID)
	require.Equal(t, expectedVote, actualVote)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckVoteByIDFailure(t *testing.T) {
	db, mock := NewMock()
	repo := NewMovieRepository(db)
	defer mock.Close()

	var movieID, userID uint = 0, 5
	var expectedVote int = 0

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(querySelectVote)).WithArgs(userID, movieID).WillReturnError(errors.New(""))
	mock.ExpectRollback()

	actualVote := repo.CheckVoteByID(movieID, userID)
	require.Equal(t, expectedVote, actualVote)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDbMovieRepository_Search(t *testing.T) {
	db, mock := NewMock()
	repo := NewMovieRepository(db)
	defer mock.Close()

	failQuery := "folm"
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(querySearch)).WithArgs(baseutils.PrepareQueryForSearch(failQuery)).WillReturnError(errors.New(""))
	mock.ExpectRollback()

	movies, err := repo.Search(failQuery)
	require.Equal(t, movie.NotFoundError, err)
	require.Equal(t, []domain.Movie(nil), movies)

	okQuery := "Film"
	var m = domain.Movie{
		ID:          1,
		Title:       "Film",
		Description: "Test data",
		IsFree:      true,
		Type:        domain.MovieT,
	}
	rows := pgxmock.NewRows([]string{"m.id", "m.title", "m.description", "m.avatar", "m.is_free", "mt.type"}).AddRow(
		cast.UintToBytes(m.ID), cast.StrToBytes(m.Title), cast.StrToBytes(m.Description), cast.StrToBytes(m.Avatar),
		cast.BoolToBytes(m.IsFree), cast.StrToBytes(string(m.Type)))
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(querySearch)).WithArgs(baseutils.PrepareQueryForSearch(okQuery)).WillReturnRows(rows)
	mock.ExpectCommit()

	movies, err = repo.Search(okQuery)
	require.NoError(t, err)
	require.Equal(t, []domain.Movie{m}, movies)
}

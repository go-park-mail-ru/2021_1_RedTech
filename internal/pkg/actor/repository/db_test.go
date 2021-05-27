package repository

import (
	actor2 "Redioteka/internal/pkg/actor"
	"Redioteka/internal/pkg/database"
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/utils/baseutils"
	"Redioteka/internal/pkg/utils/cast"
	"Redioteka/internal/pkg/utils/log"
	"errors"
	"regexp"
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

func TestDbActorRepository_GetById(t *testing.T) {
	db, mock := NewMock()
	repo := NewActorRepository(db)
	defer mock.Close()

	actor := domain.Actor{
		ID:        1,
		FirstName: "FName",
		LastName:  "LName",
		Born:      "Moscow",
		Avatar:    "path",
	}
	movie := domain.Movie{
		ID:          1,
		Title:       "Film",
		Description: "desc",
		Avatar:      "path",
	}

	actorRows := pgxmock.NewRows([]string{
		"id", "firstname", "lastname", "born", "avatar",
	}).AddRow(cast.UintToBytes(actor.ID), cast.StrToBytes(actor.FirstName),
		cast.StrToBytes(actor.LastName), cast.StrToBytes(actor.Born),
		cast.StrToBytes(actor.Avatar))

	movieRows := pgxmock.NewRows([]string{
		"m.id", "m.title", "m.description", "m.avatar",
	}).AddRow(cast.UintToBytes(movie.ID), cast.StrToBytes(movie.Title),
		cast.StrToBytes(movie.Description), cast.StrToBytes(movie.Avatar))
	id := uint(1)
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(querySelectActor)).WithArgs(id).WillReturnRows(actorRows)
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(querySelectMoviesByActor)).WithArgs(id).WillReturnRows(movieRows)
	mock.ExpectCommit()

	res, err := repo.GetById(id)
	expActor := actor
	expActor.Movies = []domain.Movie{movie}
	require.NoError(t, err)
	require.Equal(t, expActor, res)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(querySelectActor)).WithArgs(id).WillReturnError(errors.New(""))
	mock.ExpectRollback()

	_, err = repo.GetById(id)
	require.Equal(t, actor2.NotFoundError, err)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(querySelectActor)).WithArgs(id).WillReturnRows(actorRows)
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(querySelectMoviesByActor)).WithArgs(id).WillReturnError(errors.New(""))
	mock.ExpectRollback()

	_, err = repo.GetById(id)
	require.Equal(t, actor2.NotFoundError, err)
}

func TestDbActorRepository_Search(t *testing.T) {
	db, mock := NewMock()
	repo := NewActorRepository(db)
	defer mock.Close()

	actor := domain.Actor{
		ID:        1,
		FirstName: "FName",
		LastName:  "LName",
		Born:      "Moscow",
		Avatar:    "path",
	}

	actorRows := pgxmock.NewRows([]string{
		"id", "firstname", "lastname", "born", "avatar",
	}).AddRow(cast.UintToBytes(actor.ID), cast.StrToBytes(actor.FirstName),
		cast.StrToBytes(actor.LastName), cast.StrToBytes(actor.Born),
		cast.StrToBytes(actor.Avatar))

	okQuery := "query"
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(querySearchActors)).WithArgs(baseutils.PrepareQueryForSearch(okQuery)).WillReturnRows(actorRows)
	mock.ExpectCommit()

	res, err := repo.Search(okQuery)
	require.NoError(t, err)
	require.Equal(t, []domain.Actor{actor}, res)

	badQuery := "qqq"
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(querySearchActors)).WithArgs(baseutils.PrepareQueryForSearch(badQuery)).WillReturnError(errors.New(""))
	mock.ExpectRollback()

	_, err = repo.Search(badQuery)
	require.Equal(t, actor2.NotFoundError, err)
}

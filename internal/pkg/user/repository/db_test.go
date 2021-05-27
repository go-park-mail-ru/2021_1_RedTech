package repository

import (
	"Redioteka/internal/pkg/database"
	"Redioteka/internal/pkg/domain"
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

func TestGetByIDSuccess(t *testing.T) {
	db, mock := NewMock()
	repo := NewUserRepository(db)
	defer mock.Close()

	var u = domain.User{
		ID:       3,
		Username: "BestUser",
		Email:    "best@mail.ru",
		Avatar:   "defult.jpg",
	}
	rows := pgxmock.NewRows([]string{"id", "username", "email", "avatar"}).AddRow(cast.UintToBytes(u.ID),
		cast.StrToBytes(u.Username), cast.StrToBytes(u.Email), cast.StrToBytes(u.Avatar))

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(querySelectID)).WithArgs(u.ID).WillReturnRows(rows)
	mock.ExpectCommit()

	actual, err := repo.GetById(u.ID)
	require.NoError(t, err)
	require.Equal(t, u, actual)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByIDFailure(t *testing.T) {
	db, mock := NewMock()
	repo := NewUserRepository(db)
	defer mock.Close()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(querySelectID)).WithArgs(uint(0)).WillReturnError(errors.New(""))
	mock.ExpectRollback()

	actual, err := repo.GetById(0)
	require.NotNil(t, err)
	require.Equal(t, domain.User{}, actual)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByEmailSuccess(t *testing.T) {
	db, mock := NewMock()
	repo := NewUserRepository(db)
	defer mock.Close()

	var u = domain.User{
		ID:       3,
		Username: "BestUser",
		Email:    "best@mail.ru",
		Avatar:   "defult.jpg",
		Password: []byte{'p', 'a', 's', 's'},
	}
	rows := pgxmock.NewRows([]string{"id", "username", "email", "avatar", "password"}).
		AddRow(cast.UintToBytes(u.ID), cast.StrToBytes(u.Username),
			cast.StrToBytes(u.Email), cast.StrToBytes(u.Avatar), u.Password[:])

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(querySelectEmail)).WithArgs(u.Email).WillReturnRows(rows)
	mock.ExpectCommit()

	actual, err := repo.GetByEmail(u.Email)
	require.NoError(t, err)
	require.Equal(t, u, actual)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByEmailFailure(t *testing.T) {
	db, mock := NewMock()
	repo := NewUserRepository(db)
	defer mock.Close()

	email := "not email"

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(querySelectEmail)).WithArgs(email).WillReturnError(errors.New(""))
	mock.ExpectRollback()

	actual, err := repo.GetByEmail(email)
	require.NotNil(t, err)
	require.Equal(t, domain.User{}, actual)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateSuccess(t *testing.T) {
	db, mock := NewMock()
	repo := NewUserRepository(db)
	defer mock.Close()

	var u = &domain.User{
		ID:       3,
		Username: "BestUser",
		Email:    "best@mail.ru",
		Avatar:   "defult.jpg",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(queryUpdate)).WithArgs(u.Username, u.Email, u.Avatar, u.ID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()

	err := repo.Update(u)
	require.Nil(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateFailure(t *testing.T) {
	db, mock := NewMock()
	repo := NewUserRepository(db)
	defer mock.Close()

	var u = &domain.User{
		ID:       3,
		Username: "BestUser",
		Email:    "best@mail.ru",
		Avatar:   "defult.jpg",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(queryUpdate)).WithArgs(u.Username, u.Email, u.Avatar, u.ID).WillReturnError(errors.New(""))
	mock.ExpectRollback()

	err := repo.Update(u)
	require.NotNil(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreSuccess(t *testing.T) {
	db, mock := NewMock()
	repo := NewUserRepository(db)
	defer mock.Close()

	var u = &domain.User{
		Username: "BestUser",
		Email:    "best@mail.ru",
		Avatar:   "defult.jpg",
		Password: []byte{'p', 'a', 's', 's'},
	}
	var expectedID uint = 1
	rows := pgxmock.NewRows([]string{"id"}).AddRow(cast.UintToBytes(expectedID))

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(queryInsert)).WithArgs(u.Username, u.Email, u.Password[:], u.Avatar).
		WillReturnRows(rows)
	mock.ExpectCommit()

	actualID, err := repo.Store(u)
	require.Nil(t, err)
	require.Equal(t, expectedID, actualID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreFailure(t *testing.T) {
	db, mock := NewMock()
	repo := NewUserRepository(db)
	defer mock.Close()

	var u = &domain.User{
		Username: "BestUser",
		Email:    "best@mail.ru",
		Avatar:   "defult.jpg",
		Password: []byte{'p', 'a', 's', 's'},
	}
	var expectedID uint = 0

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(queryInsert)).WithArgs(u.Username, u.Email, u.Password[:], u.Avatar).WillReturnError(errors.New(""))
	mock.ExpectRollback()

	actualID, err := repo.Store(u)
	require.NotNil(t, err)
	require.Equal(t, expectedID, actualID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteSuccess(t *testing.T) {
	db, mock := NewMock()
	repo := NewUserRepository(db)
	defer mock.Close()

	var id uint = 2

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(queryDelete)).WithArgs(id).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()

	err := repo.Delete(id)
	require.Nil(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteFailure(t *testing.T) {
	db, mock := NewMock()
	repo := NewUserRepository(db)
	defer mock.Close()

	var id uint = 2

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(queryDelete)).WithArgs(id).WillReturnError(errors.New(""))
	mock.ExpectRollback()

	err := repo.Delete(id)
	require.NotNil(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetFavouritesByIDSuccess(t *testing.T) {
	db, mock := NewMock()
	repo := NewUserRepository(db)
	defer mock.Close()

	var id uint = 2
	var expected = []domain.Movie{
		{
			ID:     1,
			Title:  "Film",
			Avatar: "/default.jpg",
			Rating: 9,
			IsFree: false,
		},
		{
			ID:     13,
			Title:  "Time to Twice",
			Avatar: "",
			Rating: 666,
			IsFree: true,
		},
	}
	rows := pgxmock.NewRows([]string{"m.id", "m.title", "m.avatar", "m.rating", "m.is_free"}).
		AddRow(cast.UintToBytes(expected[0].ID), cast.StrToBytes(expected[0].Title),
			cast.StrToBytes(expected[0].Avatar), cast.FloatToBytes(expected[0].Rating), cast.BoolToBytes(expected[0].IsFree)).
		AddRow(cast.UintToBytes(expected[1].ID), cast.StrToBytes(expected[1].Title),
			cast.StrToBytes(expected[1].Avatar), cast.FloatToBytes(expected[1].Rating), cast.BoolToBytes(expected[1].IsFree))

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(querySelectFavourites)).WithArgs(id).WillReturnRows(rows)
	mock.ExpectCommit()

	actual, err := repo.GetFavouritesByID(id)
	require.Nil(t, err)
	require.Equal(t, expected, actual)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetFavouritesByIDFailure(t *testing.T) {
	db, mock := NewMock()
	repo := NewUserRepository(db)
	defer mock.Close()

	var id uint = 2

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(querySelectFavourites)).WithArgs(id).WillReturnError(errors.New(""))
	mock.ExpectRollback()

	movies, err := repo.GetFavouritesByID(id)
	require.NotNil(t, err)
	require.Nil(t, movies)
	require.NoError(t, mock.ExpectationsWereMet())
}

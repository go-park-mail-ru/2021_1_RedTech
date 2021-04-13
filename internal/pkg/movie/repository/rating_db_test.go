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

type voteTestCase struct {
	name       string
	userId     uint
	movieId    uint
	err        error
	addViewErr error
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
		name:       "add view check error",
		userId:     1,
		movieId:    1,
		err:        movie.RatingUpdateError,
		addViewErr: movie.InvalidViewCheck,
	},
	{
		name:       "add view add error",
		userId:     1,
		movieId:    1,
		err:        movie.RatingUpdateError,
		addViewErr: movie.InvalidViewAdd,
	},
}

func setUpVoteMock(dbMock pgxmock.PgxPoolIface, test voteTestCase, vote int) {
	someSqlError := errors.New("some sql error idk")
	dbMock.ExpectBegin()
	if test.err == movie.InvalidVoteError {
		dbMock.ExpectExec(regexp.QuoteMeta(queryVote)).
			WithArgs(test.userId, test.movieId, vote).
			WillReturnError(someSqlError)
		dbMock.ExpectRollback()
		return
	} else {
		dbMock.ExpectExec(regexp.QuoteMeta(queryVote)).
			WithArgs(test.userId, test.movieId, vote).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))
	}
	dbMock.ExpectCommit()

	dbMock.ExpectBegin()
	if test.addViewErr == movie.InvalidViewCheck {
		dbMock.ExpectQuery(regexp.QuoteMeta(queryCheckView)).
			WithArgs(test.userId, test.movieId).
			WillReturnError(someSqlError)
		dbMock.ExpectRollback()
		return
	} else {
		rows := pgxmock.NewRows([]string{"movie_id"})
		dbMock.ExpectQuery(regexp.QuoteMeta(queryCheckView)).
			WithArgs(test.userId, test.movieId).WillReturnRows(rows)
	}
	dbMock.ExpectCommit()

	dbMock.ExpectBegin()
	if test.addViewErr == movie.InvalidViewAdd {
		dbMock.ExpectExec(regexp.QuoteMeta(queryAddView)).
			WithArgs(test.userId, test.movieId).
			WillReturnError(someSqlError)
		dbMock.ExpectRollback()
		return
	} else {
		dbMock.ExpectExec(regexp.QuoteMeta(queryAddView)).
			WithArgs(test.userId, test.movieId).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))
	}
	dbMock.ExpectCommit()

	dbMock.ExpectBegin()
	rows := pgxmock.NewRows([]string{"count"}).AddRow(cast.Uint64ToBytes(uint64(500)))
	dbMock.ExpectQuery(regexp.QuoteMeta(queryCountLikes)).
		WithArgs(test.movieId).WillReturnRows(rows)
	dbMock.ExpectCommit()

	dbMock.ExpectBegin()
	rows = pgxmock.NewRows([]string{"count"}).AddRow(cast.Uint64ToBytes(uint64(100)))
	dbMock.ExpectQuery(regexp.QuoteMeta(queryCountDislikes)).
		WithArgs(test.movieId).WillReturnRows(rows)
	dbMock.ExpectCommit()

	dbMock.ExpectBegin()
	rows = pgxmock.NewRows([]string{"count"}).AddRow(cast.Uint64ToBytes(uint64(1000)))
	dbMock.ExpectQuery(regexp.QuoteMeta(queryCountViews)).
		WithArgs(test.movieId).WillReturnRows(rows)
	dbMock.ExpectCommit()

	dbMock.ExpectBegin()
	dbMock.ExpectExec(regexp.QuoteMeta(querySetRating)).
		WithArgs(countRating(500, 100, 1000), test.movieId).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	dbMock.ExpectCommit()
}

func TestDbMovieRepository_Like(t *testing.T) {
	mockDb, dbMock := NewMock()
	movieRepo := NewMovieRepository(mockDb)
	defer dbMock.Close()

	for testId, test := range likeTests {
		t.Run(fmt.Sprintln(testId, test.name), func(t *testing.T) {
			setUpVoteMock(dbMock, test, domain.Like)
			err := movieRepo.Like(test.userId, test.movieId)
			require.Equal(t, test.err, err)
		})
	}
}

func TestDbMovieRepository_Dislike(t *testing.T) {
	mockDb, dbMock := NewMock()
	movieRepo := NewMovieRepository(mockDb)
	defer dbMock.Close()

	for testId, test := range likeTests {
		t.Run(fmt.Sprintln(testId, test.name), func(t *testing.T) {
			setUpVoteMock(dbMock, test, domain.Dislike)
			err := movieRepo.Dislike(test.userId, test.movieId)
			require.Equal(t, test.err, err)
		})
	}
}

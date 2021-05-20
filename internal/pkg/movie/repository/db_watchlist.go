package repository

import (
	"Redioteka/internal/pkg/movie"
	"Redioteka/internal/pkg/utils/log"
	"fmt"
)

const (
	queryInsertWatchlist = `insert into user_watchlist values(default, $1, $2);`
	queryDeleteWatchlist = `delete from user_watchlist where user_id = $1 and movie_id = $2;`
	querySelectWatchlist = `select id from user_watchlist here user_id = $1 and movie_id = $2;`
)

func (mr *dbMovieRepository) AddWatchlistByID(movieID, userID uint) error {
	err := mr.db.Exec(queryInsertWatchlist, userID, movieID)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Cannot add to watchlist of movie id: %d for user id: %d", movieID, userID))
		return movie.NotFoundError
	}
	return nil
}

func (mr *dbMovieRepository) RemoveWatchlistByID(movieID, userID uint) error {
	err := mr.db.Exec(queryDeleteWatchlist, userID, movieID)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Cannot delete from watchlist movie with id: %d for user id: %d", movieID, userID))
		return movie.NotFoundError
	}

	return nil
}

func (mr *dbMovieRepository) CheckWatchlistByID(movieID, userID uint) error {
	data, err := mr.db.Query(querySelectWatchlist, userID, movieID)
	if err == nil && len(data) == 0 {
		return nil
	}
	log.Log.Warn(fmt.Sprintf("Check of watchlist failed with movie id: %d user_id: %d", movieID, userID))
	return movie.AlreadyExists
}

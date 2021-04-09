package repository

import (
	"Redioteka/internal/pkg/database"
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/movie"
	"Redioteka/internal/pkg/utils/cast"
	"Redioteka/internal/pkg/utils/log"
	"fmt"
	"strconv"
	"strings"
)

const (
	querySelectID = `select
	m.id, m.title, m.description, m.avatar, m.rating, m.countries, m.directors, m.release_year, m.is_free, mt.type,
	(select string_agg(a.firstname || ' ' || a.lastname, ';') from actors as a join movie_actors as ma on a.id = ma.actor_id join movies as m on m.id = ma.movie_id where m.id = $1) as acts,
	(select string_agg(g.name, ';') from genres as g join movie_genres as mg on g.id = mg.genre_id join movies as m on m.id = mg.movie_id where m.id = $1) as gns
	from movies as m 
	join movie_types as mt on m.type = mt.id
	where m.id = $1;`
	queryInsertFav = `insert into user_favs values(default, $1, $2);`
	queryDeleteFav = `delete from user_favs where user_id = $1 and movie_id = $2;`
	querySelectFav = `select id from user_favs where user_id = $1 and movie_id = $2;`
)

type dbMovieRepository struct {
	db *database.DBManager
}

func NewMovieRepository(db *database.DBManager) domain.MovieRepository {
	return &dbMovieRepository{db: db}
}

func (mr *dbMovieRepository) GetById(id uint) (domain.Movie, error) {
	data, err := mr.db.Query(querySelectID, id)
	if err != nil {
		log.Log.Warn(fmt.Sprint("Cannot get movie from db with id: ", id))
		return domain.Movie{}, err
	}
	if len(data) == 0 {
		log.Log.Warn(fmt.Sprintf("Movie with id: %d - not found in db", id))
		return domain.Movie{}, movie.NotFoundError
	}

	first := data[0]
	actors := cast.ToString(first[10])
	genres := cast.ToString(first[11])
	movie := domain.Movie{
		ID:          cast.ToUint(first[0]),
		Title:       cast.ToString(first[1]),
		Description: cast.ToString(first[2]),
		Avatar:      cast.ToString(first[3]),
		Rating:      cast.ToFloat(first[4]),
		Countries:   strings.Split(cast.ToString(first[5]), ", "),
		Director:    strings.Split(cast.ToString(first[6]), ", "),
		Year:        strconv.Itoa(cast.ToSmallInt(first[7])),
		IsFree:      cast.ToBool(first[8]),
		Type:        domain.MovieType(cast.ToString(first[9])),
		Actors:      strings.Split(actors, ";"),
		Genres:      strings.Split(genres, ";"),
	}
	return movie, nil
}

func (mr *dbMovieRepository) AddFavouriteByID(movieID, userID uint) error {
	err := mr.db.Exec(queryInsertFav, userID, movieID)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Cannot add fav of movie id: %d for user id: %d", movieID, userID))
		return movie.NotFoundError
	}

	return nil
}

func (mr *dbMovieRepository) RemoveFavouriteByID(movieID, userID uint) error {
	err := mr.db.Exec(queryDeleteFav, userID, movieID)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Cannot delete fav of movie id: %d for user id: %d", movieID, userID))
		return movie.NotFoundError
	}

	return nil
}

func (mr *dbMovieRepository) CheckFavouriteByID(movieID, userID uint) error {
	data, err := mr.db.Query(querySelectFav, userID, movieID)
	if err == nil && len(data) == 0 {
		return nil
	}
	log.Log.Warn(fmt.Sprintf("Check of fav failed with movie id: %d user_id: %d", movieID, userID))
	return movie.AlreadyExists
}

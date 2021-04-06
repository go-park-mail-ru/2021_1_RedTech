package repository

import (
	"Redioteka/internal/pkg/database"
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/utils/cast"
	"Redioteka/internal/pkg/utils/log"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
querySelectID = `select m.id,
    m.title,
    m.description,
    m.avatar,
    m.rating,
    m.countries,
    m.directors,
    m.release_year,
    m.price,
    mt.type,
    (
        select string_agg(a.firstname || ' ' || a.lastname, ';')
        from actors as a
            join movie_actors as ma on a.id = ma.actor_id
            join movies as m on m.id = ma.movie_id
        where m.id = $1
    ) as acts,
    (
        select string_agg(g.name, ';')
        from genres as g
            join movie_genres as mg on g.id = mg.genre_id
            join movies as m on m.id = mg.movie_id
        where m.id = $1
    ) as gns
from movies as m
    join movie_types as mt on m.type = mt.id
where m.id = $1;`
querySelectFilter = `select m.id,
    m.title,
    m.description,
    m.avatar,
    m.rating,
    m.countries,
    m.directors,
    m.release_year,
    m.price,
    mt.type,
    (
        select string_agg(a.firstname || ' ' || a.lastname, ';')
        from actors as a
            join movie_actors as ma on a.id = ma.actor_id
            join movies as m on m.id = ma.movie_id
        where m.id = $1
    ) as acts,
    (
        select string_agg(g.name, ';')
        from genres as g
            join movie_genres as mg on g.id = mg.genre_id
            join movies as m on m.id = mg.movie_id
        where m.id = $1
    ) as gns
from movies as m
    join movie_types as mt on m.type = mt.id
where m.id = $1;`

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
		return domain.Movie{}, errors.New("Movie does not exist")
	}

	first := data[0]
	price := cast.ToFloat(first[8])
	actors := cast.ToString(first[10])
	genres := cast.ToString(first[11])
	movie := domain.Movie{
		ID:          cast.ToUint(first[0]),
		Title:       cast.ToString(first[1]),
		Description: cast.ToString(first[2]),
		Avatar:      cast.ToString(first[3]),
		Rating:      cast.ToFloat(first[4]), // ?
		Countries:   strings.Split(cast.ToString(first[5]), ", "),
		Director:    strings.Split(cast.ToString(first[6]), ", "),
		Year:        strconv.Itoa(cast.ToSmallInt(first[7])),
		IsFree:      price == 0,
		Type:        domain.MovieType(cast.ToString(first[9])),
		Actors:      strings.Split(actors, ";"),
		Genres:      strings.Split(genres, ";"),
	}
	return movie, nil
}

func (mr *dbMovieRepository) GetByFilter(filter domain.MovieFilter) ([]domain.Movie, error) {
	panic("implement me")
}

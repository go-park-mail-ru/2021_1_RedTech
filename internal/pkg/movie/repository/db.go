package repository

import (
	"Redioteka/internal/pkg/database"
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/movie"
	"Redioteka/internal/pkg/utils/cast"
	"Redioteka/internal/pkg/utils/log"
	"fmt"
	sq "github.com/Masterminds/squirrel"
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
    m.is_free,
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
		Rating:      cast.ToFloat(first[4]), // ?
		Countries:   strings.Split(cast.ToString(first[5]), ", "),
		Director:    strings.Split(cast.ToString(first[6]), ", "),
		Year:        strconv.Itoa(cast.ToSmallInt(first[7])),
		IsFree:      first[8][0] != 0,
		Type:        domain.MovieType(cast.ToString(first[9])),
		Actors:      strings.Split(actors, ";"),
		Genres:      strings.Split(genres, ";"),
	}
	return movie, nil
}

func buildFilterQuery(filter domain.MovieFilter) (string, []interface{}, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	allMovies := psql.Select("distinct movies.id as movie_id, movies.title, movies.description, movies.avatar," +
		"is_free").From("movies").Join("movie_actors ma on movies.id = ma.movie_id").
		Join("movie_genres mg on movies.id = mg.movie_id").
		Join("movie_types mt on movies.type = mt.id").
		Join("genres g on mg.genre_id = g.id").
		Join("(select a.id, a.firstname || ' ' || a.lastname as full_actor_name from actors as a) full_acts on full_acts.id = ma.actor_id")
	if filter.MinRating > 0 {
		allMovies = allMovies.Where(sq.GtOrEq{"rating": filter.MinRating})
	}
	if filter.Actors != nil {
		allMovies = allMovies.Where(sq.Eq{"a.full_actor_name": filter.Actors})
	}
	if filter.Genres != nil {
		allMovies = allMovies.Where(sq.Eq{"lower(g.name)": filter.Genres})
	}
	if filter.IsFree != domain.FilterBoth {
		allMovies = allMovies.Where(sq.Eq{"is_free": filter.IsFree == domain.FilterFree})
	}
	if filter.Type != "" {
		allMovies = allMovies.Where(sq.Eq{"mt.type": filter.Type})
	}
	allMovies = allMovies.Offset(uint64(filter.Offset)).Limit(uint64(filter.Limit))
	return allMovies.ToSql()
}

func IsFilterValid(filter domain.MovieFilter) bool {
	return filter.Offset >= 0 && filter.Limit >= 0 &&
		(filter.IsFree == domain.FilterBoth ||
			filter.IsFree == domain.FilterFree ||
			filter.IsFree == domain.FilterSubscription) &&
		(filter.Type == "" ||
			filter.Type == domain.MovieT ||
			filter.Type == domain.SeriesT)
}

func (mr *dbMovieRepository) GetByFilter(filter domain.MovieFilter) ([]domain.Movie, error) {
	if !IsFilterValid(filter) {
		log.Log.Warn("Invalid filter")
		return nil, movie.InvalidFilterError
	}
	filterQuery, filterArgs, err := buildFilterQuery(filter)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Can't build filter request: %v", err))
		return nil, err
	}
	data, err := mr.db.Query(filterQuery, filterArgs...)
	if err != nil {
		log.Log.Warn(fmt.Sprint("Cannot get movies from db with filter: ", filter))
		return nil, err
	}
	var res []domain.Movie
	for _, row := range data {

		res = append(res, domain.Movie{
			ID:          cast.ToUint(row[0]),
			Title:       string(row[1]),
			Description: string(row[2]),
			Avatar:      string(row[3]),
			IsFree:      row[4][0] != 0,
		})
	}
	return res, nil
}

func (mr *dbMovieRepository) GetGenres() ([]string, error) {
	data, err := mr.db.Query(`select name from genres;`)
	if err != nil {
		log.Log.Warn(fmt.Sprint("Cannot get genres from db"))
		return nil, err
	}
	res := make([]string, len(data))
	for i, row := range data {
		res[i] = string(row[0])
	}
	return res, nil
}

func (mr *dbMovieRepository) GetStream(id uint) (domain.Stream, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select("path").From("movie_videos").Where(sq.Eq{"movie_id": id})
	sqlQuery, args, err := query.ToSql()
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Can't build stream request: %v", err))
		return domain.Stream{}, err
	}
	data, err := mr.db.Query(sqlQuery, args...)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Cannot get movie video path: %v", err))
		return domain.Stream{}, err
	} else if len(data) == 0 {
		log.Log.Warn(fmt.Sprintf("Cannot find movie with id %v", id))
		return domain.Stream{}, movie.NotFoundError
	}
	res := domain.Stream{
		Video: string(data[0][0]),
	}
	return res, nil
}

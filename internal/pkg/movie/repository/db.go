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
	queryInsertFav = `insert into user_favs values(default, $1, $2);`
	queryDeleteFav = `delete from user_favs where user_id = $1 and movie_id = $2;`
	querySelectFav = `select id from user_favs where user_id = $1 and movie_id = $2;`
	querySelectID  = `select m.id,
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

	queryVote = `INSERT INTO movie_votes
	(user_id, movie_id, value)
	VALUES ($1, $2, $3)
	on conflict (movie_id, user_id) do update
	set value=$3
	where movie_id =$2
	and user_id = $1;`

	querySetRating = `update movies set rating=$1 where id=$2;`

	queryAddView   = `insert into movie_views(user_id, movie_id) values($1, $2);`
	queryCheckView = `select movie_id from movie_views where user_id = $1 and movie_id = $2;`

	queryCountLikes    = `select count(*) from movie_votes where movie_id = $1 and rating > 0;`
	queryCountDislikes = `select count(*) from movie_votes where movie_id = $1 and rating < 0;`
	queryCountViews    = `select count(*) from movie_views where movie_id = $1;`
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
		Actors:      strings.Split(cast.ToString(first[10]), ";"),
		Genres:      strings.Split(cast.ToString(first[11]), ";"),
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

func countRating(likes, dislikes, views int) float32 {
	likeWeight := 10
	disLikeweight := -5
	viewWeight := 7
	if views == 0 {
		return 0
	}
	return 10 * float32((views-dislikes-likes)*viewWeight+
		likes*likeWeight+
		dislikes*disLikeweight) /
		float32(views*likeWeight)
}

func (mr *dbMovieRepository) updateRating(movieId uint) error {
	likes := 0
	data, err := mr.db.Query(queryCountLikes, movieId)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Can't get like count %v", err))
		return err
	}
	if len(data) > 0 {
		likes = cast.ToInt(data[0][0])
	}

	dislikes := 0
	data, err = mr.db.Query(queryCountDislikes, movieId)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Can't get dislike countpath: %v", err))
		return err
	}
	if len(data) > 0 {
		dislikes = cast.ToInt(data[0][0])
	}

	views := 0
	data, err = mr.db.Query(queryCountViews, movieId)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Can't get view count path: %v", err))
		return err
	}
	if len(data) > 0 {
		views = cast.ToInt(data[0][0])
	}

	newRating := countRating(likes, dislikes, views)
	err = mr.db.Exec(querySetRating, newRating, movieId)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Can't update rating: %v", err))
		return err
	}

	return nil
}

func (mr *dbMovieRepository) addView(userId, movieId uint) error {
	data, err := mr.db.Query(queryCheckView, userId, movieId)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("User %v can't check %v movie view: %v", userId, movieId, err))
		return movie.InvalidViewCheck
	}
	if len(data) == 0 {
		err := mr.db.Exec(queryAddView, userId, movieId)
		if err != nil {
			log.Log.Warn(fmt.Sprintf("User %v can't set %v movie view: %v", userId, movieId, err))
			return movie.InvalidViewAdd
		}
	}
	return nil
}

func (mr *dbMovieRepository) Like(userId, movieId uint) error {
	err := mr.db.Exec(queryVote, userId, movieId, domain.Like)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("User %v can't like movie %v: %v", userId, movieId, err))
		return movie.InvalidVoteError
	}
	err = mr.addView(userId, movieId)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("User %v can't add movie %v view: %v", userId, movieId, err))
		return movie.RatingUpdateError
	}
	err = mr.updateRating(movieId)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("User %v can't update movie %v rating: %v", userId, movieId, err))
		return movie.RatingUpdateError
	}
	return nil
}

func (mr *dbMovieRepository) Dislike(userId, movieId uint) error {
	err := mr.db.Exec(queryVote, userId, movieId, domain.Dislike)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("User %v can't dislike movie %v: %v", userId, movieId, err))
		return movie.InvalidVoteError
	}
	err = mr.addView(userId, movieId)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("User %v can't add movie %v view: %v", userId, movieId, err))
		return movie.RatingUpdateError
	}
	err = mr.updateRating(movieId)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("User %v can't update movie %v rating: %v", userId, movieId, err))
		return movie.RatingUpdateError
	}
	return nil
}

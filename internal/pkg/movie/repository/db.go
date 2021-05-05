package repository

import (
	"Redioteka/internal/pkg/database"
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/movie"
	"Redioteka/internal/pkg/utils/baseutils"
	"Redioteka/internal/pkg/utils/cast"
	"Redioteka/internal/pkg/utils/log"
	"fmt"
	"strconv"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

const (
	queryInsertFav  = `insert into user_favs values(default, $1, $2);`
	queryDeleteFav  = `delete from user_favs where user_id = $1 and movie_id = $2;`
	querySelectFav  = `select id from user_favs where user_id = $1 and movie_id = $2;`
	querySelectVote = `select mv.value from movie_votes as mv join movies as m on mv.movie_id = m.id 
	join users as u on mv.user_id = u.id where u.id = $1 and m.id = $2;`
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
           select string_agg(g.label_rus, ';')
           from genres as g
                    join movie_genres as mg on g.id = mg.genre_id
                    join movies as m on m.id = mg.movie_id
           where m.id = $1
       ) as gns,
       (
           select string_agg(cast(a.id as varchar), ';')
           from actors as a
                    join movie_actors as ma on a.id = ma.actor_id
                    join movies as m on m.id = ma.movie_id
           where m.id = $1
       ) as actors_ids
from movies as m
         join movie_types as mt on m.type = mt.id
where m.id = $1;`

	querySelectVideo  = `select path, season, series from movie_videos where movie_id = $1 order by season, series;`
	querySelectSeries = `select count(series) from movie_videos where movie_id = $1 group by season order by season;`

	queryVote = `insert into movie_votes (user_id, movie_id, value)
	values ($1, $2, $3)
	on conflict (user_id, movie_id) do update set value=$3;`

	querySetRating = `update movies set rating=$1 where id=$2;`

	queryAddView   = `insert into movie_views(user_id, movie_id) values($1, $2);`
	queryCheckView = `select movie_id from movie_views where user_id = $1 and movie_id = $2;`

	queryCountLikes    = `select count(*) from movie_votes where movie_id = $1 and value > 0;`
	queryCountDislikes = `select count(*) from movie_votes where movie_id = $1 and value < 0;`
	queryCountViews    = `select count(*) from movie_views where movie_id = $1;`
	querySearch        = `select id, title, description, avatar, is_free from movies where lower(title) similar to $1;`
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
	actorNames := strings.Split(cast.ToString(first[10]), ";")
	actorIds, err := baseutils.StringsToUint(strings.Split(cast.ToString(first[12]), ";"))
	if err != nil {
		return domain.Movie{}, err
	}
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
		Actors:      actorNames,
		ActorIds:    actorIds,
		Genres:      strings.Split(cast.ToString(first[11]), ";"),
	}
	return movie, nil
}

func (mr *dbMovieRepository) GetSeriesList(id uint) ([]uint, error) {
	data, err := mr.db.Query(querySelectSeries, id)
	if err != nil {
		log.Log.Warn(fmt.Sprint("Cannot get series with movie id: ", id))
		return nil, err
	}
	if len(data) == 0 {
		log.Log.Warn(fmt.Sprintf("Series for movie with id: %d - not found in db", id))
		return nil, movie.NotFoundError
	}

	seriesCount := make([]uint, 0)
	for _, row := range data {
		seriesCount = append(seriesCount, cast.ToUint64(row[0]))
	}
	return seriesCount, nil
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

func (mr *dbMovieRepository) CheckVoteByID(movieID, userID uint) int {
	data, err := mr.db.Query(querySelectVote, userID, movieID)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Check of vote failed with movie id: %d user_id: %d", movieID, userID))
		return 0
	} else if len(data) == 0 {
		return 0
	}
	return cast.ToSmallInt(data[0][0])
}

func buildFilterQuery(filter domain.MovieFilter) (string, []interface{}, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	allMovies := psql.
		Select("movies.id, movies.title, movies.description, movies.avatar, is_free").
		From("movies").
		Join("movie_actors ma on movies.id = ma.movie_id").
		Join("movie_genres mg on movies.id = mg.movie_id").
		Join("movie_types mt on movies.type = mt.id").
		Join("genres g on mg.genre_id = g.id").
		Join("(select a.id, a.firstname || ' ' || a.lastname as full_actor_name from actors as a) full_acts on full_acts.id = ma.actor_id").
		GroupBy("movies.id")
	if filter.MinRating > 0 {
		allMovies = allMovies.Where(sq.GtOrEq{"rating": filter.MinRating})
	}
	if filter.Genres != nil {
		allMovies = allMovies.Where(sq.Eq{"lower(g.name)": filter.Genres})
	}
	if filter.IsFree != domain.FilterBoth {
		allMovies = allMovies.Where(sq.Eq{"is_free": filter.IsFree == domain.FilterFree})
	}
	if filter.Type != "" {
		typeName := ""
		if filter.Type == domain.SeriesT {
			typeName = "Сериал"
		}
		if filter.Type == domain.MovieT {
			typeName = "Фильм"
		}
		allMovies = allMovies.Where(sq.Eq{"mt.type": typeName})
	}
	if filter.Order != domain.NoneOrder {
		switch filter.Order {
		case domain.DateOrder:
			allMovies = allMovies.OrderBy("movies.add_date desc")
		case domain.RatingOrder:
			allMovies = allMovies.OrderBy("movies.rating desc")
		}
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
		return nil, movie.NotFoundError
	}
	var res []domain.Movie
	for _, row := range data {
		res = append(res, domain.Movie{
			ID:          cast.ToUint(row[0]),
			Title:       cast.ToString(row[1]),
			Description: cast.ToString(row[2]),
			Avatar:      cast.ToString(row[3]),
			IsFree:      row[4][0] != 0,
		})
	}
	return res, nil
}

func (mr *dbMovieRepository) GetGenres() ([]domain.Genre, error) {
	data, err := mr.db.Query(`select name, label_rus, image from genres;`)
	if err != nil {
		log.Log.Warn("Cannot get genres from db")
		return nil, err
	}
	res := make([]domain.Genre, len(data))
	for i, row := range data {
		res[i] = domain.Genre{
			Name:     cast.ToString(row[0]),
			LabelRus: cast.ToString(row[1]),
			Image:    cast.ToString(row[2]),
		}
	}
	return res, nil
}

func (mr *dbMovieRepository) GetStream(id uint) ([]domain.Stream, error) {
	data, err := mr.db.Query(querySelectVideo, id)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Cannot get movie video path: %v", err))
		return nil, err
	} else if len(data) == 0 {
		log.Log.Warn(fmt.Sprintf("Cannot find movie with id %v", id))
		return nil, movie.NotFoundError
	}

	res := make([]domain.Stream, 0)
	for _, dataRow := range data {
		res = append(res, domain.Stream{
			Video:  cast.ToString(dataRow[0]),
			Season: cast.ToInt(dataRow[1]),
			Series: cast.ToInt(dataRow[2]),
		})
	}
	return res, nil
}

func countRating(likes, dislikes, views int) float32 {
	likeWeight := 10
	dislikeWeight := 0
	viewWeight := 7
	if views == 0 {
		return 0
	}
	rating := 10 * float32((views-dislikes-likes)*viewWeight+
		likes*likeWeight+
		dislikes*dislikeWeight) /
		float32(views*likeWeight)
	if rating < 1 {
		return 1
	}
	return rating
}

func (mr *dbMovieRepository) updateRating(movieId uint) error {
	likes := 0
	data, err := mr.db.Query(queryCountLikes, movieId)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Can't get like count %v", err))
		return err
	}
	if len(data) > 0 {
		likes = int(cast.ToUint64(data[0][0]))
	}

	dislikes := 0
	data, err = mr.db.Query(queryCountDislikes, movieId)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Can't get dislike countpath: %v", err))
		return err
	}
	if len(data) > 0 {
		dislikes = int(cast.ToUint64(data[0][0]))
	}

	views := 0
	data, err = mr.db.Query(queryCountViews, movieId)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Can't get view count path: %v", err))
		return err
	}
	if len(data) > 0 {
		views = int(cast.ToUint64(data[0][0]))
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

func (mr *dbMovieRepository) Search(query string) ([]domain.Movie, error) {
	data, err := mr.db.Query(querySearch, baseutils.PrepareQueryForSearch(query))
	if err != nil {
		log.Log.Warn(fmt.Sprint("Cannot find movies from db with search query: ", query))
		return nil, movie.NotFoundError
	}
	var res []domain.Movie
	for _, row := range data {
		res = append(res, domain.Movie{
			ID:          cast.ToUint(row[0]),
			Title:       cast.ToString(row[1]),
			Description: cast.ToString(row[2]),
			Avatar:      cast.ToString(row[3]),
			IsFree:      row[4][0] != 0,
		})
	}
	return res, nil
}

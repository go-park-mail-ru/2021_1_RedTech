package repository

import (
	"Redioteka/internal/pkg/database"
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/utils/cast"
	"Redioteka/internal/pkg/utils/log"
	"fmt"
)

const (
	querySelectActor         = `select id, firstname, lastname, born, avatar from actors where id = $1;`
	querySelectMoviesByActor = `select m.id, m.title, m.description, m.avatar
from actors a
         join movie_actors ma on a.id = ma.actor_id
         join movies m on ma.movie_id = m.id
where a.id = $1;`
)

type dbActorRepository struct {
	db *database.DBManager
}

func NewActorRepository(db *database.DBManager) domain.ActorRepository {
	return &dbActorRepository{db: db}
}

func (ar dbActorRepository) GetById(id uint) (domain.Actor, error) {
	actorData, err := ar.db.Query(querySelectActor, id)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Actor with id: #{id}  - not found in db"))
		return domain.Actor{}, err
	}
	movieData, err := ar.db.Query(querySelectMoviesByActor, id)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Movies with with actor with id: #{id}  - not found in db"))
		return domain.Actor{}, err
	}
	movies := make([]domain.Movie, len(movieData))
	for i, movie := range movieData {
		movies[i] = domain.Movie{
			ID:          cast.ToUint(movie[0]),
			Title:       cast.ToString(movie[1]),
			Description: cast.ToString(movie[2]),
			Avatar:      cast.ToString(movie[3]),
		}
	}

	actorLine := actorData[0]
	res := domain.Actor{
		ID:        cast.ToUint(actorLine[0]),
		FirstName: cast.ToString(actorLine[1]),
		LastName:  cast.ToString(actorLine[2]),
		Born:      cast.ToString(actorLine[3]),
		Avatar:    cast.ToString(actorLine[4]),
		Movies:    movies,
	}
	return res, nil
}

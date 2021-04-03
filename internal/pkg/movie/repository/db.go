package repository

import (
	"Redioteka/internal/pkg/database"
	"Redioteka/internal/pkg/domain"
	"errors"
	"fmt"
)

type dbMovieRepository struct {
	db *database.DBManager
}

func NewMovieRepository() domain.MovieRepository {
	return &dbMovieRepository{db: database.Manager}
}

func (mr *dbMovieRepository) GetById(id uint) (domain.Movie, error) {
	data, err := mr.db.Query(`select 
		m.id, m.title, m.description, m.avatar, m.rating, m.countries, m.directors, m.release_year, m.price 
		a.firstname || ' ' || a.lastname as actor, g.name, mt.type 
		from movies as m 
		join movie_genres as mg on mg.movie_id = m.id 
		join genres as g on mg.genre_id = g.id 
		join movie_actors as ma on ma.movie_id = m.id 
		join actors as a on ma.actor_id =  a.id
		join movie_types as mt on m.type = mt.id
		where id = $1`, id)
	if err != nil {
		return domain.Movie{}, err
	}
	first := data[0]
	if first == nil {
		return domain.Movie{}, errors.New("User does not exist")
	}
	fmt.Println(first)
	movie := first.(domain.Movie)
	return movie, nil
}

func (mr *dbMovieRepository) Delete(id uint) error {
	return nil
}

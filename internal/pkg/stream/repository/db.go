package repository

import (
	"Redioteka/internal/pkg/database"
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/movie"
	"Redioteka/internal/pkg/utils/cast"
	"Redioteka/internal/pkg/utils/log"
	"fmt"
)

const (
	querySelectVideo  = `select path, season, series from movie_videos where movie_id = $1 order by season, series;`
)

type dbStreamRepository struct {
	db *database.DBManager
}

func NewStreamRepository(db *database.DBManager) domain.StreamRepository {
	return &dbStreamRepository{
		db: db,
	}
}

func (sr dbStreamRepository) GetStream(id uint) ([]domain.Stream, error) {
	data, err := sr.db.Query(querySelectVideo, id)
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

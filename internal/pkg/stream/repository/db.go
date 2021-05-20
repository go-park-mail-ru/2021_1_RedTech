package repository

import (
	"Redioteka/internal/pkg/database"
	"Redioteka/internal/pkg/domain"
)

type dbStreamRepository struct {
	db *database.DBManager
}

func NewStreamRepository(db *database.DBManager) domain.StreamRepository {
	return &dbStreamRepository{
		db: db,
	}
}

func (d dbStreamRepository) GetStream(id uint) ([]domain.Stream, error) {

	panic("implement me")
}

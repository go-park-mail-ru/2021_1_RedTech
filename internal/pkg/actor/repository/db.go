package repository

import (
	"Redioteka/internal/pkg/database"
	"Redioteka/internal/pkg/domain"
)

const (
	querySelectActor = `select * from actors where id = $1;`
)

type dbActorRepository struct {
	db *database.DBManager
}

func NewActorRepository(db *database.DBManager) domain.ActorRepository {
	return &dbActorRepository{db: db}
}

func (ar dbActorRepository) GetById(id uint) (domain.Actor, error) {
	res := domain.Actor{}
	return res, nil
}


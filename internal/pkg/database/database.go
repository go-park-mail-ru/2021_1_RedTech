package database

import (
	"Redioteka/internal/pkg/utils/log"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DBManager struct {
	pool *pgxpool.Pool
}

func Connect() (*DBManager, error) {
	connString := "user=redtech password=red_tech host=localhost port=5432 dbname=netflix"
	pool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		log.Log.Error(err)
		return nil, err
	}
	log.Log.Info("Successful connection to postgres")
	return &DBManager{pool: pool}, nil
}

func Disconnect() {
	Manager.pool.Close()
}

var Manager, myErr = Connect()

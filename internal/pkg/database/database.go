package database

import (
	"Redioteka/internal/pkg/utils/log"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DBManager struct {
	pool *pgxpool.Pool
}

func (db *DBManager) Query(queryString string) ([]interface{}, error) {
	ctx := context.Background()
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		log.Log.Error(err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, queryString)
	if err != nil {
		log.Log.Error(err)
		return nil, err
	}
	defer rows.Close()

	err = tx.Commit(ctx)
	if err != nil {
		log.Log.Error(err)
		return nil, err
	}
	return rows.Values()
}

func (db *DBManager) Exec(queryString string) error {
	ctx := context.Background()
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		log.Log.Error(err)
		return err
	}
	defer tx.Rollback(ctx)

	result, err := tx.Exec(ctx, queryString)
	if err != nil {
		log.Log.Error(err)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Log.Error(err)
		return err
	}

	if result.RowsAffected() == 0 {
		log.Log.Warn("No row was changed")
	}
	return nil
}

func Connect() *DBManager {
	connString := "user=redtech password=red_tech host=localhost port=5432 dbname=netflix"
	pool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		log.Log.Error(err)
		return nil
	}
	log.Log.Info("Successful connection to postgres")
	return &DBManager{pool: pool}
}

func Disconnect() {
	Manager.pool.Close()
}

var Manager = Connect()

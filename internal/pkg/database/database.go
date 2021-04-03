package database

import (
	"Redioteka/internal/pkg/utils/log"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DBManager struct {
	pool *pgxpool.Pool
}

func (db *DBManager) Query(queryString string, params ...interface{}) ([]interface{}, error) {
	ctx := context.Background()
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		log.Log.Error(err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, queryString, params...)
	if err != nil {
		log.Log.Error(err)
		return nil, err
	}
	defer rows.Close()

	result := make([]interface{}, 0)
	for rows.Next() {
		row, err := rows.Values()
		if err != nil {
			return nil, err
		}
		result = append(result, row)
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Log.Error(err)
		return nil, err
	}
	return result, nil
}

func (db *DBManager) Exec(queryString string, params ...interface{}) error {
	ctx := context.Background()
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		log.Log.Error(err)
		return err
	}
	defer tx.Rollback(ctx)

	result, err := tx.Exec(ctx, queryString, params...)
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

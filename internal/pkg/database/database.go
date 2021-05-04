package database

import (
	"Redioteka/internal/pkg/utils/log"
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PgxPool interface {
	Begin(context.Context) (pgx.Tx, error)
	Close()
}

type DBManager struct {
	Pool PgxPool
}

func (db *DBManager) Query(queryString string, params ...interface{}) ([][][]byte, error) {
	ctx := context.Background()
	tx, err := db.Pool.Begin(ctx)
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

	result := make([][][]byte, 0)
	for rows.Next() {
		row := make([][]byte, 0)
		row = append(row, rows.RawValues()...)
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
	tx, err := db.Pool.Begin(ctx)
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
	connString := "postgres://redtech:red_tech@database:5432/netflix?sslmode=disable"
	pool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		log.Log.Error(err)
		return nil
	}
	log.Log.Info("Successful connection to postgres")
	return &DBManager{Pool: pool}
}

func Disconnect(manager *DBManager) {
	manager.Pool.Close()
	log.Log.Info("DB was disconnected")
}

package database

import (
	"Redioteka/internal/pkg/config"
	"Redioteka/internal/pkg/utils/log"
	"context"
	"fmt"
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

func Rollback(tx pgx.Tx, ctx context.Context) {
	err := tx.Rollback(ctx)
	if err != nil {
		log.Log.Error(err)
	}
}

func (db *DBManager) Query(queryString string, params ...interface{}) ([][][]byte, error) {
	ctx := context.Background()
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		log.Log.Error(err)
		return nil, err
	}
	defer Rollback(tx, ctx)

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
	defer Rollback(tx, ctx)

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
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", config.Get().Postgres.User, config.Get().Postgres.Password,
		config.Get().Postgres.Host, config.Get().Postgres.Port, config.Get().Postgres.DBName)
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

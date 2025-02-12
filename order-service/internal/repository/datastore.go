package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryxContext(context.Context, string, ...interface{}) (*sqlx.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	QueryRowxContext(context.Context, string, ...interface{}) *sqlx.Row
}

type DataStore interface {
	Atomic(ctx context.Context, fn func(DataStore) error) error
	ProductRepository() ProductRepository
}

type dataStore struct {
	conn *sqlx.DB
	db   DBTX
}

func NewDataStore(db *sqlx.DB) DataStore {
	return &dataStore{
		conn: db,
		db:   db,
	}
}

func (s *dataStore) Atomic(ctx context.Context, fn func(DataStore) error) error {
	tx, err := s.conn.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	err = fn(&dataStore{conn: s.conn, db: tx})
	if err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}

func (s *dataStore) ProductRepository() ProductRepository {
	return NewProductRepository(s.db)
}

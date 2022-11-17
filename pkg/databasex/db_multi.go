package databasex

import (
	"context"
	"database/sql"
	"fmt"
)

var (
	// check in runtime implement Databaser
	_ Adapter = (*dbMulti)(nil)
)

type dbMulti struct {
	dbRead  Adapter
	dbWrite Adapter
}

func NewMulti(dbWrite, dbRead Adapter) *dbMulti {
	return &dbMulti{
		dbWrite: dbWrite,
		dbRead:  dbRead,
	}
}

func (db *dbMulti) Ping() error {
	return db.dbWrite.Ping()
}

func (db *dbMulti) InTransaction() bool {
	return false
}

// Close closes the database connection.
func (db *dbMulti) Close() error {
	return fmt.Errorf("not implemented in multi db mode ")
}

// Exec executes a SQL statement and returns the number of rows it affected.
func (db *dbMulti) Exec(ctx context.Context, query string, args ...any) (_ int64, err error) {
	return db.dbWrite.Exec(ctx, query, args...)
}

// Query runs the DB query.
func (db *dbMulti) Query(ctx context.Context, dst any, query string, args ...any) error {
	return db.dbRead.Query(ctx, dst, query, args...)
}

// QueryRow runs the query and returns a single row.
func (db *dbMulti) QueryRow(ctx context.Context, dst any, query string, args ...any) error {
	return db.dbRead.QueryRow(ctx, dst, query, args...)
}

// QueryX runs the DB query.
func (db *dbMulti) QueryX(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return db.dbRead.QueryX(ctx, query, args...)
}

// QueryRowX runs the query and returns a single row.
func (db *dbMulti) QueryRowX(ctx context.Context, query string, args ...any) *sql.Row {
	return db.dbRead.QueryRowX(ctx, query, args...)
}

// Transact executes the given function in the context of a SQL transaction at
// the given isolation level
func (db *dbMulti) Transact(ctx context.Context, iso sql.IsolationLevel, txFunc func(*DB) error) (err error) {
	return db.dbWrite.Transact(ctx, iso, txFunc)
}

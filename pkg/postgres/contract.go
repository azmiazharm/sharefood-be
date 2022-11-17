package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

// const connStringTemplate = "postgres://%s:%s@%s:%d/%s?sslmode=disable"
const connStringTemplate = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s"

type Config struct {
	Host         string
	Port         int
	User         string
	Password     string
	Name         string
	Timeout      time.Duration
	Charset      string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
	Type         string
	TimeZone     string
}

type Adapter interface {
	QueryRow(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	QueryRows(ctx context.Context, query string, args ...interface{}) (rows *sqlx.Rows, err error)
	Fetch(ctx context.Context, dst interface{}, query string, args ...interface{}) error
	FetchRow(ctx context.Context, dst interface{}, query string, args ...interface{}) error
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Ping(ctx context.Context) error
	HealthCheck() error
}

type Transaction interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

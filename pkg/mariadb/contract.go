// Package mariadb
package mariadb

import (
	"context"
	"database/sql"
	"time"
)

// const connStringTemplate = "%s:%s@tcp(%s:%d)/%s?timeout=%v&charset=%s&parseTime=true&loc=Local"
const connStringTemplate = "%s:%s@tcp(%s:%d)/%s?%s"

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
	QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryRows(ctx context.Context, query string, args ...interface{}) (rows *sql.Rows, err error)
	Fetch(ctx context.Context, dst interface{}, query string, args ...interface{}) error
	FetchRow(ctx context.Context, dst interface{}, query string, args ...interface{}) error
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
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

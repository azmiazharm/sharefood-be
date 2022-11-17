// Package mariadb
package mariadb

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"sharefood/pkg/tracer"
	"sharefood/pkg/util"
)

type maria struct {
	db  *sqlx.DB
	cfg *Config
}

// NewMaria initialize single maria db
func NewMaria(cfg *Config) (Adapter, error) {
	x := maria{cfg: cfg}
	db, err := CreateSession(cfg)
	x.db = db

	return &x, err
}

// QueryRow select single row database will return  sql.row raw
func (d *maria) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	ctx = tracer.DBSpanStartWithOption(ctx, d.cfg.Name, "query_row",
		tracer.WithResourceNameOptions(query),
		tracer.WithOptions("sql.query", query),
		tracer.WithOptions("sql.args", util.DumpToString(args)),
	)
	defer tracer.SpanFinish(ctx)

	return d.db.QueryRowContext(ctx, query, args...)
}

// QueryRows select multiple rows of database will return  sql.rows raw
func (d *maria) QueryRows(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	ctx = tracer.DBSpanStartWithOption(ctx, d.cfg.Name, "query_rows",
		tracer.WithResourceNameOptions(query),
		tracer.WithOptions("sql.query", query),
		tracer.WithOptions("sql.args", util.DumpToString(args)),
	)
	defer tracer.SpanFinish(ctx)
	return d.db.QueryContext(ctx, query, args...)
}

// Fetch select multiple rows of database will cast data to struct passing by parameter
func (d *maria) Fetch(ctx context.Context, dst interface{}, query string, args ...interface{}) error {
	ctx = tracer.DBSpanStartWithOption(ctx, d.cfg.Name, "fetch_rows",
		tracer.WithResourceNameOptions(query),
		tracer.WithOptions("sql.query", query),
		tracer.WithOptions("sql.args", util.DumpToString(args)),
	)
	defer tracer.SpanFinish(ctx)
	return d.db.SelectContext(ctx, dst, query, args...)
}

// FetchRow fetching one row database will cast data to struct passing by parameter
func (d *maria) FetchRow(ctx context.Context, dst interface{}, query string, args ...interface{}) error {
	ctx = tracer.DBSpanStartWithOption(ctx, d.cfg.Name, "fetch_row",
		tracer.WithResourceNameOptions(query),
		tracer.WithOptions("sql.query", query),
		tracer.WithOptions("sql.args", util.DumpToString(args)),
	)
	defer tracer.SpanFinish(ctx)
	return d.db.GetContext(ctx, dst, query, args...)
}

// Exec execute mysql command query
func (d *maria) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	ctx = tracer.DBSpanStartWithOption(ctx, d.cfg.Name, "exec",
		tracer.WithResourceNameOptions(query),
		tracer.WithOptions("sql.query", query),
		tracer.WithOptions("sql.args", util.DumpToString(args)),
	)
	defer tracer.SpanFinish(ctx)
	return d.db.ExecContext(ctx, query, args...)
}

// BeginTx start new transaction session
func (d *maria) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	ctx = tracer.DBSpanStartWithOption(ctx, d.cfg.Name, "begin.transaction")
	defer tracer.SpanFinish(ctx)
	return d.db.BeginTx(ctx, opts)
}

// Ping check database connectivity
func (d *maria) Ping(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

// HealthCheck checking healthy of database connection
func (d *maria) HealthCheck() error {
	return d.Ping(context.Background())
}

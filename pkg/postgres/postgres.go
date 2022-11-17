package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"sharefood/pkg/tracer"
	"sharefood/pkg/util"
)

type postgres struct {
	db  *sqlx.DB
	cfg *Config
}

func NewPostgreSQL(cfg *Config) (Adapter, error) {
	x := postgres{cfg: cfg}
	db, err := CreateSession(cfg)
	x.db = db

	return &x, err
}

func (d *postgres) QueryRow(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	ctx = tracer.DBSpanStartWithOption(ctx, d.cfg.Name, "query_row",
		tracer.WithResourceNameOptions(query),
		tracer.WithOptions("sql.query", query),
		tracer.WithOptions("sql.args", util.DumpToString(args)),
	)
	defer tracer.SpanFinish(ctx)

	return d.db.QueryRowxContext(ctx, query, args...)
}

func (d *postgres) QueryRows(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	ctx = tracer.DBSpanStartWithOption(ctx, d.cfg.Name, "query_rows",
		tracer.WithResourceNameOptions(query),
		tracer.WithOptions("sql.query", query),
		tracer.WithOptions("sql.args", util.DumpToString(args)),
	)
	defer tracer.SpanFinish(ctx)
	return d.db.QueryxContext(ctx, query, args...)
}

func (d *postgres) Fetch(ctx context.Context, dst interface{}, query string, args ...interface{}) error {
	ctx = tracer.DBSpanStartWithOption(ctx, d.cfg.Name, "fetch_rows",
		tracer.WithResourceNameOptions(query),
		tracer.WithOptions("sql.query", query),
		tracer.WithOptions("sql.args", util.DumpToString(args)),
	)
	defer tracer.SpanFinish(ctx)
	return d.db.SelectContext(ctx, dst, query, args...)
}

func (d *postgres) FetchRow(ctx context.Context, dst interface{}, query string, args ...interface{}) error {
	ctx = tracer.DBSpanStartWithOption(ctx, d.cfg.Name, "fetch_row",
		tracer.WithResourceNameOptions(query),
		tracer.WithOptions("sql.query", query),
		tracer.WithOptions("sql.args", util.DumpToString(args)),
	)
	defer tracer.SpanFinish(ctx)
	return d.db.GetContext(ctx, dst, query, args...)
}

func (d *postgres) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	ctx = tracer.DBSpanStartWithOption(ctx, d.cfg.Name, "exec",
		tracer.WithResourceNameOptions(query),
		tracer.WithOptions("sql.query", query),
		tracer.WithOptions("sql.args", util.DumpToString(args)),
	)
	defer tracer.SpanFinish(ctx)
	return d.db.ExecContext(ctx, query, args...)
}

func (d *postgres) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) {
	ctx = tracer.DBSpanStartWithOption(ctx, d.cfg.Name, "begin.transaction")
	defer tracer.SpanFinish(ctx)
	return d.db.BeginTxx(ctx, opts)
}

func (d *postgres) Ping(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

func (d *postgres) HealthCheck() error {
	return d.Ping(context.Background())
}

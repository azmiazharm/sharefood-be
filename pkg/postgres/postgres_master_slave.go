package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"sharefood/pkg/tracer"
	"sharefood/pkg/util"
)

type postgreSQLMasterSlave struct {
	db      *sqlx.DB
	dbRead  *sqlx.DB
	cfg     *Config
	cfgRead *Config
}

func NewPostgresMasterSlave(cfgWrite *Config, cfgRead *Config) (Adapter, error) {
	x := postgreSQLMasterSlave{cfg: cfgWrite, cfgRead: cfgRead}

	e := x.initialize()

	return &x, e
}

func (d *postgreSQLMasterSlave) initialize() error {
	dbWrite, err := CreateSession(d.cfg)

	if err != nil {
		return err
	}

	err = dbWrite.Ping()
	if err != nil {
		return err
	}

	dbRead, err := CreateSession(d.cfgRead)
	if err != nil {
		return err
	}

	err = dbRead.Ping()
	if err != nil {
		return err
	}

	d.db = dbWrite
	d.dbRead = dbRead

	return nil
}

func (d *postgreSQLMasterSlave) QueryRow(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	ctx = tracer.DBSpanStartWithOption(ctx, d.cfg.Name, "query_row",
		tracer.WithResourceNameOptions(query),
		tracer.WithOptions("sql.query", query),
		tracer.WithOptions("sql.args", util.DumpToString(args)),
	)
	defer tracer.SpanFinish(ctx)
	return d.selector().QueryRowxContext(ctx, query, args...)
}

func (d *postgreSQLMasterSlave) QueryRows(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	ctx = tracer.DBSpanStartWithOption(ctx, d.cfg.Name, "query_rows",
		tracer.WithResourceNameOptions(query),
		tracer.WithOptions("sql.query", query),
		tracer.WithOptions("sql.args", util.DumpToString(args)),
	)
	defer tracer.SpanFinish(ctx)
	return d.selector().QueryxContext(ctx, query, args...)
}

func (d *postgreSQLMasterSlave) Fetch(ctx context.Context, dst interface{}, query string, args ...interface{}) error {
	ctx = tracer.DBSpanStartWithOption(ctx, d.cfg.Name, "fetch_rows",
		tracer.WithResourceNameOptions(query),
		tracer.WithOptions("sql.query", query),
		tracer.WithOptions("sql.args", util.DumpToString(args)),
	)
	defer tracer.SpanFinish(ctx)
	return d.selector().SelectContext(ctx, dst, query, args...)
}

func (d *postgreSQLMasterSlave) FetchRow(ctx context.Context, dst interface{}, query string, args ...interface{}) error {
	ctx = tracer.DBSpanStartWithOption(ctx, d.cfg.Name, "fetch_row",
		tracer.WithResourceNameOptions(query),
		tracer.WithOptions("sql.query", query),
		tracer.WithOptions("sql.args", util.DumpToString(args)),
	)
	defer tracer.SpanFinish(ctx)
	return d.selector().GetContext(ctx, dst, query, args...)
}

func (d *postgreSQLMasterSlave) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	ctx = tracer.DBSpanStartWithOption(ctx, d.cfg.Name, "exec",
		tracer.WithResourceNameOptions(query),
		tracer.WithOptions("sql.query", query),
		tracer.WithOptions("sql.args", util.DumpToString(args)),
	)
	defer tracer.SpanFinish(ctx)
	return d.db.ExecContext(ctx, query, args...)
}

func (d *postgreSQLMasterSlave) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) {
	ctx = tracer.DBSpanStartWithOption(ctx, d.cfg.Name, "begin.transaction")
	defer tracer.SpanFinish(ctx)
	return d.db.BeginTxx(ctx, opts)
}

func (d *postgreSQLMasterSlave) Ping(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

func (d *postgreSQLMasterSlave) HealthCheck() error {
	var err1, err2 error
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err1 = d.Ping(context.Background())
		wg.Done()
	}()

	if d.dbRead != nil {
		wg.Add(1)
		go func() {
			err2 = d.dbRead.PingContext(context.Background())
			wg.Done()
		}()
	}

	wg.Wait()

	if err1 != nil && err2 != nil {
		return fmt.Errorf("database write error:%s; database read error:%s; ", err1.Error(), err2.Error())
	}

	if err1 != nil {
		return fmt.Errorf("database write error:%s;", err1.Error())

	}

	if err2 != nil {
		return fmt.Errorf("database read error:%s;", err2.Error())

	}

	return nil

}

func (d *postgreSQLMasterSlave) selector() *sqlx.DB {
	if d.dbRead != nil {
		return d.dbRead
	}

	return d.db
}

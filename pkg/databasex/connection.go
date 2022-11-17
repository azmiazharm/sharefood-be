package databasex

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func mysqlDSN(cfg *Config) string {
	if len(strings.Trim(cfg.Charset, "")) == 0 {
		cfg.Charset = "UTF8"
	}

	param := url.Values{}
	param.Add("timeout", fmt.Sprintf("%v", time.Duration(cfg.Timeout)*time.Second))
	param.Add("charset", cfg.Charset)
	param.Add("parseTime", "True")
	param.Add("loc", cfg.TimeZone)

	connStr := fmt.Sprintf(connStringMysqlTemplate,
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		param.Encode(),
	)

	return connStr
}

func postgreDSN(cfg *Config) string {
	param := url.Values{}
	param.Add("connect_timeout", fmt.Sprint(cfg.Timeout))
	param.Add("timezone", cfg.TimeZone)
	param.Add("sslmode", "disable")

	connStr := fmt.Sprintf(connStringPostgresTemplate,
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		param.Encode(),
	)

	return connStr
}

// CreateSession create new session maria db
func CreateSession(cfg *Config) (*sqlx.DB, error) {

	fn, ok := dsn[cfg.Driver]
	if !ok {
		return nil, fmt.Errorf("invalid driver %s", cfg.Driver)
	}

	db, err := sqlx.Open(cfg.Driver, fn(cfg))
	if err != nil {
		return db, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return db, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.MaxLifetime)

	return db, nil
}

package postgres

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func CreateSession(cfg *Config) (*sqlx.DB, error) {
	if len(strings.Trim(cfg.Charset, "")) == 0 {
		cfg.Charset = "UTF8"
	}

	param := url.Values{}
	param.Add("timeout", fmt.Sprintf("%v", cfg.Timeout))
	param.Add("charset", cfg.Charset)
	param.Add("parseTime", "True")
	param.Add("loc", cfg.TimeZone)

	connStr := fmt.Sprintf(connStringTemplate,
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.Port,
		cfg.TimeZone,
	)

	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return db, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.MaxLifetime)

	return db, nil
}

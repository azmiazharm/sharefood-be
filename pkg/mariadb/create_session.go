// Package mariadb
package mariadb

import (
	"fmt"
	"net/url"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// CreateSession create new session maria db
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
        cfg.User,
        cfg.Password,
        cfg.Host,
        cfg.Port,
        cfg.Name,
        param.Encode(),
    )

	db, err := sqlx.Open("mysql", connStr)
	if err != nil {
		return db, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.MaxLifetime)

	return db, nil
}

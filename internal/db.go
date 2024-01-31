package internal

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/jmoiron/sqlx"
)

func InitDBConn(provider, cfg string) (*sqlx.DB, error) {
	db, err := sqlx.Connect(provider, cfg)

	if err != nil {
		return nil, err
	}
	return db, nil
}

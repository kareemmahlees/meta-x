package internal

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/kareemmahlees/meta-x/utils"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/jmoiron/sqlx"
)

func InitDBConn(provider string, cfg utils.Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect(provider, cfg.DSN())

	if err != nil {
		return nil, err
	}
	return db, nil
}

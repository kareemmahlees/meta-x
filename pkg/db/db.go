package db

import (
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/kareemmahlees/mysql-meta/utils"
)

func InitDBConn() (*sqlx.DB, error) {
	dbUsername, _ := utils.GetEnvVar("MYSQL_META_DB_USERNAME")
	dbPassword, _ := utils.GetEnvVar("MYSQL_META_DB_PASSWORD")
	dbName, _ := utils.GetEnvVar("MYSQL_META_DB_NAME")

	cfg := mysql.Config{
		User:                 dbUsername,
		Passwd:               dbPassword,
		DBName:               dbName,
		AllowNativePasswords: true,
	}

	db, err := sqlx.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}
	return db, nil
}

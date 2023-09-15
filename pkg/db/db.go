package db

import (
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/kareemmahlees/mysql-meta/utils"
)

func InitDBConn() (*sqlx.DB, error) {
	dbUsername, _ := utils.GetEnvVar("MYSQL_META_DB_USERNAME", false)
	dbPassword, _ := utils.GetEnvVar("MYSQL_META_DB_PASSWORD", false)
	dbName, _ := utils.GetEnvVar("MYSQL_META_DB_NAME", false)
	dbPort, _ := utils.GetEnvVar("MYSQL_META_DB_PORT", true)

	cfg := mysql.Config{
		User:                 dbUsername,
		Passwd:               dbPassword,
		DBName:               dbName,
		Net:                  "tcp",
		
		AllowNativePasswords: true,
	}

	// if dbPort is set it means we are in testing
	if dbPort != "" {
		cfg.Addr = fmt.Sprintf("127.0.0.1:%s", dbPort)
	}

	db, err := sqlx.Open("mysql", cfg.FormatDSN())

	db.SetMaxIdleConns(100)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Minute * 2)

	if err != nil {
		return nil, err
	}
	return db, nil
}

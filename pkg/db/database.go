package db

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/jmoiron/sqlx"
)

func ListDatabases(db *sqlx.DB) []string {

	var dbs []string
	rows, _ := db.Queryx("SHOW DATABASES")
	for rows.Next() {
		var db string
		err := rows.Scan(&db)
		if err != nil {
			log.Error(err)
		}
		dbs = append(dbs, db)
	}
	return dbs
}

func CreateDatabase(db *sqlx.DB, dbName string) (int, error) {
	res, err := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
	if err != nil {
		return 0, err
	}
	num, _ := res.RowsAffected()
	return int(num), nil
}

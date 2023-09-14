package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func ListDatabases(db *sqlx.DB) ([]string, error) {

	var dbs []string
	rows, err := db.Queryx("SHOW DATABASES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var db string
		err := rows.Scan(&db)
		if err != nil {
			return nil, err
		}
		dbs = append(dbs, db)
	}
	return dbs, nil
}

func CreateDatabase(db *sqlx.DB, dbName string) (int, error) {
	res, err := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
	if err != nil {
		return 0, err
	}
	num, _ := res.RowsAffected()
	return int(num), nil
}

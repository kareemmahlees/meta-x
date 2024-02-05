package db

import (
	"fmt"
	"meta-x/lib"

	"github.com/jmoiron/sqlx"
)

func ListDatabasesPgMySQL(db *sqlx.DB, provider string) ([]*string, error) {

	var dbs []*string
	var queryString string

	switch provider {
	case lib.PSQL:
		queryString = "SELECT DATNAME FROM PG_DATABASE"
	case lib.MYSQL:
		queryString = "SHOW DATABASES"
	}

	rows, err := db.Queryx(queryString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		db := new(string)
		err := rows.Scan(db)
		if err != nil {
			return nil, err
		}
		dbs = append(dbs, db)
	}
	return dbs, nil
}

type SqliteDatabase struct {
	File string `json:"file" db:"file"`
	Name string `json:"name" db:"name"`
}

func ListDatabasesSqlite(db *sqlx.DB) ([]*string, error) {
	dbs := []*string{}

	rows, err := db.Queryx("SELECT name,file FROM PRAGMA_DATABASE_LIST")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		db := new(SqliteDatabase)
		err = rows.StructScan(db)
		if err != nil {
			return nil, err
		}
		record := fmt.Sprintf("%s %s", db.Name, db.File)
		dbs = append(dbs, &record)
	}
	return dbs, nil
}

func CreatePgMysqlDatabase(db *sqlx.DB, provider, dbName string) error {
	var queryString string

	switch provider {
	case lib.PSQL:
		queryString = fmt.Sprintf("CREATE DATABASE %s", dbName)
	case lib.MYSQL:
		queryString = fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName)
	}

	_, err := db.Exec(queryString)
	if err != nil {
		return err
	}
	return nil
}

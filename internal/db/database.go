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

func ListDatabasesSqlite(db *sqlx.DB) ([]*SqliteDatabase, error) {
	dbs := []*SqliteDatabase{}

	rows, err := db.Queryx("SELECT name,file FROM PRAGMA_DATABASE_LIST;")
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

		dbs = append(dbs, db)
	}
	return dbs, nil
}

func AttachSqliteDatabase(db *sqlx.DB, dbName, filePath string) (int64, error) {
	res, err := db.Exec(fmt.Sprintf("ATTACH DATABASE '%s' AS %s;", filePath, dbName))
	if err != nil {
		return 0, err
	}
	fmt.Printf("%v", res)
	num, _ := res.RowsAffected()
	return num, nil
}

func CreatePgMysqlDatabase(db *sqlx.DB, dbName string) (int64, error) {
	res, err := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
	if err != nil {
		return 0, err
	}
	num, _ := res.RowsAffected()
	return num, nil
}

package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/kareemmahlees/mysql-meta/lib"
)

func ListTables(db *sqlx.DB, dbName string) (result []string, err error) {
	_, err = db.Queryx(fmt.Sprintf("use %s", dbName))
	if err != nil {
		return nil, err
	}
	rows, err := db.Queryx("show tables")
	if err != nil {
		return nil, err
	}

	var tables []string
	for rows.Next() {
		var table string
		err := rows.Scan(&table)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}
	return tables, nil
}

var dataTypesMappings = map[string]string{
	"text":   "varchar(255)",
	"number": "int",
}

func CreateTable(db *sqlx.DB, dbName string, tableName string, payload map[string]lib.TablePropsValidator) (result int64, err error) {
	_, err = db.Queryx(fmt.Sprintf("USE %s", dbName))
	if err != nil {
		return 0, err
	}
	// this long solution is made because placeholders "?" can't
	// be used for db, table or column names
	dataString := ""
	for col, props := range payload {
		if _, ok := dataTypesMappings[props.Type]; ok {
			props.Type = dataTypesMappings[props.Type]
		}
		dataString += fmt.Sprintf("%s\t%s\t", col, props.Type)
		if props.Nullable != nil && props.Nullable == false {
			dataString += "NOT NULL\t"
		}
		if props.Unique != nil && props.Unique == true {
			dataString += "UNIQUE\t"
		}
		if props.Default != nil {
			dataString += fmt.Sprintf("DEFAULT \t'%s'", props.Default)
		}
		dataString += ",\n"
	}
	res, err := db.Exec(fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		ID int NOT NULL,
		%s
		PRIMARY KEY (ID)
	)
	`, tableName, dataString))
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

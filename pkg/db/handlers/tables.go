package handlers

import (
	"fmt"

	"github.com/jmoiron/sqlx"
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

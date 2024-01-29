package db

import (
	"fmt"
	"strings"

	"meta-x/lib"

	"github.com/jmoiron/sqlx"
)

type tableInfoStruct struct {
	Field   string `db:"Field" json:"field"`
	Type    string `db:"Type" json:"type"`
	Null    string `db:"Null" json:"null"`
	Key     string `db:"Key" json:"key"`
	Default any    `db:"Default" json:"default"`
	Extra   any    `db:"Extra" json:"extra"`
}

func GetTableInfo(db *sqlx.DB, tableName string) (result []tableInfoStruct, err error) {
	rows, err := db.Queryx(fmt.Sprintf("desc %s", tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tablesDescriptions = []tableInfoStruct{}
	for rows.Next() {
		var tableDesc tableInfoStruct
		err := rows.StructScan(&tableDesc)
		if err != nil {
			return nil, err
		}
		tablesDescriptions = append(tablesDescriptions, tableDesc)
	}
	return tablesDescriptions, nil
}

func ListTables(db *sqlx.DB) (result []string, err error) {
	rows, err := db.Queryx("show tables")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tables = []string{}
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

func CreateTable(db *sqlx.DB, tableName string, payload map[string]lib.CreateTableProps) error {
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
		id int NOT NULL,
		%s
		PRIMARY KEY (id)
	)
	`, tableName, dataString))
	if err != nil {
		return err
	}
	_, err = res.LastInsertId()
	return err

}

func UpdateTable(db *sqlx.DB, tableName string, payload lib.UpdateTableProps) error {
	dataString := ""
	switch payload.Operation.Type {
	case "add":
		for col, dataType := range payload.Operation.Data.(map[string]interface{}) {
			dataString += fmt.Sprintf("ADD %s %s,\n", col, dataType)
		}
	case "modify":
		for col, dataType := range payload.Operation.Data.(map[string]interface{}) {
			dataString += fmt.Sprintf("MODIFY COLUMN %s %s,\n", col, dataType)
		}
	case "delete":
		for _, col := range payload.Operation.Data.([]interface{}) {
			dataString += fmt.Sprintf("DROP COLUMN %s,\n", col)
		}
	}
	dataString, _ = strings.CutSuffix(dataString, ",\n")
	_, err := db.Exec(fmt.Sprintf(`
	ALTER TABLE %s 
		%s
	`, tableName, dataString))
	if err != nil {
		return err
	}
	return nil
}

func DeleteTable(db *sqlx.DB, tableName string) error {
	_, err := db.Exec(fmt.Sprintf(`DROP TABLE %s`, tableName))
	if err != nil {
		return err
	}
	return nil
}

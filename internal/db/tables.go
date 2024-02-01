package db

import (
	"fmt"
	"strings"

	"meta-x/lib"

	"github.com/jmoiron/sqlx"
)

type tableDescriptionStruct struct {
	Name     string `db:"name" json:"name"`
	Type     string `db:"type" json:"type"`
	Nullable string `db:"nullable" json:"nullable"`
	Key      any    `db:"key" json:"key"`
	Default  any    `db:"default" json:"default"`
}

func GetTableInfo(db *sqlx.DB, tableName, provider string) (result []*tableDescriptionStruct, err error) {
	var queryString string
	switch provider {
	case lib.SQLITE3:
		queryString = `
		SELECT name,type,
			CASE when 'notnull' = 1
			THEN 'NO'
			ELSE 'YES'
			END AS nullable,

			CASE WHEN pk = 1
			THEN 'PRI'
			END AS key,
			dflt_value AS 'default'
		FROM pragma_table_info('%s');`
	case lib.PSQL:
		queryString = `
		SELECT col.column_name AS name,
			col.data_type AS type,
			col.is_nullable AS nullable,
			kcu.constraint_name AS key,
			col.column_default AS default
		FROM information_schema.columns AS col
		LEFT JOIN information_schema.key_column_usage AS kcu ON col.column_name = kcu.column_name
		WHERE col.table_name = '%s';`
	case lib.MYSQL:
		queryString = `
		SELECT column_name AS name,
			column_type AS type,
			is_nullable AS nullable,
			column_key AS key,
			column_default AS default
		FROM information_schema.columns
		WHERE table_name='%s';
		`
	}
	rows, err := db.Queryx(fmt.Sprintf(queryString, tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tablesDescriptions := []*tableDescriptionStruct{}
	for rows.Next() {
		tableDesc := new(tableDescriptionStruct)
		err := rows.StructScan(tableDesc)
		if err != nil {
			return nil, err
		}
		tablesDescriptions = append(tablesDescriptions, tableDesc)
	}
	return tablesDescriptions, nil
}

func ListTables(db *sqlx.DB, provider string) (result []*string, err error) {
	var queryString string
	switch provider {
	case lib.SQLITE3:
		queryString = "SELECT tbl_name FROM sqlite_master where type='table'"
	case lib.PSQL:
		queryString = "SELECT tablename FROM pg_catalog.pg_tables where schemaname='public';"
	case lib.MYSQL:
		queryString = "SHOW TABLES"
	}

	rows, err := db.Queryx(queryString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables = []*string{}
	for rows.Next() {
		table := new(string)
		err := rows.Scan(table)
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

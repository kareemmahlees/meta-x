package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/kareemmahlees/meta-x/models"
)

type pgProvider struct {
	db *sqlx.DB
}

func NewPgProvider(db *sqlx.DB) *pgProvider {
	return &pgProvider{db}
}

func (p *pgProvider) ListDBs() ([]*string, error) {
	var dbs []*string

	queryString := "SELECT DATNAME FROM PG_DATABASE"

	rows, err := p.db.Queryx(queryString)
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

func (p *pgProvider) CreateDB(dbName string) error {
	queryString := fmt.Sprintf("CREATE DATABASE %s", dbName)

	_, err := p.db.Exec(queryString)
	if err != nil {
		return err
	}
	return nil
}

func (p *pgProvider) GetTable(tableName string) ([]*models.TableColumnInfo, error) {
	queryString := `
		SELECT col.column_name AS name,
			col.data_type AS type,
			col.is_nullable AS nullable,
			kcu.constraint_name AS key,
			col.column_default AS default
		FROM information_schema.columns AS col
		LEFT JOIN information_schema.key_column_usage AS kcu ON col.column_name = kcu.column_name
		WHERE col.table_name = '%s';`
	rows, err := p.db.Queryx(fmt.Sprintf(queryString, tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tablesDescriptions := []*models.TableColumnInfo{}
	for rows.Next() {
		tableDesc := new(models.TableColumnInfo)
		err := rows.StructScan(tableDesc)
		if err != nil {
			return nil, err
		}
		tablesDescriptions = append(tablesDescriptions, tableDesc)
	}
	return tablesDescriptions, nil
}
func (p *pgProvider) ListTables() ([]*string, error) {
	queryString := "SELECT tablename FROM pg_catalog.pg_tables where schemaname='public';"

	rows, err := p.db.Queryx(queryString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables = []*string{}
	for rows.Next() {
		table := new(string)
		err := rows.Scan(&table)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}
	return tables, nil
}
func (p *pgProvider) CreateTable(tableName string, data []models.CreateTablePayload) error {
	// this long solution is made because placeholders "?" can't
	// be used for db, table or column names
	dataString := ""
	for _, props := range data {
		dataString += fmt.Sprintf("%s\t%s\t", props.ColName, props.Type)
		if props.Nullable != nil && props.Nullable == false {
			dataString += "NOT NULL\t"
		}
		if props.Unique != nil && props.Unique == true {
			dataString += "UNIQUE\t"
		}
		if props.Default != nil {
			dataString += fmt.Sprintf("DEFAULT\t'%s'", props.Default)
		}
		dataString += ",\n"
	}
	_, err := p.db.Exec(fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		id int NOT NULL,
		%s
		PRIMARY KEY (id)
	)
	`, tableName, dataString))
	if err != nil {
		return err
	}
	return nil

}
func (p *pgProvider) DeleteTable(tableName string) error {
	_, err := p.db.Exec(fmt.Sprintf(`DROP TABLE %s`, tableName))
	if err != nil {
		return err
	}
	return nil
}
func (p *pgProvider) AddColumn(tableName string, data models.AddModifyColumnPayload) error {
	dataString := fmt.Sprintf("ADD %s %s\n", data.ColName, data.Type)

	return alterTable(p.db, tableName, dataString)
}
func (p *pgProvider) UpdateColumn(tableName string, data models.AddModifyColumnPayload) error {
	dataString := fmt.Sprintf("ALTER COLUMN %s TYPE %s\n", data.ColName, data.Type)

	return alterTable(p.db, tableName, dataString)
}
func (p *pgProvider) DeleteColumn(tableName string, data models.DeleteColumnPayload) error {
	dataString := fmt.Sprintf("DROP COLUMN %s\n", data.ColName)

	return alterTable(p.db, tableName, dataString)
}

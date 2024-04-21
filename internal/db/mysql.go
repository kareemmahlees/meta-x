package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/kareemmahlees/meta-x/models"
)

type mysqlProvider struct {
	db *sqlx.DB
}

func NewMySQLProvider(db *sqlx.DB) *mysqlProvider {
	return &mysqlProvider{db}
}

func (p *mysqlProvider) ListDBs() ([]*string, error) {
	var dbs []*string

	rows, err := p.db.Queryx("SHOW DATABASES")
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

func (p *mysqlProvider) CreateDB(dbName string) error {
	queryString := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName)

	_, err := p.db.Exec(queryString)
	if err != nil {
		return err
	}
	return nil
}
func (p *mysqlProvider) GetTable(tableName string) ([]*models.TableInfoResp, error) {
	queryString := `
		SELECT column_name AS name,
			column_type AS type,
			is_nullable AS nullable,
			column_key AS 'key',
			column_default AS 'default'
		FROM information_schema.columns
		WHERE table_name='%s';
		`
	rows, err := p.db.Queryx(fmt.Sprintf(queryString, tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tablesDescriptions := []*models.TableInfoResp{}
	for rows.Next() {
		tableDesc := new(models.TableInfoResp)
		err := rows.StructScan(tableDesc)
		if err != nil {
			return nil, err
		}
		tablesDescriptions = append(tablesDescriptions, tableDesc)
	}
	return tablesDescriptions, nil
}
func (p *mysqlProvider) ListTables() ([]*string, error) {
	queryString := "SHOW TABLES"

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
func (p *mysqlProvider) CreateTable(tableName string, data []models.CreateTablePayload) error {
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
func (p *mysqlProvider) DeleteTable(tableName string) error {
	_, err := p.db.Exec(fmt.Sprintf(`DROP TABLE %s`, tableName))
	if err != nil {
		return err
	}
	return nil
}
func (p *mysqlProvider) AddColumn(tableName string, data models.AddModifyColumnPayload) error {
	dataString := fmt.Sprintf("ADD %s %s\n", data.ColName, data.Type)

	return alterTable(p.db, tableName, dataString)
}
func (p *mysqlProvider) UpdateColumn(tableName string, data models.AddModifyColumnPayload) error {
	dataString := fmt.Sprintf("MODIFY COLUMN %s %s\n", data.ColName, data.Type)

	return alterTable(p.db, tableName, dataString)
}
func (p *mysqlProvider) DeleteColumn(tableName string, data models.DeleteColumnPayload) error {
	dataString := fmt.Sprintf("DROP COLUMN %s\n", data.ColName)

	return alterTable(p.db, tableName, dataString)
}

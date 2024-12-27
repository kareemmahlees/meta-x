package db

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/kareemmahlees/meta-x/models"
)

type SqliteProvider struct {
	db *sqlx.DB
}

func NewSQLiteProvider(db *sqlx.DB) *SqliteProvider {
	return &SqliteProvider{db}
}

type sqliteDatabase struct {
	File string `json:"file" db:"file"`
	Name string `json:"name" db:"name"`
}

func (p *SqliteProvider) ListDBs() ([]*string, error) {
	dbs := []*string{}

	rows, err := p.db.Queryx("SELECT name,file FROM PRAGMA_DATABASE_LIST")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		db := new(sqliteDatabase)
		err = rows.StructScan(db)
		if err != nil {
			return nil, err
		}
		record := fmt.Sprintf("%s %s", db.Name, db.File)
		dbs = append(dbs, &record)
	}
	return dbs, nil
}

func (p *SqliteProvider) CreateDB(dbName string) error {
	return errors.New("Unsupported")
}

func (p *SqliteProvider) GetTable(tableName string) ([]*models.TableColumnInfo, error) {
	queryString := `
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
func (p *SqliteProvider) ListTables() ([]*string, error) {
	queryString := "SELECT tbl_name FROM sqlite_master where type='table'"

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
func (p *SqliteProvider) CreateTable(tableName string, data []models.CreateTablePayload) error {
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
func (p *SqliteProvider) DeleteTable(tableName string) error {
	_, err := p.db.Exec(fmt.Sprintf(`DROP TABLE %s`, tableName))
	if err != nil {
		return err
	}
	return nil
}
func (p *SqliteProvider) AddColumn(tableName string, data models.AddModifyColumnPayload) error {
	dataString := ""
	dataString += fmt.Sprintf("ADD %s %s\n", data.ColName, data.Type)

	return alterTable(p.db, tableName, dataString)
}
func (p *SqliteProvider) UpdateColumn(tableName string, data models.AddModifyColumnPayload) error {
	return errors.New("Unsupported")
}
func (p *SqliteProvider) DeleteColumn(tableName string, data models.DeleteColumnPayload) error {
	dataString := fmt.Sprintf("DROP COLUMN %s\n", data.ColName)

	return alterTable(p.db, tableName, dataString)
}

func alterTable(db *sqlx.DB, tableName string, dataString string) error {
	_, err := db.Exec(fmt.Sprintf(`
	ALTER TABLE %s 
		%s
	`, tableName, dataString))
	if err != nil {
		return err
	}
	return nil
}

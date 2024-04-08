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

func (p *SqliteProvider) GetTable(tableName string) ([]*models.TableInfoResp, error) {
	return nil, nil
}
func (p *SqliteProvider) ListTables() ([]*string, error) {
	return nil, nil
}
func (p *SqliteProvider) CreateTable(tableName string, data []models.CreateTablePayload) error {
	return nil
}
func (p *SqliteProvider) DeleteTable(tableName string) error {
	return nil
}
func (p *SqliteProvider) AddColumn(tableName string, data models.AddModifyColumnPayload) error {
	return nil
}
func (p *SqliteProvider) UpdateColumn(tableName string, data models.AddModifyColumnPayload) error {
	return nil
}

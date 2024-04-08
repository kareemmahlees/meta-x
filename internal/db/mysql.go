package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
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

package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
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

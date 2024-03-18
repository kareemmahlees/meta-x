package db

import "github.com/kareemmahlees/meta-x/models"

// Interface that must be implemented by any
// storage driver dealing with database logic
type DatabaseExecuter interface {
	ListDBs() ([]*string, error)
	CreateDB(dbName string) error
}

// Interface that must be implemented by any
// storage driver dealing with table logic
type TableExecuter interface {
	GetTable(tableName string) ([]*models.TableInfoResp, error)
	ListTables() ([]*string, error)
	CreateTable(tableName string, data []models.CreateTablePayload) error
	DeleteTable(tableName string) error
	AddColumn(tableName string, data models.AddModifyColumnPayload) error
	UpdateColumn(tableName string, data models.AddModifyColumnPayload) error
}

type Storage interface {
	DatabaseExecuter
	TableExecuter
}

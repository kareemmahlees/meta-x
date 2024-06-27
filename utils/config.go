package utils

import (
	"fmt"

	"github.com/go-sql-driver/mysql"
)

type Config interface {
	DSN() string
}

type SQLiteConfig struct {
	filePath string
}

func NewSQLiteConfig(filePath string) *SQLiteConfig {
	return &SQLiteConfig{filePath}
}

func (sc SQLiteConfig) DSN() string {
	return sc.filePath
}

type PgConfig struct {
	ConnUrl    *string
	ConnParams *PgConnectionParams
}

type PgConnectionParams struct {
	DBUsername string
	DBHost     string
	DBName     string
	DBSslMode  string
	DBPassword string
	DBPort     int
}

func NewPGConfig(connUrl *string, pgConnParams *PgConnectionParams) *PgConfig {
	return &PgConfig{connUrl, pgConnParams}
}

func (pc PgConfig) DSN() string {
	if pc.ConnUrl != nil {
		return *pc.ConnUrl
	} else {
		cfg := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
			pc.ConnParams.DBUsername,
			pc.ConnParams.DBPassword,
			pc.ConnParams.DBHost,
			pc.ConnParams.DBPort,
			pc.ConnParams.DBName,
			pc.ConnParams.DBSslMode,
		)
		return cfg
	}
}

type MySQLConfig struct {
	ConnUrl    *string
	ConnParams *MySQLConnectionParams
}

type MySQLConnectionParams struct {
	DBUsername string
	DBHost     string
	DBName     string
	DBPassword string
	DBPort     int
}

func NewMySQLConfig(connUrl *string, pgConnParams *MySQLConnectionParams) *MySQLConfig {
	return &MySQLConfig{connUrl, pgConnParams}
}

func (mc *MySQLConfig) DSN() string {
	if mc.ConnUrl != nil {
		return *mc.ConnUrl
	} else {
		cfg := mysql.Config{
			User:   mc.ConnParams.DBUsername,
			Passwd: mc.ConnParams.DBPassword,
			DBName: mc.ConnParams.DBName,
			Net:    "tcp",
			Addr:   fmt.Sprintf("%s:%d", mc.ConnParams.DBHost, mc.ConnParams.DBPort),
		}
		return cfg.FormatDSN()
	}
}

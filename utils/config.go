package utils

import (
	"fmt"

	"github.com/lib/pq"
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
	DbUsername string
	DbHost     string
	DbPort     int
	DbName     string
	DbSslMode  string
	DbPassword string
}

func NewPGConfig(connUrl *string, pgConnParams *PgConnectionParams) *PgConfig {
	return &PgConfig{connUrl, pgConnParams}
}

func (pc PgConfig) DSN() string {
	if *pc.ConnUrl != "" {
		return *pc.ConnUrl
	} else {
		cfg, _ := pq.ParseURL(fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
			pc.ConnParams.DbUsername,
			pc.ConnParams.DbPassword,
			pc.ConnParams.DbHost,
			pc.ConnParams.DbPort,
			pc.ConnParams.DbName,
			pc.ConnParams.DbSslMode,
		))
		return cfg
	}
}

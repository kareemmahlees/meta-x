package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSQLiteConfig(t *testing.T) {
	filePath := "any"
	conf := NewSQLiteConfig(filePath)

	assert.Equal(t, conf.DSN(), filePath)
}

func TestPostgresConfig(t *testing.T) {
	connUrl := "any"

	conf := NewPGConfig(&connUrl, nil)
	assert.Equal(t, conf.DSN(), connUrl)

	conf = NewPGConfig(nil, &PgConnectionParams{
		DBUsername: "test",
		DBPassword: "test",
		DBHost:     "localhost",
		DBPort:     5432,
		DBName:     "test",
	})

	assert.Regexp(t, `postgres:\/\/\w+:\w+@.+:\d+\/\w+`, conf.DSN())
}

func TestMySQLConfig(t *testing.T) {
	connUrl := "any"

	conf := NewMySQLConfig(&connUrl, nil)
	assert.Equal(t, conf.DSN(), connUrl)

	conf = NewMySQLConfig(nil, &MySQLConnectionParams{
		DBUsername: "test",
		DBPassword: "test",
		DBHost:     "localhost",
		DBPort:     5432,
		DBName:     "test",
	})

	assert.Regexp(t, `\w+:\w+@tcp\(.+:\d+\)/\w+`, conf.DSN())
}

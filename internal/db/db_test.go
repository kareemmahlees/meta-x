package db

import (
	"testing"

	"github.com/kareemmahlees/meta-x/lib"
	"github.com/kareemmahlees/meta-x/utils"
	"github.com/stretchr/testify/assert"
)

func TestInitDBConn(t *testing.T) {
	var cfg utils.Config

	cfg = utils.NewSQLiteConfig(":memory:")

	conn, err := InitDBConn(lib.SQLITE3, cfg)
	assert.Nil(t, err)
	defer conn.Close()

	malformedConnUrl := "postgres://malformed"
	cfg = utils.NewPGConfig(&malformedConnUrl, nil)

	conn, err = InitDBConn(lib.PSQL, cfg)
	assert.Error(t, err)
}

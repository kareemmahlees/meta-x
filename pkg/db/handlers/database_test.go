package handlers

import (
	"os"
	"testing"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/kareemmahlees/mysql-meta/pkg/db"
	"github.com/stretchr/testify/assert"
)

var con *sqlx.DB

func TestMain(m *testing.M){
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatal(err)
	}
	con,err =db.InitDBConn()
	if err != nil {
		log.Error(err)
	}
	defer con.Close()

	// Run tests
    exitVal := m.Run()
    
    // Exit with exit value from tests
    os.Exit(exitVal)
}

func TestListDatabases(t *testing.T) {
	dbs := ListDatabases(con)
	assert.Greater(t,len(dbs),0)
}

func TestCreateDatabase(t *testing.T){
	rowsAffected,err := CreateDatabase(con,"mysqlmeta")
	assert.Nil(t,err)
	assert.Equal(t,rowsAffected,1)
}
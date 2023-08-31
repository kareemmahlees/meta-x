package db

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	os.Exit(m.Run())
}

func TestInitDBConn(t *testing.T) {
	db, err := InitDBConn()
	assert.Nil(t, err)

	defer db.Close()

	err = db.Ping()
	assert.Nil(t, err)
}

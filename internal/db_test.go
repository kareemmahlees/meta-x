package internal

import (
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func init() {
	err := godotenv.Load("../../.env.test")
	if err != nil {
		log.Fatal(err)
	}
}

func TestInitDBConn(t *testing.T) {
	db, err := InitDBConn()
	assert.Nil(t, err)

	defer db.Close()

	err = db.Ping()
	assert.Nil(t, err)
}

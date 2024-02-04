package utils_test

import (
	"context"
	"io"
	"meta-x/lib"
	"meta-x/utils"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func TestCreatePostgresContainer(t *testing.T) {
	ctx := context.Background()
	pgContainer, err := utils.CreatePostgresContainer(ctx)
	defer pgContainer.Terminate(ctx)

	assert.Nil(t, err)

	con, err := sqlx.Open(lib.PSQL, pgContainer.ConnectionString)
	assert.Nil(t, err)

	defer con.Close()

	err = con.Ping()
	assert.Nil(t, err)
}

func TestCreateMySQLContainer(t *testing.T) {
	ctx := context.Background()
	mysqlContainer, err := utils.CreateMySQLContainer(ctx)
	defer mysqlContainer.Terminate(ctx)

	assert.Nil(t, err)

	con, err := sqlx.Open(lib.MYSQL, mysqlContainer.ConnectionString)
	assert.Nil(t, err)

	defer con.Close()

	err = con.Ping()
	assert.Nil(t, err)
}

type mockBody struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestDecodeBody(t *testing.T) {
	testBody := io.NopCloser(strings.NewReader(`
	{
		"name":"foo",
		"age": 123
	}`))

	decodedBody := utils.DecodeBody[mockBody](testBody)

	name := decodedBody.Name
	assert.Equal(t, name, "foo")

	age := decodedBody.Age
	assert.Equal(t, age, 123)

	testBody = io.NopCloser(strings.NewReader(`
	[ 
		{
			"name":"foo",
			"age" : 123
		},
		{
			"name":"bar",
			"age" : 123
		}
	 ]`))

	decodedBody2 := utils.DecodeBody[[]mockBody](testBody)

	name = decodedBody2[1].Name
	assert.Equal(t, name, "bar")

	age = decodedBody2[1].Age
	assert.Equal(t, age, 123)
}

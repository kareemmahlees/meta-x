package utils

import (
	"context"
	"io"
	"log"
	"meta-x/lib"
	"reflect"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func TestCreatePostgresContainer(t *testing.T) {
	ctx := context.Background()
	pgContainer, err := CreatePostgresContainer(ctx)
	defer func() {
		_ = pgContainer.Terminate(ctx)
	}()

	assert.Nil(t, err)

	con, err := sqlx.Open(lib.PSQL, pgContainer.ConnectionString)
	assert.Nil(t, err)

	defer con.Close()

	err = con.Ping()
	assert.Nil(t, err)
}

func TestCreateMySQLContainer(t *testing.T) {
	ctx := context.Background()
	mysqlContainer, err := CreateMySQLContainer(ctx)
	defer func() {
		_ = mysqlContainer.Terminate(ctx)
	}()

	assert.Nil(t, err)

	con, err := sqlx.Open(lib.MYSQL, mysqlContainer.ConnectionString)
	assert.Nil(t, err)

	defer con.Close()
}

func TestNewTestingFiberApp(t *testing.T) {
	app := NewTestingFiberApp(lib.SQLITE3)
	defer func() {
		err := app.Shutdown()
		if err != nil {
			log.Fatal(err)
		}
	}()

	listenCh := make(chan bool)

	app.Hooks().OnListen(func(ld fiber.ListenData) error {
		listenCh <- true
		return nil
	})

	go func() {
		if err := app.Listen(":5522"); err != nil {
			listenCh <- false
			log.Fatal(err)
		}
	}()

	startedListening := <-listenCh

	assert.True(t, startedListening)
}

func TestEncodeBody(t *testing.T) {
	mockBody := "test"
	encodedBody := EncodeBody(mockBody)

	assert.Equal(t, 7, encodedBody.Len())
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

	decodedBody := DecodeBody[mockBody](testBody)

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

	decodedBody2 := DecodeBody[[]mockBody](testBody)

	name = decodedBody2[1].Name
	assert.Equal(t, name, "bar")

	age = decodedBody2[1].Age
	assert.Equal(t, age, 123)
}

func TestSliceOfPointersToSliceOfValues(t *testing.T) {
	var testSlice []*string
	testSlice = append(testSlice, new(string))

	soptsov := SliceOfPointersToSliceOfValues(testSlice)

	assert.IsType(t, reflect.SliceOf(reflect.TypeOf("")), reflect.TypeOf(soptsov))
}

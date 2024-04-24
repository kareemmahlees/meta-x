package utils

import (
	"context"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kareemmahlees/meta-x/lib"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func TestCreatePostgresContainer(t *testing.T) {
	t.Run("should pass", func(t *testing.T) {
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
	})

	t.Run("should fail timeout exceeded", func(t *testing.T) {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Millisecond)
		defer cancel()

		_, err := CreatePostgresContainer(ctx)

		assert.Error(t, err)

	})
}

func TestCreateMySQLContainer(t *testing.T) {
	t.Run("should pass", func(t *testing.T) {
		ctx := context.Background()
		mysqlContainer, err := CreateMySQLContainer(ctx)
		defer func() {
			_ = mysqlContainer.Terminate(ctx)
		}()

		assert.Nil(t, err)

		con, err := sqlx.Open(lib.MYSQL, mysqlContainer.ConnectionString)
		assert.Nil(t, err)

		defer con.Close()

	})

	t.Run("should fail timetout exceeded", func(t *testing.T) {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Millisecond)
		defer cancel()

		_, err := CreateMySQLContainer(ctx)

		assert.Error(t, err)

	})
}

func TestEncodeBody(t *testing.T) {
	mockBody := "test"
	encodedBody, err := EncodeBody(mockBody)

	assert.Equal(t, 7, encodedBody.Len())
	assert.Nil(t, err)

	_, err = EncodeBody(make(chan any)) // make encoding fail
	assert.NotNil(t, err)
}

type mockBody struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// A mocked Reader that implements io.ReadClose which always returns an error
type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func (errReader) Close() error {
	return errors.New("test error")
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

	decodedBody3 := DecodeBody[string](errReader(0))

	assert.Empty(t, decodedBody3)
}

func TestSliceOfPointersToSliceOfValues(t *testing.T) {
	var testSlice []*string
	testSlice = append(testSlice, new(string))

	soptsov := SliceOfPointersToSliceOfValues(testSlice)

	assert.IsType(t, reflect.SliceOf(reflect.TypeOf("")), reflect.TypeOf(soptsov))
}

func TestRequestTest(t *testing.T) {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})

	rr := TestRequest(r, http.MethodGet, "/", http.NoBody)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, rr.Body.String(), "Hello")
}

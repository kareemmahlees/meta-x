package utils

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestRouteHandler(t *testing.T) {

	mockRouteHandlerFunc := func(*fiber.Ctx, *sqlx.DB) error {
		return nil
	}

	con, _ := sqlx.Open("sqlite3", ":memory:")
	defer con.Close()

	err := RouteHandler(con, mockRouteHandlerFunc)(&fiber.Ctx{})

	assert.Nil(t, err)

}

func TestGraphQLHandler(t *testing.T) {
	app := fiber.New()
	c := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer func() {
		_ = app.Shutdown()
	}()

	mockGraphQLHandlerFunc := func(http.ResponseWriter, *http.Request) {}

	handler := GraphQLHandler(mockGraphQLHandlerFunc)

	assert.NotPanics(t, func() { handler(c) })
}

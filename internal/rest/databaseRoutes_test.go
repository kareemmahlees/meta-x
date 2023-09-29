package routes

import (
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/kareemmahlees/mysql-meta/internal/db"
	"github.com/kareemmahlees/mysql-meta/utils"
	"github.com/stretchr/testify/assert"
)

func init() {
	err := godotenv.Load("../../.env.test")
	if err != nil {
		log.Fatal(err)
	}
}

func TestRegisterDatabaseRoutes(t *testing.T) {
	app := fiber.New()

	RegisterDatabasesRoutes(app, nil)

	var routes []struct {
		method string
		params []string
		path   string
	}

	for _, route := range app.GetRoutes() {
		routes = append(routes, struct {
			method string
			params []string
			path   string
		}{
			method: route.Method,
			params: route.Params,
			path:   route.Path,
		})
	}

	assert.Contains(t, routes, struct {
		method string
		params []string
		path   string
	}{
		method: "GET",
		params: []string(nil),
		path:   "/databases",
	})

}

func TestHandleListDatabases(t *testing.T) {
	app := fiber.New()

	con, err := db.InitDBConn()
	if err != nil {
		log.Fatal(err)
	}
	defer con.Close()

	RegisterDatabasesRoutes(app, con)
	req := httptest.NewRequest("GET", "http://localhost:4000/databases", nil)

	resp, _ := app.Test(req)
	payload := utils.ReadBody[map[string]any](resp.Body)

	assert.Equal(t, resp.StatusCode, fiber.StatusOK)

	dbs, ok := payload["databases"]
	assert.True(t, ok)
	assert.Greater(t, len(dbs.([]interface{})), 0)
}

func TestHandleCreateDatabase(t *testing.T) {
	app := fiber.New()

	con, err := db.InitDBConn()
	if err != nil {
		log.Fatal(err)
	}
	defer con.Close()

	RegisterDatabasesRoutes(app, con)
	req := httptest.NewRequest("POST", "http://localhost:4000/databases/mysqlmeta", nil)

	resp, _ := app.Test(req)
	payload := utils.ReadBody[map[string]any](resp.Body)

	assert.Equal(t, resp.StatusCode, fiber.StatusCreated)
	var foo float64
	assert.IsType(t, foo, payload["created"])
	assert.Equal(t, int(payload["created"].(float64)), 1)

}

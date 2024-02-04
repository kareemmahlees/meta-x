package routes_test

import (
	"context"
	"log"
	"net/http/httptest"
	"strings"
	"testing"

	"meta-x/internal"
	routes "meta-x/internal/rest"
	"meta-x/lib"
	"meta-x/models"
	"meta-x/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DatabaseRoutesTestSuite struct {
	suite.Suite
	providers        []string
	pgContainer      *utils.PostgresContainer
	pgConnection     *sqlx.DB
	mysqlContainer   *utils.MySQLContainer
	mysqlConnection  *sqlx.DB
	sqliteConnection *sqlx.DB
	ctx              context.Context
}

func (suite *DatabaseRoutesTestSuite) getConnection(provider string) *sqlx.DB {
	switch provider {
	case lib.SQLITE3:
		return suite.sqliteConnection
	case lib.PSQL:
		return suite.pgConnection
	case lib.MYSQL:
		return suite.mysqlConnection
	default:
		return suite.sqliteConnection
	}
}

func (suite *DatabaseRoutesTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	suite.sqliteConnection, _ = internal.InitDBConn(lib.SQLITE3, ":memory:")

	pgContainer, err := utils.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.pgContainer = pgContainer
	suite.pgConnection, _ = internal.InitDBConn(lib.PSQL, pgContainer.ConnectionString)

	mysqlContainer, err := utils.CreateMySQLContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.mysqlContainer = mysqlContainer
	suite.mysqlConnection, _ = internal.InitDBConn(lib.MYSQL, mysqlContainer.ConnectionString)

	suite.providers = []string{lib.SQLITE3, lib.PSQL, lib.MYSQL}
}

func (suite *DatabaseRoutesTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}

	if err := suite.mysqlContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating mysql container: %s", err)
	}
	suite.sqliteConnection.Close()
	suite.pgConnection.Close()
	suite.mysqlConnection.Close()
}

func (suite *DatabaseRoutesTestSuite) TestRegisterDatabaseRoutes() {
	t := suite.T()
	app := fiber.New()
	routes.RegisterDatabasesRoutes(app, nil)

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
		path:   "/database",
	})

}

func (suite *DatabaseRoutesTestSuite) TestHandleListDatabases() {
	t := suite.T()
	for _, provider := range suite.providers {
		app := utils.NewTestingFiberApp(provider)
		con := suite.getConnection(provider)
		routes.RegisterDatabasesRoutes(app, con)

		req := httptest.NewRequest("GET", "http://localhost:5522/database", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()
		payload := utils.DecodeBody[models.ListDatabasesResp](resp.Body)

		assert.Equal(t, resp.StatusCode, fiber.StatusOK)
		assert.NotZero(t, len(payload.Databases))
	}
}

func (suite *DatabaseRoutesTestSuite) TestHandleCreateDatabase() {
	t := suite.T()
	for _, provider := range suite.providers {
		app := utils.NewTestingFiberApp(provider)
		con := suite.getConnection(provider)
		routes.RegisterDatabasesRoutes(app, con)

		req := httptest.NewRequest("POST", "http://localhost:5522/database", strings.NewReader(`
			{
				"name":"test"
			}
		`))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		defer resp.Body.Close()
		payload := utils.DecodeBody[models.SuccessResp](resp.Body)

		assert.Equal(t, resp.StatusCode, fiber.StatusCreated)
		assert.True(t, payload.Success)
	}

}

func TestDatabaseRoutesTestSuite(t *testing.T) {
	suite.Run(t, new(DatabaseRoutesTestSuite))
}

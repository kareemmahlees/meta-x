package routes_test

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kareemmahlees/meta-x/internal"
	routes "github.com/kareemmahlees/meta-x/internal/rest"
	"github.com/kareemmahlees/meta-x/lib"
	"github.com/kareemmahlees/meta-x/models"
	"github.com/kareemmahlees/meta-x/utils"

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

	var routes []utils.FiberRoute
	for _, route := range app.GetRoutes() {
		routes = append(routes, utils.FiberRoute{
			Method: route.Method,
			Path:   route.Path,
		})
	}

	assert.Contains(t, routes, utils.FiberRoute{
		Method: "GET",
		Path:   "/database",
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

func (suite *DatabaseRoutesTestSuite) TestHandleCreateDatabasePassing() {
	t := suite.T()
	for _, provider := range suite.providers {
		app := utils.NewTestingFiberApp(provider)
		con := suite.getConnection(provider)
		routes.RegisterDatabasesRoutes(app, con)

		passingBody, _ := utils.EncodeBody(models.CreatePgMySqlDBPayload{
			Name: "testing",
		})
		passing := utils.RequestTesting[models.SuccessResp]{
			ReqMethod: http.MethodPost,
			ReqUrl:    "/database",
			ReqBody:   passingBody,
		}

		decodedRes, rawRes := passing.RunRequest(app)
		assert.Equal(t, rawRes.StatusCode, fiber.StatusCreated)
		assert.True(t, decodedRes.Success)
	}

}

func (suite *DatabaseRoutesTestSuite) TestHandleCreateDatabaseFailing() {

	t := suite.T()
	for _, provider := range suite.providers {
		app := utils.NewTestingFiberApp(provider)
		con := suite.getConnection(provider)
		routes.RegisterDatabasesRoutes(app, con)

		failingUnprocessableEntity := utils.RequestTesting[models.ErrResp]{
			ReqMethod: http.MethodPost,
			ReqUrl:    "/database",
		}
		decodedRes, rawRes := failingUnprocessableEntity.RunRequest(app)
		assert.Equal(t, http.StatusUnprocessableEntity, rawRes.StatusCode)
		assert.Contains(t, decodedRes.Message, "Unprocessable Entity")

		failingBadRequestBody, _ := utils.EncodeBody(models.CreatePgMySqlDBPayload{
			Name: "",
		})
		failingBadRequest := utils.RequestTesting[models.ErrResp]{
			ReqMethod: http.MethodPost,
			ReqUrl:    "/database",
			ReqBody:   failingBadRequestBody,
		}
		decodedRes, rawRes = failingBadRequest.RunRequest(app)
		assert.Equal(t, http.StatusBadRequest, rawRes.StatusCode)
		assert.Len(t, decodedRes.Message, 1)
	}
}

func TestDatabaseRoutesTestSuite(t *testing.T) {
	suite.Run(t, new(DatabaseRoutesTestSuite))
}

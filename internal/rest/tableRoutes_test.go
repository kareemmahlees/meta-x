package routes_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"meta-x/internal"
	"meta-x/internal/db"
	routes "meta-x/internal/rest"
	"meta-x/lib"
	"meta-x/models"
	"meta-x/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TableRoutesTestSuite struct {
	suite.Suite
	providers        []string
	pgContainer      *utils.PostgresContainer
	pgConnection     *sqlx.DB
	mysqlContainer   *utils.MySQLContainer
	mysqlConnection  *sqlx.DB
	sqliteConnection *sqlx.DB
	ctx              context.Context
}

func (suite *TableRoutesTestSuite) getConnection(provider string) *sqlx.DB {
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

func (suite *TableRoutesTestSuite) SetupSuite() {
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

func (suite *TableRoutesTestSuite) TearDownSuite() {
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

func (suite *TableRoutesTestSuite) BeforeTest(suiteName, testName string) {
	queryString := "CREATE TABLE test ( name varchar(255) );"
	_, _ = suite.sqliteConnection.Exec(queryString)
	_, _ = suite.pgConnection.Exec(queryString)
	_, _ = suite.mysqlConnection.Exec(queryString)
}

func (suite *TableRoutesTestSuite) AfterTest(suiteName, testName string) {
	queryString := "DROP TABLE test;"
	_, _ = suite.sqliteConnection.Exec(queryString)
	_, _ = suite.pgConnection.Exec(queryString)
	_, _ = suite.mysqlConnection.Exec(queryString)
}

func (suite *TableRoutesTestSuite) TestRegisterTablesRoutes() {
	t := suite.T()
	app := fiber.New()

	routes.RegisterTablesRoutes(app, nil)

	var routes []utils.FiberRoute
	for _, route := range app.GetRoutes() {
		routes = append(routes, utils.FiberRoute{
			Method: route.Method,
			Path:   route.Path,
		})
	}

	assert.Contains(t, routes, utils.FiberRoute{
		Method: "GET",
		Path:   "/table",
	})

}

func (suite *TableRoutesTestSuite) TestHandleDescribeTable() {
	t := suite.T()
	for _, provider := range suite.providers {
		app := utils.NewTestingFiberApp(provider)
		con := suite.getConnection(provider)
		routes.RegisterTablesRoutes(app, con)

		req := httptest.NewRequest("GET", "http://localhost:5522/table/test/describe", nil)

		resp, _ := app.Test(req)

		tableFields := utils.DecodeBody[[]models.TableInfoResp](resp.Body)
		assert.NotEmpty(t, tableFields)
		assert.Equal(t, tableFields[0].Name, "name")
	}

}

func (suite *TableRoutesTestSuite) TestHandleListTables() {
	t := suite.T()
	for _, provider := range suite.providers {
		app := utils.NewTestingFiberApp(provider)
		con := suite.getConnection(provider)
		routes.RegisterTablesRoutes(app, con)

		req := httptest.NewRequest("GET", "http://localhost:5522/table", nil)

		resp, _ := app.Test(req)

		body := utils.DecodeBody[models.ListTablesResp](resp.Body)
		tables := utils.SliceOfPointersToSliceOfValues(body.Tables)
		assert.NotEmpty(t, tables)
		assert.Contains(t, tables, "test")
	}
}

func (suite *TableRoutesTestSuite) TestHandleCreateTable() {
	t := suite.T()
	for idx, provider := range suite.providers {
		app := utils.NewTestingFiberApp(provider)
		con := suite.getConnection(provider)
		routes.RegisterTablesRoutes(app, con)

		req := httptest.NewRequest("POST", fmt.Sprintf("http://localhost:5522/table/test%d", idx), strings.NewReader(fmt.Sprintf(`
		[{	
				"column_name":"test%d",
				"type": "varchar(255)",
				"nullable": true,
				"default": "kareem",
				"unique": true
		}]`, idx)))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		body := utils.DecodeBody[models.CreateTableResp](resp.Body)
		assert.Equal(t, body.Created, fmt.Sprintf("test%d", idx))

		tables, _ := db.ListTables(con, provider)
		convertedTables := utils.SliceOfPointersToSliceOfValues(tables)
		assert.NotEmpty(t, convertedTables)
		assert.Contains(t, convertedTables, fmt.Sprintf("test%d", idx))
	}

}

func (suite *TableRoutesTestSuite) TestHandleAddColumn() {
	t := suite.T()
	for idx, provider := range suite.providers {
		app := utils.NewTestingFiberApp(provider)
		con := suite.getConnection(provider)
		routes.RegisterTablesRoutes(app, con)

		reqBody := utils.EncodeBody(models.AddModifyColumnPayload{ColName: fmt.Sprintf("test%d", idx), Type: "varchar(255)"})

		req := httptest.NewRequest(http.MethodPost, "http://localhost:5522/table/test/column/add", reqBody)
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		respBody := utils.DecodeBody[models.SuccessResp](resp.Body)
		assert.True(t, respBody.Success)

		tableInfo, _ := db.GetTableInfo(con, "test", provider)
		convertedTableInfo := utils.SliceOfPointersToSliceOfValues(tableInfo)

		var columnType string
		var key any

		switch provider {
		case lib.PSQL:
			columnType = "character varying"
		case lib.MYSQL:
			key = []uint8{}
			columnType = "varchar(255)"
		case lib.SQLITE3:
			columnType = "varchar(255)"
		}

		assert.Contains(t, convertedTableInfo,
			models.TableInfoResp{Name: fmt.Sprintf("test%d", idx),
				Type:     columnType,
				Nullable: "YES",
				Key:      key,
				Default:  nil})
	}
}

func (suite *TableRoutesTestSuite) TestHandleModifyColumn() {
	t := suite.T()
	for idx, provider := range suite.providers {
		app := utils.NewTestingFiberApp(provider)
		con := suite.getConnection(provider)
		routes.RegisterTablesRoutes(app, con)

		reqBody := utils.EncodeBody(models.AddModifyColumnPayload{ColName: "name", Type: fmt.Sprintf("varchar(5%d)", idx)})

		req := httptest.NewRequest(http.MethodPut, "http://localhost:5522/table/test/column/modify", reqBody)
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		respBody := utils.DecodeBody[models.SuccessResp](resp.Body)

		switch provider {
		case lib.SQLITE3:
			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		default:
			assert.True(t, respBody.Success)
		}

	}
}

func (suite *TableRoutesTestSuite) TestHandleDeleteColumn() {
	t := suite.T()
	for _, provider := range suite.providers {
		app := utils.NewTestingFiberApp(provider)
		con := suite.getConnection(provider)
		routes.RegisterTablesRoutes(app, con)

		reqBody := utils.EncodeBody(models.DeleteColumnPayload{ColName: "name"})

		req := httptest.NewRequest(http.MethodDelete, "http://localhost:5522/table/test/column/delete", reqBody)
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		switch provider {
		// those two forbid to delete all columns of table
		case lib.SQLITE3, lib.MYSQL:
			assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		case lib.PSQL:
			respBody := utils.DecodeBody[models.SuccessResp](resp.Body)
			assert.True(t, respBody.Success)

			tableInfo, _ := db.GetTableInfo(con, "test", provider)
			assert.Empty(t, tableInfo)
		}

	}

}

func (suite *TableRoutesTestSuite) TestHandleDeleteTable() {
	t := suite.T()
	for _, provider := range suite.providers {
		app := utils.NewTestingFiberApp(provider)
		con := suite.getConnection(provider)
		routes.RegisterTablesRoutes(app, con)

		req := httptest.NewRequest("DELETE", "http://localhost:5522/table/test", nil)

		resp, _ := app.Test(req)

		payload := utils.DecodeBody[models.SuccessResp](resp.Body)
		assert.True(t, payload.Success)
	}

}

func TestTableRoutesTestSuite(t *testing.T) {
	suite.Run(t, new(TableRoutesTestSuite))
}

package db_test

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"reflect"
// 	"testing"

// 	"github.com/kareemmahlees/meta-x/internal"
// 	"github.com/kareemmahlees/meta-x/internal/db"
// 	"github.com/kareemmahlees/meta-x/lib"
// 	"github.com/kareemmahlees/meta-x/models"
// 	"github.com/kareemmahlees/meta-x/utils"

// 	"github.com/jmoiron/sqlx"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/suite"
// )

// type TablesTestSuite struct {
// 	suite.Suite
// 	providers        []string
// 	pgContainer      *utils.PostgresContainer
// 	pgConnection     *sqlx.DB
// 	mysqlContainer   *utils.MySQLContainer
// 	mysqlConnection  *sqlx.DB
// 	sqliteConnection *sqlx.DB
// 	ctx              context.Context
// }

// func (suite *TablesTestSuite) getConnection(provider string) *sqlx.DB {
// 	switch provider {
// 	case lib.SQLITE3:
// 		return suite.sqliteConnection
// 	case lib.PSQL:
// 		return suite.pgConnection
// 	case lib.MYSQL:
// 		return suite.mysqlConnection
// 	default:
// 		return suite.sqliteConnection
// 	}
// }

// func (suite *TablesTestSuite) SetupSuite() {
// 	suite.ctx = context.Background()

// 	suite.sqliteConnection, _ = internal.InitDBConn(lib.SQLITE3, ":memory:")

// 	pgContainer, err := utils.CreatePostgresContainer(suite.ctx)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	suite.pgContainer = pgContainer
// 	suite.pgConnection, _ = internal.InitDBConn(lib.PSQL, pgContainer.ConnectionString)

// 	mysqlContainer, err := utils.CreateMySQLContainer(suite.ctx)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	suite.mysqlContainer = mysqlContainer
// 	suite.mysqlConnection, _ = internal.InitDBConn(lib.MYSQL, mysqlContainer.ConnectionString)

// 	suite.providers = []string{lib.SQLITE3, lib.PSQL, lib.MYSQL}
// }

// func (suite *TablesTestSuite) TearDownSuite() {
// 	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
// 		log.Fatalf("error terminating postgres container: %s", err)
// 	}

// 	if err := suite.mysqlContainer.Terminate(suite.ctx); err != nil {
// 		log.Fatalf("error terminating mysql container: %s", err)
// 	}
// 	suite.sqliteConnection.Close()
// 	suite.pgConnection.Close()
// 	suite.mysqlConnection.Close()

// }

// func (suite *TablesTestSuite) BeforeTest(suiteName, testName string) {
// 	queryString := "CREATE TABLE test ( name varchar(255) );"
// 	_, _ = suite.sqliteConnection.Exec(queryString)
// 	_, _ = suite.pgConnection.Exec(queryString)
// 	_, _ = suite.mysqlConnection.Exec(queryString)
// }

// func (suite *TablesTestSuite) AfterTest(suiteName, testName string) {
// 	queryString := "DROP TABLE test;"
// 	_, _ = suite.sqliteConnection.Exec(queryString)
// 	_, _ = suite.pgConnection.Exec(queryString)
// 	_, _ = suite.mysqlConnection.Exec(queryString)
// }

// func (suite *TablesTestSuite) TestDescribeTable() {
// 	t := suite.T()
// 	for _, provider := range suite.providers {
// 		con := suite.getConnection(provider)
// 		tableInfo, err := db.GetTableInfo(con, "test", provider)

// 		assert.Nil(t, err)
// 		assert.NotEmpty(t, tableInfo)
// 		assert.Equal(t, tableInfo[0].Name, "name")
// 		switch provider {
// 		case lib.PSQL:
// 			assert.Equal(t, tableInfo[0].Type, "character varying")
// 		default:
// 			assert.Equal(t, tableInfo[0].Type, "varchar(255)")
// 		}
// 	}
// }

// func (suite *TablesTestSuite) TestListTable() {
// 	t := suite.T()
// 	for _, provider := range suite.providers {
// 		con := suite.getConnection(provider)
// 		result, err := db.ListTables(con, provider)
// 		assert.Nil(t, err)
// 		assert.NotEmpty(t, result)
// 		assert.IsType(t, reflect.Slice, reflect.TypeOf(result).Kind())
// 	}

// }

// func (suite *TablesTestSuite) TestCreateTable() {
// 	t := suite.T()
// 	for _, provider := range suite.providers {
// 		con := suite.getConnection(provider)

// 		err := db.CreateTable(con, "test_create_table", []models.CreateTablePayload{
// 			{
// 				ColName:  "test",
// 				Type:     "varchar(255)",
// 				Default:  "defaultText",
// 				Unique:   true,
// 				Nullable: false,
// 			},
// 		})
// 		assert.Nil(t, err)

// 		result, _ := db.ListTables(con, provider)
// 		convertedResult := utils.SliceOfPointersToSliceOfValues(result)
// 		assert.Contains(t, convertedResult, "test_create_table")

// 		err = db.CreateTable(con, "123456", nil)
// 		assert.NotNil(t, err)
// 	}
// }

// func (suite *TablesTestSuite) TestAddColumn() {
// 	t := suite.T()
// 	for idx, provider := range suite.providers {
// 		con := suite.getConnection(provider)
// 		err := db.AddColumn(con, "test", models.AddModifyColumnPayload{
// 			ColName: fmt.Sprintf("test%d", idx),
// 			Type:    "int",
// 		})
// 		assert.Nil(t, err)

// 		result, _ := db.GetTableInfo(con, "test", provider)
// 		assert.Greater(t, len(result), 1)
// 	}
// }

// func (suite *TablesTestSuite) TestModifyColumn() {
// 	t := suite.T()
// 	for _, provider := range suite.providers {
// 		con := suite.getConnection(provider)
// 		err := db.UpdateColumn(con, provider, "test", models.AddModifyColumnPayload{
// 			ColName: "name",
// 			Type:    "varchar(255)",
// 		})
// 		switch provider {
// 		case lib.MYSQL, lib.PSQL:
// 			assert.Nil(t, err)
// 		case lib.SQLITE3:
// 			assert.Error(t, err)
// 		}
// 	}
// }

// func (suite *TablesTestSuite) TestDeleteColumn() {
// 	t := suite.T()
// 	for _, provider := range suite.providers {
// 		con := suite.getConnection(provider)
// 		err := db.DeleteColumn(con, "test", models.DeleteColumnPayload{
// 			ColName: "name",
// 		})
// 		switch provider {
// 		case lib.SQLITE3, lib.MYSQL:
// 			assert.Error(t, err)
// 		default:
// 			assert.Nil(t, err)
// 			result, _ := db.GetTableInfo(con, "test", provider)
// 			assert.Zero(t, len(result))
// 		}

// 	}
// }

// func (suite *TablesTestSuite) TestDeleteTable() {
// 	t := suite.T()
// 	for _, provider := range suite.providers {
// 		con := suite.getConnection(provider)
// 		err := db.DeleteTable(con, "test")
// 		assert.Nil(t, err)

// 		result, _ := db.ListTables(con, provider)
// 		assert.NotContains(t, result, "test")
// 	}

// }

// func TestTablesTestSuite(t *testing.T) {
// 	suite.Run(t, new(TablesTestSuite))
// }

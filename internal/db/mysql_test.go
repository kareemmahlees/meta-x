package db

import (
	"context"
	"log"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/kareemmahlees/meta-x/lib"
	"github.com/kareemmahlees/meta-x/models"
	"github.com/kareemmahlees/meta-x/utils"
	"github.com/stretchr/testify/suite"
)

type MySQLTestSuite struct {
	suite.Suite
	mysqlContainer  *utils.MySQLContainer
	mysqlConnection *sqlx.DB
	mysqlProvider   *mysqlProvider
	ctx             context.Context
}

func (suite *MySQLTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	mysqlContainer, err := utils.CreateMySQLContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.mysqlContainer = mysqlContainer

	mysqlConfig := utils.NewMySQLConfig(&mysqlContainer.ConnectionString, nil)
	conn, err := InitDBConn(lib.MYSQL, mysqlConfig)
	if err != nil {
		log.Fatal(err)
	}
	suite.mysqlConnection = conn

	mysqlProvider := NewMySQLProvider(conn)
	suite.mysqlProvider = mysqlProvider
}

func (suite *MySQLTestSuite) TearDownSuite() {
	if err := suite.mysqlContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
	suite.mysqlConnection.Close()
}

func (suite *MySQLTestSuite) BeforeTest(suiteName, testName string) {
	queryString := "CREATE TABLE test ( name varchar(255) );"
	_, _ = suite.mysqlConnection.Exec(queryString)
}

func (suite *MySQLTestSuite) AfterTest(suiteName, testName string) {
	queryString := "DROP TABLE test;"
	_, _ = suite.mysqlConnection.Exec(queryString)
}

func (suite *MySQLTestSuite) TestListDBs() {
	assert := suite.Assert()

	dbs, err := suite.mysqlProvider.ListDBs()

	assert.Nil(err)
	assert.Greater(len(dbs), 0)

}
func (suite *MySQLTestSuite) TestCreateDB() {
	assert := suite.Assert()

	err := suite.mysqlProvider.CreateDB("metax")
	assert.Nil(err)

	err = suite.mysqlProvider.CreateDB("true")
	assert.NotNil(err)

}
func (suite *MySQLTestSuite) TestGetTable() {
	assert := suite.Assert()

	tableInfo, err := suite.mysqlProvider.GetTable("test")

	assert.Nil(err)
	assert.NotEmpty(tableInfo)
	assert.Equal(tableInfo[0].Name, "name")
	assert.Equal(tableInfo[0].Type, "varchar(255)")

}
func (suite *MySQLTestSuite) TestListTables() {
	assert := suite.Assert()

	result, err := suite.mysqlProvider.ListTables()

	assert.Nil(err)
	assert.NotEmpty(result)

}
func (suite *MySQLTestSuite) TestCreateTable() {
	assert := suite.Assert()

	err := suite.mysqlProvider.CreateTable("test_create_table", []models.CreateTablePayload{
		{
			ColName:  "test",
			Type:     "varchar(255)",
			Default:  "defaultText",
			Unique:   true,
			Nullable: false,
		},
	})
	assert.Nil(err)

	result, _ := suite.mysqlProvider.ListTables()
	convertedResult := utils.SliceOfPointersToSliceOfValues(result)
	assert.Contains(convertedResult, "test_create_table")

	err = suite.mysqlProvider.CreateTable("123456", nil)
	assert.NotNil(err)

}
func (suite *MySQLTestSuite) TestDeleteTable() {
	assert := suite.Assert()
	err := suite.mysqlProvider.DeleteTable("test")
	assert.Nil(err)

	result, _ := suite.mysqlProvider.ListTables()
	assert.NotContains(result, "test")

}
func (suite *MySQLTestSuite) TestAddColumn() {
	assert := suite.Assert()

	err := suite.mysqlProvider.AddColumn("test", models.AddModifyColumnPayload{
		ColName: "test1",
		Type:    "int",
	})
	assert.Nil(err)

	result, _ := suite.mysqlProvider.GetTable("test")
	assert.NotEmpty(result)

}
func (suite *MySQLTestSuite) TestUpdateColumn() {
	assert := suite.Assert()

	err := suite.mysqlProvider.UpdateColumn("test", models.AddModifyColumnPayload{
		ColName: "name",
		Type:    "varchar(255)",
	})
	assert.Nil(err)

}
func (suite *MySQLTestSuite) TestDeleteColumn() {
	assert := suite.Assert()

	err := suite.mysqlProvider.DeleteColumn("test", models.DeleteColumnPayload{
		ColName: "name",
	})
	assert.Error(err)

}

func TestMySQLTestSuite(t *testing.T) {
	suite.Run(t, new(MySQLTestSuite))
}

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

type SQLiteTestSuite struct {
	suite.Suite
	sqliteConnection *sqlx.DB
	sqliteProvider   *SqliteProvider
	ctx              context.Context
}

func (suite *SQLiteTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	sqliteConfig := utils.NewSQLiteConfig(":memory:")

	conn, err := InitDBConn(lib.SQLITE3, sqliteConfig)
	if err != nil {
		log.Fatal(err)
	}
	suite.sqliteConnection = conn

	sqliteProvider := NewSQLiteProvider(conn)
	suite.sqliteProvider = sqliteProvider
}

func (suite *SQLiteTestSuite) TearDownSuite() {
	suite.sqliteConnection.Close()
}

func (suite *SQLiteTestSuite) BeforeTest(suiteName, testName string) {
	queryString := "CREATE TABLE test ( name varchar(255) );"
	_, _ = suite.sqliteConnection.Exec(queryString)
}

func (suite *SQLiteTestSuite) AfterTest(suiteName, testName string) {
	queryString := "DROP TABLE test;"
	_, _ = suite.sqliteConnection.Exec(queryString)
}

func (suite *SQLiteTestSuite) TestListDBs() {
	assert := suite.Assert()

	dbs, err := suite.sqliteProvider.ListDBs()

	assert.Nil(err)
	assert.NotEmpty(dbs)
}

func (suite *SQLiteTestSuite) TestCreateDatabase() {
	assert := suite.Assert()

	err := suite.sqliteProvider.CreateDB("something")
	assert.Error(err)
}

func (suite *SQLiteTestSuite) TestGetTable() {
	assert := suite.Assert()

	tableInfo, err := suite.sqliteProvider.GetTable("test")

	assert.Nil(err)
	assert.NotEmpty(tableInfo)

	tableInfo, _ = suite.sqliteProvider.GetTable("notexist")
	assert.Empty(tableInfo)
}
func (suite *SQLiteTestSuite) TestListTables() {
	assert := suite.Assert()

	result, err := suite.sqliteProvider.ListTables()

	assert.Nil(err)
	assert.NotEmpty(result)
}
func (suite *SQLiteTestSuite) TestCreateTable() {
	assert := suite.Assert()
	err := suite.sqliteProvider.CreateTable("test_create_table", []models.CreateTablePayload{
		{
			ColName:  "test",
			Type:     "varchar(255)",
			Default:  "defaultText",
			Unique:   true,
			Nullable: false,
		},
	})
	assert.Nil(err)

	result, _ := suite.sqliteProvider.ListTables()
	convertedResult := utils.SliceOfPointersToSliceOfValues(result)
	assert.Contains(convertedResult, "test_create_table")

	err = suite.sqliteProvider.CreateTable("123456", nil)
	assert.NotNil(err)
}
func (suite *SQLiteTestSuite) TestDeleteTable() {
	assert := suite.Assert()
	err := suite.sqliteProvider.DeleteTable("test")
	assert.Nil(err)
	err = suite.sqliteProvider.DeleteTable("not exist")
	assert.Error(err)
}
func (suite *SQLiteTestSuite) TestAddColumn() {
	assert := suite.Assert()
	err := suite.sqliteProvider.AddColumn("test", models.AddModifyColumnPayload{
		ColName: "test1",
		Type:    "int",
	})
	assert.Nil(err)

	result, _ := suite.sqliteProvider.GetTable("test")
	assert.NotEmpty(len(result))
}
func (suite *SQLiteTestSuite) TestUpdateColumn() {
	assert := suite.Assert()
	err := suite.sqliteProvider.UpdateColumn("test", models.AddModifyColumnPayload{
		ColName: "name",
		Type:    "varchar(255)",
	})
	assert.Error(err)
}
func (suite *SQLiteTestSuite) TestDeleteColumn() {
	assert := suite.Assert()
	err := suite.sqliteProvider.DeleteColumn("test", models.DeleteColumnPayload{
		ColName: "name",
	})
	assert.Error(err)
}

func TestSqliteTestSuite(t *testing.T) {
	suite.Run(t, new(SQLiteTestSuite))
}

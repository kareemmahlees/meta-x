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

type PgTestSuite struct {
	suite.Suite
	pgContainer  *utils.PostgresContainer
	pgConnection *sqlx.DB
	pgProvider   *pgProvider
	ctx          context.Context
}

func (suite *PgTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	pgContainer, err := utils.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.pgContainer = pgContainer

	pgConfig := utils.NewPGConfig(&pgContainer.ConnectionString, nil)
	conn, err := InitDBConn(lib.PSQL, pgConfig)
	if err != nil {
		log.Fatal(err)
	}
	suite.pgConnection = conn

	pgProvider := NewPgProvider(conn)
	suite.pgProvider = pgProvider
}

func (suite *PgTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
	suite.pgConnection.Close()
}

func (suite *PgTestSuite) BeforeTest(suiteName, testName string) {
	queryString := "CREATE TABLE test ( name varchar(255) );"
	_, _ = suite.pgConnection.Exec(queryString)
}

func (suite *PgTestSuite) AfterTest(suiteName, testName string) {
	queryString := "DROP TABLE test;"
	_, _ = suite.pgConnection.Exec(queryString)
}

func (suite *PgTestSuite) TestListDBs() {
	assert := suite.Assert()

	dbs, err := suite.pgProvider.ListDBs()

	assert.Nil(err)
	assert.NotEmpty(dbs)

}
func (suite *PgTestSuite) TestCreateDB() {
	assert := suite.Assert()

	err := suite.pgProvider.CreateDB("testing")
	assert.Nil(err)

	dbs, _ := suite.pgProvider.ListDBs()
	conv := utils.SliceOfPointersToSliceOfValues(dbs)
	assert.Contains(conv, "testing")
}
func (suite *PgTestSuite) TestGetTable() {
	assert := suite.Assert()

	tableInfo, err := suite.pgProvider.GetTable("test")

	assert.Nil(err)
	assert.NotEmpty(tableInfo)
	assert.Equal(tableInfo[0].Name, "name")
	assert.Equal(tableInfo[0].Type, "character varying")

}
func (suite *PgTestSuite) TestListTables() {
	assert := suite.Assert()

	result, err := suite.pgProvider.ListTables()

	assert.Nil(err)
	assert.NotEmpty(result)
}
func (suite *PgTestSuite) TestCreateTable() {
	assert := suite.Assert()

	err := suite.pgProvider.CreateTable("test_create_table", []models.CreateTablePayload{
		{
			ColName:  "test",
			Type:     "varchar(255)",
			Default:  "defaultText",
			Unique:   true,
			Nullable: false,
		},
	})
	assert.Nil(err)

	result, _ := suite.pgProvider.ListTables()
	convertedResult := utils.SliceOfPointersToSliceOfValues(result)
	assert.Contains(convertedResult, "test_create_table")

	err = suite.pgProvider.CreateTable("123456", nil)
	assert.Error(err)

}
func (suite *PgTestSuite) TestDeleteTable() {
	assert := suite.Assert()

	err := suite.pgProvider.DeleteTable("test")
	assert.Nil(err)

	result, _ := suite.pgProvider.ListTables()
	assert.NotContains(result, "test")

}
func (suite *PgTestSuite) TestAddColumn() {
	assert := suite.Assert()

	err := suite.pgProvider.AddColumn("test", models.AddModifyColumnPayload{
		ColName: "test1",
		Type:    "int",
	})
	assert.Nil(err)

	result, _ := suite.pgProvider.GetTable("test")
	assert.NotEmpty(result)

}
func (suite *PgTestSuite) TestUpdateColumn() {
	assert := suite.Assert()

	err := suite.pgProvider.UpdateColumn("test", models.AddModifyColumnPayload{
		ColName: "name",
		Type:    "varchar(255)",
	})
	assert.Nil(err)

}
func (suite *PgTestSuite) TestDeleteColumn() {
	assert := suite.Assert()
	err := suite.pgProvider.DeleteColumn("test", models.DeleteColumnPayload{
		ColName: "name",
	})
	assert.Nil(err)
	result, _ := suite.pgProvider.GetTable("test")
	assert.Zero(len(result))

}

func TestPgTestSuite(t *testing.T) {
	suite.Run(t, new(PgTestSuite))
}

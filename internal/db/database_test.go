package db_test

// import (
// 	"context"
// 	"log"
// 	"github.com/kareemmahlees/meta-x/internal"
// 	"github.com/kareemmahlees/meta-x/internal/db"
// 	"github.com/kareemmahlees/meta-x/lib"
// 	"github.com/kareemmahlees/meta-x/utils"
// 	"testing"

// 	"github.com/jmoiron/sqlx"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/suite"
// )

// type DatabaseTestSuite struct {
// 	suite.Suite
// 	providers        []string
// 	pgContainer      *utils.PostgresContainer
// 	pgConnection     *sqlx.DB
// 	mysqlContainer   *utils.MySQLContainer
// 	mysqlConnection  *sqlx.DB
// 	sqliteConnection *sqlx.DB
// 	ctx              context.Context
// }

// func (suite *DatabaseTestSuite) getConnection(provider string) *sqlx.DB {
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

// func (suite *DatabaseTestSuite) SetupSuite() {
// 	suite.ctx = context.Background()

// 	var err error

// 	suite.sqliteConnection, err = internal.InitDBConn(lib.SQLITE3, ":memory:")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	pgContainer, err := utils.CreatePostgresContainer(suite.ctx)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	suite.pgContainer = pgContainer
// 	suite.pgConnection, err = internal.InitDBConn(lib.PSQL, pgContainer.ConnectionString)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	mysqlContainer, err := utils.CreateMySQLContainer(suite.ctx)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	suite.mysqlContainer = mysqlContainer
// 	suite.mysqlConnection, err = internal.InitDBConn(lib.MYSQL, mysqlContainer.ConnectionString)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	suite.providers = []string{lib.SQLITE3, lib.PSQL, lib.MYSQL}
// }

// func (suite *DatabaseTestSuite) TearDownSuite() {
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

// func (suite *DatabaseTestSuite) TestListDatabases() {
// 	t := suite.T()

// 	for _, provider := range suite.providers {
// 		con := suite.getConnection(provider)

// 		var dbs []*string
// 		var err error

// 		switch provider {
// 		case lib.SQLITE3:
// 			dbs, err = db.ListDatabasesSqlite(suite.sqliteConnection)
// 		case lib.PSQL, lib.MYSQL:
// 			dbs, err = db.ListDatabasesPgMySQL(con, provider)
// 		}

// 		assert.Nil(t, err)
// 		assert.Greater(t, len(dbs), 0)

// 		switch provider {
// 		// we reverse the order here to intentionally make the query fail
// 		case lib.SQLITE3:
// 			_, err = db.ListDatabasesPgMySQL(suite.sqliteConnection, provider)
// 		case lib.PSQL, lib.MYSQL:
// 			_, err = db.ListDatabasesSqlite(con)
// 		}
// 		assert.NotNil(t, err)
// 	}
// }

// func (suite *DatabaseTestSuite) TestCreateDatabase() {
// 	t := suite.T()

// 	for _, provider := range suite.providers {
// 		con := suite.getConnection(provider)

// 		var err error

// 		switch provider {
// 		case lib.PSQL, lib.MYSQL:
// 			err = db.CreatePgMysqlDatabase(con, provider, "metax")
// 			assert.Nil(t, err)
// 		}

// 		switch provider {
// 		case lib.PSQL, lib.MYSQL:
// 			err = db.CreatePgMysqlDatabase(con, provider, "true")
// 			assert.NotNil(t, err)
// 		}
// 	}
// }

// func TestDatabaseTestSuite(t *testing.T) {
// 	suite.Run(t, new(DatabaseTestSuite))
// }

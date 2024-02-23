package internal_test

import (
	"context"
	"log"
	"github.com/kareemmahlees/meta-x/internal"
	"github.com/kareemmahlees/meta-x/lib"
	"github.com/kareemmahlees/meta-x/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DBTestSuite struct {
	suite.Suite
	pgContainer    *utils.PostgresContainer
	mysqlContainer *utils.MySQLContainer
	ctx            context.Context
}

func (suite *DBTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	pgContainer, err := utils.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.pgContainer = pgContainer

	mysqlContainer, err := utils.CreateMySQLContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.mysqlContainer = mysqlContainer
}

func (suite *DBTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}

	if err := suite.mysqlContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating mysql container: %s", err)
	}

}

func (suite *DBTestSuite) TestInitDBConn() {
	providers := []string{lib.SQLITE3, lib.PSQL, lib.MYSQL}
	t := suite.T()
	for _, provider := range providers {

		var cfg string
		switch provider {
		case lib.PSQL:
			cfg = suite.pgContainer.ConnectionString
		case lib.MYSQL:
			cfg = suite.mysqlContainer.ConnectionString
		case lib.SQLITE3:
			cfg = ":memory:"
		}

		conn, err := internal.InitDBConn(provider, cfg)

		conn.Close()

		assert.Nil(t, err)

	}
}

func TestDBTestSuite(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}

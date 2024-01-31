package utils

import (
	"context"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresContainer struct {
	*postgres.PostgresContainer
	ConnectionString string
}

type MySQLContainer struct {
	*mysql.MySQLContainer
	ConnectionString string
}

func CreatePostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15"),
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	return &PostgresContainer{
		PostgresContainer: pgContainer,
		ConnectionString:  connStr,
	}, nil
}

func CreateMySQLContainer(ctx context.Context) (*MySQLContainer, error) {
	mysqlContainer, err := mysql.RunContainer(ctx,
		testcontainers.WithImage("mysql:latest"),
		mysql.WithDatabase("test-db"),
		mysql.WithUsername("root"),
		mysql.WithPassword("root"),
	)
	if err != nil {
		return nil, err
	}
	connStr, err := mysqlContainer.ConnectionString(ctx)
	if err != nil {
		return nil, err
	}

	return &MySQLContainer{
		MySQLContainer:   mysqlContainer,
		ConnectionString: connStr,
	}, nil
}

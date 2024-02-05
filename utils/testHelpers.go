package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
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
		postgres.WithPassword("postgres"), // pragma: allowlist secret
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
		mysql.WithPassword("root"), // pragma: allowlist secret
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

func NewTestingFiberApp(provider string) *fiber.App {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("provider", provider)
		return c.Next()
	})
	return app
}

func EncodeBody[T any](body T) *bytes.Buffer {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		log.Fatal(err)
	}
	return &buf
}

func DecodeBody[K any](body io.ReadCloser) K {
	var parsedPayload K
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&parsedPayload)
	if err != nil {
		log.Fatal(err)
	}
	return parsedPayload
}

func SliceOfPointersToSliceOfValues[T any](s []*T) []T {
	v := make([]T, len(s))
	for i, p := range s {
		v[i] = *p
	}
	return v
}

type FiberRoute struct {
	Method string
	Path   string
}

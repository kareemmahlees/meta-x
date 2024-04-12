package cmd

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/kareemmahlees/meta-x/models"
	"github.com/kareemmahlees/meta-x/utils"
	"github.com/stretchr/testify/assert"
)

type MockStorage struct{}

func NewMockStorage() *MockStorage {
	return &MockStorage{}
}

func (ms *MockStorage) ListDBs() ([]*string, error) {
	return nil, nil
}
func (ms *MockStorage) CreateDB(dbName string) error {
	return nil
}
func (ms *MockStorage) GetTable(tableName string) ([]*models.TableInfoResp, error) {
	return nil, nil
}
func (ms *MockStorage) ListTables() ([]*string, error) {
	return nil, nil
}
func (ms *MockStorage) CreateTable(tableName string, data []models.CreateTablePayload) error {
	return nil
}
func (ms *MockStorage) DeleteTable(tableName string) error {
	return nil
}
func (ms *MockStorage) AddColumn(tableName string, data models.AddModifyColumnPayload) error {
	return nil
}
func (ms *MockStorage) UpdateColumn(tableName string, data models.AddModifyColumnPayload) error {
	return nil
}
func (ms *MockStorage) DeleteColumn(tableName string, data models.DeleteColumnPayload) error {
	return nil
}

func TestServe(t *testing.T) {
	listenCh := make(chan bool, 1)
	server := NewServer(NewMockStorage(), 5522, listenCh)
	defer server.Shutdown()

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatal(err)
		}
	}()

	assert.True(t, <-listenCh)

	testRoutes := []string{"graphql", "playground", "swagger"}

	for _, route := range testRoutes {
		foundRoute := server.app.GetRoute(route)
		assert.NotEmpty(t, foundRoute)

		request := utils.RequestTesting[any]{
			ReqMethod: http.MethodGet,
			ReqUrl:    fmt.Sprintf("/%s", route),
		}
		_, res := request.RunRequest(server.app)

		assert.NotEqual(t, http.StatusNotFound, res.StatusCode)

	}
}

func TestShutDown(t *testing.T) {
	listenCh := make(chan bool, 1)
	server := NewServer(NewMockStorage(), 5522, listenCh)

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatal(err)
		}
	}()

	assert.True(t, <-server.listenCh)

	err := server.Shutdown()

	assert.Nil(t, err)
	assert.False(t, <-server.listenCh)
}

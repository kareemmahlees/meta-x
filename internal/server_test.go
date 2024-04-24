package internal

import (
	"log"
	"testing"

	"github.com/kareemmahlees/meta-x/models"
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
	server := NewServer(NewMockStorage(), 5522)

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatal(err)
		}
	}()

	testRoutes := []string{"/graphql/*", "/playground/*", "/swagger/*"}
	registerdRoutes := []string{}

	for _, route := range server.router.Routes() {
		registerdRoutes = append(registerdRoutes, route.Pattern)
	}

	for _, route := range testRoutes {
		assert.Contains(t, registerdRoutes, route)
	}

}

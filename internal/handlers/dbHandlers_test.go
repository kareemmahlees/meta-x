package handlers

import (
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/kareemmahlees/meta-x/models"
	"github.com/kareemmahlees/meta-x/utils"
	"github.com/stretchr/testify/suite"
)

type MockDBExecutor struct{}

func NewMockDBExecutor() *MockDBExecutor {
	return &MockDBExecutor{}
}

func (md *MockDBExecutor) ListDBs() ([]*string, error) {
	db := "test"
	return []*string{&db}, nil
}
func (md *MockDBExecutor) CreateDB(dbName string) error {
	return nil
}

type DBHandlerTestSuite struct {
	suite.Suite
	api     humatest.TestAPI
	handler *DBHandler
}

func (suite *DBHandlerTestSuite) SetupSuite() {
	_, api := humatest.New(suite.T())
	storage := NewMockDBExecutor()

	handler := NewDBHandler(storage)
	handler.RegisterRoutes(api)

	suite.api = api
	suite.handler = handler
}

func (suite *DBHandlerTestSuite) TestHandleListDatabases() {
	assert := suite.Assert()
	resp := suite.api.Get("/database")

	assert.Equal(http.StatusOK, resp.Code)
	assert.Contains(resp.Body.String(), "test")
}

func (suite *DBHandlerTestSuite) TestHandleCreateDatabase() {
	assert := suite.Assert()
	t := suite.T()

	t.Run("should pass", func(t *testing.T) {
		resp := suite.api.Post("/database", map[string]any{
			"name": "testing",
		})

		assert.Equal(resp.Code, http.StatusCreated)

		decodedRes := utils.DecodeBody[models.SuccessResp](resp.Result().Body)
		assert.True(decodedRes.Success)

	})

	t.Run("should fail name required", func(t *testing.T) {
		resp := suite.api.Post("/database", map[string]any{
			"name": "",
		})
		assert.Equal(resp.Code, http.StatusBadRequest)

		decodedRes := utils.DecodeBody[huma.ErrorModel](resp.Result().Body)
		assert.Contains(decodedRes.Detail, "required")
	})
}

func TestDBHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(DBHandlerTestSuite))
}

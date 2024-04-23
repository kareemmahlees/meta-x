package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
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
	r       *chi.Mux
	handler *DBHandler
}

func (suite *DBHandlerTestSuite) SetupSuite() {
	r := chi.NewRouter()
	storage := NewMockDBExecutor()

	handler := NewDBHandler(storage)
	handler.RegisterRoutes(r)

	suite.r = r
	suite.handler = handler
}

func (suite *DBHandlerTestSuite) TestRegisterRoutes() {
	assert := suite.Assert()

	var routes []string
	for _, route := range suite.r.Routes() {
		routes = append(routes, route.Pattern)
	}

	assert.Contains(routes, "/database/*")
}

func (suite *DBHandlerTestSuite) TestHandleListDatabases() {
	assert := suite.Assert()

	assert.HTTPSuccess(suite.handler.handleListDatabases, http.MethodGet, "/database", nil)
	assert.HTTPBodyContains(suite.handler.handleListDatabases, http.MethodGet, "/database", nil, "test")
}

func (suite *DBHandlerTestSuite) TestHandleCreateDatabase() {
	assert := suite.Assert()
	t := suite.T()

	t.Run("should pass", func(t *testing.T) {
		passingBody, _ := utils.EncodeBody(models.CreatePgMySqlDBPayload{
			Name: "testing",
		})
		req, _ := http.NewRequest(http.MethodPost, "/database", passingBody)
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(suite.handler.handleCreateDatabase)
		handler.ServeHTTP(rr, req)
		assert.Equal(rr.Code, http.StatusCreated)

		decodedRes := utils.DecodeBody[models.SuccessResp](rr.Result().Body)
		assert.True(decodedRes.Success)

	})

	t.Run("should fail unproccessable entity", func(t *testing.T) {
		assert.HTTPError(suite.handler.handleCreateDatabase, http.MethodPost, "/database", nil)
	})

	t.Run("should fail bad request", func(t *testing.T) {
		failingBadRequestBody, _ := utils.EncodeBody(models.CreatePgMySqlDBPayload{
			Name: "",
		})
		req, _ := http.NewRequest(http.MethodPost, "/database", failingBadRequestBody)
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(suite.handler.handleCreateDatabase)
		handler.ServeHTTP(rr, req)
		assert.Equal(rr.Code, http.StatusBadRequest)

		decodedRes := utils.DecodeBody[models.ErrResp](rr.Result().Body)
		assert.Len(decodedRes.Message, 1)
	})
}

func TestDBHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(DBHandlerTestSuite))
}

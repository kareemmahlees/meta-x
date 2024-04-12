package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
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
	app *fiber.App
}

func (suite *DBHandlerTestSuite) SetupSuite() {
	app := fiber.New()
	storage := NewMockDBExecutor()

	handler := NewDBHandler(storage)
	handler.RegisterRoutes(app)

	suite.app = app
}

func (suite *DBHandlerTestSuite) TestRegisterRoutes() {
	assert := suite.Assert()

	var routes []utils.FiberRoute
	for _, route := range suite.app.GetRoutes() {
		routes = append(routes, utils.FiberRoute{
			Method: route.Method,
			Path:   route.Path,
		})
	}

	assert.Contains(routes, utils.FiberRoute{
		Method: "GET",
		Path:   "/database",
	})

}

func (suite *DBHandlerTestSuite) TestHandleListDatabases() {
	assert := suite.Assert()

	req := httptest.NewRequest("GET", "http://localhost:5522/database", nil)

	resp, _ := suite.app.Test(req)
	defer resp.Body.Close()
	payload := utils.DecodeBody[models.ListDatabasesResp](resp.Body)

	assert.Equal(resp.StatusCode, fiber.StatusOK)
	assert.NotEmpty(payload.Databases)

}

func (suite *DBHandlerTestSuite) TestHandleCreateDatabase() {
	assert := suite.Assert()
	t := suite.T()

	t.Run("should pass", func(t *testing.T) {
		passingBody, _ := utils.EncodeBody(models.CreatePgMySqlDBPayload{
			Name: "testing",
		})
		passing := utils.RequestTesting[models.SuccessResp]{
			ReqMethod: http.MethodPost,
			ReqUrl:    "/database",
			ReqBody:   passingBody,
		}

		decodedRes, rawRes := passing.RunRequest(suite.app)
		assert.Equal(rawRes.StatusCode, fiber.StatusCreated)
		assert.True(decodedRes.Success)
	})

	t.Run("should fail unproccessable entity", func(t *testing.T) {
		failingUnprocessableEntity := utils.RequestTesting[models.ErrResp]{
			ReqMethod: http.MethodPost,
			ReqUrl:    "/database",
		}
		decodedRes, rawRes := failingUnprocessableEntity.RunRequest(suite.app)
		assert.Equal(http.StatusUnprocessableEntity, rawRes.StatusCode)
		assert.Contains(decodedRes.Message, "Unprocessable Entity")

		failingBadRequestBody, _ := utils.EncodeBody(models.CreatePgMySqlDBPayload{
			Name: "",
		})
		failingBadRequest := utils.RequestTesting[models.ErrResp]{
			ReqMethod: http.MethodPost,
			ReqUrl:    "/database",
			ReqBody:   failingBadRequestBody,
		}
		decodedRes, rawRes = failingBadRequest.RunRequest(suite.app)
		assert.Equal(http.StatusBadRequest, rawRes.StatusCode)
		assert.Len(decodedRes.Message, 1)

	})

}

func TestDBHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(DBHandlerTestSuite))
}

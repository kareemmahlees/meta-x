package handlers

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kareemmahlees/meta-x/models"
	"github.com/kareemmahlees/meta-x/utils"
	"github.com/stretchr/testify/suite"
)

type MockTableExecutor struct{}

func NewMockTableExecutor() *MockTableExecutor {
	return &MockTableExecutor{}
}

func (ms *MockTableExecutor) GetTable(tableName string) ([]*models.TableInfoResp, error) {
	tableInfo := models.TableInfoResp{
		Name:     "name",
		Type:     "varchar(255)",
		Nullable: "Yes",
	}
	return []*models.TableInfoResp{&tableInfo}, nil
}
func (ms *MockTableExecutor) ListTables() ([]*string, error) {
	table := "test"
	return []*string{&table}, nil
}
func (ms *MockTableExecutor) CreateTable(tableName string, data []models.CreateTablePayload) error {
	return nil
}
func (ms *MockTableExecutor) DeleteTable(tableName string) error {
	return nil
}
func (ms *MockTableExecutor) AddColumn(tableName string, data models.AddModifyColumnPayload) error {
	return nil
}
func (ms *MockTableExecutor) UpdateColumn(tableName string, data models.AddModifyColumnPayload) error {
	return nil
}
func (ms *MockTableExecutor) DeleteColumn(tableName string, data models.DeleteColumnPayload) error {
	return nil
}

type TableHandlerTestSuite struct {
	suite.Suite
	app *fiber.App
}

func (suite *TableHandlerTestSuite) SetupSuite() {
	app := fiber.New()
	storage := NewMockTableExecutor()

	handler := NewTableHandler(storage)
	handler.RegisterRoutes(app)

	suite.app = app
}

func (suite *TableHandlerTestSuite) TestRegisterRoutes() {
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
		Path:   "/table",
	})

}

func (suite *TableHandlerTestSuite) TestHandleGetTableInfo() {
	assert := suite.Assert()
	t := suite.T()

	t.Run("should pass", func(t *testing.T) {
		passing := utils.RequestTesting[[]models.TableInfoResp]{
			ReqMethod: http.MethodGet,
			ReqUrl:    "/table/test/describe",
		}

		tableInfo, _ := passing.RunRequest(suite.app)
		assert.NotEmpty(tableInfo)
		assert.Equal(tableInfo[0].Name, "name")
	})

	t.Run("should fail bad request", func(t *testing.T) {
		failingBadRequest := utils.RequestTesting[models.ErrResp]{
			ReqMethod: http.MethodGet,
			ReqUrl:    "/table/12345/describe",
		}
		decoedResp, rawResp := failingBadRequest.RunRequest(suite.app)
		assert.Equal(http.StatusBadRequest, rawResp.StatusCode)
		assert.Len(decoedResp.Message, 1)
	})
}

func (suite *TableHandlerTestSuite) TestHandleListTables() {
	assert := suite.Assert()
	passing := utils.RequestTesting[models.ListTablesResp]{
		ReqMethod: http.MethodGet,
		ReqUrl:    "/table",
	}
	decoedRes, _ := passing.RunRequest(suite.app)

	tables := utils.SliceOfPointersToSliceOfValues(decoedRes.Tables)
	assert.NotEmpty(tables)
	assert.Contains(tables, "test")
}

func (suite *TableHandlerTestSuite) TestHandleCreateTable() {
	assert := suite.Assert()
	t := suite.T()

	t.Run("should pass", func(t *testing.T) {
		passingBody, _ := utils.EncodeBody([]models.CreateTablePayload{{ColName: "test1",
			Type:     "varchar(255)",
			Nullable: true,
			Default:  "kareem",
			Unique:   true,
		}})
		passing := utils.RequestTesting[models.CreateTableResp]{
			ReqMethod: http.MethodPost,
			ReqUrl:    "/table/test1",
			ReqBody:   passingBody,
		}
		decodedResp, rawResp := passing.RunRequest(suite.app)
		assert.Equal(http.StatusCreated, rawResp.StatusCode)
		assert.Equal(decodedResp.Created, "test1")

	})

	t.Run("should fail unprocessable entitiy", func(t *testing.T) {
		failingUnprocessableEntitiy := utils.RequestTesting[models.ErrResp]{
			ReqMethod: http.MethodPost,
			ReqUrl:    "/table/anything",
		}
		decodedResp, rawResp := failingUnprocessableEntitiy.RunRequest(suite.app)
		assert.Equal(http.StatusUnprocessableEntity, rawResp.StatusCode)
		assert.Contains(decodedResp.Message, "Unprocessable Entity")
	})

	t.Run("should fail bad request", func(t *testing.T) {
		failingBadRequest := utils.RequestTesting[models.ErrResp]{
			ReqMethod: http.MethodPost,
			ReqUrl:    "/table/1.1",
		}
		decodedResp, rawResp := failingBadRequest.RunRequest(suite.app)
		assert.Equal(http.StatusBadRequest, rawResp.StatusCode)
		assert.NotZero(decodedResp.Message)

		failingBadRequestBody, _ := utils.EncodeBody([]models.CreateTablePayload{{
			ColName:  "test2",
			Type:     "varchar(255)",
			Nullable: "should fail",
			Default:  nil,
			Unique:   nil,
		}})
		failingBadRequest = utils.RequestTesting[models.ErrResp]{
			ReqMethod: http.MethodPost,
			ReqUrl:    "/table/anything",
			ReqBody:   failingBadRequestBody,
		}
		decodedResp, rawResp = failingBadRequest.RunRequest(suite.app)
		assert.Equal(http.StatusBadRequest, rawResp.StatusCode)
		assert.NotZero(decodedResp.Message)

	})

}

func (suite *TableHandlerTestSuite) TestHandleAddColumn() {
	assert := suite.Assert()
	t := suite.T()

	t.Run("should pass", func(t *testing.T) {
		passingBody, _ := utils.EncodeBody(models.AddModifyColumnPayload{ColName: "test3", Type: "varchar(255)"})
		passing := utils.RequestTesting[models.SuccessResp]{
			ReqMethod: http.MethodPost,
			ReqUrl:    "/table/test/column/add",
			ReqBody:   passingBody,
		}
		decoedBody, _ := passing.RunRequest(suite.app)
		assert.True(decoedBody.Success)

	})

	t.Run("should fail unproccessable entity", func(t *testing.T) {
		failingUnprocessableEntitiy := utils.RequestTesting[models.ErrResp]{
			ReqMethod: http.MethodPost,
			ReqUrl:    "/table/test/column/add",
		}
		decodedResp, rawResp := failingUnprocessableEntitiy.RunRequest(suite.app)
		assert.Equal(http.StatusUnprocessableEntity, rawResp.StatusCode)
		assert.Contains(decodedResp.Message, "Unprocessable Entity")

		failingBadRequestParam := utils.RequestTesting[models.ErrResp]{
			ReqMethod: http.MethodPost,
			ReqUrl:    "/table/1.1/column/add",
		}
		decodedRes, rawRes := failingBadRequestParam.RunRequest(suite.app)
		assert.Equal(http.StatusBadRequest, rawRes.StatusCode)
		assert.Len(decodedRes.Message, 1)
	})

	t.Run("should fail bad request", func(t *testing.T) {
		failingBadRequestBody, _ := utils.EncodeBody(models.AddModifyColumnPayload{
			ColName: "",
			Type:    "varchar(255)",
		})
		failingBadRequest := utils.RequestTesting[models.ErrResp]{
			ReqMethod: http.MethodPost,
			ReqUrl:    "/table/test/column/add",
			ReqBody:   failingBadRequestBody,
		}
		decodedResp, rawResp := failingBadRequest.RunRequest(suite.app)
		assert.Equal(http.StatusBadRequest, rawResp.StatusCode)
		assert.Len(decodedResp.Message, 1)
	})

}

func (suite *TableHandlerTestSuite) TestHandleUpdateColumn() {
	assert := suite.Assert()
	t := suite.T()

	t.Run("should pass", func(t *testing.T) {
		passingBody, _ := utils.EncodeBody(models.AddModifyColumnPayload{ColName: "name", Type: "varchar(255)"})
		passing := utils.RequestTesting[models.SuccessResp]{
			ReqMethod: http.MethodPut,
			ReqUrl:    "/table/test/column/modify",
			ReqBody:   passingBody,
		}
		decodedRes, _ := passing.RunRequest(suite.app)

		assert.True(decodedRes.Success)
	})

	t.Run("should fail unproccessable entity", func(t *testing.T) {
		failingUnprocessableEntity := utils.RequestTesting[models.ErrResp]{
			ReqMethod: http.MethodPut,
			ReqUrl:    "/table/test/column/modify",
		}
		decodedRes, rawRes := failingUnprocessableEntity.RunRequest(suite.app)
		assert.Equal(http.StatusUnprocessableEntity, rawRes.StatusCode)
		assert.Contains(decodedRes.Message, "Unprocessable Entity")
	})
	t.Run("should fail bad request param", func(t *testing.T) {
		failingBadRequestParam := utils.RequestTesting[models.ErrResp]{
			ReqMethod: http.MethodPut,
			ReqUrl:    "/table/1.1/column/modify",
		}
		decodedRes, rawRes := failingBadRequestParam.RunRequest(suite.app)
		assert.Equal(http.StatusBadRequest, rawRes.StatusCode)
		assert.Len(decodedRes.Message, 1)
	})

	t.Run("should fail bad request body", func(t *testing.T) {
		failingBadRequestBody, _ := utils.EncodeBody(models.AddModifyColumnPayload{
			ColName: "",
			Type:    "varchar(255)",
		})
		failingBadRequest := utils.RequestTesting[models.ErrResp]{
			ReqMethod: http.MethodPut,
			ReqUrl:    "/table/test/column/modify",
			ReqBody:   failingBadRequestBody,
		}
		decodedRes, rawRes := failingBadRequest.RunRequest(suite.app)
		assert.Equal(http.StatusBadRequest, rawRes.StatusCode)
		assert.Len(decodedRes.Message, 1)
	})
}

func (suite *TableHandlerTestSuite) TestHandleDeleteColumn() {
	assert := suite.Assert()
	t := suite.T()

	t.Run("should pass", func(t *testing.T) {
		passingBody, _ := utils.EncodeBody(models.DeleteColumnPayload{ColName: "name"})
		passing := utils.RequestTesting[models.SuccessResp]{
			ReqMethod: http.MethodDelete,
			ReqUrl:    "/table/test/column/delete",
			ReqBody:   passingBody,
		}

		decoedRes, _ := passing.RunRequest(suite.app)

		assert.True(decoedRes.Success)
	})
	t.Run("should fail bad request param", func(t *testing.T) {
		failingBadRequestParam := utils.RequestTesting[models.ErrResp]{
			ReqMethod: http.MethodDelete,
			ReqUrl:    "/table/1.1/column/delete",
		}
		decodedRes, rawRes := failingBadRequestParam.RunRequest(suite.app)
		assert.Equal(http.StatusBadRequest, rawRes.StatusCode)
		assert.Len(decodedRes.Message, 1)
	})

	t.Run("should fail bad request", func(t *testing.T) {
		failingBadRequestBody, _ := utils.EncodeBody(models.DeleteColumnPayload{
			ColName: "",
		})
		failingBadRequest := utils.RequestTesting[models.ErrResp]{
			ReqMethod: http.MethodDelete,
			ReqUrl:    "/table/test/column/delete",
			ReqBody:   failingBadRequestBody,
		}
		decodedRes, rawRes := failingBadRequest.RunRequest(suite.app)
		assert.Equal(http.StatusBadRequest, rawRes.StatusCode)
		assert.Len(decodedRes.Message, 1)
	})

	t.Run("should fail unproccessable entity", func(t *testing.T) {
		failingUnprocessableEntity := utils.RequestTesting[models.ErrResp]{
			ReqMethod: http.MethodDelete,
			ReqUrl:    "/table/test/column/delete",
		}
		decodedRes, rawRes := failingUnprocessableEntity.RunRequest(suite.app)
		assert.Equal(http.StatusUnprocessableEntity, rawRes.StatusCode)
		assert.Contains(decodedRes.Message, "Unprocessable Entity")
	})

}

func (suite *TableHandlerTestSuite) TestHandleDeleteTable() {
	assert := suite.Assert()
	t := suite.T()

	t.Run("should pass", func(t *testing.T) {
		passing := utils.RequestTesting[models.SuccessResp]{
			ReqMethod: http.MethodDelete,
			ReqUrl:    "/table/test",
		}
		decodedRes, _ := passing.RunRequest(suite.app)
		assert.True(decodedRes.Success)
	})

	t.Run("should fail bad request param", func(t *testing.T) {
		failingBadRequestParams := utils.RequestTesting[models.ErrResp]{
			ReqMethod: http.MethodDelete,
			ReqUrl:    "/table/1.1",
		}
		decodedRes, rawRes := failingBadRequestParams.RunRequest(suite.app)
		assert.Equal(http.StatusBadRequest, rawRes.StatusCode)
		assert.Len(decodedRes.Message, 1)
	})
}

func TestTableHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(TableHandlerTestSuite))
}

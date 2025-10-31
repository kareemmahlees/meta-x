package handlers

import (
	"errors"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/kareemmahlees/meta-x/models"
	"github.com/kareemmahlees/meta-x/utils"
	"github.com/stretchr/testify/suite"
)

type TableHandlerTestSuite struct {
	suite.Suite
	api     humatest.TestAPI
	handler *TableHandler
}

func (suite *TableHandlerTestSuite) SetupSuite() {
	_, api := humatest.New(suite.T())
	storage := NewMockTableExecutor()

	handler := NewTableHandler(storage)
	handler.RegisterRoutes(api)

	suite.api = api
	suite.handler = handler
}

func (suite *TableHandlerTestSuite) TestHandleGetTableInfo() {
	assert := suite.Assert()
	t := suite.T()

	t.Run("should pass", func(t *testing.T) {
		rr := suite.api.Get("/table/test/describe")
		assert.Equal(rr.Code, http.StatusOK)
	})

	t.Run("should fail bad request", func(t *testing.T) {
		rr := suite.api.Get("/table/12345/describe")
		assert.Equal(rr.Code, http.StatusBadRequest)

		decodedResp := utils.DecodeBody[models.ErrResp](rr.Result().Body)
		assert.NotEmpty(decodedResp.Message)
	})

	t.Run("should fail internal server", func(t *testing.T) {
		_, testAPI := humatest.New(t)
		storage := NewFaultyTableExecutor()

		handler := NewTableHandler(storage)
		handler.RegisterRoutes(testAPI)

		rr := testAPI.Get("/table/test/describe")
		assert.Equal(rr.Code, http.StatusInternalServerError)
	})
}

func (suite *TableHandlerTestSuite) TestHandleListTables() {
	assert := suite.Assert()
	t := suite.T()

	t.Run("should pass", func(t *testing.T) {
		rr := suite.api.Get("/table")
		assert.Equal(http.StatusOK, rr.Code)
		assert.Contains(rr.Body.String(), "test")
	})

	t.Run("should fail internal server", func(t *testing.T) {
		_, testAPI := humatest.New(t)
		storage := NewFaultyTableExecutor()

		handler := NewTableHandler(storage)
		handler.RegisterRoutes(testAPI)
		rr := testAPI.Get("/table")

		assert.Equal(http.StatusInternalServerError, rr.Code)
	})
}

func (suite *TableHandlerTestSuite) TestHandleCreateTable() {
	assert := suite.Assert()
	t := suite.T()

	t.Run("should pass", func(t *testing.T) {
		rr := suite.api.Post("/table/test1", []models.CreateTablePayload{{ColName: "test1",
			Type:     "varchar(255)",
			Nullable: true,
			Default:  "kareem",
			Unique:   true,
		}})
		assert.Equal(http.StatusCreated, rr.Code)

		decodedRes := utils.DecodeBody[models.CreateTableResp](rr.Result().Body)
		assert.Equal(decodedRes.Created, "test1")
	})

	t.Run("should fail unprocessable entitiy", func(t *testing.T) {
		rr := suite.api.Post("/table/anything")
		assert.Equal(http.StatusUnprocessableEntity, rr.Code)
	})

	t.Run("should fail bad request params", func(t *testing.T) {
		rr := suite.api.Post("/table/1.1")
		assert.Equal(http.StatusBadRequest, rr.Code)

	})

	t.Run("should fail bad request body", func(t *testing.T) {
		rr := suite.api.Post("/table/anything", []models.CreateTablePayload{{
			ColName:  "test2",
			Type:     "varchar(255)",
			Nullable: "should fail",
			Default:  nil,
			Unique:   nil,
		}})
		assert.Equal(http.StatusBadRequest, rr.Code)
	})

	t.Run("should fail internal server", func(t *testing.T) {
		_, testAPI := humatest.New(t)
		storage := NewFaultyTableExecutor()

		handler := NewTableHandler(storage)
		handler.RegisterRoutes(testAPI)

		rr := testAPI.Post("/table/test1", []models.CreateTablePayload{{ColName: "test1",
			Type:     "varchar(255)",
			Nullable: true,
			Default:  "kareem",
			Unique:   true,
		}})

		assert.Equal(http.StatusInternalServerError, rr.Code)
	})

}

func (suite *TableHandlerTestSuite) TestHandleAddColumn() {
	assert := suite.Assert()
	t := suite.T()

	t.Run("should pass", func(t *testing.T) {
		rr := suite.api.Post("/table/test/column/add", models.AddModifyColumnPayload{ColName: "test3", Type: "varchar(255)"})
		decodedBody := utils.DecodeBody[models.SuccessResp](rr.Result().Body)
		assert.True(decodedBody.Success)
	})

	t.Run("should fail unproccessable entity", func(t *testing.T) {
		rr := suite.api.Post("/table/test/column/add")
		assert.Equal(http.StatusUnprocessableEntity, rr.Code)
	})

	t.Run("should fail bad request param", func(t *testing.T) {
		rr := suite.api.Post("/table/1.1/column/add")
		assert.Equal(http.StatusBadRequest, rr.Code)
	})

	t.Run("should fail bad request body", func(t *testing.T) {
		rr := suite.api.Post("/table/test/column/add", models.AddModifyColumnPayload{
			ColName: "",
			Type:    "varchar(255)",
		})
		assert.Equal(http.StatusBadRequest, rr.Code)
	})

	t.Run("should fail internal server", func(t *testing.T) {
		_, testAPI := humatest.New(t)
		storage := NewFaultyTableExecutor()
		handler := NewTableHandler(storage)
		handler.RegisterRoutes(testAPI)

		rr := testAPI.Post("/table/test/column/add", models.AddModifyColumnPayload{ColName: "test3", Type: "varchar(255)"})

		assert.Equal(http.StatusInternalServerError, rr.Code)
	})
}

func (suite *TableHandlerTestSuite) TestHandleUpdateColumn() {
	assert := suite.Assert()
	t := suite.T()

	t.Run("should pass", func(t *testing.T) {
		rr := suite.api.Put("/table/test/column/modify", models.AddModifyColumnPayload{ColName: "name", Type: "varchar(255)"})
		decodedRes := utils.DecodeBody[models.SuccessResp](rr.Result().Body)

		assert.Equal(rr.Code, http.StatusOK)
		assert.True(decodedRes.Success)
	})

	t.Run("should fail unproccessable entity", func(t *testing.T) {
		rr := suite.api.Put("/table/test/column/modify")
		assert.Equal(http.StatusUnprocessableEntity, rr.Code)

		decodedRes := utils.DecodeBody[models.ErrResp](rr.Result().Body)
		assert.NotEmpty(decodedRes.Message)
	})
	t.Run("should fail bad request param", func(t *testing.T) {
		rr := suite.api.Put("/table/1.1/column/modify")
		assert.Equal(http.StatusBadRequest, rr.Code)

		decodedRes := utils.DecodeBody[models.ErrResp](rr.Result().Body)
		assert.NotEmpty(decodedRes.Message)
	})

	t.Run("should fail bad request body", func(t *testing.T) {
		rr := suite.api.Put("/table/test/column/modify", models.AddModifyColumnPayload{
			ColName: "",
			Type:    "varchar(255)",
		})
		assert.Equal(http.StatusBadRequest, rr.Code)

		decodedRes := utils.DecodeBody[models.ErrResp](rr.Result().Body)
		assert.NotEmpty(decodedRes.Message)
	})

	t.Run("should fail internal server", func(t *testing.T) {
		_, testAPI := humatest.New(t)
		storage := NewFaultyTableExecutor()
		handler := NewTableHandler(storage)
		handler.RegisterRoutes(testAPI)

		rr := testAPI.Put("/table/test/column/modify", models.AddModifyColumnPayload{ColName: "name", Type: "varchar(255)"})
		assert.Equal(http.StatusInternalServerError, rr.Code)
	})
}

func (suite *TableHandlerTestSuite) TestHandleDeleteColumn() {
	assert := suite.Assert()
	t := suite.T()

	t.Run("should pass", func(t *testing.T) {
		rr := suite.api.Delete("/table/test/column/delete", models.DeleteColumnPayload{ColName: "name"})

		decodedRes := utils.DecodeBody[models.SuccessResp](rr.Result().Body)
		assert.True(decodedRes.Success)
	})
	t.Run("should fail bad request param", func(t *testing.T) {
		rr := suite.api.Delete("/table/1.1/column/delete")
		decodedRes := utils.DecodeBody[models.ErrResp](rr.Result().Body)

		assert.Equal(http.StatusBadRequest, rr.Code)
		assert.NotEmpty(decodedRes.Message)
	})

	t.Run("should fail bad request", func(t *testing.T) {
		rr := suite.api.Delete("/table/test/column/delete", models.DeleteColumnPayload{
			ColName: "",
		})
		decodedRes := utils.DecodeBody[models.ErrResp](rr.Result().Body)

		assert.Equal(http.StatusBadRequest, rr.Code)
		assert.NotEmpty(decodedRes.Message)
	})

	t.Run("should fail unproccessable entity", func(t *testing.T) {
		rr := suite.api.Delete("/table/test/column/delete")
		decodedRes := utils.DecodeBody[models.ErrResp](rr.Result().Body)

		assert.Equal(http.StatusUnprocessableEntity, rr.Code)
		assert.NotEmpty(decodedRes.Message)
	})

	t.Run("should fail internal server", func(t *testing.T) {
		_, testAPI := humatest.New(t)
		storage := NewFaultyTableExecutor()
		handler := NewTableHandler(storage)
		handler.RegisterRoutes(testAPI)

		rr := testAPI.Delete("/table/test/column/delete", models.DeleteColumnPayload{ColName: "name"})

		assert.Equal(http.StatusInternalServerError, rr.Code)
	})
}

func (suite *TableHandlerTestSuite) TestHandleDeleteTable() {
	assert := suite.Assert()
	t := suite.T()

	t.Run("should pass", func(t *testing.T) {
		rr := suite.api.Delete("/table/test")
		decodedRes := utils.DecodeBody[models.SuccessResp](rr.Result().Body)

		assert.True(decodedRes.Success)
	})

	t.Run("should fail bad request param", func(t *testing.T) {
		rr := suite.api.Delete("/table/1.1")
		decodedRes := utils.DecodeBody[models.ErrResp](rr.Result().Body)

		assert.Equal(http.StatusBadRequest, rr.Code)
		assert.NotEmpty(decodedRes.Message)
	})

	t.Run("should fail internal server", func(t *testing.T) {
		_, testAPI := humatest.New(t)
		storage := NewFaultyTableExecutor()
		handler := NewTableHandler(storage)
		handler.RegisterRoutes(testAPI)

		rr := testAPI.Delete("/table/test")

		assert.Equal(http.StatusInternalServerError, rr.Code)
	})
}

func TestTableHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(TableHandlerTestSuite))
}

type MockTableExecutor struct{}

func NewMockTableExecutor() *MockTableExecutor {
	return &MockTableExecutor{}
}

func (ms *MockTableExecutor) GetTable(tableName string) ([]*models.TableColumnInfo, error) {
	tableInfo := models.TableColumnInfo{
		Name:     "name",
		Type:     "varchar(255)",
		Nullable: "Yes",
	}
	return []*models.TableColumnInfo{&tableInfo}, nil
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

type FaultyTableExecutor struct{}

func NewFaultyTableExecutor() *FaultyTableExecutor {
	return &FaultyTableExecutor{}
}

var err = errors.New("error")

func (ms *FaultyTableExecutor) GetTable(tableName string) ([]*models.TableColumnInfo, error) {
	return nil, err
}
func (ms *FaultyTableExecutor) ListTables() ([]*string, error) {
	return nil, err
}
func (ms *FaultyTableExecutor) CreateTable(tableName string, data []models.CreateTablePayload) error {
	return err
}
func (ms *FaultyTableExecutor) DeleteTable(tableName string) error {
	return err
}
func (ms *FaultyTableExecutor) AddColumn(tableName string, data models.AddModifyColumnPayload) error {
	return err
}
func (ms *FaultyTableExecutor) UpdateColumn(tableName string, data models.AddModifyColumnPayload) error {
	return err
}
func (ms *FaultyTableExecutor) DeleteColumn(tableName string, data models.DeleteColumnPayload) error {
	return err
}

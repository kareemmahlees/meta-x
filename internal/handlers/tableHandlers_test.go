package handlers

import (
	"errors"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/kareemmahlees/meta-x/models"
	"github.com/kareemmahlees/meta-x/utils"
	"github.com/stretchr/testify/suite"
)

type TableHandlerTestSuite struct {
	suite.Suite
	r       *chi.Mux
	handler *TableHandler
}

func (suite *TableHandlerTestSuite) SetupSuite() {
	r := chi.NewRouter()
	storage := NewMockTableExecutor()

	handler := NewTableHandler(storage)
	handler.RegisterRoutes(r)

	suite.r = r
	suite.handler = handler
}

func (suite *TableHandlerTestSuite) TestRegisterRoutes() {
	assert := suite.Assert()

	var routes []string
	for _, route := range suite.r.Routes() {
		routes = append(routes, route.Pattern)
	}

	assert.Contains(routes, "/table/*")
}

func (suite *TableHandlerTestSuite) TestHandleGetTableInfo() {
	assert := suite.Assert()
	t := suite.T()

	t.Run("should pass", func(t *testing.T) {
		rr := utils.TestRequest(suite.r, http.MethodGet, "/table/test/describe", http.NoBody)
		assert.Equal(rr.Code, http.StatusOK)
	})

	t.Run("should fail bad request", func(t *testing.T) {
		rr := utils.TestRequest(suite.r, http.MethodGet, "/table/12345/describe", http.NoBody)
		assert.Equal(rr.Code, http.StatusBadRequest)

		decodedResp := utils.DecodeBody[models.ErrResp](rr.Result().Body)
		assert.NotEmpty(decodedResp.Message)
	})

	t.Run("should fail internal server", func(t *testing.T) {

		r := chi.NewRouter()
		storage := NewFaultyTableExecutor()

		handler := NewTableHandler(storage)
		handler.RegisterRoutes(r)

		rr := utils.TestRequest(r, http.MethodGet, "/table/test/describe", http.NoBody)
		assert.Equal(rr.Code, http.StatusInternalServerError)
	})
}

func (suite *TableHandlerTestSuite) TestHandleListTables() {
	assert := suite.Assert()
	t := suite.T()

	t.Run("should pass", func(t *testing.T) {
		assert.HTTPSuccess(suite.handler.handleListTables, http.MethodGet, "/table", nil)
		assert.HTTPBodyContains(suite.handler.handleListTables, http.MethodGet, "/table", nil, "test")
	})

	t.Run("should fail internal server", func(t *testing.T) {
		r := chi.NewRouter()
		storage := NewFaultyTableExecutor()

		handler := NewTableHandler(storage)
		handler.RegisterRoutes(r)

		rr := utils.TestRequest(r, http.MethodGet, "/table", http.NoBody)

		assert.Equal(http.StatusInternalServerError, rr.Code)
	})
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
		rr := utils.TestRequest(suite.r, http.MethodPost, "/table/test1", passingBody)
		assert.Equal(http.StatusCreated, rr.Code)

		decodedRes := utils.DecodeBody[models.CreateTableResp](rr.Result().Body)
		assert.Equal(decodedRes.Created, "test1")
	})

	t.Run("should fail unprocessable entitiy", func(t *testing.T) {
		rr := utils.TestRequest(suite.r, http.MethodPost, "/table/anything", http.NoBody)
		assert.Equal(http.StatusUnprocessableEntity, rr.Code)
	})

	t.Run("should fail bad request params", func(t *testing.T) {
		rr := utils.TestRequest(suite.r, http.MethodPost, "/table/1.1", http.NoBody)
		assert.Equal(http.StatusBadRequest, rr.Code)

	})

	t.Run("should fail bad request body", func(t *testing.T) {
		failingBadRequestBody, _ := utils.EncodeBody([]models.CreateTablePayload{{
			ColName:  "test2",
			Type:     "varchar(255)",
			Nullable: "should fail",
			Default:  nil,
			Unique:   nil,
		}})
		rr := utils.TestRequest(suite.r, http.MethodPost, "/table/anything", failingBadRequestBody)
		assert.Equal(http.StatusBadRequest, rr.Code)
	})

	t.Run("should fail internal server", func(t *testing.T) {
		failingInternalServerBody, _ := utils.EncodeBody([]models.CreateTablePayload{{ColName: "test1",
			Type:     "varchar(255)",
			Nullable: true,
			Default:  "kareem",
			Unique:   true,
		}})
		r := chi.NewRouter()
		storage := NewFaultyTableExecutor()

		handler := NewTableHandler(storage)
		handler.RegisterRoutes(r)

		rr := utils.TestRequest(r, http.MethodPost, "/table/test1", failingInternalServerBody)

		assert.Equal(http.StatusInternalServerError, rr.Code)
	})

}

func (suite *TableHandlerTestSuite) TestHandleAddColumn() {
	assert := suite.Assert()
	t := suite.T()

	t.Run("should pass", func(t *testing.T) {
		passingBody, _ := utils.EncodeBody(models.AddModifyColumnPayload{ColName: "test3", Type: "varchar(255)"})
		rr := utils.TestRequest(suite.r, http.MethodPost, "/table/test/column/add", passingBody)
		decoedBody := utils.DecodeBody[models.SuccessResp](rr.Result().Body)
		assert.True(decoedBody.Success)
	})

	t.Run("should fail unproccessable entity", func(t *testing.T) {
		rr := utils.TestRequest(suite.r, http.MethodPost, "/table/test/column/add", http.NoBody)
		assert.Equal(http.StatusUnprocessableEntity, rr.Code)
	})

	t.Run("should fail bad request param", func(t *testing.T) {
		rr := utils.TestRequest(suite.r, http.MethodPost, "/table/1.1/column/add", http.NoBody)
		assert.Equal(http.StatusBadRequest, rr.Code)
	})

	t.Run("should fail bad request body", func(t *testing.T) {
		failingBadRequestBody, _ := utils.EncodeBody(models.AddModifyColumnPayload{
			ColName: "",
			Type:    "varchar(255)",
		})
		rr := utils.TestRequest(suite.r, http.MethodPost, "/table/test/column/add", failingBadRequestBody)
		assert.Equal(http.StatusBadRequest, rr.Code)
	})

	t.Run("should fail internal server", func(t *testing.T) {
		failingInternalServerBody, _ := utils.EncodeBody(models.AddModifyColumnPayload{ColName: "test3", Type: "varchar(255)"})

		r := chi.NewRouter()
		storage := NewFaultyTableExecutor()
		handler := NewTableHandler(storage)
		handler.RegisterRoutes(r)

		rr := utils.TestRequest(r, http.MethodPost, "/table/test/column/add", failingInternalServerBody)

		assert.Equal(http.StatusInternalServerError, rr.Code)
	})
}

func (suite *TableHandlerTestSuite) TestHandleUpdateColumn() {
	assert := suite.Assert()
	t := suite.T()

	t.Run("should pass", func(t *testing.T) {
		passingBody, _ := utils.EncodeBody(models.AddModifyColumnPayload{ColName: "name", Type: "varchar(255)"})
		rr := utils.TestRequest(suite.r, http.MethodPut, "/table/test/column/modify", passingBody)
		decodedRes := utils.DecodeBody[models.SuccessResp](rr.Result().Body)

		assert.Equal(rr.Code, http.StatusOK)
		assert.True(decodedRes.Success)
	})

	t.Run("should fail unproccessable entity", func(t *testing.T) {
		rr := utils.TestRequest(suite.r, http.MethodPut, "/table/test/column/modify", http.NoBody)
		assert.Equal(http.StatusUnprocessableEntity, rr.Code)

		decodedRes := utils.DecodeBody[models.ErrResp](rr.Result().Body)
		assert.NotEmpty(decodedRes.Message)
	})
	t.Run("should fail bad request param", func(t *testing.T) {
		rr := utils.TestRequest(suite.r, http.MethodPut, "/table/1.1/column/modify", http.NoBody)
		assert.Equal(http.StatusBadRequest, rr.Code)

		decodedRes := utils.DecodeBody[models.ErrResp](rr.Result().Body)
		assert.NotEmpty(decodedRes.Message)
	})

	t.Run("should fail bad request body", func(t *testing.T) {
		failingBadRequestBody, _ := utils.EncodeBody(models.AddModifyColumnPayload{
			ColName: "",
			Type:    "varchar(255)",
		})
		rr := utils.TestRequest(suite.r, http.MethodPut, "/table/test/column/modify", failingBadRequestBody)
		assert.Equal(http.StatusBadRequest, rr.Code)

		decodedRes := utils.DecodeBody[models.ErrResp](rr.Result().Body)
		assert.NotEmpty(decodedRes.Message)
	})

	t.Run("should fail internal server", func(t *testing.T) {
		failingInternalServerBody, _ := utils.EncodeBody(models.AddModifyColumnPayload{ColName: "name", Type: "varchar(255)"})

		r := chi.NewRouter()
		storage := NewFaultyTableExecutor()
		handler := NewTableHandler(storage)
		handler.RegisterRoutes(r)

		rr := utils.TestRequest(r, http.MethodPut, "/table/test/column/modify", failingInternalServerBody)

		assert.Equal(http.StatusInternalServerError, rr.Code)
	})
}

func (suite *TableHandlerTestSuite) TestHandleDeleteColumn() {
	assert := suite.Assert()
	t := suite.T()

	t.Run("should pass", func(t *testing.T) {
		passingBody, _ := utils.EncodeBody(models.DeleteColumnPayload{ColName: "name"})
		rr := utils.TestRequest(suite.r, http.MethodDelete, "/table/test/column/delete", passingBody)

		decodedRes := utils.DecodeBody[models.SuccessResp](rr.Result().Body)
		assert.True(decodedRes.Success)
	})
	t.Run("should fail bad request param", func(t *testing.T) {
		rr := utils.TestRequest(suite.r, http.MethodDelete, "/table/1.1/column/delete", http.NoBody)
		decodedRes := utils.DecodeBody[models.ErrResp](rr.Result().Body)

		assert.Equal(http.StatusBadRequest, rr.Code)
		assert.NotEmpty(decodedRes.Message)
	})

	t.Run("should fail bad request", func(t *testing.T) {
		failingBadRequestBody, _ := utils.EncodeBody(models.DeleteColumnPayload{
			ColName: "",
		})
		rr := utils.TestRequest(suite.r, http.MethodDelete, "/table/test/column/delete", failingBadRequestBody)
		decodedRes := utils.DecodeBody[models.ErrResp](rr.Result().Body)

		assert.Equal(http.StatusBadRequest, rr.Code)
		assert.NotEmpty(decodedRes.Message)
	})

	t.Run("should fail unproccessable entity", func(t *testing.T) {
		rr := utils.TestRequest(suite.r, http.MethodDelete, "/table/test/column/delete", http.NoBody)
		decodedRes := utils.DecodeBody[models.ErrResp](rr.Result().Body)

		assert.Equal(http.StatusUnprocessableEntity, rr.Code)
		assert.NotEmpty(decodedRes.Message)
	})

	t.Run("should fail internal server", func(t *testing.T) {
		failingInternalServerBody, _ := utils.EncodeBody(models.DeleteColumnPayload{ColName: "name"})

		r := chi.NewRouter()
		storage := NewFaultyTableExecutor()
		handler := NewTableHandler(storage)
		handler.RegisterRoutes(r)

		rr := utils.TestRequest(r, http.MethodDelete, "/table/test/column/delete", failingInternalServerBody)

		assert.Equal(http.StatusInternalServerError, rr.Code)
	})
}

func (suite *TableHandlerTestSuite) TestHandleDeleteTable() {
	assert := suite.Assert()
	t := suite.T()

	t.Run("should pass", func(t *testing.T) {
		rr := utils.TestRequest(suite.r, http.MethodDelete, "/table/test", http.NoBody)
		decodedRes := utils.DecodeBody[models.SuccessResp](rr.Result().Body)

		assert.True(decodedRes.Success)
	})

	t.Run("should fail bad request param", func(t *testing.T) {
		rr := utils.TestRequest(suite.r, http.MethodDelete, "/table/1.1", http.NoBody)
		decodedRes := utils.DecodeBody[models.ErrResp](rr.Result().Body)

		assert.Equal(http.StatusBadRequest, rr.Code)
		assert.NotEmpty(decodedRes.Message)
	})

	t.Run("should fail internal server", func(t *testing.T) {
		r := chi.NewRouter()
		storage := NewFaultyTableExecutor()
		handler := NewTableHandler(storage)
		handler.RegisterRoutes(r)

		rr := utils.TestRequest(r, http.MethodDelete, "/table/test", http.NoBody)

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

type FaultyTableExecutor struct{}

func NewFaultyTableExecutor() *FaultyTableExecutor {
	return &FaultyTableExecutor{}
}

var err = errors.New("error")

func (ms *FaultyTableExecutor) GetTable(tableName string) ([]*models.TableInfoResp, error) {
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

package handlers

import (
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/stretchr/testify/suite"
)

type DefaultHandlerTestSuite struct {
	suite.Suite
	api     humatest.TestAPI
	handler *DefaultHandler
}

func (suite *DefaultHandlerTestSuite) SetupSuite() {
	_, api := humatest.New(suite.T())
	handler := NewDefaultHandler()
	handler.RegisterRoutes(api)

	suite.api = api
	suite.handler = handler
}

func (suite *DefaultHandlerTestSuite) TestAPIInfo() {
	assert := suite.Assert()
	resp := suite.api.Get("/")

	assert.Equal(http.StatusOK, resp.Code)
}

func TestDefaultHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(DefaultHandlerTestSuite))
}

package handlers

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/suite"
)

type DefaultHandlerTestSuite struct {
	suite.Suite
	r       *chi.Mux
	handler *DefaultHandler
}

func (suite *DefaultHandlerTestSuite) SetupSuite() {
	r := chi.NewRouter()
	handler := NewDefaultHandler()
	handler.RegisterRoutes(r)

	suite.r = r
	suite.handler = handler
}

func (suite *DefaultHandlerTestSuite) TestRegisterDefaultRoutes() {
	assert := suite.Assert()

	var routes []string
	for _, route := range suite.r.Routes() {
		routes = append(routes, route.Pattern)
	}

	assert.Contains(routes, "/health")
	assert.Contains(routes, "/")
}

func (suite *DefaultHandlerTestSuite) TestHealthCheck() {
	assert := suite.Assert()

	assert.HTTPSuccess(suite.handler.healthCheck, http.MethodGet, "/health", nil)
	assert.HTTPBodyContains(suite.handler.healthCheck, http.MethodGet, "/health", nil, "date")
}

func (suite *DefaultHandlerTestSuite) TestAPIInfo() {
	assert := suite.Assert()

	assert.HTTPSuccess(suite.handler.apiInfo, http.MethodGet, "/", nil)
}

func TestDefaultHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(DefaultHandlerTestSuite))
}

package handlers

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/stretchr/testify/suite"
)

type DefaultHandlerTestSuite struct {
	suite.Suite
	api     huma.API
	handler *DefaultHandler
}

// func (suite *DefaultHandlerTestSuite) SetupSuite() {
// 	r := chi.NewRouter()
// 	handler := NewDefaultHandler()
// 	handler.RegisterRoutes(api)

// 	suite.r = r
// 	suite.handler = handler
// }

// func (suite *DefaultHandlerTestSuite) TestRegisterDefaultRoutes() {
// 	assert := suite.Assert()

// 	var routes []string
// 	for _, route := range suite.r.Routes() {
// 		routes = append(routes, route.Pattern)
// 	}

// 	assert.Contains(routes, "/")
// }

// func (suite *DefaultHandlerTestSuite) TestAPIInfo() {
// 	assert := suite.Assert()

// 	assert.HTTPSuccess(suite.handler.apiInfo, http.MethodGet, "/", nil)
// }

// func TestDefaultHandlerTestSuite(t *testing.T) {
// 	suite.Run(t, new(DefaultHandlerTestSuite))
// }

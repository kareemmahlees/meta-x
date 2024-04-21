package lib

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
)

type ErrorTestSuite struct {
	suite.Suite
	fiberCtx *fiber.Ctx
}

func (suite *ErrorTestSuite) SetupSuite() {
	app := fiber.New()
	suite.fiberCtx = app.AcquireCtx(&fasthttp.RequestCtx{})
}

func (suite *ErrorTestSuite) TestBadRequestErr() {
	t := suite.T()

	err := BadRequestErr(suite.fiberCtx, "anything")
	assert.Nil(t, err)

	err = BadRequestErr(suite.fiberCtx, make(chan any))
	assert.NotNil(t, err)
}

func (suite *ErrorTestSuite) TestUnprocessableEntityErr() {
	t := suite.T()

	err := UnprocessableEntityErr(suite.fiberCtx, "anything")
	assert.Nil(t, err)

	err = UnprocessableEntityErr(suite.fiberCtx, make(chan any))
	assert.NotNil(t, err)
}

func (suite *ErrorTestSuite) TestForbiddenErr() {
	t := suite.T()

	err := ForbiddenErr(suite.fiberCtx, "anything")
	assert.Nil(t, err)

	err = ForbiddenErr(suite.fiberCtx, make(chan any))
	assert.NotNil(t, err)
}
func (suite *ErrorTestSuite) TestInternalServerErr() {
	t := suite.T()

	err := InternalServerErr(suite.fiberCtx, "anything")
	assert.Nil(t, err)

	err = InternalServerErr(suite.fiberCtx, make(chan any))
	assert.NotNil(t, err)
}

func TestErrorsTestSuite(t *testing.T) {
	suite.Run(t, new(ErrorTestSuite))
}

package handlers

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kareemmahlees/meta-x/utils"
	"github.com/stretchr/testify/assert"
)

func TestRegisterDefaultRoutes(t *testing.T) {

	app := fiber.New()
	handler := NewDefaultHandler(nil)

	handler.RegisterRoutes(app)

	var routes []utils.FiberRoute
	for _, route := range app.GetRoutes() {
		routes = append(routes, utils.FiberRoute{
			Method: route.Method,
			Path:   route.Path,
		})
	}

	assert.Contains(t, routes, utils.FiberRoute{
		Method: "GET",
		Path:   "/health",
	})

}

func TestHealthCheck(t *testing.T) {
	app := fiber.New()
	handler := NewDefaultHandler(nil)
	handler.RegisterRoutes(app)

	req := httptest.NewRequest("GET", "http://localhost:4000/health", nil)

	resp, _ := app.Test(req)
	payload := utils.DecodeBody[map[string]any](resp.Body)

	assert.Equal(t, resp.StatusCode, fiber.StatusOK)

	_, ok := payload["date"]
	assert.True(t, ok)
}

func TestBaseUrl(t *testing.T) {
	app := fiber.New()
	handler := NewDefaultHandler(nil)
	handler.RegisterRoutes(app)

	req := httptest.NewRequest("GET", "http://localhost:4000", nil)

	resp, err := app.Test(req)
	assert.Nil(t, err)
	_ = utils.DecodeBody[map[string]any](resp.Body)

	assert.Equal(t, resp.StatusCode, fiber.StatusOK)
}

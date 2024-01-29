package routes

import (
	"net/http/httptest"
	"testing"

	"meta-x/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestRegisterDefaultRoutes(t *testing.T) {

	app := fiber.New()

	RegisterDefaultRoutes(app)

	var routes []struct {
		method string
		params []string
		path   string
	}

	for _, route := range app.GetRoutes() {
		routes = append(routes, struct {
			method string
			params []string
			path   string
		}{
			method: route.Method,
			params: route.Params,
			path:   route.Path,
		})
	}

	assert.Contains(t, routes, struct {
		method string
		params []string
		path   string
	}{
		method: "GET",
		params: []string(nil),
		path:   "/health",
	})

}

func TestHealthCheck(t *testing.T) {
	app := fiber.New()
	RegisterDefaultRoutes(app)

	req := httptest.NewRequest("GET", "http://localhost:4000/health", nil)

	resp, _ := app.Test(req)
	payload := utils.ReadBody[map[string]any](resp.Body)

	assert.Equal(t, resp.StatusCode, fiber.StatusOK)

	_, ok := payload["date"]
	assert.True(t, ok)
}

func TestBaseUrl(t *testing.T) {
	app := fiber.New()
	RegisterDefaultRoutes(app)

	req := httptest.NewRequest("GET", "http://localhost:4000", nil)

	resp, err := app.Test(req)
	assert.Nil(t, err)
	_ = utils.ReadBody[map[string]any](resp.Body)

	assert.Equal(t, resp.StatusCode, fiber.StatusOK)
}

package routes

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kareemmahlees/mysql-meta/utils"
	"github.com/stretchr/testify/assert"
)

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

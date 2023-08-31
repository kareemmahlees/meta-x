package routes

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kareemmahlees/mysql-meta/utils"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	app := fiber.New()
	RegisterDefaultRoutes(app, nil)

	req := httptest.NewRequest("GET", "http://localhost:4000/health", nil)

	resp, _ := app.Test(req)
	payload := utils.ReadBody(resp.Body)

	assert.Equal(t, resp.StatusCode, fiber.StatusOK)

	_, ok := payload["date"]
	assert.True(t, ok)
}

func TestBaseUrl(t *testing.T) {
	app := fiber.New()
	RegisterDefaultRoutes(app, nil)

	req := httptest.NewRequest("GET", "http://localhost:4000", nil)

	resp, _ := app.Test(req)
	body, _ := io.ReadAll(resp.Body)

	assert.Greater(t, len(string(body)), 0)
	assert.Equal(t, resp.StatusCode, fiber.StatusOK)
}
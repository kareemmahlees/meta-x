package default_routes_test

import (
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	default_routes "github.com/kareemmahlees/mysql-meta/pkg/routes/default"
	"github.com/kareemmahlees/mysql-meta/utils"
	"github.com/stretchr/testify/assert"
)

var app *fiber.App

func TestMain(m *testing.M) {
	app = fiber.New()
	default_routes.RegisterDefaultRoutes(app)

	// Run tests
    exitVal := m.Run()
    
    // Exit with exit value from tests
    os.Exit(exitVal)
}

func TestHealthCheck(t *testing.T){
	req := httptest.NewRequest("GET", "http://localhost:4000/health", nil)

	resp, _ := app.Test(req)
	payload := utils.ReadBody(resp.Body)

	assert.Equal(t,resp.StatusCode,fiber.StatusOK)

	_,ok:= payload["date"]
	assert.True(t,ok)
}

func TestBaseUrl(t *testing.T){
	req := httptest.NewRequest("GET", "http://localhost:4000", nil)

	resp, _ := app.Test(req)
	body,_ := io.ReadAll(resp.Body)

	assert.Greater(t,len(string(body)),0)
	assert.Equal(t,resp.StatusCode,fiber.StatusOK)
}
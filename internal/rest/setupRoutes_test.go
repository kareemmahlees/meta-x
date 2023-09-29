package routes

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestSetup(t *testing.T) {
	app := fiber.New()
	Setup(app, nil)

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

	testData := []struct {
		method string
		params []string
		path   string
	}{
		{method: "GET", params: []string(nil), path: "/health"},
		{method: "GET", params: []string(nil), path: "/databases"},
		{method: "GET", params: []string(nil), path: "/tables"},
	}
	for _, test := range testData {
		assert.Contains(t, routes, test)
	}
}

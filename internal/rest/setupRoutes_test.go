package routes_test

// import (
// 	routes "github.com/kareemmahlees/meta-x/internal/rest"
// 	"github.com/kareemmahlees/meta-x/utils"
// 	"testing"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/stretchr/testify/assert"
// )

// func TestSetup(t *testing.T) {
// 	app := fiber.New()
// 	routes.Setup(app, nil)

// 	var routes []utils.FiberRoute

// 	for _, route := range app.GetRoutes() {
// 		routes = append(routes, utils.FiberRoute{
// 			Method: route.Method,
// 			Path:   route.Path,
// 		})
// 	}

// 	testData := []utils.FiberRoute{
// 		{Method: "GET", Path: "/health"},
// 		{Method: "GET", Path: "/database"},
// 		{Method: "GET", Path: "/table"},
// 	}
// 	for _, test := range testData {
// 		assert.Contains(t, routes, test)
// 	}
// }

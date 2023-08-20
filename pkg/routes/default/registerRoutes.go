package default_routes

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterDefaultRoutes(app *fiber.App){
	app.Get("/health",healthCheck)
	app.Get("/",apiInfo)
}
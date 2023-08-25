package default_routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func RegisterDefaultRoutes(app *fiber.App,db *sqlx.DB){
	app.Get("/health",healthCheck)
	app.Get("/",apiInfo)
}
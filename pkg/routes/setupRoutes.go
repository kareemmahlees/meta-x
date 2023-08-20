package routes

import (
	"github.com/gofiber/fiber/v2"
	default_routes "github.com/kareemmahlees/mysql-meta/pkg/routes/default"
)

func Setup(app *fiber.App){
	default_routes.RegisterDefaultRoutes(app)
}
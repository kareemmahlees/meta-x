package routes

import (
	"github.com/gofiber/fiber/v2"
	routes "github.com/kareemmahlees/mysql-meta/pkg/routes/default"
)

func Setup(app *fiber.App){
	routes.RegisterDefaultRoutes(app)
}
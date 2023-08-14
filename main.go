package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kareemmahlees/mysql-meta/pkg/routes"
)
var app *fiber.App

func main() {
	app = fiber.New()

	routes.Setup(app)

	app.Listen(":4000")
}
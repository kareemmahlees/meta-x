package handlers

import "github.com/gofiber/fiber/v2"

type Handler interface {
	RegisterRoutes(app *fiber.App)
}

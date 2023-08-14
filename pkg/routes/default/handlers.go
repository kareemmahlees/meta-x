package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func healthCheck(c *fiber.Ctx)error {
	return c.JSON(fiber.Map{"date":time.Now()})
}

func apiInfo(c *fiber.Ctx)error {
	return c.JSON(fiber.Map{
		"author":"Kareem Ebrahim",
		"year":2023,
		"contact":"kareemmahlees@gmail.com",
		"repo":"https://github.com/kareemmahlees/mysql-meta",
	})
}
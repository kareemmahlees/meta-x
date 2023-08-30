package lib

import "github.com/gofiber/fiber/v2"

func ResponseError500(errMsg string) fiber.Map {
	return fiber.Map{"status": 500, "error": errMsg}
}

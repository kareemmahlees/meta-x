package lib

import "github.com/gofiber/fiber/v2"

func ResponseError400(errMsg any) fiber.Map {
	return fiber.Map{"status": 400, "error": errMsg}
}

func ResponseError500(errMsg any) fiber.Map {
	return fiber.Map{"status": 500, "error": errMsg}
}

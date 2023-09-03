package lib

import "github.com/gofiber/fiber/v2"

// this utility is typically used for validation erros
// because fiber.Error doesn't accept list of errors
func ResponseError400(errMsg any) fiber.Map {
	return fiber.Map{"status": 400, "error": errMsg}
}

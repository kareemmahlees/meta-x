package lib

import "github.com/gofiber/fiber/v2"

func BadRequestErr(c *fiber.Ctx, errMsg any) error {
	return c.Status(fiber.StatusBadRequest).JSON(errMsg)
}

func UnprocessableEntityErr(c *fiber.Ctx, errMsg any) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(errMsg)
}

func ForbiddenErr(c *fiber.Ctx, errMsg any) error {
	return c.Status(fiber.StatusForbidden).JSON(errMsg)
}

func InternalServerErr(c *fiber.Ctx, errMsg any) error {
	return c.Status(fiber.StatusInternalServerError).JSON(errMsg)
}

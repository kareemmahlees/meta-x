package lib

import "github.com/gofiber/fiber/v2"

func BadRequestErr(c *fiber.Ctx, errMsg any) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"code":    fiber.StatusBadRequest,
		"message": errMsg,
	})
}

func UnprocessableEntityErr(c *fiber.Ctx, errMsg any) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
		"code":    fiber.StatusUnprocessableEntity,
		"message": errMsg,
	})
}

func ForbiddenErr(c *fiber.Ctx, errMsg any) error {
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		"code":    fiber.StatusForbidden,
		"message": errMsg,
	})
}

func InternalServerErr(c *fiber.Ctx, errMsg any) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"code":    fiber.StatusInternalServerError,
		"message": errMsg,
	})
}

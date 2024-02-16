package lib

import (
	"meta-x/models"

	"github.com/gofiber/fiber/v2"
)

func BadRequestErr(c *fiber.Ctx, errMsg any) error {
	return c.Status(fiber.StatusBadRequest).JSON(models.ErrResp{
		Code:    fiber.StatusBadRequest,
		Message: errMsg,
	})
}

func UnprocessableEntityErr(c *fiber.Ctx, errMsg any) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(models.ErrResp{
		Code:    fiber.StatusUnprocessableEntity,
		Message: errMsg,
	})
}

func ForbiddenErr(c *fiber.Ctx, errMsg any) error {
	return c.Status(fiber.StatusForbidden).JSON(models.ErrResp{
		Code:    fiber.StatusForbidden,
		Message: errMsg,
	})
}

func InternalServerErr(c *fiber.Ctx, errMsg any) error {
	return c.Status(fiber.StatusInternalServerError).JSON(models.ErrResp{
		Code:    fiber.StatusInternalServerError,
		Message: errMsg,
	})
}

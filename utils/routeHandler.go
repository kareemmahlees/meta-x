package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type RouteHandlerFunc func(*fiber.Ctx,*sqlx.DB) error

func RouteHandler(db *sqlx.DB,routeHandlerFunc RouteHandlerFunc) fiber.Handler {
    fn := func(c *fiber.Ctx) error {
		return routeHandlerFunc(c,db)
    }

	return fn
}
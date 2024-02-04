package utils

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

type RouteHandlerFunc func(*fiber.Ctx, *sqlx.DB) error

func RouteHandler(db *sqlx.DB, routeHandlerFunc RouteHandlerFunc) fiber.Handler {
	fn := func(c *fiber.Ctx) error {
		return routeHandlerFunc(c, db)
	}

	return fn
}

func GraphQLHandler(f func(http.ResponseWriter, *http.Request)) func(ctx *fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		fasthttpadaptor.NewFastHTTPHandler(http.HandlerFunc(f))(ctx.Context())
	}
}

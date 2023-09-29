package utils

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

func RouteHandler(db *sqlx.DB, routeHandlerFunc func(*fiber.Ctx, *sqlx.DB) error) fiber.Handler {
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

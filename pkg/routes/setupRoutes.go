package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func Setup(app *fiber.App, db *sqlx.DB) {
	RegisterDefaultRoutes(app)
	RegisterDatabasesRoutes(app, db)
	RegisterTablesRoutes(app, db)
}

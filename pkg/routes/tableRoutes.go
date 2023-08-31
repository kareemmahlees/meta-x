package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/kareemmahlees/mysql-meta/lib"
	db_handlers "github.com/kareemmahlees/mysql-meta/pkg/db"
	"github.com/kareemmahlees/mysql-meta/utils"
)

func RegisterTablesRoutes(app *fiber.App, db *sqlx.DB) {
	tableGroup := app.Group("tables")
	tableGroup.Get("/:dbName", utils.RouteHandler(db, handleListTables))
}

func handleListTables(c *fiber.Ctx, db *sqlx.DB) error {
	tables, err := db_handlers.ListTables(db, c.Params("dbName"))
	if err != nil {
		return c.JSON(lib.ResponseError500(err.Error()))
	}
	return c.JSON(fiber.Map{"tables": tables})
}

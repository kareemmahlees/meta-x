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
	tableGroup.Post("/:dbName/:tableName", utils.RouteHandler(db, handleCreateTable))
}

func handleListTables(c *fiber.Ctx, db *sqlx.DB) error {
	tables, err := db_handlers.ListTables(db, c.Params("dbName"))
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return c.JSON(fiber.Map{"tables": tables})
}

func handleCreateTable(c *fiber.Ctx, db *sqlx.DB) error {
	var payload map[string]lib.TablePropsValidator
	if err := c.BodyParser(&payload); err != nil {
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	for _, v := range payload {
		errs := lib.Validate(v)
		if len(errs) > 0 {
			return c.Status(400).JSON(lib.ResponseError400(errs))
		}
	}
	_, err := db_handlers.CreateTable(db, c.Params("dbName"), c.Params("tableName"), payload)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"created": c.Params("tableName")})
}

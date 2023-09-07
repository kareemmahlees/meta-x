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
	tableGroup.Get("/", utils.RouteHandler(db, handleListTables))
	tableGroup.Post("/:tableName", utils.RouteHandler(db, handleCreateTable))
	tableGroup.Delete("/:tableName", utils.RouteHandler(db, handleDeleteTable))
	tableGroup.Put("/:tableName", utils.RouteHandler(db, handleUpdateTable))
}

func handleListTables(c *fiber.Ctx, db *sqlx.DB) error {
	tables, err := db_handlers.ListTables(db)
	if err != nil {
		return c.JSON(lib.ResponseError500(err.Error()))
	}
	return c.JSON(fiber.Map{"tables": tables})
}

func handleCreateTable(c *fiber.Ctx, db *sqlx.DB) error {
	if err := lib.ValidateVar(c.Params("tableName"), "required,alpha"); err != nil {
		return c.JSON(lib.ResponseError400(err.Error()))
	}
	var payload map[string]lib.CreateTableProps
	if err := c.BodyParser(&payload); err != nil {
		return c.JSON(lib.ResponseError500(err.Error()))
	}
	for _, v := range payload {
		errs := lib.ValidateStruct(v)
		if len(errs) > 0 {
			return c.Status(400).JSON(lib.ResponseError400(errs))
		}
	}
	err := db_handlers.CreateTable(db, c.Params("tableName"), payload)
	if err != nil {
		return c.JSON(lib.ResponseError500(err.Error()))
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"created": c.Params("tableName")})
}

func handleUpdateTable(c *fiber.Ctx, db *sqlx.DB) error {
	if err := lib.ValidateVar(c.Params("tableName"), "required,alpha"); err != nil {
		return c.JSON(lib.ResponseError400(err.Error()))
	}
	var payload lib.UpdateTableProps
	if err := c.BodyParser(&payload); err != nil {
		return c.JSON(lib.ResponseError500(err.Error()))
	}
	errs := lib.ValidateStruct(payload)
	if len(errs) > 0 {
		return c.JSON(lib.ResponseError400(errs))
	}
	err := db_handlers.UpdateTable(db, c.Params("tableName"), payload)
	if err != nil {
		return c.JSON(lib.ResponseError500(err.Error()))
	}
	return c.JSON(fiber.Map{"success": true})
}

func handleDeleteTable(c *fiber.Ctx, db *sqlx.DB) error {
	if err := lib.ValidateVar(c.Params("tableName"), "required,alpha"); err != nil {
		return c.JSON(lib.ResponseError400(err.Error()))
	}
	err := db_handlers.DeleteTable(db, c.Params("tableName"))
	if err != nil {
		return c.JSON(lib.ResponseError500(err.Error()))
	}
	return c.JSON(fiber.Map{"success": true})
}

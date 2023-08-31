package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/kareemmahlees/mysql-meta/docs"
	db_handlers "github.com/kareemmahlees/mysql-meta/pkg/db"
	"github.com/kareemmahlees/mysql-meta/utils"
)

func RegisterDatabasesRoutes(app *fiber.App, db *sqlx.DB) {
	dbGroup := app.Group("databases")
	dbGroup.Get("", utils.RouteHandler(db, handleListDatabases))
	dbGroup.Post("/:name", utils.RouteHandler(db, handlerCreateDatabase))
}

type HandleListDatabasesResult struct {
	Databases []string
}

// @tags		Databases
// @decription	list databases
// @router		/databases [get]
// @produce	json
// @success	200	{object}	HandleListDatabasesResult
func handleListDatabases(c *fiber.Ctx, db *sqlx.DB) error {
	dbs := db_handlers.ListDatabases(db)
	return c.JSON(fiber.Map{"databases": dbs})
}

type HandleCreateDatabaseResult struct {
	Created int
}

// @tags			Databases
// @description	create database
// @router			/databases [post]
// @param			name	path	string	true	"database name"
// @prduce			json
// @success		201	{object}	HandleCreateDatabaseResult
func handlerCreateDatabase(c *fiber.Ctx, db *sqlx.DB) error {
	rowsAffected, err := db_handlers.CreateDatabase(db, c.Params("name"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err})
	}
	return c.Status(201).JSON(fiber.Map{"created": rowsAffected})
}

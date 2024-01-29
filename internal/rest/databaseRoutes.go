package routes

import (
	_ "meta-x/docs"
	db_handlers "meta-x/internal/db"
	"meta-x/lib"
	"meta-x/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func RegisterDatabasesRoutes(app *fiber.App, db *sqlx.DB) {
	dbGroup := app.Group("databases")
	dbGroup.Get("", utils.RouteHandler(db, handleListDatabases))
	dbGroup.Post("/:name", utils.RouteHandler(db, handlerCreateDatabase))
}

type HandleListDatabasesResult struct {
	Databases []string
}

// Lists databases
//
//	@tags		Databases
//	@decription	list databases
//	@router		/databases [get]
//	@produce	json
//	@success	200	{object}	HandleListDatabasesResult
func handleListDatabases(c *fiber.Ctx, db *sqlx.DB) error {
	dbs, err := db_handlers.ListDatabases(db)
	if err != nil {
		return c.JSON(lib.ResponseError500(err.Error()))
	}
	return c.JSON(fiber.Map{"databases": dbs})
}

type HandleCreateDatabaseResult struct {
	Created int
}

// Creates database
//
//	@tags			Databases
//	@description	create database
//	@router			/databases [post]
//	@param			name	path	string	true	"database name"
//	@prduce			json
//	@success		201	{object}	HandleCreateDatabaseResult
func handlerCreateDatabase(c *fiber.Ctx, db *sqlx.DB) error {
	rowsAffected, err := db_handlers.CreateDatabase(db, c.Params("name"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err})
	}
	return c.Status(201).JSON(fiber.Map{"created": rowsAffected})
}

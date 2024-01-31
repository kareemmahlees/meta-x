package routes

import (
	db_handlers "meta-x/internal/db"
	"meta-x/lib"
	"meta-x/models"
	"meta-x/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func RegisterDatabasesRoutes(app *fiber.App, db *sqlx.DB) {
	dbGroup := app.Group("databases")
	dbGroup.Get("", utils.RouteHandler(db, handleListDatabases))
	dbGroup.Post("/", utils.RouteHandler(db, handleCreateDatabase))
}

// Lists databases
//
//	@tags			Databases
//	@description	list databases
//	@router			/databases [get]
//	@produce		json
//	@success		200	{object}	models.ListDatabasesResult
func handleListDatabases(c *fiber.Ctx, db *sqlx.DB) error {
	var dbs any
	var err error

	provider := c.Locals("provider")

	switch provider {
	case lib.SQLITE3:
		dbs, err = db_handlers.ListDatabasesSqlite(db)
	case lib.PSQL, lib.MYSQL:
		dbs, err = db_handlers.ListDatabasesPgMySQL(db, provider.(string))
	}

	if err != nil {
		return c.JSON(lib.ResponseError500(err.Error()))
	}
	return c.JSON(fiber.Map{"databases": dbs})
}

// Creates new pg/mysql database
//
//	@tags			Databases
//	@description	create pg/mysql database
//	@router			/databases [post]
//	@produce		json
//	@accept			json
//	@param			pg_mysql_db_data	body		models.CreatePgMySqlDBPayload	true	"only supported for pg and mysql, because attached sqlite dbs are temporary"
//	@success		201					{object}	models.CreateDatabaseResult
func handleCreateDatabase(c *fiber.Ctx, db *sqlx.DB) error {
	payload := new(models.CreatePgMySqlDBPayload)

	if err := c.BodyParser(payload); err != nil {
		return c.JSON(lib.ResponseError400(err))
	}

	if err := lib.ValidateStruct(payload); err != nil {
		return c.JSON(lib.ResponseError400(err))
	}

	rowsAffected, err := db_handlers.CreatePgMysqlDatabase(db, payload.Name)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err})
	}
	return c.Status(201).JSON(fiber.Map{"created": rowsAffected})
}

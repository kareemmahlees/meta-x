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
	dbGroup := app.Group("database")
	dbGroup.Get("", utils.RouteHandler(db, handleListDatabases))
	dbGroup.Post("", utils.RouteHandler(db, handleCreateDatabase))
}

// Lists databases
//
//	@tags			Databases
//	@description	list databases
//	@router			/database [get]
//	@produce		json
//	@success		200	{object}	models.ListDatabasesResp
func handleListDatabases(c *fiber.Ctx, db *sqlx.DB) error {
	var dbs []*string
	var err error

	provider := c.Locals("provider")

	switch provider {
	case lib.SQLITE3:
		dbs, err = db_handlers.ListDatabasesSqlite(db)
	case lib.PSQL, lib.MYSQL:
		dbs, err = db_handlers.ListDatabasesPgMySQL(db, provider.(string))
	}

	if err != nil {
		return lib.InternalServerErr(c, err.Error())
	}
	return c.JSON(models.ListDatabasesResp{Databases: dbs})
}

// Creates new pg/mysql database
//
//	@tags			Databases
//	@description	create pg/mysql database
//	@router			/database [post]
//	@produce		json
//	@accept			json
//	@param			pg_mysql_db_data	body		models.CreatePgMySqlDBPayload	true	"only supported for pg and mysql, because attached sqlite dbs are temporary"
//	@success		201					{object}	models.SuccessResp
func handleCreateDatabase(c *fiber.Ctx, db *sqlx.DB) error {
	var payload models.CreatePgMySqlDBPayload

	if err := c.BodyParser(&payload); err != nil {
		return lib.UnprocessableEntityErr(c, err.Error())
	}

	if errs := lib.ValidateStruct(payload); len(errs) > 0 {
		return lib.BadRequestErr(c, errs)
	}

	err := db_handlers.CreatePgMysqlDatabase(db, c.Locals("provider").(string), payload.Name)

	if err != nil {
		return lib.InternalServerErr(c, err.Error())
	}
	return c.Status(201).JSON(models.SuccessResp{Success: true})
}

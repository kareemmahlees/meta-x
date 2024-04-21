package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kareemmahlees/meta-x/internal/db"
	"github.com/kareemmahlees/meta-x/lib"
	"github.com/kareemmahlees/meta-x/models"
)

type DBHandler struct {
	storage db.DatabaseExecuter
}

// TODO: change the interface
func NewDBHandler(storage db.DatabaseExecuter) *DBHandler {
	return &DBHandler{storage}
}

func (dh *DBHandler) RegisterRoutes(app *fiber.App) {
	dbGroup := app.Group("database")
	dbGroup.Get("", dh.handleListDatabases)
	dbGroup.Post("", dh.handleCreateDatabase)
}

// Lists databases
//
//	@tags			Databases
//	@description	list databases
//	@router			/database [get]
//	@produce		json
//	@success		200	{object}	models.ListDatabasesResp
func (dh *DBHandler) handleListDatabases(c *fiber.Ctx) error {
	dbs, err := dh.storage.ListDBs()

	if err != nil {
		return lib.InternalServerErr(c, err.Error())
	}
	return c.JSON(models.ListDatabasesResp{Databases: dbs})
}

func (dh *DBHandler) handleCreateDatabase(c *fiber.Ctx) error {
	var payload models.CreatePgMySqlDBPayload

	if err := c.BodyParser(&payload); err != nil {
		return lib.UnprocessableEntityErr(c, err.Error())
	}

	if errs := lib.ValidateStruct(payload); len(errs) > 0 {
		return lib.BadRequestErr(c, errs)
	}

	err := dh.storage.CreateDB(payload.Name)

	if err != nil {
		return lib.InternalServerErr(c, err.Error())
	}
	return c.Status(201).JSON(models.SuccessResp{Success: true})
}

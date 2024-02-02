package routes

import (
	_ "meta-x/docs"
	db_handlers "meta-x/internal/db"
	"meta-x/lib"
	"meta-x/models"
	"meta-x/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func RegisterTablesRoutes(app *fiber.App, db *sqlx.DB) {
	tableGroup := app.Group("tables")
	tableGroup.Get("", utils.RouteHandler(db, handleListTables))
	tableGroup.Get("/:tableName/describe", utils.RouteHandler(db, handeGetTableInfo))
	tableGroup.Post("/:tableName", utils.RouteHandler(db, handleCreateTable))
	tableGroup.Delete("/:tableName", utils.RouteHandler(db, handleDeleteTable))
	tableGroup.Post("/:tableName/add", utils.RouteHandler(db, handleAddColumn))
	tableGroup.Put("/:tableName/modify", utils.RouteHandler(db, handleUpdateColumn))
	tableGroup.Delete("/:tableName/delete", utils.RouteHandler(db, handleDeleteColumn))
}

func handeGetTableInfo(c *fiber.Ctx, db *sqlx.DB) error {
	tableName := c.Params("tableName")

	if err := lib.ValidateVar(tableName, "required,alpha"); err != nil {
		return lib.BadRequestErr(c, err.Error())
	}
	tableInfo, err := db_handlers.GetTableInfo(db, tableName, c.Locals("provider").(string))
	if err != nil {
		return lib.InternalServerErr(c, err.Error())
	}
	return c.JSON(tableInfo)
}

// Lists all tables in the database
//
//	@tags			Tables
//	@description	list tables
//	@router			/tables [get]
//	@produce		json
//	@success		200	{object}	models.ListTablesResp
func handleListTables(c *fiber.Ctx, db *sqlx.DB) error {
	tables, err := db_handlers.ListTables(db, c.Locals("provider").(string))
	if err != nil {
		return lib.InternalServerErr(c, err.Error())
	}
	return c.JSON(models.ListTablesResp{Tables: tables})
}

// Creates a Table
//
//	@tags			Tables
//	@description	create table
//	@router			/tables/{tableName} [post]
//	@param			tableName	path	string					true	"table name"
//	@param			tableData	body	models.CreateTablePayload	true	"create table data"
//	@accept			json
//	@produce		json
//	@success		201	{object}	models.CreateTableResp
func handleCreateTable(c *fiber.Ctx, db *sqlx.DB) error {
	tableName := c.Params("tableName")

	if err := lib.ValidateVar(tableName, "required,alphanum"); err != nil {
		return lib.BadRequestErr(c, err.Error())
	}
	var payload []models.CreateTablePayload
	if err := c.BodyParser(&payload); err != nil {
		return lib.UnprocessableEntityErr(c, err)
	}
	for _, v := range payload {
		if errs := lib.ValidateStruct(v); len(errs) > 0 {
			return lib.BadRequestErr(c, errs)
		}
	}
	err := db_handlers.CreateTable(db, tableName, payload)
	if err != nil {
		return c.JSON(lib.InternalServerErr(c, err.Error()))
	}
	return c.Status(fiber.StatusCreated).JSON(models.CreateTableResp{Created: tableName})
}

type HandleUpdateDeleteResp struct {
	Success bool
}

// Updates a Table either by add,modify or delete
//
//	@tags			Tables
//	@description	update table
//	@router			/tables/{tableName} [put]
//	@param			tableName	path	string					true	"table name"
//	@param			tableData	body	lib.UpdateTableProps	true	"update table data"
//	@accept			json
//	@produce		json
//	@success		200	{object}	models.UpdateDeleteResp

func handleAddColumn(c *fiber.Ctx, db *sqlx.DB) error {
	if err := lib.ValidateVar(c.Params("tableName"), "required,alphanum"); err != nil {
		return lib.BadRequestErr(c, err.Error())
	}
	var payload models.AddUpdateColumnPayload
	if err := c.BodyParser(&payload); err != nil {
		return lib.UnprocessableEntityErr(c, err.Error())
	}
	if errs := lib.ValidateStruct(payload); len(errs) > 0 {
		return lib.BadRequestErr(c, errs)
	}
	err := db_handlers.AddColumn(db, c.Params("tableName"), payload)
	if err != nil {
		return lib.InternalServerErr(c, err.Error())
	}
	return c.JSON(models.UpdateDeleteResp{Success: true})
}

func handleUpdateColumn(c *fiber.Ctx, db *sqlx.DB) error {
	if c.Locals("provider") == lib.SQLITE3 {
		return fiber.NewError(fiber.StatusForbidden, "MODIFY COLUMN not supported by sqlite")
	}

	if err := lib.ValidateVar(c.Params("tableName"), "required,alphanum"); err != nil {
		return lib.BadRequestErr(c, err.Error())
	}
	var payload models.AddUpdateColumnPayload
	if err := c.BodyParser(&payload); err != nil {
		return lib.UnprocessableEntityErr(c, err.Error())
	}
	if errs := lib.ValidateStruct(payload); len(errs) > 0 {
		return lib.BadRequestErr(c, errs)
	}
	err := db_handlers.UpdateColumn(db, c.Params("tableName"), payload)
	if err != nil {
		return lib.InternalServerErr(c, err.Error())
	}
	return c.JSON(models.UpdateDeleteResp{Success: true})
}

func handleDeleteColumn(c *fiber.Ctx, db *sqlx.DB) error {
	if err := lib.ValidateVar(c.Params("tableName"), "required,alphanum"); err != nil {
		return lib.BadRequestErr(c, err.Error())
	}
	var payload models.DeleteColumnPayload
	if err := c.BodyParser(&payload); err != nil {
		return lib.UnprocessableEntityErr(c, err.Error())
	}
	if errs := lib.ValidateStruct(payload); len(errs) > 0 {
		return lib.BadRequestErr(c, errs)
	}
	err := db_handlers.DeleteColumn(db, c.Params("tableName"), payload)
	if err != nil {
		return lib.InternalServerErr(c, err.Error())
	}
	return c.JSON(models.UpdateDeleteResp{Success: true})
}

// Deletes a table
//
//	@tags		Tables
//	@decription	delete table
//	@router		/tables/{tableName} [delete]
//	@param		tableName	path	string	true	"table name"
//	@accept		json
//	@produce	json
//	@success	200	{object}	HandleUpdateDeleteResp
func handleDeleteTable(c *fiber.Ctx, db *sqlx.DB) error {
	if err := lib.ValidateVar(c.Params("tableName"), "required,alpha"); err != nil {
		return lib.BadRequestErr(c, err.Error())
	}
	err := db_handlers.DeleteTable(db, c.Params("tableName"))
	if err != nil {
		return lib.InternalServerErr(c, err.Error())
	}
	return c.JSON(fiber.Map{"success": true})
}

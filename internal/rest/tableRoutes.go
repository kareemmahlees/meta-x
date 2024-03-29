package routes

import (
	_ "github.com/kareemmahlees/meta-x/docs"
	db_handlers "github.com/kareemmahlees/meta-x/internal/db"
	"github.com/kareemmahlees/meta-x/lib"
	"github.com/kareemmahlees/meta-x/models"
	"github.com/kareemmahlees/meta-x/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func RegisterTablesRoutes(app *fiber.App, db *sqlx.DB) {
	tableGroup := app.Group("table")
	tableGroup.Get("", utils.RouteHandler(db, handleListTables))
	tableGroup.Get("/:tableName/describe", utils.RouteHandler(db, handleGetTableInfo))
	tableGroup.Post("/:tableName", utils.RouteHandler(db, handleCreateTable))
	tableGroup.Delete("/:tableName", utils.RouteHandler(db, handleDeleteTable))
	tableGroup.Post("/:tableName/column/add", utils.RouteHandler(db, handleAddColumn))
	tableGroup.Put("/:tableName/column/modify", utils.RouteHandler(db, handleModifyColumn))
	tableGroup.Delete("/:tableName/column/delete", utils.RouteHandler(db, handleDeleteColumn))
}

// Get detailed info about the specified table
//
//	@tags			Tables
//	@description	Get detailed info about a specific table
//	@router			/table/{tableName}/describe [get]
//	@produce		json
//	@success		200	{object}	[]models.TableInfoResp
func handleGetTableInfo(c *fiber.Ctx, db *sqlx.DB) error {
	params := struct {
		TableName string `params:"tableName" validate:"required,alpha"`
	}{}
	_ = c.ParamsParser(&params)

	if errs := lib.ValidateStruct(params); len(errs) > 0 {
		return lib.BadRequestErr(c, errs)
	}
	tableInfo, err := db_handlers.GetTableInfo(db, params.TableName, c.Locals("provider").(string))
	if err != nil {
		return lib.InternalServerErr(c, err.Error())
	}
	return c.JSON(tableInfo)
}

// Lists all tables in the database
//
//	@tags			Tables
//	@description	list tables
//	@router			/table [get]
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
//	@router			/table/{tableName} [post]
//	@param			tableName	path	string						true	"table name"
//	@param			tableData	body	models.CreateTablePayload	true	"create table data"
//	@accept			json
//	@produce		json
//	@success		201	{object}	models.CreateTableResp
func handleCreateTable(c *fiber.Ctx, db *sqlx.DB) error {
	params := struct {
		TableName string `params:"tableName" validate:"required,alphanum"`
	}{}
	_ = c.ParamsParser(&params)

	if errs := lib.ValidateStruct(params); len(errs) > 0 {
		return lib.BadRequestErr(c, errs)
	}
	var payload []models.CreateTablePayload
	if err := c.BodyParser(&payload); err != nil {
		return lib.UnprocessableEntityErr(c, err.Error())
	}
	for _, v := range payload {
		if errs := lib.ValidateStruct(v); len(errs) > 0 {
			return lib.BadRequestErr(c, errs)
		}
	}
	err := db_handlers.CreateTable(db, params.TableName, payload)
	if err != nil {
		return lib.InternalServerErr(c, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(models.CreateTableResp{Created: params.TableName})
}

// Updates a table by adding a column
//
//	@tags			Tables
//	@description	Add column to table
//	@router			/table/{tableName}/column/add [post]
//	@param			tableName	path	string							true	"table name"
//	@param			columnData	body	models.AddModifyColumnPayload	true	"column data"
//	@accept			json
//	@produce		json
//	@success		201	{object}	models.SuccessResp
func handleAddColumn(c *fiber.Ctx, db *sqlx.DB) error {
	params := struct {
		TableName string `params:"tableName" validate:"required,alphanum"`
	}{}
	_ = c.ParamsParser(&params)
	if errs := lib.ValidateStruct(params); len(errs) > 0 {
		return lib.BadRequestErr(c, errs)
	}
	var payload models.AddModifyColumnPayload
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
	return c.Status(fiber.StatusCreated).JSON(models.SuccessResp{Success: true})
}

// Updates a table by modifying a column
//
//	@tags			Tables
//	@description	Update table column
//	@router			/table/{tableName}/column/modify [put]
//	@param			tableName	path	string							true	"table name"
//	@param			columnData	body	models.AddModifyColumnPayload	true	"column data"
//	@accept			json
//	@produce		json
//	@success		200	{object}	models.SuccessResp
func handleModifyColumn(c *fiber.Ctx, db *sqlx.DB) error {
	if c.Locals("provider") == lib.SQLITE3 {
		return lib.ForbiddenErr(c, "MODIFY COLUMN not supported by sqlite")
	}
	params := struct {
		TableName string `params:"tableName" validate:"required,alphanum"`
	}{}

	_ = c.ParamsParser(&params)
	if errs := lib.ValidateStruct(params); len(errs) > 0 {
		return lib.BadRequestErr(c, errs)
	}
	var payload models.AddModifyColumnPayload
	if err := c.BodyParser(&payload); err != nil {
		return lib.UnprocessableEntityErr(c, err.Error())
	}
	if errs := lib.ValidateStruct(payload); len(errs) > 0 {
		return lib.BadRequestErr(c, errs)
	}
	err := db_handlers.UpdateColumn(db, c.Locals("provider").(string), c.Params("tableName"), payload)
	if err != nil {
		return lib.InternalServerErr(c, err.Error())
	}
	return c.JSON(models.SuccessResp{Success: true})
}

// Updates a table by deleting/dropping a column
//
//	@tags			Tables
//	@description	Delete/Drop table column
//	@router			/table/{tableName}/column/delete [delete]
//	@param			tableName	path	string						true	"table name"
//	@param			columnData	body	models.DeleteColumnPayload	true	"column name"
//	@accept			json
//	@produce		json
//	@success		200	{object}	models.SuccessResp
func handleDeleteColumn(c *fiber.Ctx, db *sqlx.DB) error {
	params := struct {
		TableName string `params:"tableName" validate:"required,alphanum"`
	}{}
	_ = c.ParamsParser(&params)
	if errs := lib.ValidateStruct(params); len(errs) > 0 {
		return lib.BadRequestErr(c, errs)
	}

	var payload models.DeleteColumnPayload
	if err := c.BodyParser(&payload); err != nil {
		return lib.UnprocessableEntityErr(c, err.Error())
	}
	if errs := lib.ValidateStruct(payload); len(errs) > 0 {
		return lib.BadRequestErr(c, errs)
	}
	err := db_handlers.DeleteColumn(db, params.TableName, payload)
	if err != nil {
		return lib.InternalServerErr(c, err.Error())
	}
	return c.JSON(models.SuccessResp{Success: true})
}

// Deletes a table
//
//	@tags		Tables
//	@decription	delete table
//	@router		/table/{tableName} [delete]
//	@param		tableName	path	string	true	"table name"
//	@accept		json
//	@produce	json
//	@success	200	{object}	models.SuccessResp
func handleDeleteTable(c *fiber.Ctx, db *sqlx.DB) error {
	params := struct {
		TableName string `params:"tableName" validate:"required,alpha"`
	}{}
	_ = c.ParamsParser(&params)
	if errs := lib.ValidateStruct(params); len(errs) > 0 {
		return lib.BadRequestErr(c, errs)
	}

	err := db_handlers.DeleteTable(db, c.Params("tableName"))
	if err != nil {
		return lib.InternalServerErr(c, err.Error())
	}
	return c.JSON(models.SuccessResp{Success: true})
}

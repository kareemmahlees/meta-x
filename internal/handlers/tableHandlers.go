package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kareemmahlees/meta-x/internal/db"
	"github.com/kareemmahlees/meta-x/lib"
	"github.com/kareemmahlees/meta-x/models"
)

type TableHandler struct {
	storage db.TableExecuter
}

func NewTableHandler(storage db.TableExecuter) *TableHandler {
	return &TableHandler{storage}
}

func (th *TableHandler) RegisterRoutes(app *fiber.App) {
	tableGroup := app.Group("table")
	tableGroup.Get("", th.handleListTables)
	tableGroup.Get("/:tableName/describe", th.handleGetTableInfo)
	tableGroup.Post("/:tableName", th.handleCreateTable)
	tableGroup.Delete("/:tableName", th.handleDeleteTable)
	tableGroup.Post("/:tableName/column/add", th.handleAddColumn)
	tableGroup.Put("/:tableName/column/modify", th.handleModifyColumn)
	tableGroup.Delete("/:tableName/column/delete", th.handleDeleteColumn)
}

// Get detailed info about the specified table
//
//	@tags			Tables
//	@description	Get detailed info about a specific table
//	@router			/table/{tableName}/describe [get]
//	@produce		json
//	@success		200	{object}	[]models.TableInfoResp
func (th *TableHandler) handleGetTableInfo(c *fiber.Ctx) error {
	params := struct {
		TableName string `params:"tableName" validate:"required,alpha"`
	}{}
	_ = c.ParamsParser(&params)

	if errs := lib.ValidateStruct(params); len(errs) > 0 {
		return lib.BadRequestErr(c, errs)
	}
	tableInfo, err := th.storage.GetTable(params.TableName)
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
func (th *TableHandler) handleListTables(c *fiber.Ctx) error {
	tables, err := th.storage.ListTables()
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
func (th *TableHandler) handleCreateTable(c *fiber.Ctx) error {
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
	err := th.storage.CreateTable(params.TableName, payload)
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
func (th *TableHandler) handleAddColumn(c *fiber.Ctx) error {
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
	err := th.storage.AddColumn(c.Params("tableName"), payload)
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
func (th *TableHandler) handleModifyColumn(c *fiber.Ctx) error {
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
	err := th.storage.UpdateColumn(c.Params("tableName"), payload)
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
func (th *TableHandler) handleDeleteColumn(c *fiber.Ctx) error {
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
	err := th.storage.DeleteColumn(params.TableName, payload)
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
func (th *TableHandler) handleDeleteTable(c *fiber.Ctx) error {
	params := struct {
		TableName string `params:"tableName" validate:"required,alpha"`
	}{}
	_ = c.ParamsParser(&params)
	if errs := lib.ValidateStruct(params); len(errs) > 0 {
		return lib.BadRequestErr(c, errs)
	}

	err := th.storage.DeleteTable(c.Params("tableName"))
	if err != nil {
		return lib.InternalServerErr(c, err.Error())
	}
	return c.JSON(models.SuccessResp{Success: true})
}

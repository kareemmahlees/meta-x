package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
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

func (h *TableHandler) RegisterRoutes(r *chi.Mux) {
	r.Route("/table", func(r chi.Router) {
		r.Get("/", h.handleListTables)
		r.Get("/{tableName}/describe", h.handleGetTableInfo)
		r.Post("/{tableName}", h.handleCreateTable)
		r.Delete("/{tableName}", h.handleDeleteTable)
		r.Post("/{tableName}/column/add", h.handleAddColumn)
		r.Put("/{tableName}/column/modify", h.handleModifyColumn)
		r.Delete("/{tableName}/column/delete", h.handleDeleteColumn)
	})
}

// Lists all tables in the database
//
//	@summary		List tables.
//	@description	Get a list of the available tables in the database.
//	@tags			Table
//	@router			/table [get]
//	@produce		json
//	@success		200	{object}	models.ListTablesResp
//	@failure		500	{object}	models.InternalServerError
func (h *TableHandler) handleListTables(w http.ResponseWriter, r *http.Request) {
	tables, err := h.storage.ListTables()
	if err != nil {
		httpError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJson(w, models.ListTablesResp{Tables: tables})
}

// Get detailed info about the specified table
//
//	@summary		Get table info
//	@description	Get detailed info about a table's fields.
//	@tags			Table
//	@router			/table/{table_name}/describe [get]
//	@param			table_name	path	string	true	"Table Name"
//	@produce		json
//	@success		200	{array}		models.TableColumnInfo
//	@failure		500	{object}	models.InternalServerError
func (h *TableHandler) handleGetTableInfo(w http.ResponseWriter, r *http.Request) {
	params := struct {
		TableName string `validate:"required,alpha"`
	}{
		TableName: chi.URLParam(r, "tableName"),
	}

	if errs := lib.ValidateStruct(&params); len(errs) > 0 {
		httpError(w, http.StatusBadRequest, errs)
		return
	}
	tableInfo, err := h.storage.GetTable(params.TableName)
	if err != nil {
		httpError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJson(w, tableInfo)
}

// Creates a Table
//
//	@tags			Table
//	@summary		Create a Table.
//	@description	Creates a new table with the specified columns.
//	@router			/table/{table_name} [post]
//	@param			table_name	path	string						true	"Table Name"
//	@param			table_data	body	[]models.CreateTablePayload	true	"Table Data"
//	@accept			json
//	@produce		json
//	@success		201	{object}	models.CreateTableResp
func (h *TableHandler) handleCreateTable(w http.ResponseWriter, r *http.Request) {
	params := struct {
		TableName string `validate:"required,alphanum"`
	}{
		TableName: chi.URLParam(r, "tableName"),
	}
	if errs := lib.ValidateStruct(&params); len(errs) > 0 {
		httpError(w, http.StatusBadRequest, errs)
		return
	}
	var payload []models.CreateTablePayload
	if err := parseBody(r.Body, &payload); err != nil {
		httpError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	for _, v := range payload {
		if errs := lib.ValidateStruct(v); len(errs) > 0 {
			httpError(w, http.StatusBadRequest, errs)
			return
		}
	}
	err := h.storage.CreateTable(params.TableName, payload)
	if err != nil {
		httpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	writeJson(w, models.CreateTableResp{Created: params.TableName})
}

// Updates a table by adding a column
//
//	@tags			Table
//	@description	Add column to table
//	@router			/table/{tableName}/column/add [post]
//	@param			tableName	path	string							true	"table name"
//	@param			columnData	body	models.AddModifyColumnPayload	true	"column data"
//	@accept			json
//	@produce		json
//	@success		201	{object}	models.SuccessResp
func (h *TableHandler) handleAddColumn(w http.ResponseWriter, r *http.Request) {
	params := struct {
		TableName string `validate:"required,alphanum"`
	}{
		TableName: chi.URLParam(r, "tableName"),
	}
	if errs := lib.ValidateStruct(&params); len(errs) > 0 {
		httpError(w, http.StatusBadRequest, errs)
		return
	}
	var payload models.AddModifyColumnPayload
	if err := parseBody(r.Body, &payload); err != nil {
		httpError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	if errs := lib.ValidateStruct(&payload); len(errs) > 0 {
		httpError(w, http.StatusBadRequest, errs)
		return
	}
	err := h.storage.AddColumn(params.TableName, payload)
	if err != nil {
		httpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	writeJson(w, models.SuccessResp{Success: true})
}

// Updates a table by modifying a column
//
//	@tags			Table
//	@description	Update table column
//	@router			/table/{tableName}/column/modify [put]
//	@param			tableName	path	string							true	"table name"
//	@param			columnData	body	models.AddModifyColumnPayload	true	"column data"
//	@accept			json
//	@produce		json
//	@success		200	{object}	models.SuccessResp
func (h *TableHandler) handleModifyColumn(w http.ResponseWriter, r *http.Request) {
	params := struct {
		TableName string `validate:"required,alphanum"`
	}{
		TableName: chi.URLParam(r, "tableName"),
	}

	if errs := lib.ValidateStruct(&params); len(errs) > 0 {
		httpError(w, http.StatusBadRequest, errs)
		return
	}
	var payload models.AddModifyColumnPayload
	if err := parseBody(r.Body, &payload); err != nil {
		httpError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	if errs := lib.ValidateStruct(&payload); len(errs) > 0 {
		httpError(w, http.StatusBadRequest, errs)
		return
	}
	err := h.storage.UpdateColumn(params.TableName, payload)
	if err != nil {
		httpError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJson(w, models.SuccessResp{Success: true})
}

// Updates a table by deleting/dropping a column
//
//	@tags			Table
//	@description	Delete/Drop table column
//	@router			/table/{tableName}/column/delete [delete]
//	@param			tableName	path	string						true	"table name"
//	@param			columnData	body	models.DeleteColumnPayload	true	"column name"
//	@accept			json
//	@produce		json
//	@success		200	{object}	models.SuccessResp
func (h *TableHandler) handleDeleteColumn(w http.ResponseWriter, r *http.Request) {
	params := struct {
		TableName string `params:"tableName" validate:"required,alphanum"`
	}{
		TableName: chi.URLParam(r, "tableName"),
	}
	if errs := lib.ValidateStruct(&params); len(errs) > 0 {
		httpError(w, http.StatusBadRequest, errs)
		return
	}

	var payload models.DeleteColumnPayload
	if err := parseBody(r.Body, &payload); err != nil {
		httpError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	if errs := lib.ValidateStruct(&payload); len(errs) > 0 {
		httpError(w, http.StatusBadRequest, errs)
		return
	}
	err := h.storage.DeleteColumn(params.TableName, payload)
	if err != nil {
		httpError(w, http.StatusInternalServerError, err)
		return
	}
	writeJson(w, models.SuccessResp{Success: true})
}

// Deletes a table
//
//	@tags		Table
//	@decription	delete table
//	@router		/table/{tableName} [delete]
//	@param		tableName	path	string	true	"table name"
//	@accept		json
//	@produce	json
//	@success	200	{object}	models.SuccessResp
func (h *TableHandler) handleDeleteTable(w http.ResponseWriter, r *http.Request) {
	params := struct {
		TableName string `params:"tableName" validate:"required,alpha"`
	}{
		TableName: chi.URLParam(r, "tableName"),
	}
	if errs := lib.ValidateStruct(&params); len(errs) > 0 {
		httpError(w, http.StatusBadRequest, errs)
		return
	}

	err := h.storage.DeleteTable(params.TableName)
	if err != nil {
		httpError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJson(w, models.SuccessResp{Success: true})
}

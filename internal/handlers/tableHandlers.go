package handlers

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/kareemmahlees/meta-x/internal/db"
	"github.com/kareemmahlees/meta-x/models"
)

type TableHandler struct {
	storage db.TableExecuter
}

func NewTableHandler(storage db.TableExecuter) *TableHandler {
	return &TableHandler{storage}
}

func (h *TableHandler) RegisterRoutes(api huma.API) {
	group := huma.NewGroup(api, "/tables")
	huma.Register(group, huma.Operation{
		OperationID: "list-tables",
		Method:      http.MethodGet,
		Path:        "/",
		Summary:     "List tables",
		Description: "Get a list of the available tables in the database",
		Tags:        []string{"Table"},
	}, h.handleListTables)
	huma.Register(group, huma.Operation{
		OperationID: "get-table-info",
		Method:      http.MethodGet,
		Path:        "/{tableName}",
		Summary:     "Get table info",
		Description: "Get detailed info about a table's fields",
		Tags:        []string{"Table"},
	}, h.handleGetTableInfo)
	huma.Register(group, huma.Operation{
		OperationID:   "create-table",
		Method:        http.MethodPost,
		Path:          "/{tableName}",
		Summary:       "Create a table",
		Description:   "Creates a new table with the specified columns",
		Tags:          []string{"Table"},
		DefaultStatus: http.StatusCreated,
	}, h.handleCreateTable)
	huma.Register(group, huma.Operation{
		OperationID: "delete-table",
		Method:      http.MethodDelete,
		Path:        "/{tableName}",
		Summary:     "Delete table",
		Description: "Delete the given table",
		Tags:        []string{"Table"},
	}, h.handleDeleteTable)
	huma.Register(group, huma.Operation{
		OperationID:   "add-column",
		Method:        http.MethodPost,
		Path:          "/{tableName}/column/add",
		Summary:       "Add a column",
		Description:   "Adds a column with the provided data to the given table",
		Tags:          []string{"Table", "Column"},
		DefaultStatus: http.StatusCreated,
	}, h.handleAddColumn)
	huma.Register(group, huma.Operation{
		OperationID: "alter-column",
		Method:      http.MethodPut,
		Path:        "/{tableName}/column/alter",
		Summary:     "Update column",
		Description: "Update table column properties, **Only supports updating the column type for now**",
		Tags:        []string{"Table", "Column"},
	}, h.handleAlterColumn)
	huma.Register(group, huma.Operation{
		OperationID: "drop-column",
		Method:      http.MethodDelete,
		Path:        "/{tableName}/column/drop",
		Summary:     "Delete Column",
		Tags:        []string{"Table", "Column"},
	}, h.handleDropColumn)
}

type SuccessOutput struct {
	Body models.SuccessResp
}

type ListTablesOutput struct {
	Body models.ListTablesResp
}

func (h *TableHandler) handleListTables(ctx context.Context, input *struct{}) (*ListTablesOutput, error) {
	tables, err := h.storage.ListTables()
	if err != nil {
		return nil, huma.Error500InternalServerError("Something went wrong")
	}
	return &ListTablesOutput{
		Body: models.ListTablesResp{Tables: tables},
	}, nil
}

type GetTableInfoInput struct {
	TableName string `path:"tableName"`
}
type GetTableInfoOutput struct {
	Body []*models.TableColumnInfo
}

func (h *TableHandler) handleGetTableInfo(ctx context.Context, input *GetTableInfoInput) (*GetTableInfoOutput, error) {
	tableInfo, err := h.storage.GetTable(input.TableName)
	if err != nil {
		return nil, huma.Error500InternalServerError("Something went wrong")
	}
	return &GetTableInfoOutput{
		Body: tableInfo,
	}, nil
}

type CreateTableInput struct {
	GetTableInfoInput
	Body []models.CreateTablePayload
}
type CreateTableOutput struct {
	Body models.CreateTableResp
}

func (h *TableHandler) handleCreateTable(ctx context.Context, input *CreateTableInput) (*CreateTableOutput, error) {
	err := h.storage.CreateTable(input.TableName, input.Body)
	if err != nil {
		return nil, huma.Error500InternalServerError("Something went wrong")
	}

	return &CreateTableOutput{
		Body: models.CreateTableResp{Created: input.TableName},
	}, nil
}

type AddColumnInput struct {
	GetTableInfoInput
	Body models.AddModifyColumnPayload
}

func (h *TableHandler) handleAddColumn(ctx context.Context, input *AddColumnInput) (*SuccessOutput, error) {
	err := h.storage.AddColumn(input.TableName, input.Body)
	if err != nil {
		return nil, huma.Error500InternalServerError("Something went wrong")
	}

	return &SuccessOutput{Body: models.SuccessResp{Success: true}}, nil
}

type AlterColumnInput struct {
	GetTableInfoInput
	Body models.AddModifyColumnPayload
}

func (h *TableHandler) handleAlterColumn(ctx context.Context, input *AlterColumnInput) (*SuccessOutput, error) {
	err := h.storage.UpdateColumn(input.TableName, input.Body)
	if err != nil {
		return nil, huma.Error500InternalServerError("Something went wrong")
	}
	return &SuccessOutput{Body: models.SuccessResp{Success: true}}, nil
}

type DropColumnInput struct {
	GetTableInfoInput
	Body models.DeleteColumnPayload
}

func (h *TableHandler) handleDropColumn(ctx context.Context, input *DropColumnInput) (*SuccessOutput, error) {
	err := h.storage.DeleteColumn(input.TableName, input.Body)
	if err != nil {
		return nil, huma.Error500InternalServerError("Something went wrong")
	}
	return &SuccessOutput{models.SuccessResp{Success: true}}, nil
}

type DeleteTableInput struct {
	GetTableInfoInput
}

func (h *TableHandler) handleDeleteTable(ctx context.Context, input *DeleteTableInput) (*SuccessOutput, error) {
	err := h.storage.DeleteTable(input.TableName)
	if err != nil {
		return nil, huma.Error500InternalServerError("Something went wrong")
	}
	return &SuccessOutput{
		Body: models.SuccessResp{Success: true},
	}, nil
}

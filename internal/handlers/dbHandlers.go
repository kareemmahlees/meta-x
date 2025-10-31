package handlers

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/kareemmahlees/meta-x/internal/db"
	"github.com/kareemmahlees/meta-x/models"
)

type DBHandler struct {
	storage db.DatabaseExecuter
}

func NewDBHandler(storage db.DatabaseExecuter) *DBHandler {
	return &DBHandler{storage}
}

func (h *DBHandler) RegisterRoutes(api huma.API) {
	group := huma.NewGroup(api, "/database")
	huma.Register(group, huma.Operation{
		OperationID: "list-databases",
		Method:      http.MethodGet,
		Path:        "/",
		Summary:     "List Databases",
		Description: "Get all the available databases",
		Tags:        []string{"Database"},
	}, h.handleListDatabases)
	huma.Register(group, huma.Operation{
		OperationID:   "create-database",
		Method:        http.MethodPost,
		DefaultStatus: http.StatusCreated,
		Path:          "/",
		Summary:       "Creates a new Database",
		Description:   "Creates a new database. SQLite is **not** supported",
		Tags:          []string{"Database"},
	}, h.handleCreateDatabase)
}

type ListDatabasesOutput struct {
	Body models.ListDatabasesResp
}

func (dh *DBHandler) handleListDatabases(ctx context.Context, input *struct{}) (*ListDatabasesOutput, error) {
	dbs, err := dh.storage.ListDBs()

	if err != nil {
		return nil, huma.Error500InternalServerError("Something went wrong")
	}
	return &ListDatabasesOutput{
		Body: models.ListDatabasesResp{Databases: dbs},
	}, nil
}

type CreateDatabaseOutput struct {
	Body models.SuccessResp
}

func (dh *DBHandler) handleCreateDatabase(ctx context.Context, input *struct {
	Body models.CreatePgMySqlDBPayload
}) (*CreateDatabaseOutput, error) {
	if err := dh.storage.CreateDB(input.Body.Name); err != nil {
		return nil, huma.Error500InternalServerError("Something went wrong")
	}

	return &CreateDatabaseOutput{
		Body: models.SuccessResp{Success: true},
	}, nil
}

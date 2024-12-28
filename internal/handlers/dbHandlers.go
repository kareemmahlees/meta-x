package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kareemmahlees/meta-x/internal/db"
	"github.com/kareemmahlees/meta-x/lib"
	"github.com/kareemmahlees/meta-x/models"
)

type DBHandler struct {
	storage db.DatabaseExecuter
}

func NewDBHandler(storage db.DatabaseExecuter) *DBHandler {
	return &DBHandler{storage}
}

func (dh *DBHandler) RegisterRoutes(r *chi.Mux) {
	r.Route("/database", func(r chi.Router) {
		r.Get("/", dh.handleListDatabases)
		r.Post("/", dh.handleCreateDatabase)
	})
}

// Lists databases
//
//	@summary		List Databases
//	@description	Get all the available databases.
//	@tags			Database
//	@router			/database [get]
//	@produce		json
//	@success		200	{object}	models.ListDatabasesResp
//	@failure		500	{object}	models.InternalServerError
func (dh *DBHandler) handleListDatabases(w http.ResponseWriter, r *http.Request) {
	dbs, err := dh.storage.ListDBs()

	if err != nil {
		httpError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJson(w, models.ListDatabasesResp{Databases: dbs})
}

// Creates a new Database
//
//	@summary		Create a new Database
//	@description	Creates a new database with the specified table schema.
//	@description
//	@description	**NOTE**: SQLite is Unsupported.
//	@tags			Database
//	@router			/database [post]
//	@accept			json
//	@produce		json
//	@param			payload	body		models.CreatePgMySqlDBPayload	true	"Database Info"
//	@success		200		{object}	models.SuccessResp
//	@failure		400		{object}	models.ErrResp
//	@failure		500		{object}	models.InternalServerError
func (dh *DBHandler) handleCreateDatabase(w http.ResponseWriter, r *http.Request) {
	var payload models.CreatePgMySqlDBPayload

	if err := parseBody(r.Body, &payload); err != nil {
		httpError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if errs := lib.ValidateStruct(payload); len(errs) > 0 {
		httpError(w, http.StatusBadRequest, errs)
		return
	}

	if err := dh.storage.CreateDB(payload.Name); err != nil {
		httpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	writeJson(w, models.SuccessResp{Success: true})
}

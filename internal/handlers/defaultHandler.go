package handlers

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type DefaultHandler struct{}

func NewDefaultHandler() *DefaultHandler {
	return &DefaultHandler{}
}

func (h *DefaultHandler) RegisterRoutes(r *chi.Mux) {
	r.Get("/health", h.healthCheck)
	r.Get("/", h.apiInfo)
}

type HealthCheckResult struct {
	Date string
}

// Checks the health
//
//	@description	check application health by getting current date
//	@produce		json
//	@tags			default
//	@router			/health [get]
//	@success		200	{object}	HealthCheckResult
func (h *DefaultHandler) healthCheck(w http.ResponseWriter, r *http.Request) {
	writeJson(w, map[string]time.Time{
		"date": time.Now(),
	})
}

type APIInfoResult struct {
	Author  string `json:"author"`
	Contact string `json:"contact"`
	Repo    string `json:"repo"`
	Year    int    `json:"yeaer"`
}

// Get info about the api
//
//	@description	get info about the api
//	@produce		json
//	@tags			default
//	@router			/ [get]
//	@success		200	{object}	APIInfoResult
func (h *DefaultHandler) apiInfo(w http.ResponseWriter, r *http.Request) {
	writeJson(w, APIInfoResult{
		Author:  "Kareem Ebrahim",
		Year:    2024,
		Contact: "kareemmahlees@gmail.com",
		Repo:    "https://github.com/kareemmahlees/meta-x",
	})
}

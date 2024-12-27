package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type DefaultHandler struct{}

func NewDefaultHandler() *DefaultHandler {
	return &DefaultHandler{}
}

func (h *DefaultHandler) RegisterRoutes(r *chi.Mux) {
	r.Get("/", h.apiInfo)
}

// General API Info
//
//	@description	General Info about the API, author name, how to contact, etc.
type APIInfo struct {
	Author  string `json:"author" example:"Kareem Ebrahim"`                        // Author name.
	Contact string `json:"contact" example:"kareemmahlees@gmail.com"`              // How to contact the author.
	Repo    string `json:"repo" example:"https://github.com/kareemmahlees/meta-x"` // Git repository name for contributions.
	Year    int    `json:"year" example:"2024"`                                    // Year of launch
}

// Gets general info about the API
//
//	@summary		Get API Info
//	@description	Get general info about the API
//	@produce		json
//	@tags			General
//	@router			/ [get]
//	@success		200	{object}	APIInfo
func (h *DefaultHandler) apiInfo(w http.ResponseWriter, r *http.Request) {
	writeJson(w, APIInfo{
		Author:  "Kareem Ebrahim",
		Year:    2024,
		Contact: "kareemmahlees@gmail.com",
		Repo:    "https://github.com/kareemmahlees/meta-x",
	})
}

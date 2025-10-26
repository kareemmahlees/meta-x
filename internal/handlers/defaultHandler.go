package handlers

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type DefaultHandler struct{}

func NewDefaultHandler() *DefaultHandler {
	return &DefaultHandler{}
}

func (h *DefaultHandler) RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "api-info",
		Method:      http.MethodGet,
		Path:        "/",
		Summary:     "API Info",
		Description: "General Info about the API, author name, how to contact, etc.",
	}, h.apiInfo)
}

type APIInfo struct {
	Author  string `json:"author"`
	Contact string `json:"contact"`
	Repo    string `json:"repo"`
	Year    int16  `json:"year"`
}

type APIInfoOutput struct {
	Body APIInfo
}

func (h *DefaultHandler) apiInfo(ctx context.Context, input *struct{}) (*APIInfoOutput, error) {
	return &APIInfoOutput{
		Body: APIInfo{
			Author:  "Kareem Ebrahim",
			Year:    int16(2024),
			Contact: "kareemmahlees@gmail.com",
			Repo:    "https://github.com/kareemmahlees/meta-x",
		},
	}, nil
}

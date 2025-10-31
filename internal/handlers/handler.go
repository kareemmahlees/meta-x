package handlers

import "github.com/danielgtaylor/huma/v2"

type Handler interface {
	RegisterRoutes(r huma.API)
}

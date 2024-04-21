package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kareemmahlees/meta-x/internal/db"
)

type DefaultHandler struct {
	storage *db.Storage
}

func NewDefaultHandler(storage *db.Storage) *DefaultHandler {
	return &DefaultHandler{storage}
}

func (h *DefaultHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/health", h.healthCheck)
	app.Get("/", h.apiInfo)
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
func (h *DefaultHandler) healthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"date": time.Now()})
}

type APIInfoResult struct {
	Author  string
	Year    int
	Contact string
	Repo    string
}

// Get info about the api
//
//	@description	get info about the api
//	@produce		json
//	@tags			default
//	@router			/ [get]
//	@success		200	{object}	APIInfoResult
func (h *DefaultHandler) apiInfo(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"author":  "Kareem Ebrahim",
		"year":    2023,
		"contact": "kareemmahlees@gmail.com",
		"repo":    "https://github.com/kareemmahlees/meta-x",
	})
}

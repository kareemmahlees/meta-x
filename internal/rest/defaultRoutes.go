package routes

import (
	"time"

	_ "github.com/kareemmahlees/meta-x/docs"

	"github.com/gofiber/fiber/v2"
)

func RegisterDefaultRoutes(app *fiber.App) {
	app.Get("/health", healthCheck)
	app.Get("/", apiInfo)
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
func healthCheck(c *fiber.Ctx) error {
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
func apiInfo(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"author":  "Kareem Ebrahim",
		"year":    2023,
		"contact": "kareemmahlees@gmail.com",
		"repo":    "https://github.com/kareemmahlees/mysql-meta",
	})
}

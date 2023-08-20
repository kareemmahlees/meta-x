package default_routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/kareemmahlees/mysql-meta/docs"
)

type HealthCheckResult struct {
	Date string `json:"date"`
}

//	@description	check application health by getting current date
//	@produce		json
//	@tag			Health Check
//	@router			/health [get]
//  @success 200 {object} HealthCheckResult 
func healthCheck(c *fiber.Ctx)error {
	return c.JSON(fiber.Map{"date":time.Now()})
}

type APIInfoResult struct{
	Author string `json:"author"`
	Year int `json:"year"`
	Contact string `json:"contact"`
	Repo string `json:"repo"`
}

// @description get info about the api
// @produce json
// @tag about
// @router / [get]
// @success 200 {object} APIInfoResult
func apiInfo(c *fiber.Ctx)error {
	return c.JSON(fiber.Map{
		"author":"Kareem Ebrahim",
		"year":2023,
		"contact":"kareemmahlees@gmail.com",
		"repo":"https://github.com/kareemmahlees/mysql-meta",
	})
}
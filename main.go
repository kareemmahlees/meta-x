package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	_ "github.com/kareemmahlees/mysql-meta/docs"
	"github.com/kareemmahlees/mysql-meta/pkg/routes"
)


var app *fiber.App

//	@title			MySQL Meta
//	@version		1.0
//	@description	A RESTFull and GraphQL API to manage your MySQL DB 
//	@contact.name	Kareem Ebrahim
//	@contact.email	kareemmahlees@gmail.com
//	@host			localhost:4000
//	@BasePath		/
func main() {
	app = fiber.New()

	app.Get("/swagger/*", swagger.HandlerDefault) 

	routes.Setup(app)

	port,exists := os.LookupEnv("PORT")
	if (exists){
		app.Listen(fmt.Sprintf(":%s",port))
	}
	app.Listen(":4000")
}
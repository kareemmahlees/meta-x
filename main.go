package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/kareemmahlees/mysql-meta/docs"
	"github.com/kareemmahlees/mysql-meta/pkg/db"
	"github.com/kareemmahlees/mysql-meta/pkg/routes"
)

var DB *sqlx.DB

func init(){
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
	}
	DB,err = db.InitDBConn()
	if err != nil {
		log.Error("Error connecting to DB")
	}
	if err := DB.Ping(); err!=nil{
		log.Error("Something wrong with DB" + err.Error())
	}
}

//	@title			MySQL Meta
//	@version		1.0
//	@description	A RESTFull and GraphQL API to manage your MySQL DB 
//	@contact.name	Kareem Ebrahim
//	@contact.email	kareemmahlees@gmail.com
//	@host			localhost:4000
//	@BasePath		/
func main() {
	app := fiber.New()
	defer DB.Close()

	app.Get("/swagger/*", swagger.HandlerDefault) 
	app.Use(logger.New())

	routes.Setup(app,DB)

	port,exists := os.LookupEnv("PORT")
	if (exists){
		app.Listen(fmt.Sprintf(":%s",port))
	}
	app.Listen(":4000")
}
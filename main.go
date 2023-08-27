package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
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

var con *sqlx.DB

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
	}
	con, err = db.InitDBConn()
	if err != nil {
		log.Error("Error connecting to DB")
	}
	if err := con.Ping(); err != nil {
		log.Error("Something wrong with DB" + err.Error())
	}
}

var RestStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FFFFFF")).
	Background(lipgloss.Color("#4B87FF")).
	MarginTop(1)

var GraphQLStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FFFFFF")).
	Background(lipgloss.Color("#FF70FD")).
	MarginTop(1)

var SwaggerDocsStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FFFFFF")).
	Background(lipgloss.Color("#0EEBA1")).
	MarginTop(1)

//	@title			MySQL Meta
//	@version		1.0
//	@description	A RESTFull and GraphQL API to manage your MySQL DB
//	@contact.name	Kareem Ebrahim
//	@contact.email	kareemmahlees@gmail.com
//	@host			localhost:4000
//	@BasePath		/
func main() {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	defer con.Close()

	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Use(logger.New())

	routes.Setup(app, con)

	var port string
	var exists bool
	port, exists = os.LookupEnv("PORT")
	if exists {
		err := app.Listen(fmt.Sprintf(":%s", port))
		if err != nil {
			log.Error(err)
		}
	} else {
		port = "4000"

		fmt.Println(RestStyle.Render("REST"), fmt.Sprintf("http://localhost:%s", port))
		fmt.Println(SwaggerDocsStyle.Render("SWAGGER"),fmt.Sprintf("http://localhost:%s/swagger", port))
		fmt.Println(GraphQLStyle.Render("GraphQL"), fmt.Sprintf("http://localhost:%s/graph\n", port))
		err :=app.Listen(fmt.Sprintf(":%s", port))
		if err != nil {
			log.Error(err)
		}
	}
}

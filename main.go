package main

import (
	"flag"
	"fmt"

	"github.com/kareemmahlees/mysql-meta/internal/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/kareemmahlees/mysql-meta/docs"
	"github.com/kareemmahlees/mysql-meta/internal/db"
	routes "github.com/kareemmahlees/mysql-meta/internal/rest"
	"github.com/kareemmahlees/mysql-meta/utils"
)

var con *sqlx.DB
var port int
var err error

func init() {
	flag.IntVar(&port, "port", 4000, "port to listen on")
	flag.Parse()

	_ = godotenv.Load()
	con, err = db.InitDBConn()
	if err != nil {
		log.Error("Error connecting to DB")
	}
	if err := con.Ping(); err != nil {
		log.Error("Something wrong with DB" + err.Error())
	}
}

// Main
//
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

	// see https://github.com/99designs/gqlgen/issues/1664#issuecomment-1616620967
	// Create a gqlgen handler
	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	app.All("/graphql", func(c *fiber.Ctx) error {
		utils.GraphQLHandler(h.ServeHTTP)(c)
		return nil
	})

	app.All("/playground", func(c *fiber.Ctx) error {
		utils.GraphQLHandler(playground.Handler("GraphQL", "/graphql"))(c)
		return nil
	})

	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Use(logger.New())

	routes.Setup(app, con)

	fmt.Println(utils.NewStyle("REST", "#4B87FF"), fmt.Sprintf("http://localhost:%d", port))
	fmt.Println(utils.NewStyle("Swagger", "#0EEBA1"), fmt.Sprintf("http://localhost:%d/swagger", port))
	fmt.Println(utils.NewStyle("GraphQl", "#FF70FD"), fmt.Sprintf("http://localhost:%d/graphql\n", port))

	if err := app.Listen(fmt.Sprintf(":%d", port)); err != nil {
		log.Error(err)
	}
}

package internal

import (
	"fmt"
	"github.com/kareemmahlees/meta-x/internal/graph"
	routes "github.com/kareemmahlees/meta-x/internal/rest"
	"github.com/kareemmahlees/meta-x/utils"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
)

// Initializes a database connection and starts the fiber
// server on the supplied port.
//
// listenCh is only used for the sake of testing, and app is passed as an argument instead of initializing
// inside the function to make it easier for testing
func InitDBAndServer(app *fiber.App, provider, cfg string, port int, listenCh chan bool) error {

	con, err := InitDBConn(provider, cfg)
	if err != nil {
		return err
	}
	defer con.Close()

	// see https://github.com/99designs/gqlgen/issues/1664#issuecomment-1616620967
	// Create a gqlgen handler
	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{DB: con, Provider: provider}}))

	app.All("/graphql", func(c *fiber.Ctx) error {
		utils.GraphQLHandler(h.ServeHTTP)(c)
		return nil
	}).Name("graphql")

	app.All("/playground", func(c *fiber.Ctx) error {
		utils.GraphQLHandler(playground.Handler("GraphQL", "/graphql"))(c)
		return nil
	}).Name("playground")

	app.Get("/swagger/*", swagger.HandlerDefault).Name("swagger")
	app.Use(logger.New())
	// set the provider to pass it to other handlers
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("provider", provider)
		return c.Next()
	})

	routes.Setup(app, con)

	fmt.Println(utils.NewStyle("REST", "#4B87FF"), fmt.Sprintf("http://localhost:%d", port))
	fmt.Println(utils.NewStyle("Swagger", "#0EEBA1"), fmt.Sprintf("http://localhost:%d/swagger", port))
	fmt.Println(utils.NewStyle("GraphQl", "#FF70FD"), fmt.Sprintf("http://localhost:%d/graphql", port))
	fmt.Println(utils.NewStyle("Playground", "#B6B5B5"), fmt.Sprintf("http://localhost:%d/playground\n", port))

	app.Hooks().OnListen(func(ld fiber.ListenData) error {
		listenCh <- true
		return nil
	})

	if err := app.Listen(fmt.Sprintf(":%d", port)); err != nil {
		listenCh <- false
		return err
	}
	return nil
}

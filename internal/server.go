package internal

import (
	"fmt"
	"meta-x/internal/graph"
	routes "meta-x/internal/rest"
	"meta-x/utils"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
)

// Initializes a database connection and starts the fiber
// server on the supplied port
func InitDBAndServer(provider, cfg string, port int) error {

	con, err := InitDBConn(provider, cfg)
	if err != nil {
		// log.Error("Error connecting to DB")
		return err
	}
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	defer con.Close()

	// see https://github.com/99designs/gqlgen/issues/1664#issuecomment-1616620967
	// Create a gqlgen handler
	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{DB: con, Provider: provider}}))

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

	if err := app.Listen(fmt.Sprintf(":%d", port)); err != nil {
		return err
	}
	return nil
}

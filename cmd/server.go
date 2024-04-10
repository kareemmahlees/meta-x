package cmd

import (
	"fmt"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/kareemmahlees/meta-x/internal/db"
	"github.com/kareemmahlees/meta-x/internal/graph"
	"github.com/kareemmahlees/meta-x/internal/handlers"
	"github.com/kareemmahlees/meta-x/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
)

type Server struct {
	storage  db.Storage
	port     int
	listenCh chan<- bool
}

func NewServer(storage db.Storage, port int, listenCh chan<- bool) *Server {
	return &Server{storage, port, listenCh}
}

func (s *Server) Serve() error {
	// see https://github.com/99designs/gqlgen/issues/1664#issuecomment-1616620967
	// Create a gqlgen handler
	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Storage: s.storage}}))

	app := fiber.New(fiber.Config{DisableStartupMessage: true})

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

	defaultHandler := handlers.NewDefaultHandler(nil)
	dbHandler := handlers.NewDBHandler(s.storage)
	tableHandler := handlers.NewTableHandler(s.storage)

	defaultHandler.RegisterRoutes(app)
	dbHandler.RegisterRoutes(app)
	tableHandler.RegisterRoutes(app)

	app.Hooks().OnListen(func(ld fiber.ListenData) error {

		fmt.Println(utils.NewStyle("REST", "#4B87FF"), fmt.Sprintf("http://localhost:%d", s.port))
		fmt.Println(utils.NewStyle("Swagger", "#0EEBA1"), fmt.Sprintf("http://localhost:%d/swagger", s.port))
		fmt.Println(utils.NewStyle("GraphQl", "#FF70FD"), fmt.Sprintf("http://localhost:%d/graphql", s.port))
		fmt.Println(utils.NewStyle("Playground", "#B6B5B5"), fmt.Sprintf("http://localhost:%d/playground\n", s.port))

		// s.listenCh <- true
		return nil
	})

	if err := app.Listen(fmt.Sprintf(":%d", s.port)); err != nil {
		s.listenCh <- false
		return err
	}
	return nil

}

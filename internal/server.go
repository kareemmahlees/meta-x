package internal

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kareemmahlees/meta-x/internal/db"
	"github.com/kareemmahlees/meta-x/internal/graph"
	"github.com/kareemmahlees/meta-x/internal/handlers"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	storage  db.Storage
	port     int
	listenCh chan bool
	app      *fiber.App
}

func NewServer(storage db.Storage, port int, listenCh chan bool) *Server {
	return &Server{storage, port, listenCh, nil}
}

func (s *Server) Serve() error {
	// see https://github.com/99designs/gqlgen/issues/1664#issuecomment-1616620967
	// Create a gqlgen handler
	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Storage: s.storage}}))
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Post("/graphql", h.ServeHTTP)
	r.Get("/playground", playground.ApolloSandboxHandler("GraphQL", "/graphql"))
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", s.port)),
	))

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	s.app = app

	defaultHandler := handlers.NewDefaultHandler()
	dbHandler := handlers.NewDBHandler(s.storage)
	tableHandler := handlers.NewTableHandler(s.storage)

	defaultHandler.RegisterRoutes(r)
	dbHandler.RegisterRoutes(r)
	tableHandler.RegisterRoutes(app)

	slog.Info("Server started listening", "port", s.port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), r); err != nil {
		return err
	}
	return nil
}

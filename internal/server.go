package internal

import (
	"fmt"
	"log/slog"
	"net/http"

	graphQlHandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kareemmahlees/meta-x/internal/db"
	"github.com/kareemmahlees/meta-x/internal/graph"
	"github.com/kareemmahlees/meta-x/internal/handlers"
)

type Server struct {
	storage  db.Storage
	router   *chi.Mux
	listenCh chan bool
	port     int
}

func NewServer(storage db.Storage, port int, listenCh chan bool) *Server {
	router := chi.NewMux()
	router.Use(middleware.Logger)
	router.Use(middleware.Heartbeat("/health"))
	return &Server{storage, router, listenCh, port}
}

func (s *Server) Serve() error {
	api := humachi.New(s.router, huma.DefaultConfig("MetaX", "1.0.0"))
	gqlHandler := graphQlHandler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Storage: s.storage}}))

	s.router.Post("/graphql", gqlHandler.ServeHTTP)
	s.router.Get("/playground", playground.ApolloSandboxHandler("GraphQL", "/graphql"))
	s.router.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!doctype html>
<html>
  <head>
    <title>API Reference</title>
    <meta charset="utf-8" />
    <meta
      name="viewport"
      content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <script
      id="api-reference"
      data-url="/openapi.json"></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>`))
	})

	defaultHandler := handlers.NewDefaultHandler()
	// dbHandler := handlers.NewDBHandler(s.storage)
	// tableHandler := handlers.NewTableHandler(s.storage)
	defaultHandler.RegisterRoutes(api)

	// s.registerRoutes(defaultHandler, dbHandler, tableHandler)

	slog.Info("Server started listening", "port", s.port)

	s.listenCh <- true

	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.router); err != nil {
		s.listenCh <- false
		return err
	}
	return nil
}

// func (s *Server) registerRoutes(handlers ...handlers.Handler) {
// 	for _, h := range handlers {
// 		h.RegisterRoutes(s.router)
// 	}
// }

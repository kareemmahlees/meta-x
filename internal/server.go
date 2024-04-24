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
)

type Server struct {
	storage db.Storage
	port    int
	router  *chi.Mux
}

func NewServer(storage db.Storage, port int) *Server {
	return &Server{storage, port, nil}
}

func (s *Server) Serve() error {
	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Storage: s.storage}}))
	r := chi.NewRouter()
	s.router = r

	r.Use(middleware.Logger)
	r.Post("/graphql", h.ServeHTTP)
	r.Get("/playground", playground.ApolloSandboxHandler("GraphQL", "/graphql"))
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", s.port)),
	))

	defaultHandler := handlers.NewDefaultHandler()
	dbHandler := handlers.NewDBHandler(s.storage)
	tableHandler := handlers.NewTableHandler(s.storage)

	defaultHandler.RegisterRoutes(r)
	dbHandler.RegisterRoutes(r)
	tableHandler.RegisterRoutes(r)

	slog.Info("Server started listening", "port", s.port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), r); err != nil {
		return err
	}
	return nil
}

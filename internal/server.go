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
	storage  db.Storage
	router   *chi.Mux
	listenCh chan bool
	port     int
}

func NewServer(storage db.Storage, port int, listenCh chan bool) *Server {
	r := chi.NewRouter()
	return &Server{storage, r, listenCh, port}
}

func (s *Server) Serve() error {
	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Storage: s.storage}}))

	s.router.Use(middleware.Logger)
	s.router.Post("/graphql", h.ServeHTTP)
	s.router.Get("/playground", playground.ApolloSandboxHandler("GraphQL", "/graphql"))
	s.router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", s.port)),
	))

	defaultHandler := handlers.NewDefaultHandler()
	dbHandler := handlers.NewDBHandler(s.storage)
	tableHandler := handlers.NewTableHandler(s.storage)

	defaultHandler.RegisterRoutes(s.router)
	dbHandler.RegisterRoutes(s.router)
	tableHandler.RegisterRoutes(s.router)

	slog.Info("Server started listening", "port", s.port)

	s.listenCh <- true

	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.router); err != nil {
		s.listenCh <- false
		return err
	}
	return nil
}

package internal

import (
	"fmt"
	"log/slog"
	"net/http"

	graphQlHandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/MarceloPetrucio/go-scalar-api-reference"
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
	r := chi.NewRouter()
	return &Server{storage, r, listenCh, port}
}

func (s *Server) Serve() error {
	h := graphQlHandler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Storage: s.storage}}))

	s.router.Use(middleware.Logger)
	s.router.Post("/graphql", h.ServeHTTP)
	s.router.Get("/playground", playground.ApolloSandboxHandler("GraphQL", "/graphql"))
	s.router.HandleFunc("/spec", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/swagger.json")
	})
	s.router.Get("/reference", func(w http.ResponseWriter, r *http.Request) {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			// Because the scalar client has some problems with path resolution on windows
			// we use the direct spec url.
			SpecURL: fmt.Sprintf("http://localhost:%d/spec", s.port),
			CustomOptions: scalar.CustomOptions{
				PageTitle: "MetaX",
			},
			DarkMode: true,
		})

		if err != nil {
			fmt.Printf("%v", err)
		}

		fmt.Fprintln(w, htmlContent)
	})

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

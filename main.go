package main

import (
	"meta-x/cmd"
	_ "meta-x/docs"

	"github.com/jmoiron/sqlx"
)

var con *sqlx.DB
var port int
var err error

// func init() {
// 	flag.IntVar(&port, "port", 4000, "port to listen on")
// 	provider := flag.String("provider", "", "database provider to connect with")
// 	flag.Parse()

// 	_ = godotenv.Load()
// 	switch *provider {
// 	case "sqlite":
// 		// con, err = db.InitSQLiteConn()
// 	}
// 	con, err = db.InitDBConn()
// 	if err != nil {
// 		log.Error("Error connecting to DB")
// 	}
// 	if err := con.Ping(); err != nil {
// 		log.Error("Something wrong with DB" + err.Error())
// 	}
// }

// Main
//
//	@title			MySQL Meta
//	@version		1.0
//	@description	A RESTFull and GraphQL API to manage your MySQL DB
//	@contact.name	Kareem Ebrahim
//	@contact.email	kareemmahlees@gmail.com
//	@host			localhost:5522
//	@BasePath		/
func main() {
	// app := fiber.New(fiber.Config{
	// 	DisableStartupMessage: true,
	// })
	// defer con.Close()

	// see https://github.com/99designs/gqlgen/issues/1664#issuecomment-1616620967
	// Create a gqlgen handler
	// h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{DB: con}}))

	// app.All("/graphql", func(c *fiber.Ctx) error {
	// 	utils.GraphQLHandler(h.ServeHTTP)(c)
	// 	return nil
	// })

	// app.All("/playground", func(c *fiber.Ctx) error {
	// 	utils.GraphQLHandler(playground.Handler("GraphQL", "/graphql"))(c)
	// 	return nil
	// })

	// app.Get("/swagger/*", swagger.HandlerDefault)
	// app.Use(logger.New())

	// routes.Setup(app, con)

	// fmt.Println(utils.NewStyle("REST", "#4B87FF"), fmt.Sprintf("http://localhost:%d", port))
	// fmt.Println(utils.NewStyle("Swagger", "#0EEBA1"), fmt.Sprintf("http://localhost:%d/swagger", port))
	// fmt.Println(utils.NewStyle("GraphQl", "#FF70FD"), fmt.Sprintf("http://localhost:%d/graphql", port))
	// fmt.Println(utils.NewStyle("Playground", "#B6B5B5"), fmt.Sprintf("http://localhost:%d/playground\n", port))

	// if err := app.Listen(fmt.Sprintf(":%d", port)); err != nil {
	// 	log.Error(err)
	// }

	cmd.Execute()
}

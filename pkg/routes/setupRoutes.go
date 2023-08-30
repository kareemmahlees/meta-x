package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	databases_routes "github.com/kareemmahlees/mysql-meta/pkg/routes/databases"
	default_routes "github.com/kareemmahlees/mysql-meta/pkg/routes/default"
	tables_routes "github.com/kareemmahlees/mysql-meta/pkg/routes/tables"
)

func Setup(app *fiber.App, db *sqlx.DB) {
	default_routes.RegisterDefaultRoutes(app, db)
	databases_routes.RegisterDatabasesRoutes(app, db)
	tables_routes.RegisterTablesRoutes(app, db)
}

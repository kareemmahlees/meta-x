package databases

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/kareemmahlees/mysql-meta/utils"
)

func RegisterDatabasesRoutes(app *fiber.App,db *sqlx.DB) {
	dbGroup := app.Group("databases")
	dbGroup.Get("",utils.RouteHandler(db,handleListDatabases))
	dbGroup.Post("/:name",utils.RouteHandler(db,handlerCreateDatabase))
}
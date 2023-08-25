package databases

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/kareemmahlees/mysql-meta/pkg/db/handlers"
)

func handleListDatabases(c *fiber.Ctx,db *sqlx.DB) error {
	dbs := handlers.ListDatabases(db)
	return c.JSON(fiber.Map{"databases":dbs})
}

func handlerCreateDatabase(c *fiber.Ctx,db *sqlx.DB)error{
	rowsAffected,err := handlers.CreateDatabase(db,"mysqlmeta")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error":err})
	}
	return c.Status(201).JSON(fiber.Map{"created":rowsAffected})
}
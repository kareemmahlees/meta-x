package internal

import (
	"log"
	"meta-x/lib"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestInitDBAndServer(t *testing.T) {
	listenCh := make(chan bool)

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	err := app.Shutdown()
	if err != nil {
		log.Fatal(err)
	}

	err = InitDBAndServer(app, "anything", "mallformed", 5522, listenCh)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "unknown driver")

	go func(app *fiber.App) {
		err = InitDBAndServer(app, lib.SQLITE3, ":memory:", 5522, listenCh)
		if err != nil {
			log.Fatal(err)
		}
	}(app)

	listenting := <-listenCh
	assert.True(t, listenting)

	testRoutes := []string{"graphql", "playground", "swagger"}

	for _, route := range testRoutes {
		foundRoute := app.GetRoute(route)
		assert.NotEmpty(t, foundRoute)
	}
}

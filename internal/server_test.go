package internal

import (
	"fmt"
	"log"
	"meta-x/lib"
	"meta-x/utils"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestInitDBAndServer(t *testing.T) {
	listenCh := make(chan bool)

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	defer func() {
		err := app.Shutdown()
		if err != nil {
			log.Fatal(err)
		}
	}()

	err := InitDBAndServer(app, "anything", "mallformed", 5522, listenCh)
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

		request := utils.RequestTesting[any]{
			ReqMethod: http.MethodGet,
			ReqUrl:    fmt.Sprintf("/%s", route),
		}
		_, res := request.RunRequest(app)

		assert.NotEqual(t, http.StatusNotFound, res.StatusCode)

	}

	go func(app *fiber.App) {
		err = InitDBAndServer(app, lib.SQLITE3, ":memory:", 100000, listenCh)
		assert.NotNil(t, err)
	}(app)

	listenting = <-listenCh
	assert.False(t, listenting)

}

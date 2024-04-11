package cmd

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/kareemmahlees/meta-x/utils"
	"github.com/stretchr/testify/assert"
)

// import (
// 	"fmt"
// 	"log"
// 	"github.com/kareemmahlees/meta-x/lib"
// 	"github.com/kareemmahlees/meta-x/utils"
// 	"net/http"
// 	"testing"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/stretchr/testify/assert"
// )

// func TestInitDBAndServerPassing(t *testing.T) {
// 	listenCh := make(chan bool)

// 	app := fiber.New(fiber.Config{DisableStartupMessage: true})
// 	defer func() {
// 		err := app.Shutdown()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}()

// 	err := InitDBAndServer(app, "anything", "mallformed", 5522, listenCh)
// 	assert.NotNil(t, err)
// 	assert.ErrorContains(t, err, "unknown driver")

// 	go func(app *fiber.App) {
// 		err = InitDBAndServer(app, lib.SQLITE3, ":memory:", 5522, listenCh)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}(app)

// 	listenting := <-listenCh
// 	assert.True(t, listenting)

// 	testRoutes := []string{"graphql", "playground", "swagger"}

// 	for _, route := range testRoutes {
// 		foundRoute := app.GetRoute(route)
// 		assert.NotEmpty(t, foundRoute)

// 		request := utils.RequestTesting[any]{
// 			ReqMethod: http.MethodGet,
// 			ReqUrl:    fmt.Sprintf("/%s", route),
// 		}
// 		_, res := request.RunRequest(app)

// 		assert.NotEqual(t, http.StatusNotFound, res.StatusCode)

// 	}

// }

// func TestInitDBAndServerFailing(t *testing.T) {
// 	listenCh := make(chan bool)

// 	app := fiber.New(fiber.Config{DisableStartupMessage: true})
// 	defer func() {
// 		err := app.Shutdown()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}()

// 	go func(app *fiber.App) {
// 		err := InitDBAndServer(app, lib.SQLITE3, ":memory:", 100000, listenCh)
// 		assert.NotNil(t, err)
// 	}(app)

// 	listenting := <-listenCh
// 	assert.False(t, listenting)
// }

func TestServe(t *testing.T) {
	listenCh := make(chan bool, 1)
	server := NewServer(utils.NewMockStorage(), 5522, listenCh)
	defer server.Shutdown()

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatal(err)
		}
	}()

	assert.True(t, <-listenCh)

	testRoutes := []string{"graphql", "playground", "swagger"}

	for _, route := range testRoutes {
		foundRoute := server.app.GetRoute(route)
		assert.NotEmpty(t, foundRoute)

		request := utils.RequestTesting[any]{
			ReqMethod: http.MethodGet,
			ReqUrl:    fmt.Sprintf("/%s", route),
		}
		_, res := request.RunRequest(server.app)

		assert.NotEqual(t, http.StatusNotFound, res.StatusCode)

	}
}

func TestShutDown(t *testing.T) {
	listenCh := make(chan bool, 1)
	server := NewServer(utils.NewMockStorage(), 5522, listenCh)

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatal(err)
		}
	}()

	assert.True(t, <-server.listenCh)

	err := server.Shutdown()

	assert.Nil(t, err)
	assert.False(t, <-server.listenCh)
}

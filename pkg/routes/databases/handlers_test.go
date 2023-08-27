package databases

import (
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/kareemmahlees/mysql-meta/pkg/db"
	"github.com/kareemmahlees/mysql-meta/utils"
	"github.com/stretchr/testify/assert"
)

var app *fiber.App



func TestMain(m *testing.M) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatal(err)
	}
	app = fiber.New()

	con,err :=db.InitDBConn()
	if err != nil {
		log.Fatal(err)
	}
	defer con.Close()

	RegisterDatabasesRoutes(app,con)

	// Run tests
	exitVal := m.Run()

	// Exit with exit value from tests
	os.Exit(exitVal)
}


func TestHandleListDatabases(t *testing.T){
	req := httptest.NewRequest("GET", "http://localhost:4000/databases", nil)

	resp, _ := app.Test(req)
	payload := utils.ReadBody(resp.Body)

	assert.Equal(t,resp.StatusCode,fiber.StatusOK)

	dbs,ok := payload["databases"]
	assert.True(t,ok)
	assert.Greater(t,len(dbs.([]interface{})),0)
}

func TestHandleCreateDatabase(t *testing.T){
	req := httptest.NewRequest("POST", "http://localhost:4000/databases/mysqlmeta", nil)

	resp, _ := app.Test(req)
	payload := utils.ReadBody(resp.Body)

	assert.Equal(t,resp.StatusCode,fiber.StatusCreated)
	var foo float64
	assert.IsType(t,foo,payload["created"])
	assert.Equal(t,int(payload["created"].(float64)),1)

}
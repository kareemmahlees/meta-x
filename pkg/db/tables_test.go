package db

import (
	"log"
	"reflect"
	"testing"

	"github.com/joho/godotenv"
	"github.com/kareemmahlees/mysql-meta/lib"
	"github.com/stretchr/testify/assert"
)

func init() {
	err := godotenv.Load("../../.env.test")
	if err != nil {
		log.Fatal(err)
	}
}

func TestDescribeTable(t *testing.T) {
	con, erro := InitDBConn()
	if erro != nil {
		t.Fatal(erro)
	}
	defer con.Close()

	// create testing table
	err := CreateTable(con, "testDescribeTable", map[string]lib.CreateTableProps{"name": {
		Type:     "varchar(255)",
		Default:  "defaultText",
		Unique:   true,
		Nullable: false,
	}})
	assert.Nil(t, err)

	result, _ := ListTables(con)
	assert.Greater(t, len(result), 0)
	assert.Contains(t, result, "testDescribeTable")

	tableInfo, err := GetTableInfo(con, "testDescribeTable")
	assert.Nil(t, err)
	assert.Greater(t, len(tableInfo), 0)
	assert.Equal(t, tableInfo[1].Type, "varchar(255)")
}

func TestListTable(t *testing.T) {
	con, erro := InitDBConn()
	if erro != nil {
		t.Fatal(erro)
	}
	defer con.Close()

	result, err := ListTables(con)
	assert.Nil(t, err)
	assert.IsType(t, reflect.Slice, reflect.TypeOf(result).Kind())
}

func TestCreateTable(t *testing.T) {
	con, erro := InitDBConn()
	if erro != nil {
		t.Fatal(erro)
	}
	defer con.Close()
	defer func() {
		err := DeleteTable(con, "testCreateTable")
		assert.Nil(t, err)
	}()

	// create testing table
	err := CreateTable(con, "testCreateTable", map[string]lib.CreateTableProps{"name": {
		Type:     "text",
		Default:  "defaultText",
		Unique:   true,
		Nullable: false,
	}})
	assert.Nil(t, err)

	result, _ := ListTables(con)
	assert.Greater(t, len(result), 0)
	assert.Contains(t, result, "testCreateTable")
}

func TestUpdateTable(t *testing.T) {
	con, erro := InitDBConn()
	if erro != nil {
		t.Fatal(erro)
	}
	defer con.Close()
	defer func() {
		err := DeleteTable(con, "testUpdateTable")
		assert.Nil(t, err)
	}()

	// create testing table
	err := CreateTable(con, "testUpdateTable", map[string]lib.CreateTableProps{"name": {
		Type:     "text",
		Default:  "defaultText",
		Unique:   true,
		Nullable: false,
	}})
	assert.Nil(t, err)

	updateTableProps := lib.UpdateTableProps{}
	updateTableProps.Opertaion.Type = "add"
	updateTableProps.Opertaion.Data = map[string]any{
		"age": "int",
	}
	err = UpdateTable(con, "testUpdateTable", updateTableProps)
	assert.Nil(t, err)

	updateTableProps.Opertaion.Type = "modify"
	updateTableProps.Opertaion.Data = map[string]any{
		"name": "varchar(55)",
	}
	err = UpdateTable(con, "testUpdateTable", updateTableProps)
	assert.Nil(t, err)

	updateTableProps.Opertaion.Type = "delete"
	updateTableProps.Opertaion.Data = []interface{}{"name", "age"}
	err = UpdateTable(con, "testUpdateTable", updateTableProps)
	assert.Nil(t, err)
}

func TestDeleteTable(t *testing.T) {
	con, erro := InitDBConn()
	if erro != nil {
		t.Fatal(erro)
	}
	defer con.Close()

	// create testing table
	err := CreateTable(con, "testDeleteTable", map[string]lib.CreateTableProps{"name": {
		Type:     "text",
		Default:  "defaultText",
		Unique:   true,
		Nullable: false,
	}})
	assert.Nil(t, err)

	err = DeleteTable(con, "testDeleteTable")
	assert.Nil(t, err)

	result, _ := ListTables(con)
	assert.NotContains(t, result, "testDeleteTable")
}

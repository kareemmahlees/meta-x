//go:generate go run cmd/gen/main.go

package main

import (
	"github.com/kareemmahlees/meta-x/cmd"
	_ "github.com/kareemmahlees/meta-x/docs"
)

// Main
//
//	@title					MetaX
//	@version				0.1.1
//	@description			A RESTFull and GraphQL API to supercharge your database
//	@contact.name			Kareem Ebrahim
//	@contact.email			kareemmahlees@gmail.com
//
//	@tag.name				General
//	@tag.description		general info about the API.
//	@tag.name				Database
//	@tag.description		Database level related operations.
//	@tag.name				Table
//	@tag.description		Table Level related operations.
//
//	@servers.url			http://localhost:5522
//	@servers.description	Home town of Meta-X
func main() {
	cmd.Execute()
}

package main

import (
	"github.com/kareemmahlees/meta-x/cmd"
	_ "github.com/kareemmahlees/meta-x/docs"
)

// Main
//
//	@title			MetaX
//	@version		0.1.1
//	@description	A RESTFull and GraphQL API to supercharge your database
//	@contact.name	Kareem Ebrahim
//	@contact.email	kareemmahlees@gmail.com
//	@host			localhost:5522
//	@BasePath		/
func main() {
	cmd.Execute()
}

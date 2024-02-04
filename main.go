package main

import (
	"meta-x/cmd"
	_ "meta-x/docs"
)

// Main
//
//	@title			MySQL Meta
//	@version		1.0
//	@description	A RESTFull and GraphQL API to manage your MySQL DB
//	@contact.name	Kareem Ebrahim
//	@contact.email	kareemmahlees@gmail.com
//	@host			localhost:5522
//	@BasePath		/
func main() {
	cmd.Execute()
}

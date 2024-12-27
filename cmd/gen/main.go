//go:build ignore
// +build ignore

package main

import (
	"log"

	"github.com/swaggo/swag/v2/format"
	"github.com/swaggo/swag/v2/gen"
)

func main() {
	log.Println("Formatting swagger comments...")
	err := format.New().Build(&format.Config{
		SearchDir: "./",
		Excludes:  "cmd,.vscode,.github,docs",
	})
	if err != nil {
		log.Fatalf("Error formatting swagger comments: %v", err)
	}

	config := &gen.Config{
		OutputDir:           "./docs",
		MainAPIFile:         "./main.go",
		SearchDir:           "./",
		PropNamingStrategy:  "camelcase",
		GenerateOpenAPI3Doc: true,
	}

	err = gen.New().Build(config)
	if err != nil {
		log.Fatalf("Error generating swagger docs: %v", err)
	}

	log.Println("Swagger documentation generated successfully!")
}

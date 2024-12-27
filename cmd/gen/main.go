package main

import (
	"errors"
	"io/fs"
	"log"

	graphqlAPI "github.com/99designs/gqlgen/api"
	graphqlConfig "github.com/99designs/gqlgen/codegen/config"
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

	cfg, err := graphqlConfig.LoadConfigFromDefaultLocations()
	if errors.Is(err, fs.ErrNotExist) {
		cfg, err = graphqlConfig.LoadDefaultConfig()
	}
	if err != nil {
		log.Fatalf("Error while loading GraphQL config: %v", err)
		return
	}

	log.Println("Generating GraphQL code...")

	err = graphqlAPI.Generate(cfg)
	if err != nil {
		log.Fatalf("Error while generating GraphQL code: %v", err)
		return
	}

	log.Println("Successfully generated GraphQL code!")
}

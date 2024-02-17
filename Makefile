build:
	@go build -o bin/
run: build
	@./bin/meta-x.exe --help

# make sure to run docker first
test:
	@go test ./... -race
swag:
	@swag fmt
	@swag init 
graphql:
	@go run github.com/99designs/gqlgen generate
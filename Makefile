build:
	@go build -o bin/ -ldflags "-w -s"
run: build
	@./bin/meta-x.exe --help

# WARNING: make sure to run docker first
test:
	@go test ./... -race
swag:
	@swag fmt
	@swag init 
graphql:
	@go run github.com/99designs/gqlgen generate
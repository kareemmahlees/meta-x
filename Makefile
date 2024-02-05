build:
	@go build -o bin/
run: build
	@./bin/meta-x.exe
testv:
	@go test ./... -race
swag:
	@swag fmt
	@swag init 
graphql:
	@go run github.com/99designs/gqlgen generate
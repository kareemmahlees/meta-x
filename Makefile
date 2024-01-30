build:
	@go build -o bin/
run: build
	@./bin/mysql-meta.exe
watch:
	@docker compose up -d mysql-meta
	@air
testv:
	@go test -v ./...
setup_test:
	@docker compose up -d test
test:
	-@go test ./...
cleanup_test:
	@docker compose down
swag:
	@swag fmt
	@swag init 
generate:
	@go run github.com/99designs/gqlgen generate
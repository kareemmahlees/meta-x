build:
	@go build -o bin/
run: build
	@./bin/mysql-meta.exe
watch:
	@docker compose up -d mysql-meta
	@air
testv:
	@go test -v ./...
test:
	@docker compose up -d test
	-@go test ./...
	docker compose down
swag:
	@swag fmt
	@swag init
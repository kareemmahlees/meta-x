build:
	@go build -o bin/
run: build
	@./bin/mysql-meta.exe
watch:
	@air
testv:
	@go test -v ./...
test:
	@go test ./...
swag:
	@swag fmt
	@swag init
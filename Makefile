run:
	@go run main.go

watch:
	@air

testv:
	@go test -v ./...
test:
	@go test ./...
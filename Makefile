run:
	@go run main.go

watch:
	@air

testv:
	@go test -v ./pkg/...
test:
	@go test ./pkg/...
build:
	@go build -o bin/ 
build_prod:
	@go build -o bin/ -ldflags "-w -s"
run: generate
	@air

# WARNING: make sure to run docker first
test:
	@go test ./... -race
generate:
	@go generate
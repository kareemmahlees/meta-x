set shell := ["powershell.exe", "-c"]

build:
	@go build -o bin/ 
build_prod:
	@go build -o bin/ -ldflags "-w -s"
run:
	@air

# WARNING: make sure to run docker first
test:
	@go test ./... -race
generate_graphql:
	@go tool gqlgen generate
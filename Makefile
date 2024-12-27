build_dev:
	@go build -o bin/ 
build_prod:
	@go build -o bin/ -ldflags "-w -s"
run: generate build_dev
	@./bin/meta-x --help

# WARNING: make sure to run docker first
test:
	@go test ./... -race
generate:
	@go generate
# graphql:
#	 @go run github.com/99designs/gqlgen generate
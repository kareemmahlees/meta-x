build_dev:
	@go build -o bin/ -tags dev ./cmd/dev 
build_prod:
	@go build ./cmd/prod -o bin/ -ldflags "-w -s" -tags prod
run: build_dev
	@./bin/meta-x.exe --help

# WARNING: make sure to run docker first
test:
	@go test ./... -race
swag:
	@swag fmt
	@swag init 
graphql:
	@go run github.com/99designs/gqlgen generate
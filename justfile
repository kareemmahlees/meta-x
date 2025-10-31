set shell := ["powershell.exe", "-c"]

# Produce a development binary.
build_dev:
	@go build -o bin/ 
# Produce a production binary.
build_prod:
	@go build -o bin/ -ldflags "-w -s"
# Run the app in watch mode.
run:
	@air
# Run the entire test suite with race detection enabled.
test:
	@go test ./... -race
# Generate graphql resolvers and types.
generate_graphql:
	@go tool gqlgen generate
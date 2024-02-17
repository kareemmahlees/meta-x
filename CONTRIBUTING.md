Thank you for being interested in contributing to MetaX ✨

## Getting Started

Start by cloning the repo:

```shell
git clone --depth=1 https://github.com/kareemmahlees/mysql-meta.git
```

the `--depth` flag specifies the number of commits you want to clone with the repo.

You will find in the [Makefile](Makefile) some useful scripts:

- `build`: builds the code and produce a binary.
- `run`: builds the code and runs it.
- `test`: runs the test suite, make sure to startup docker first because we use **TestContainers**.
- `swag`: will generate the required swagger documentation.
- `generate`: will generate the required graphql code, e.g resolvers.

## Codebase overview

### CLI

**MetaX** uses [cobra](https://github.com/spf13/cobra) to parse CLI arguments, commands, etc.
All **MetaX's** commands are located under the `cmd` directory.

### REST API

We use [fiber](https://gofiber.io) as our framework. All the routes are located under the `internal/rest` directory.

### GraphQL

We use [gqlgen](https://gqlgen.com/) as our graphql server. All GraphQL related stuff are located under `internal/graph`.

Each time you make a change in the schema, make sure you run `make generate` to generate the resolvers and models.

### Swagger

We use [swag](https://github.com/swaggo/swag) to generate swagger documentation for the REST endpoints.

Each time you make a change in the endpoints documentation, make sure your run `make swag` to formate the comments and generate the required data.

### Other

#### /lib

Contains logic related to validation, constants, and http errors.

#### /models

Contains structs that define the shape of requests and responses.

#### /utils

Some utilities for encoding/decoding http payloads, test helpers and styled printing

## Testing

We use [TestContainers](https://testcontainers.com/) to create and terminate docker containers in the `setup`&`teardown` of our test suites.

So if you want to run the tests after making any change, make sure to startup docker first.

Please make sure that any change you make is accompanied with sufficient test cases.

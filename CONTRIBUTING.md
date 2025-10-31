Thank you for being interested in contributing to MetaX ✨

## Getting Started

Start by cloning the repo:

```shell
git clone --depth=1 https://github.com/kareemmahlees/mysql-meta.git
```

> [!TIP]
>
> the `--depth` flag specifies the number of commits you want to clone with the repo.

In the [justfile](justfile) you will find all the scripts you would need regarding building, running or testing the app.

You can list all the available commands by running:

```shell
just --list
```

## Codebase overview

### CLI

**MetaX** uses [cobra](https://github.com/spf13/cobra) to parse CLI arguments, commands, etc.

All **MetaX's** commands are located under the `cmd` directory.

### REST API

We use [Huma](https://huma.rocks) as our framework. All the routes are located under the `internal/handlers` directory.

### GraphQL

We use [gqlgen](https://gqlgen.com/) as our graphql server. All GraphQL related stuff are located under `internal/graph`.

Each time you make a change in the schema, make sure you run `just generate_graphql` to generate the resolvers and models.

### Other

#### /lib

Contains logic related to validation, constants, and http errors.

#### /models

Contains structs that define the shape of requests and responses.

#### /utils

Some utilities for encoding/decoding http payloads, test helpers and styled printing

## Testing

We use [TestContainers](https://testcontainers.com/) to create and terminate docker containers in the `setup`&`teardown` of our test suites.

> [!IMPORTANT]
>
> Make sure to start the docker daemon before running the tests.

Please make sure that any change you make is accompanied with sufficient test cases.

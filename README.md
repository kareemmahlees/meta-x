<h1 align='center'>
    MySQL Meta
</h1>
<p align='center'>
    RESTfull & GraphQL API to manage your MySQL Database
</p>

<div align='center'>

![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/kareemmahlees/mysql-meta)
![Codecov](https://img.shields.io/codecov/c/github/kareemmahlees/mysql-meta)
![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/kareemmahlees/mysql-meta/lint.yml)

</div>

## What does it introduce

MySQL Meta introduces a RESTFull and GraphQL API that allows you to control your MySQL database through http requests, i.e creating tables and databases, adding columns and deleting columns.

![Screen shot of running application](./docs/screenshot.png)

## Installation

- ### Download the binary

  You can download the binary from the [releases page](https://github.com/kareemmahlees/mysql-meta/releases).

  Make sure you add the executable to your `PATH` environment variable.

- ### Go Install

  You can use the `go install` command to install like so:

  ```shell
    go install github.com/kareemmahlees/mysql-meta
  ```

## Running

#### Environment Variables Setup

Before starting the application, MySQL Meta expects some env vars to be set, you can find them also in the [.env.example](.env.example) file

Once you [installed](#installation) the application you can run it like so:

```shell
mysql-meta
```

MySQL Meta by default serves on port 4000, you can configure the port by passing the `--port` flag:

```shell
mysql-meta --port 4444
```

## Documentation

The API is fully documented, the REST version is documented using **Swagger Docs** and is served on `http://localhost:4000/swagger`

Regarding the GraphQL version, you can run the application and then use the GraphQL endpoint `http://localhost:4000/graphql` to introspect the schema with your favorite tool, e.g postman, insomnia, hoppscotch.

Additionally, you can playaround with the GraphQL version by jumping into the playground at `http://localhost:4000/playground`

## With Postman

MySQL Meta is well integrated with postman, you can import the collection along side it's documentation from the REST specification in the [postman collection folder](postman/collections/)

## Progress Track

- Databases
  - [x] list databases
  - [x] create database
- Tables
  - [x] list tables
  - [x] table info
  - [x] create table
  - [x] delete table
  - [x] update table
- Queries
  - [ ] execute single query
  - [ ] execute multiple queries in transaction
- Views
  - [ ] list views
  - [ ] create views
  - [ ] delete views
  - [ ] query by views
- Config
  - [ ] get mysql version

## Contributing

We strongly encourage anyone who wants to contribute to go ahead, not matter what skill level your are.

Contributions can be as small as suggesting a feature, reporting a bug or enhancing the docs.

For more details, please visit [CONTRIBUTING.md]()

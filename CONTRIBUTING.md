Thank you for being interested in contributing to MySQL Meta ✨

## Local Development

Start by cloning the repo:

```shell
git clone --depth=1 https://github.com/kareemmahlees/mysql-meta.git
```

the `--depth` flag specifies the number of commits you want to clone with the repo.

You will find in the [Makefile](Makefile) some useful scripts:

- watch: will start the database container and run the application in watch mode
- swag: will generate the required swagger documentation
- generate: will generate the required graphql code, e.g resolvers.

## Testing

You can setup the testing containers by running the following:

```shell
make setup_test
```

then run the tests:

```shell
make test
```

or, for verbose output:

```shell
make testv
```

after your are done, remember to cleanup afterwards:

```shell
make cleanup_test
```

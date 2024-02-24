# golang-sqlc-use-example

Example of using `sqlc` library with the Go.

## sqlc

"**sqlc** generates fully type-safe idiomatic Go code from SQL."

Docs: https://docs.sqlc.dev/en/latest/#

### commands

## golang-migrate

**Migrate** reads migrations from sources and applies them in correct order to
a database.

Docs: https://github.com/golang-migrate/migrate

### commands:

```shell
# create a migration
migrate create -ext=sql -dir=sql/migrations -seq init

# execute migrations
migrate -path=sql/migrations -database "mysql://root:root@tcp(localhost:3306)/courses" -verbose up

# restore migrations
migrate -path=sql/migrations -database "mysql://root:root@tcp(localhost:3306)/courses" -verbose down
```
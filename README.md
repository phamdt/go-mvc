# GO-MVC

*** **Disclaimer** ***
There are no guarantees of API stability until the first major version is released. Think of the current state as an alpha-phase product looking for validation from users about what features make sense and what the best UX is for those features.

## Dependencies
[goimports](https://godoc.org/golang.org/x/tools/cmd/goimports)

## Installation

- set up your GOPATH
- make sure your regular PATH includes the go binary folder e.g. `/whereveryouinstalled/go/bin`
- `go get golang.org/x/tools/cmd/goimports`
- `go get -u github.com/go-generation/go-mvc/cmd/gomvc`

## Development Workflow

In its current state, using gomvc assumes you will do some combination of the following:
1) use the `gomvc application {{yourAppName}}` command and be done with gomvc
2) use the `gomvc oa` command in the directory where your openapi.yml file is
3) use one of the two previous commands and then also use `gomvc resource {{someRESTfulResource}}` subsequently

As of now, gomvc assumes you will want to use the following dependencies:
- gin
- postgres
- https://github.com/golang-migrate/migrate (will need to be installed by you)
- https://github.com/volatiletech/sqlboiler


### Example steps: using sqlboiler
1. Generate application: `gomvc application api --dest ~/Code/api`
1. Design your SQL schema or at least one table
1. Create a migration per table: `migrate create -ext sql -dir migrations -seq create_users_table`
1. Fill the migration file with your SQL commands e.g. `CREATE USERS (etc...)`
1. Migrate your local database: `migrate -database YOUR_DATBASE_URL -path PATH_TO_YOUR_MIGRATIONS up`
1. Install sqlboiler by running `make dev-dependencies`
1. Create a sqlboiler configuration per https://github.com/volatiletech/sqlboiler#getting-started (will be automatically done for you in future GO-MVC release)
1. Start your application so that you have postgres running: `docker-compose up`
1. Run sqlboiler: `sqlboiler psql`
1. Run `gomvc resource {{tableName}}` for each of the tables you want to create an endpoint for.
1. Verify naming conventions align with what sqlboiler generated. You might need to edit the generated controllers.
1. Generate seeder to fill DB with test data (assumes models dir exists in your app directory): `gomvc seed`.
1. Continue dev-ing as you would normally

If you're managing your schema independently, you can completely remove the migrate dependency from both your workflow and the app but you can still use sqlboiler regardless.

### Example steps: generating application with OpenAPI 3.x spec file
1. `gomvc application petstore`
1. `gomvc oa --spec ./petstore.yaml`

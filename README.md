# GO-MVC
![](https://github.com/go-generation/go-mvc/workflows/test/badge.svg?branch=main)

*** **Disclaimer** ***
There are no guarantees of API stability until the first major version is released. Think of the current state as an alpha-phase product looking for validation from users about what features make sense and what the best UX is for those features.

## Installation

- set up your GOPATH
- make sure your regular PATH includes the go binary folder e.g. `/whereveryouinstalled/go/bin`
- download dependencies
```
go get golang.org/x/tools/cmd/goimports
```
- download cli
```
go get -u github.com/go-generation/go-mvc/cmd/gomvc
```

## Commands

  `gomvc application` Generate application files
  `gomvc help`        Help about any command
  `gomvc oa`          Generate controllers from an OpenAPI yml file
  `gomvc resource`    Generate resource files
  `gomvc seed`        Generate seed files
  `gomvc swagger`     Generate controllers from a v2 Swagger yml file

## Usage

Create an application:
```
$ mkdir yourapplication && cd yourapplication
$ gomvc application yourapplication
```

Optionally create controllers from your OpenAPI or Swagger 2.0 yaml spec (json support is coming):
```
$ gomvc oa --spec path/to/openapi.yml
```
or

```
$ gomvc swagger --spec path/to/swagger.yml
```

Create more resources:
```
$ gomvc resource dogs
```

## Swagger vs OpenAPI Disclaimer
For ease of parsing, the `gomvc swagger` command currently converts the Swagger spec into an OpenAPI 3 spec. This is done by an underlying library gomvc uses which attempts to convert backwards compatible aspects of the spec. As such there may be some side effects or ignored properties. For example, a Swagger spec lists its reused parameters in a top level section called "parameters". In OA3, everything is nested under the "components" section. When resolving object/schema refs, you may find missing elements because it will try to look for parameters under "#/components/parameters" and not #/parameters". 

## Generated Application Dependencies
As of now, gomvc assumes you will want to use the following dependencies:
- gin
- postgres
- zap
- new relic agent
- https://github.com/golang-migrate/migrate (will need to be installed by you)
- https://github.com/volatiletech/sqlboiler


### Example steps: using sqlboiler
1. Generate application: `gomvc application petstore --dest ~/Code/petstore`
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


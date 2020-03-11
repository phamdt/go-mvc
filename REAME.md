# GO-MVC

## Development Workflow

In its current state, using gomvc assumes you will do the following:
1) use the `gomvc application {{yourAppName}}` command and be done with gomvc
2) use the above command and then also use `gomvc resource {{someRESTfulResource}}` subsequently

As of now, using the first command assumes you will want to use the following dependencies:
- gin
- postgres
- https://github.com/golang-migrate/migrate (will need to be installed by you)


As of now, using the second command assumes you will use the following dependencies as part of your workflow: 
- https://github.com/volatiletech/sqlboiler


An example workflow after creating the application:

1) Design your SQL schema or at least one table
2) Create a migration per table: `migrate create -ext sql -dir migrations -seq create_users_table`
3) Fill the migration file with your SQL commands e.g. `CREATE USERS (etc...)`
4) Migrate your local database: `migrate -database YOUR_DATBASE_URL -path PATH_TO_YOUR_MIGRATIONS up`
5) Install sqlboiler by running `make dev-dependencies`
6) Create a sqlboiler configuration per https://github.com/volatiletech/sqlboiler#getting-started (will be automatically done for you in future GO-MVC release)
7) Run sqlboiler: `sqlboiler psql`
8) Run `gomvc resource {{tableName}}` for each of the tables you want to create an endpoint for.
9) Verify naming conventions align with what sqlboiler generated. You might need to edit the generated controllers.
10) Continuing dev-ing as you would normally

If you're managing your schema independently, you can completely remove the migrate dependency from both your workflow and the app but you can still use sqlboiler regardless.

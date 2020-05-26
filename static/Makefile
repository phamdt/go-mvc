.PHONY: models

# Go parameters
GOBUILD=go build
GOCLEAN=go clean
GOTEST=go test
GOGET=go get

all: test build

dev-dependencies:
	go get -u -t github.com/volatiletech/sqlboiler
	go get github.com/volatiletech/sqlboiler/drivers/sqlboiler-psql

build: 
	$(GOBUILD) -tags=jsoniter .

test: 
	$(GOTEST) -v ./...

start:
	make build
	go run main.go

# usage: make migration N=tableName
migration:
	migrate create -ext sql -dir ./migrations -seq $(N)

migratedb:
	migrate up

dropdb:
	migrate drop

models:
	sqlboiler psql --no-tests --no-hooks --no-context

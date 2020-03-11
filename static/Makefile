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
	go build .
	go run main.go
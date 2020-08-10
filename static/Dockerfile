FROM golang:1.13-alpine

WORKDIR /app

ADD go.sum .
ADD go.mod .
RUN go mod vendor

ADD . /app

CMD go run main.go

FROM golang:1.13

RUN go get -u github.com/smartystreets/goconvey

ADD . /app
WORKDIR /app

RUN go install -v

CMD goconvey -host 0.0.0.0 -port=9999

EXPOSE 9999
FROM golang:1.19.4-bullseye

ENV DEBIAN_FRONTEND=noninteractive

RUN apt update && apt install -y git curl jq
    
RUN go get github.com/nsf/jsondiff \
    go get github.com/mattn/go-sqlite3 \
    go get -u github.com/go-sql-driver/mysql \
    go get -u go.uber.org/multierr

WORKDIR /work
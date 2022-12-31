FROM golang:1.19.4-bullseye

ENV DEBIAN_FRONTEND=noninteractive

RUN apt update && apt install -y git curl jq
    
RUN go get -u github.com/nsf/jsondiff \
    go get -u github.com/mattn/go-sqlite3 \
    go get -u github.com/go-sql-driver/mysql \
    go get -u github.com/lib/pq \
    go get -u github.com/docopt/docopt-go \
    go get -u go.uber.org/multierr

WORKDIR /work
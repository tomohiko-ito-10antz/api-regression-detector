FROM golang:1.19.4-bullseye

ENV DEBIAN_FRONTEND=noninteractive

RUN apt update && apt install -y git curl jq
    
RUN go get github.com/nsf/jsondiff
WORKDIR /work
FROM golang:1.20.0-bullseye

ENV DEBIAN_FRONTEND=noninteractive

RUN apt update && apt install -y git curl

WORKDIR /work/server

# go tools
RUN go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2 && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway/v2 && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 
    GO111MODULE=on GOBIN=/usr/local/bin go install github.com/bufbuild/buf/cmd/buf@v1.14.0
# go modules
COPY go.mod /work/go.mod
RUN go mod download


CMD ["make", "run"]
FROM golang:1.20.0-bullseye

ENV DEBIAN_FRONTEND=noninteractive

RUN apt update && apt install -y git curl jq

WORKDIR /work/server

COPY . /work/server/

# go modules
RUN go mod download

# go tools
RUN go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2 && \
    go install github.com/bufbuild/buf/cmd/buf@v1.14.0
# && \
#    go get github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor@v2.15.0 && go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway

CMD ["make", "serve"]
# API Regression Detector server

## Overview

Start GRPC server on port 50051 and GRPC gateway on port 80 by the following command:

```sh
make serve
```

The GRPC server serves a GRPC API of GreetingService having SayHello method.

The GRPC gateway, which translates HTTP and GRPC, behaves as REST API version of the GRPC API.

## Example

```sh
make serve &

ENDPOINT=localhost

curl "${ENDPOINT}:80/say/hello/MyName" | jq

grpcurl -plaintext -d '{"name":"MyName"}' "${ENDPOINT}:50051" api.GreetingService/SayHello | jq
```

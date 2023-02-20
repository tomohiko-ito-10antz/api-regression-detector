#!/bin/sh

set -eux

ENDPOINT='api:50051'
METHOD='api.GreetingService/SayHello'
REQUEST='test/data/call/grpc-success/request.json'
ACTUAL_RESPONSE='test/data/call/grpc-success/actual.json'
EXPECTED_RESPONSE='test/data/call/grpc-success/expected.json'

go run cmd/call/call-grpc/main.go "${ENDPOINT}" "${METHOD}" < "${REQUEST}" > "${ACTUAL_RESPONSE}"

go run cmd/compare/main.go "${EXPECTED_RESPONSE}" "${ACTUAL_RESPONSE}"
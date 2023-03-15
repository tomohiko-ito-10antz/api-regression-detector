#!/bin/sh

set -eux

ENDPOINT='api:50051'
METHOD='api.GreetingService/Error'
REQUEST='test/data/call/grpc-error/request.json'
ACTUAL_RESPONSE='test/data/call/grpc-error/actual.json'
EXPECTED_RESPONSE='test/data/call/grpc-error/expected.json'

go run cmd/call-grpc/main.go "${ENDPOINT}" "${METHOD}" < "${REQUEST}" > "${ACTUAL_RESPONSE}"

go run cmd/compare/main.go "${EXPECTED_RESPONSE}" "${ACTUAL_RESPONSE}"
#!/bin/sh

set -eux

ENDPOINT='api:50051'
METHOD='api.GreetingService/GetHello'
REQUEST='test/data/call/get/request.json'
ACTUAL_RESPONSE='test/data/call/get/actual.json'
EXPECTED_RESPONSE='test/data/call/get/expected.json'

go run cmd/call-grpc/main.go "${ENDPOINT}" "${METHOD}" < "${REQUEST}" > "${ACTUAL_RESPONSE}"

go run cmd/compare/main.go "${EXPECTED_RESPONSE}" "${ACTUAL_RESPONSE}"
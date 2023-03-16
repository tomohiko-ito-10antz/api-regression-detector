#!/bin/sh

set -eux

ENDPOINT='api:50051'
METHOD='api.GreetingService/PostHello'
REQUEST='test/data/call/post/request.json'
ACTUAL_RESPONSE='test/data/call/post/actual.json'
EXPECTED_RESPONSE='test/data/call/post/expected.json'

go run cmd/call-grpc/main.go "${ENDPOINT}" "${METHOD}" < "${REQUEST}" > "${ACTUAL_RESPONSE}"

go run cmd/compare/main.go "${EXPECTED_RESPONSE}" "${ACTUAL_RESPONSE}"
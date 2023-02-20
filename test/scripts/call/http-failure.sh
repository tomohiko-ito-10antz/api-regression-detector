#!/bin/sh

set -eux

ENDPOINT='http://api:80/error'
METHOD='GET'
REQUEST='test/data/call/http-failure/request.json'
ACTUAL_RESPONSE='test/data/call/http-failure/actual.json'
EXPECTED_RESPONSE='test/data/call/http-failure/expected.json'

go run cmd/call/call-http/main.go "${ENDPOINT}" "${METHOD}" < "${REQUEST}" > "${ACTUAL_RESPONSE}"

go run cmd/compare/main.go "${EXPECTED_RESPONSE}" "${ACTUAL_RESPONSE}"
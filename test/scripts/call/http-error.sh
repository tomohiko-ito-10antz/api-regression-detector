#!/bin/sh

set -eux

ENDPOINT='http://api:80/error'
METHOD='GET'
REQUEST='test/data/call/error/request.json'
ACTUAL_RESPONSE='test/data/call/error/actual.json'
EXPECTED_RESPONSE='test/data/call/error/expected.json'

go run cmd/call-http/main.go "${ENDPOINT}" "${METHOD}" < "${REQUEST}" > "${ACTUAL_RESPONSE}"

go run cmd/compare/main.go "${EXPECTED_RESPONSE}" "${ACTUAL_RESPONSE}"
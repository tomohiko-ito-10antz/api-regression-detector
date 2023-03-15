#!/bin/sh

set -eux

ENDPOINT='http://api:80/hello/[name]'
METHOD='PATCH'
REQUEST='test/data/call/patch/request.json'
ACTUAL_RESPONSE='test/data/call/patch/actual.json'
EXPECTED_RESPONSE='test/data/call/patch/expected.json'

go run cmd/call-http/main.go "${ENDPOINT}" "${METHOD}" < "${REQUEST}" > "${ACTUAL_RESPONSE}"

go run cmd/compare/main.go "${EXPECTED_RESPONSE}" "${ACTUAL_RESPONSE}"
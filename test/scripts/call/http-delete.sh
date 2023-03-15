#!/bin/sh

set -eux

ENDPOINT='http://api:80/hello/[name]'
METHOD='DELETE'
REQUEST='test/data/call/http-delete/request.json'
ACTUAL_RESPONSE='test/data/call/http-delete/actual.json'
EXPECTED_RESPONSE='test/data/call/http-delete/expected.json'

go run cmd/call-http/main.go "${ENDPOINT}" "${METHOD}" < "${REQUEST}" > "${ACTUAL_RESPONSE}"

go run cmd/compare/main.go "${EXPECTED_RESPONSE}" "${ACTUAL_RESPONSE}"
#!/bin/sh

set -eux

ENDPOINT='http://api:80/hello/[name]'
METHOD='POST'
REQUEST='test/data/call/post/request.json'
ACTUAL_RESPONSE='test/data/call/post/actual.json'
EXPECTED_RESPONSE='test/data/call/post/expected.json'

go run cmd/call-http/main.go "${ENDPOINT}" "${METHOD}" < "${REQUEST}" > "${ACTUAL_RESPONSE}"

go run cmd/compare/main.go "${EXPECTED_RESPONSE}" "${ACTUAL_RESPONSE}"
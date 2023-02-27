#!/bin/sh

set -eux

DRIVER='sqlite3'
CONNECT='file:examples/sqlite/sqlite.db'
INIT_TABLES='test/data/db/sqlite/init.json'
DUMP_TABLES='test/data/db/sqlite/dump.json'
EXPECTED_TABLES='test/data/db/sqlite/expected.json'

go run cmd/db-init/main.go "${DRIVER}" "${CONNECT}" < "${INIT_TABLES}"
jq '. | keys' <  "${EXPECTED_TABLES}" \
	| go run cmd/db-dump/main.go "${DRIVER}" "${CONNECT}" > "${DUMP_TABLES}"

go run cmd/compare/main.go "${EXPECTED_TABLES}" "${DUMP_TABLES}"
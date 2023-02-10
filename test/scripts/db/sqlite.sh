#!/bin/sh

set -eux

DRIVER='sqlite3'
CONNECT='file:examples/sqlite/sqlite.db'
INIT_TABLES='test/data/db/sqlite/init.json'
DUMP_TABLES='test/data/db/sqlite/dump.json'
EXPECTED_TABLES='test/data/db/sqlite/expected.json'

make init-sqlite
go run main.go init "${DRIVER}" "${CONNECT}" < "${INIT_TABLES}"
jq '. | keys' <  "${EXPECTED_TABLES}" \
	| go run main.go dump "${DRIVER}" "${CONNECT}" > "${DUMP_TABLES}"

go run main.go compare "${EXPECTED_TABLES}" "${DUMP_TABLES}"
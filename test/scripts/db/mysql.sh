#!/bin/sh

set -eux

DRIVER='mysql'
CONNECT='root:password@(mysql)/main'
INIT_TABLES='test/data/db/mysql/init.json'
DUMP_TABLES='test/data/db/mysql/dump.json'
EXPECTED_TABLES='test/data/db/mysql/expected.json'

go run cmd/db/db-init/main.go "${DRIVER}" "${CONNECT}" < "${INIT_TABLES}"
jq '. | keys' <  "${EXPECTED_TABLES}" \
	| go run cmd/db/db-dump/main.go "${DRIVER}" "${CONNECT}" > "${DUMP_TABLES}"

go run cmd/compare/main.go "${EXPECTED_TABLES}" "${DUMP_TABLES}"
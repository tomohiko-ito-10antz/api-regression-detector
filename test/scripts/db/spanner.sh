#!/bin/sh

set -eux

DRIVER='spanner'
CONNECT='projects/regression-detector/instances/example/databases/main'
INIT_TABLES='test/data/db/spanner/init.json'
DUMP_TABLES='test/data/db/spanner/dump.json'
EXPECTED_TABLES='test/data/db/spanner/expected.json'

go run cmd/db-init/main.go "${DRIVER}" "${CONNECT}" < "${INIT_TABLES}"
jq '. | keys' <  "${EXPECTED_TABLES}" \
	| go run cmd/db-dump/main.go "${DRIVER}" "${CONNECT}" > "${DUMP_TABLES}"

go run cmd/compare/main.go "${EXPECTED_TABLES}" "${DUMP_TABLES}"
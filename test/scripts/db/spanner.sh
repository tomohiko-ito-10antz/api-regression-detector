#!/bin/sh

set -eux

DRIVER='spanner'
CONNECT='projects/regression-detector/instances/example/databases/main'
INIT_TABLES='test/data/db/spanner/init.json'
DUMP_TABLES='test/data/db/spanner/dump.json'
EXPECTED_TABLES='test/data/db/spanner/expected.json'

go run main.go init "${DRIVER}" "${CONNECT}" < "${INIT_TABLES}"
jq '. | keys' <  "${EXPECTED_TABLES}" \
	| go run main.go dump "${DRIVER}" "${CONNECT}" > "${DUMP_TABLES}"

go run main.go compare "${EXPECTED_TABLES}" "${DUMP_TABLES}"
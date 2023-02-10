#!/bin/sh

set -eux

DRIVER='mysql'
CONNECT='root:password@(mysql)/main'
INIT_TABLES='test/data/db/mysql/init.json'
DUMP_TABLES='test/data/db/mysql/dump.json'
EXPECTED_TABLES='test/data/db/mysql/expected.json'

go run main.go init "${DRIVER}" "${CONNECT}" < "${INIT_TABLES}"
jq '. | keys' <  "${EXPECTED_TABLES}" \
	| go run main.go dump "${DRIVER}" "${CONNECT}" > "${DUMP_TABLES}"

go run main.go compare "${EXPECTED_TABLES}" "${DUMP_TABLES}"
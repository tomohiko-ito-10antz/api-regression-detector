#!/bin/sh

set -eux

DRIVER='postgres'
CONNECT='user=root password=password host=postgres dbname=main sslmode=disable'
INIT_TABLES='test/data/db/postgres/init.json'
DUMP_TABLES='test/data/db/postgres/dump.json'
EXPECTED_TABLES='test/data/db/postgres/expected.json'

go run cmd/db-init/main.go "${DRIVER}" "${CONNECT}" < "${INIT_TABLES}"
jq '. | keys' <  "${EXPECTED_TABLES}" \
	| go run cmd/db-dump/main.go "${DRIVER}" "${CONNECT}" > "${DUMP_TABLES}"

go run cmd/compare/main.go "${EXPECTED_TABLES}" "${DUMP_TABLES}"
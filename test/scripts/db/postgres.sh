#!/bin/sh

set -eux

DRIVER='postgres'
CONNECT='user=root password=password host=postgres dbname=main sslmode=disable'
INIT_TABLES='test/data/db/postgres/init.json'
DUMP_TABLES='test/data/db/postgres/dump.json'
EXPECTED_TABLES='test/data/db/postgres/expected.json'

go run main.go init "${DRIVER}" "${CONNECT}" < "${INIT_TABLES}"
jq '. | keys' <  "${EXPECTED_TABLES}" \
	| go run main.go dump "${DRIVER}" "${CONNECT}" > "${DUMP_TABLES}"

go run main.go compare "${EXPECTED_TABLES}" "${DUMP_TABLES}"
.DEFAULT_GOAL := help
.PHONY: help
help: ## Show help
	@grep -E '^[a-zA-Z_-]+:.*?##.*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?##"}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'


.PHONY: init-mysql
init-mysql: ## Initialize MySQL database for develop
	mysql --host=mysql --password=password main < examples/mysql/create.sql

.PHONY: init-postgres
init-postgres: ## Initialize PostgreSQL database for develop
	psql --host=postgres --username=root --dbname=main < examples/postgres/create.sql

.PHONY: init-spanner
init-spanner: ## Initialize Spanner emulator database for develop
	gcloud config set project regression-detector
	gcloud config set auth/disable_credentials true
	gcloud config set api_endpoint_overrides/spanner http://spanner:9020/
	gcloud spanner instances describe example || gcloud spanner instances create example --config=emulator-config --description="Instance for example using spanner"
	gcloud spanner databases describe main --instance=example \
		&& gcloud spanner databases delete main --instance=example \
		|| true
	gcloud spanner databases create main --instance=example
	spanner-cli -p regression-detector -i example -d main --file=examples/spanner/create.sql

.PHONY: init-sqlite
init-sqlite: ## Initialize SQLite3 database for develop
	sqlite3 examples/sqlite/sqlite.db <examples/sqlite/create.sql

.PHONY: lint
lint: ## Perform lint
	find . -print | grep --regex '.*\.go' | xargs goimports -w
	golangci-lint run --fix --enable-all ./...


.PHONY: test
test: ## Perform test
	go clean -testcache
	go test ./...
	./test/scripts/db/mysql.sh 2> /dev/null
	./test/scripts/db/postgres.sh 2> /dev/null
	./test/scripts/db/spanner.sh 2> /dev/null
	./test/scripts/db/sqlite.sh 2> /dev/null

	./test/scripts/call/grpc-error.sh 2> /dev/null
	./test/scripts/call/grpc-get.sh 2> /dev/null
	./test/scripts/call/grpc-post.sh 2> /dev/null
	./test/scripts/call/grpc-put.sh 2> /dev/null
	./test/scripts/call/grpc-patch.sh 2> /dev/null
	./test/scripts/call/grpc-delete.sh 2> /dev/null
	./test/scripts/call/http-error.sh 2> /dev/null
	./test/scripts/call/http-get.sh 2> /dev/null
	./test/scripts/call/http-post.sh 2> /dev/null
	./test/scripts/call/http-put.sh 2> /dev/null
	./test/scripts/call/http-patch.sh 2> /dev/null
	./test/scripts/call/http-delete.sh 2> /dev/null

	go test -cover ./... -coverprofile=test/cover.out && go tool cover -html=test/cover.out -o test/cover.html

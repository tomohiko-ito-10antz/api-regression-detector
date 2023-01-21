.PHONY: init-spanner
init-spanner:
	gcloud config set project regression-detector
	gcloud config set auth/disable_credentials true
	gcloud config set api_endpoint_overrides/spanner http://spanner:9020/
	gcloud spanner instances describe example || gcloud spanner instances create example --config=emulator-config --description="Instance for example using spanner"
	gcloud spanner databases describe main --instance=example || gcloud spanner databases create main --instance=example
	spanner-cli -p regression-detector -i example -d main --file=examples/spanner/create.sql

.PHONY: init-mysql
init-mysql:
	mysql --host=mysql --password=password main < examples/mysql/create.sql
FROM golang:1.19.4-bullseye

ENV DEBIAN_FRONTEND=noninteractive
# https://cloud.google.com/spanner/docs/emulator#client-libraries
ENV SPANNER_EMULATOR_HOST=spanner:9010

RUN apt update && apt install -y git curl jq default-mysql-client
# gcloud to use spanner https://cloud.google.com/sdk/docs/install-sdk#installing_the_latest_version
RUN echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] http://packages.cloud.google.com/apt cloud-sdk main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list && curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | tee /usr/share/keyrings/cloud.google.gpg && apt-get update -y && apt-get install google-cloud-sdk -y

# go modules
RUN go get -u github.com/nsf/jsondiff && \
    go get -u github.com/mattn/go-sqlite3 && \
    go get -u github.com/go-sql-driver/mysql && \
    go get -u github.com/lib/pq && \
    go get -u github.com/googleapis/go-sql-spanner && \
    go get -u github.com/docopt/docopt-go &&\
    go get -u go.uber.org/multierr && \
    go get -u github.com/cloudspannerecosystem/spanner-cli



WORKDIR /work
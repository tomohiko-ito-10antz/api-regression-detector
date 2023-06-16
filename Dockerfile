FROM golang:1.20.0-bullseye

ENV DEBIAN_FRONTEND=noninteractive
# https://cloud.google.com/spanner/docs/emulator#client-libraries
ENV SPANNER_EMULATOR_HOST=spanner:9010

RUN apt update && apt install -y git curl jq default-mysql-client postgresql-client sqlite3
# gcloud to use spanner https://cloud.google.com/sdk/docs/install-sdk#installing_the_latest_version
RUN echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] http://packages.cloud.google.com/apt cloud-sdk main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list && curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key --keyring /usr/share/keyrings/cloud.google.gpg  add - && apt-get update -y && apt-get install google-cloud-cli -y

WORKDIR /work

# go tools
RUN go install github.com/cloudspannerecosystem/spanner-cli@latest && \
    go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest && \
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest && \
    go install github.com/bufbuild/buf/cmd/buf@v1.14.0 && \
    go install golang.org/x/tools/cmd/goimports@latest

# go modules
COPY go.mod /work/go.mod
RUN go mod download

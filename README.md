# api-regression-detector

## Install

```sh
go install github.com/Jumpaku/api-regression-detector/cmd/call-http@latest
go install github.com/Jumpaku/api-regression-detector/cmd/call-grpc@latest
go install github.com/Jumpaku/api-regression-detector/cmd/compare@latest
go install github.com/Jumpaku/api-regression-detector/cmd/db-init@latest
go install github.com/Jumpaku/api-regression-detector/cmd/db-dump@latest
```

## Usage

### db-init

```
Regression detector db-init.
db-init initializes tables according to JSON data.

Usage:
        program init <database-driver> <connection-string>
        program -h | --help
        program --version

Options:
        <database-driver>   Supported database driver name which is one of mysql, spanner, sqlite3, or postgres
        <connection-string> Connection string corresponding to the database driver.
        -h --help          Show this screen.
        --version          Show version.
```

### db-dump

```
Regression detector db-dump.
db-dump outputs data within tables in JSON format.

Usage:
        program dump <database-driver> <connection-string>
        program -h | --help
        program --version

Options:
        <database-driver>   Supported database driver name which is one of mysql, spanner, sqlite3, or postgres
        <connection-string> Connection string corresponding to the database driver.
        -h --help           Show this screen.
        --version           Show version.
```

### compare

```
Regression detector compare.
compare compares two JSON files.

Usage:
        program compare [--verbose] [--strict] <expected-json> <actual-json>
        program -h | --help
        program --version

Options:
        <expected-json>    JSON file path of expected value.
        <actual-json>      JSON file path of actual value.
        --verbose          Show verbose difference. [default: false]
        --strict           Disallow superset match. [default: false]
        -h --help          Show this screen.
        --version          Show version.
```

### call-http

```
Regression detector call-http.
call-http calls HTTP API: sending JSON request and receiving JSON response.

Usage:
        program call http <endpoint-url> <http-method>
        program -h | --help
        program --version

Options:
        <endpoint-url>     The URL of the HTTP endpoint which may has path parameters enclosed in '[' and ']'.
        <http-method>      One of GET, HEAD, POST, PUT, DELETE, CONNECT, OPTIONS, TRACE, or PATCH.
        -h --help          Show this screen.
        --version          Show version.
```

### call-grpc

```
Regression detector call-grpc.
call-grpc calls GRPC API: sending JSON request and receiving JSON response.

Usage:
        program <grpc-endpoint> <grpc-full-method>
        program -h | --help
        program --version

Options:
        <grpc-endpoint>    host and port joined by ':'.
        <grpc-full-method> full method in the form 'package.name.ServiceName/MethodName'.
        -h --help          Show this screen.
        --version          Show version.
```

## Examples

### db-init

```sh
go run main.go init mysql "root:password@(mysql)/main" <examples/init.json
```

```sh
go run main.go init postgres "user=root password=password host=postgres dbname=main sslmode=disable" <examples/init.json
```

```sh
go run main.go init sqlite3 "file:examples/sqlite/sqlite.db" <examples/init.json
```

```sh
go run main.go init spanner "projects/regression-detector/instances/example/databases/main" <examples/init.json
```

### db-dump

```sh
jq '. | keys' <examples/init.json | go run main.go dump mysql "root:password@(mysql)/main"
```

```sh
jq '. | keys' <examples/init.json | go run main.go dump postgres "user=root password=password host=postgres dbname=main sslmode=disable" <examples/init.json
```

```sh
jq '. | keys' <examples/init.json | go run main.go dump sqlite3 "file:examples/sqlite/sqlite.db" <examples/init.json
```

```sh
jq '. | keys' <examples/init.json | go run main.go dump spanner "projects/regression-detector/instances/example/databases/main" <examples/init.json
```

### compare

```sh
go run main.go compare --strict --verbose examples/expected.json examples/actual.json
```

### call-http

```sh
go run main.go call http 'http://api:80/say/hello/[name]' 'GET' < examples/http/say/hello/[name]/request.json
```

### call-grpc

```sh
go run main.go call grpc 'api:50051' 'api.GreetingService/SayHello' < examples/grpc/api/GreetingService/SayHello/request.json
```

## Development

### Execution

#### Init

#### Dump

#### Compare

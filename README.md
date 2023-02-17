# api-regression-detector

## Usage

```
Regression detector.
The following commands are available:
* init: It initializes tables according to JSON data.
* dump: It outputs data within tables in JSON format.
* call: It calls RPC of HTTP or GRPC: sending JSON request and receiving JSON response.
* compare: It compares two JSON files.

Usage:
	program init <database-driver> <connection-string>
	program dump <database-driver> <connection-string>
	program call http <endpoint-url> <http-method>
	program call grpc <grpc-endpoint> <grpc-full-method>
	program compare [--verbose] [--strict] <expected-json> <actual-json>
	program -h | --help
	program --version

Options:
	-h --help          Show this screen.
	--version          Show version.
	--verbose          Show verbose difference. [default: false]
	--strict           Disallow superset match. [default: false]
```

## Examples

### Init

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

### Dump

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

### Call

```sh
go run main.go call grpc 'api:50051' 'api.GreetingService/SayHello' < examples/grpc/api/GreetingService/SayHello/request.json
```

```sh
go run main.go call http 'http://api:80/say/hello/[name]' 'GET' < examples/http/say/hello/[name]/request.json
```



### Compare

```sh
go run main.go compare --strict --verbose examples/expected.json examples/actual.json
```

## Development

### Execution

#### Init

#### Dump

#### Compare

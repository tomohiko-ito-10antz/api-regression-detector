# api-regression-detector

## Usage

```
Regression detector.
The following commands are available:
* init: It initializes database according to json provided from stdin.
* dump: It outputs database according to json provided from stdin.
* compare: It compares two JSON files and outputs the comparison result to stdout.

Usage:
  program init <database-driver> <connection-string>
  program dump <database-driver> <connection-string>
  program compare [--verbose] [--strict] <expected-json> <actual-json>
  program -h | --help
  program --version

Options:
  -h --help          Show this screen.
  --version          Show version.
  --verbose          Show verbose difference. [default: false]
  --strict           Disallow superset match. [default: false]
  --strict        Disallow superset match. [default: false]
```

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
go run main.go dump sqlite3 "file:examples/sqlite/sqlite.db" <examples/init.json
```

```sh
go run main.go dump spanner "projects/regression-detector/instances/example/databases/main" <examples/init.json
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

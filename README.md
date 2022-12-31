# api-regression-detector

## Usage

```
Regression detector.

Usage:
  program prepare <database-driver> <connection-string>
  program dump <database-driver> <connection-string>
  program compare [--verbose] [--strict] <expected-json> <actual-json>
  program -h | --help
  program --version

Options:
  -h --help       Show this screen.
  --version       Show version.
  --verbose       Show verbose difference. [default: false]
  --strict        Disallow superset match. [default: false]
```

### Options:

  -h --help       Show this screen.
  --version       Show version.
  --verbose       Show verbose difference. [default: false]
  --strict        Disallow superset match. [default: false]

## Development

### Execution

#### Prepare

```sh
go run main.go prepare mysql "root:password@(mysql)/main" <examples/prepare.json
```

```sh
go run main.go prepare postgres "user=postgres password=password host=postgres dbname=main sslmode=disable" <examples/prepare.json
```

```sh
go run main.go prepare sqlite3 "file:examples/sqlite/sqlite.db" <examples/prepare.json
```

#### Dump

#### Compare

```sh
go run main.go compare --strict-match --verbose examples/expected.json examples/actual.json
```

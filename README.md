# api-regression-detector

CLI tools that can be used for the following API regression test:

1. Initialize tables in database, which a target API uses.
2. Send request to the target API.
3. Receive response from the target API.
4. Dump tables in database, which the target API modified.
5. Compare expected response and actual response.
6. Compare expected tables and actual tables.

Table of contents:
1. [Install](#install)
2. [Usage](#usage)
     1. [db-init](#db-init)
     1. [db-dump](#db-dump)
     1. [compare](#compare)
     1. [call-http](#call-http)
     1. [call-grpc](#call-grpc)

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

#### Input

To initialize tables in the database, `db-init` receives the JSON data from stdin. It is assumed that the JSON data is represented as the following type `DBInitInput`.

```ts
type DBInitInput = InitTable[];
type InitTable = {
        /** name of table to be initialized */
        name: string,
        /** rows to be inserted */
        rows: Row[] 
}
type Row = { [columnName: string]: ColumnValue };
type ColumnValue = boolean | string | number | null
```

Example:
```json
[
    {
        "name": "example_table",
        "rows": [
            {
                "c0": "abc",
                "c1": 123,
                "c2": true,
                "c3": "2022-12-25T00:45:17Z",
                "id": 1
            },
            {
                "c0": "",
                "c1": 0,
                "c2": false,
                "c3": "2022-12-24T15:54:17Z",
                "id": 2
            }
        ]
    },
    {
        "name": "child_example_table_1",
        "rows": [
            {
                "id": 1,
                "example_table_id": 1
            }
        ]
    },
    {
        "name": "child_example_table_2",
        "rows": [
            {
                "id": 2,
                "example_table_id": 2
            }
        ]
    }
]
```

#### Output

Nothing.

<span id="usage-db-dump"></span>

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

#### Input

To dump specified tables in the database, `db-dump` receives the JSON data from stdin. It is assumed that the JSON data is represented as the following type `DBDumpInput`.

```ts
/** Array of table names to be dumped */
type DBDumpInput = string[];
```

Example:
```json
[
    "example_table",
    "child_example_table_1",
    "child_example_table_2"
]
```

#### Output

`db-dump` outputs data of specified tables as JSON to stdout. The JSON data is represented as the following type `DBDumpOutput`.

```ts
type DBDumpOutput = { [tableName: string]: Row[] };
type Row = { [columnName: string]: ColumnValue };
type ColumnValue = boolean | string | number | null
```

Example:
```json
{
    "example_table": [
        {
            "c0": "abc",
            "c1": 123,
            "c2": true,
            "c3": "2022-12-25T00:45:17Z",
            "id": 1
        },
        {
            "c0": "",
            "c1": 0,
            "c2": false,
            "c3": "2022-12-24T15:54:17Z",
            "id": 2
        }
    ],
    "child_example_table_1": [
        {
            "example_table_id": 1,
            "id": 1
        }
    ],
    "child_example_table_2": [
        {
            "example_table_id": 2,
            "id": 2
        }
    ]
}
```

<span id="usage-compare"></span>

### compare

```
Regression detector compare.
compare compares two JSON files.

Usage:
	compare [--show-diff] [--no-superset] <expected-json> <actual-json>
	compare -h | --help
	compare --version

Options:
	<expected-json>    JSON file path of expected value.
	<actual-json>      JSON file path of actual value.
	--show-diff        Show difference. [default: false]
	--no-superset      Disallow superset match. [default: false]
	-h --help          Show this screen.
	--version          Show version.
```

#### Input

Nothing.

#### Output

`compare` outputs comparison result of specified two JSON data to stdout. The result is represented as the following type `CompareOutput`.

```ts
type CompareOutput = 
        | "FullMatch" /** two JSON data match exactly */
        | "SupersetMatch" /** second JSON data is an extension of first JSON data */
        | "NoMatch" /** two JSON data are incompatible exactly */;
```

With `--show-diff` option, `compare` outputs difference of specified two JSON data after the result as follows:

```
NoMatch
 |[
 |    {
 |        "add": {
 |            ...skipped 2 object properties...,
+|            "c": 3,
+|            "d": 4
 |        }
 |    },
 |    {
 |        "remove": {
-|            "a": 1
 |        }
 |    },
 |    {
 |        "change": {
~|            "a": 1 => 2
 |        }
 |    }
 |]
```

expected.json
```json
[
    {
        "add": {
            "a": 1,
            "b": 2
        }
    },
    {
        "remove": {
            "a": 1
        }
    },
    {
        "change": {
            "a": 1
        }
    }
]
```

actual.json
```json
[
    {
        "add": {
            "d": 4,
            "a": 1,
            "b": 2,
            "c": 3
        }
    },
    {
        "remove": {}
    },
    {
        "change": {
            "a": 2
        }
    }
]
```

<span id="call-http"></span>

### call-http

```
Regression detector call-http.
call-http calls HTTP API: sending JSON request and receiving JSON response.

Usage:
        program call http <endpoint-url> <http-method>
        program -h | --help
        program --version

Options:
        <endpoint-url>     The URL of the HTTP endpoint which may have path parameters enclosed in '[' and ']'.
        <http-method>      One of GET, HEAD, POST, PUT, DELETE, CONNECT, OPTIONS, TRACE, or PATCH.
        -h --help          Show this screen.
        --version          Show version.
```

#### Input

`call-http` sends request with JSON data given from stdin.

#### Output

`call-http` outputs response body as JSON data to stdout.

<span id="call-grpc"></span>

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

#### Input

`call-grpc` sends request with JSON data given from stdin.

#### Output

`call-grpc` outputs response body as JSON data to stdout.

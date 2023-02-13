package cli

func GetDoc() string {
	return `Regression detector.
The following commands are available:
* init: It initializes tables according to JSON data.
* dump: It outputs data within tables in JSON format.
* call: It calls RPC of HTTP or GRPC: sending JSON request and receiving JSON response.
* compare: It compares two JSON files.

Usage:
	program init <database-driver> <connection-string>
	program dump <database-driver> <connection-string>
	program call <rpc-endpoint> <method>
	program compare [--verbose] [--strict] <expected-json> <actual-json>
	program -h | --help
	program --version

Options:
	-h --help          Show this screen.
	--version          Show version.
	--verbose          Show verbose difference. [default: false]
	--strict           Disallow superset match. [default: false]`
}

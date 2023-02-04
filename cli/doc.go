package cli

func GetDoc() string {
	return `Regression detector.
The following commands are available:
* init: It initializes tables according to JSON data.
* dump: It outputs data within tables in JSON format.
* compare: It compares two JSON files.

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
	--strict           Disallow superset match. [default: false]`
}

package cli

func GetDoc() string {
	return `Regression detector.
The following commands are available:
* init: It initializes database according to json provided by stdin.
* dump: It outputs database according to json provided by stdin.
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
	--strict           Disallow superset match. [default: false]`
}

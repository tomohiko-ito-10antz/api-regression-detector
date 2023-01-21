package db

type Row map[string]any
type Rows []Row
type Tables map[string]Rows
type ColumnTypes map[string]columnType
type columnType string

const (
	ColumnTypeUnknown   columnType = "UNKNOWN"
	ColumnTypeBoolean   columnType = "BOOL"
	ColumnTypeInteger   columnType = "INTEGER"
	ColumnTypeFloat     columnType = "FLOAT"
	ColumnTypeString    columnType = "STRING"
	ColumnTypeTimestamp columnType = "TIMESTAMP"
)

package db

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

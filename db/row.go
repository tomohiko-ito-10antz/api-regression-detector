package db

type Row map[string]*ColumnValue

func (row Row) GetColumnValue(columnName string) (*ColumnValue, bool) {
	val, exists := row[columnName]
	return val, exists
}
func (row Row) SetColumnValue(columnName string, val any, typ ColumnType) {
	row[columnName] = &ColumnValue{Type: typ, value: val}
}

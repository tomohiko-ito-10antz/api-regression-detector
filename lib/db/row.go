package db

type Row map[string]*ColumnValue

func (row Row) GetColumnValue(columnName string) (*ColumnValue, bool) {
	val, exists := row[columnName]
	return val, exists
}
func (row Row) SetColumnValue(columnName string, val any, typ ColumnType) {
	switch val := val.(type) {
	case nil:
		row[columnName] = &ColumnValue{Type: typ, value: nil}
	case *ColumnValue:
		if val == nil {
			row[columnName] = &ColumnValue{Type: typ, value: nil}
		} else {
			row[columnName] = &ColumnValue{Type: typ, value: val.value}
		}
	case ColumnValue:
		row[columnName] = &ColumnValue{Type: typ, value: val.value}
	default:
		row[columnName] = &ColumnValue{Type: typ, value: val}
	}
}

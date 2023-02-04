package mysql

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/db"
)

type selectOperation struct{}

func ListRows() selectOperation {
	return selectOperation{}
}

var _ cmd.RowLister = selectOperation{}

func (o selectOperation) ListRows(ctx context.Context, tx db.Tx, tableName string, schema db.Schema) (rows []db.Row, err error) {
	stmt := fmt.Sprintf(`SELECT * FROM %s ORDER BY %s`, tableName, strings.Join(schema.PrimaryKeys, ", "))
	rows, err = tx.Read(ctx, stmt, nil)
	if err != nil {
		return nil, err
	}
	out := db.Table{}
	for _, row := range rows {
		outRow := db.Row{}
		for _, columnName := range schema.GetColumnNames() {
			col, ok := row[columnName]
			if !ok {
				return nil, fmt.Errorf("column %s not found", columnName)
			}
			colBytes, err := col.AsBytes()
			if err != nil {
				return nil, err
			}
			typ, exists := schema.ColumnTypes[columnName]
			if !exists {
				return nil, fmt.Errorf("column %s not found", columnName)
			}
			var val any
			switch typ {
			case db.ColumnTypeBoolean:
				if colBytes.Valid {
					v, err := strconv.ParseBool(string(colBytes.Bytes))
					if err != nil {
						return nil, err
					}
					val = v
				}
			case db.ColumnTypeInteger:
				if colBytes.Valid {
					v, err := strconv.ParseInt(string(colBytes.Bytes), 10, 64)
					if err != nil {
						return nil, err
					}
					val = v
				}
			case db.ColumnTypeFloat:
				if colBytes.Valid {
					v, err := strconv.ParseFloat(string(colBytes.Bytes), 64)
					if err != nil {
						return nil, err
					}
					val = v
				}
			case db.ColumnTypeTime, db.ColumnTypeString:
				if colBytes.Valid {
					val = string(colBytes.Bytes)
					typ = db.ColumnTypeString
				}
			default:
				return nil, fmt.Errorf("unexpected type %v of column %s not found", typ, columnName)
			}
			outRow[columnName] = db.NewColumnValue(val, typ)
		}
		out.Rows = append(out.Rows, outRow)
	}
	return out.Rows, nil
}

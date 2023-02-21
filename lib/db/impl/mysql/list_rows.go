package mysql

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/db"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
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
		return nil, errors.Wrap(
			errors.DBFailure,
			"fail to select all rows (stmt=%v)", stmt)
	}

	out := db.Table{}
	for _, row := range rows {
		outRow := db.Row{}
		for _, columnName := range schema.GetColumnNames() {
			col, ok := row[columnName]
			if !ok {
				return nil, errors.Wrap(
					errors.BadKeyAccess,
					"column %s not found in schema of table %s", columnName, tableName)
			}

			colBytes, err := col.AsBytes()
			if err != nil {
				return nil, errors.Wrap(
					err,
					"fail to parse value of column %s as []byte", columnName)
			}

			typ, exists := schema.ColumnTypes[columnName]
			if !exists {
				return nil, errors.Wrap(
					errors.BadKeyAccess,
					"column %s not found in schema of table %s", columnName, tableName)
			}

			var val any

			switch typ {
			case db.ColumnTypeBoolean:
				if colBytes.Valid {
					val, err = strconv.ParseBool(string(colBytes.Bytes))
					if err != nil {
						return nil, errors.Wrap(
							errors.Join(err, errors.BadConversion),
							"fail to convert value %v:%T of column %s as bool", colBytes, colBytes, columnName)
					}
				}
			case db.ColumnTypeInteger:
				if colBytes.Valid {
					val, err = strconv.ParseInt(string(colBytes.Bytes), 10, 64)
					if err != nil {
						return nil, errors.Wrap(
							errors.Join(err, errors.BadConversion),
							"fail to convert value %v:%T of column %s as integer", colBytes, colBytes, columnName)
					}
				}
			case db.ColumnTypeFloat:
				if colBytes.Valid {
					val, err = strconv.ParseFloat(string(colBytes.Bytes), 64)
					if err != nil {
						return nil, errors.Wrap(
							errors.Join(err, errors.BadConversion),
							"fail to convert value %v:%T of column %s as float", colBytes, colBytes, columnName)
					}
				}
			case db.ColumnTypeTime, db.ColumnTypeString:
				if colBytes.Valid {
					val = string(colBytes.Bytes)
					typ = db.ColumnTypeString
				}
			default:
				return nil, errors.Wrap(
					errors.Join(err, errors.Unexpected),
					"unexpected type %v of column %s", typ, columnName)
			}

			outRow[columnName] = db.NewColumnValue(val, typ)
		}

		out.Rows = append(out.Rows, outRow)
	}

	return out.Rows, nil
}

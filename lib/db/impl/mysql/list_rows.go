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

	errInfo := errors.Info{"stmt": stmt}

	rows, err = tx.Read(ctx, stmt, nil)
	if err != nil {
		return nil, errors.Wrap(
			errors.DBFailure.Err(err),
			errInfo.AppendTo("fail to select all rows"))
	}

	out := db.Table{}
	for _, row := range rows {
		outRow := db.Row{}
		for _, columnName := range schema.GetColumnNames() {
			errInfo := errInfo.With("columnName", columnName)

			col, ok := row[columnName]
			if !ok {
				return nil, errors.BadKeyAccess.New(
					errInfo.AppendTo("column not found in row"))
			}

			colBytes, err := col.AsBytes()
			if err != nil {
				return nil, errors.Wrap(
					errors.BadConversion.Err(err),
					errInfo.With("columnValue", col).AppendTo("fail to parse column value to []byte"))
			}

			typ, exists := schema.ColumnTypes[columnName]
			if !exists {
				return nil, errors.BadKeyAccess.New(
					errInfo.AppendTo("column not found in table schema"))
			}

			errInfo = errInfo.
				With("colBytes", string(colBytes.Bytes)).
				With("columnType", typ)

			var val any

			switch typ {
			case db.ColumnTypeBoolean:
				if colBytes.Valid {
					val, err = strconv.ParseBool(string(colBytes.Bytes))
					if err != nil {
						return nil, errors.Wrap(
							errors.BadConversion.Err(err),
							errInfo.AppendTo("fail to convert column value to bool"))
					}
				}
			case db.ColumnTypeInteger:
				if colBytes.Valid {
					val, err = strconv.ParseInt(string(colBytes.Bytes), 10, 64)
					if err != nil {
						return nil, errors.Wrap(
							errors.BadConversion.Err(err),
							errInfo.AppendTo("fail to convert column value to integer"))
					}
				}
			case db.ColumnTypeFloat:
				if colBytes.Valid {
					val, err = strconv.ParseFloat(string(colBytes.Bytes), 64)
					if err != nil {
						return nil, errors.Wrap(
							errors.BadConversion.Err(err),
							errInfo.AppendTo("fail to convert column value to float"))
					}
				}
			case db.ColumnTypeTime, db.ColumnTypeString:
				if colBytes.Valid {
					val = string(colBytes.Bytes)
					typ = db.ColumnTypeString
				}
			default:
				return nil, errors.Unsupported.New(errInfo.AppendTo("unsupported DB column type"))
			}

			outRow[columnName] = db.NewColumnValue(val, typ)
		}

		out.Rows = append(out.Rows, outRow)
	}

	return out.Rows, nil
}

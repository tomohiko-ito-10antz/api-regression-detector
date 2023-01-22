package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Jumpaku/api-regression-detector/cmd"
	"github.com/Jumpaku/api-regression-detector/db"
)

type insertOperation struct {
}

func Insert() insertOperation {
	return insertOperation{}
}

var _ cmd.Insert = insertOperation{}

func (o insertOperation) Insert(ctx context.Context, exec db.Exec, table string, rows db.Rows) (err error) {
	columns := getColumns(rows)
	if len(columns) == 0 {
		return nil
	}
	stmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES", table, strings.Join(columns, ", "))
	values := []any{}
	n := 0
	for i, row := range rows {
		if i > 0 {
			stmt += ","
		}
		stmt += "("
		for j, column := range columns {
			if j > 0 {
				stmt += ","
			}
			n++
			stmt += fmt.Sprintf(`$%d`, n)
			value, ok := row[column]
			if !ok || value == nil {
				values = append(values, nil)
			} else {
				switch value := value.(type) {
				case json.Number, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
					values = append(values, value)
				case string:
					values = append(values, value)
				case bool:
					values = append(values, value)
				default:
					return fmt.Errorf("unexpected value %v", value)
				}
			}
		}
		stmt += ")"
	}
	err = exec.Write(ctx, stmt, values)
	if err != nil {
		return err
	}
	return nil
}

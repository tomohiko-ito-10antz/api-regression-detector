package sqlite

import (
	"context"
	"strings"

	"github.com/Jumpaku/api-regression-detector/db"
)

func getColumns(rows db.Rows) (columns []string) {
	columnAdded := map[string]bool{}
	for _, row := range rows {
		for column := range row {
			if _, added := columnAdded[column]; !added {
				columnAdded[column] = true
				columns = append(columns, column)
			}
		}
	}
	return columns
}

func getPrimaryKeys(ctx context.Context, tx db.Exec, table string) (primaryKeys []string, err error) {
	rows, err := tx.Read(ctx, `SELECT name FROM pragma_table_info(?) WHERE pk > 0 ORDER BY pk`, []any{table})
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		col, err := row.GetString("name")
		if err != nil {
			return nil, err
		}
		primaryKeys = append(primaryKeys, strings.ToLower(col))
	}
	return primaryKeys, nil
}

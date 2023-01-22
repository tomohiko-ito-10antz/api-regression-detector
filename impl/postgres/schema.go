package postgres

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

func getColumnNames(ctx context.Context, tx db.Exec, table string) (primaryKeys []string, err error) {
	trimmed := table
	if strings.HasPrefix(trimmed, `"`) && strings.HasSuffix(trimmed, `"`) {
		trimmed = strings.TrimPrefix(strings.TrimSuffix(trimmed, `"`), `"`)
	}
	rows, err := tx.Read(ctx, `SELECT column_name FROM information_schema.columns WHERE table_name = $1 ORDER BY ordinal_position`, []any{trimmed})
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		col, err := row.GetByteArray("column_name")
		if err != nil {
			return nil, err
		}
		primaryKeys = append(primaryKeys, strings.ToLower(string(col)))
	}
	return primaryKeys, nil
}

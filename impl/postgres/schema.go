package postgres

import (
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

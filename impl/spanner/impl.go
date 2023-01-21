package spanner

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/Jumpaku/api-regression-detector/cmd"
	"github.com/Jumpaku/api-regression-detector/db"
)

type op struct {
}

func Truncate() interface {
	Truncate(ctx context.Context, tx db.Exec, table string) (err error)
} {
	return &op{}
}

func Insert() interface {
	Insert(ctx context.Context, tx db.Exec, table string, rows db.Rows) (err error)
} {
	return &op{}
}

func Select() interface {
	Select(ctx context.Context, exec db.Exec, table string) (rows db.Rows, err error)
} {
	return &op{}
}

var _ cmd.Truncate = (*op)(nil)
var _ cmd.Insert = (*op)(nil)
var _ cmd.Select = (*op)(nil)

func (o *op) Select(ctx context.Context, exec db.Exec, table string) (rows db.Rows, err error) {
	rows, err = exec.Read(ctx, fmt.Sprintf(`SELECT * FROM %s`, table), nil)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
func (o *op) Truncate(ctx context.Context, tx db.Exec, table string) (err error) {
	err = tx.Write(ctx, fmt.Sprintf(`DELETE FROM %s WHERE TRUE`, table), nil)
	if err != nil {
		return err
	}
	return nil
}
func (o *op) Insert(ctx context.Context, tx db.Exec, table string, rows db.Rows) (err error) {
	columnTypes, err := getColumnSchema(ctx, tx, table)
	if err != nil {
		return err
	}
	columns := getColumns(rows)
	if len(columns) == 0 {
		return nil
	}
	fmt.Printf("%v\n", columnTypes)
	fmt.Printf("%v\n", columns)
	stmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES", table, strings.Join(columns, ", "))
	values := []any{}
	for i, row := range rows {
		if i > 0 {
			stmt += ","
		}
		stmt += "("
		for j, column := range columns {
			if j > 0 {
				stmt += ","
			}
			stmt += "?"
			columnType, ok := columnTypes[column]
			if !ok {
				return err
			}
			value, ok := row[column]
			if !ok {
				value = nil
			}
			switch columnType {
			case db.ColumnTypeBoolean:
				v, err := toBoolean(value)
				if err != nil {
					return err
				}
				values = append(values, v)
			case db.ColumnTypeInteger:
				v, err := toInteger(value)
				if err != nil {
					return err
				}
				values = append(values, v)
			case db.ColumnTypeFloat:
				v, err := toFloat(value)
				if err != nil {
					return err
				}
				values = append(values, v)
			default:
				v, err := toString(value)
				if err != nil {
					return err
				}
				values = append(values, v)
			}
		}
		stmt += ")"
	}
	err = tx.Write(ctx, stmt, values)
	if err != nil {
		return err
	}
	return nil
}

func getColumns(rows db.Rows) (columns []string) {
	columnAdded := map[string]bool{}
	for _, row := range rows {
		for column := range row {
			if _, added := columnAdded[column]; !added {
				columnAdded[column] = true
				columns = append(columns, strings.ToLower(column))
			}
		}
	}
	return columns
}

func getColumnSchema(ctx context.Context, tx db.Exec, table string) (columnTypes db.ColumnTypes, err error) {
	rows, err := tx.Read(ctx, `SELECT column_name, spanner_type FROM information_schema.columns WHERE table_name = ?`, []any{table})
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v\n", rows)
	columnTypes = db.ColumnTypes{}
	for _, row := range rows {
		colAny, ok := row["column_name"]
		if !ok || colAny == nil {
			return nil, err
		}
		col, ok := colAny.(string)
		if !ok || colAny == nil {
			return nil, err
		}
		typAny, ok := row["spanner_type"]
		if !ok || typAny == nil {
			return nil, err
		}
		typ, ok := typAny.(string)
		if !ok || typAny == nil {
			return nil, err
		}
		fmt.Println("5")
		s := strings.ToUpper(typ)
		startsWith := func(prefix string) bool {
			return strings.HasPrefix(s, prefix)
		}
		switch {
		case startsWith("BOOL"):
			columnTypes[strings.ToLower(col)] = db.ColumnTypeBoolean
		case startsWith("INT64"):
			columnTypes[strings.ToLower(col)] = db.ColumnTypeInteger
		case startsWith("FLOAT64"):
			columnTypes[strings.ToLower(col)] = db.ColumnTypeFloat
		default:
			columnTypes[strings.ToLower(col)] = db.ColumnTypeString
		}
	}
	return columnTypes, nil
}

func toString(value any) (p *string, err error) {
	var s string
	switch value := value.(type) {
	case nil:
		return nil, nil
	case json.Number:
		s = value.String()
	case int:
		s = strconv.FormatInt(int64(value), 64)
	case int8:
		s = strconv.FormatInt(int64(value), 64)
	case int16:
		s = strconv.FormatInt(int64(value), 64)
	case int32:
		s = strconv.FormatInt(int64(value), 64)
	case int64:
		s = strconv.FormatInt(int64(value), 64)
	case uint:
		s = strconv.FormatUint(uint64(value), 64)
	case uint8:
		s = strconv.FormatUint(uint64(value), 64)
	case uint16:
		s = strconv.FormatUint(uint64(value), 64)
	case uint32:
		s = strconv.FormatUint(uint64(value), 64)
	case uint64:
		s = strconv.FormatUint(uint64(value), 64)
	case string:
		s = value
	case bool:
		s = strconv.FormatBool(value)
	default:
		return nil, fmt.Errorf("unexpected value %v", value)
	}
	return &s, nil
}

func toInteger(value any) (p *int64, err error) {
	var i int64
	switch value := value.(type) {
	case nil:
		return nil, nil
	case json.Number:
		i, err = value.Int64()
		if err != nil {
			return nil, fmt.Errorf("unexpected value %v", value)
		}
	case int:
		i = int64(value)
	case int8:
		i = int64(value)
	case int16:
		i = int64(value)
	case int32:
		i = int64(value)
	case int64:
		i = int64(value)
	case uint:
		i = int64(value)
	case uint8:
		i = int64(value)
	case uint16:
		i = int64(value)
	case uint32:
		i = int64(value)
	case uint64:
		i = int64(value)
	case string:
		i, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("unexpected value %v", value)
		}
	default:
		return nil, fmt.Errorf("unexpected value %v", value)
	}
	return &i, nil
}

func toFloat(value any) (p *float64, err error) {
	var f float64
	switch value := value.(type) {
	case nil:
		return nil, nil
	case json.Number:
		f, err = value.Float64()
		if err != nil {
			return nil, fmt.Errorf("unexpected value %v", value)
		}
	case int:
		f = float64(value)
	case int8:
		f = float64(value)
	case int16:
		f = float64(value)
	case int32:
		f = float64(value)
	case int64:
		f = float64(value)
	case uint:
		f = float64(value)
	case uint8:
		f = float64(value)
	case uint16:
		f = float64(value)
	case uint32:
		f = float64(value)
	case uint64:
		f = float64(value)
	case string:
		f, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("unexpected value %v", value)
		}
	default:
		return nil, fmt.Errorf("unexpected value %v", value)
	}
	return &f, nil
}

func toBoolean(value any) (p *bool, err error) {
	var b bool
	switch value := value.(type) {
	case nil:
		return nil, nil
	case json.Number:
		i, err := value.Int64()
		if err != nil {
			return nil, fmt.Errorf("unexpected value %v", value)
		}
		b = i != 0
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		b = value != 0
	case string:
		if s := strings.ToLower(value); s == "true" {
			b = true
		} else if s == "false" {
			b = false
		} else {
			return nil, fmt.Errorf("unexpected value %v", value)
		}
	case bool:
		b = value
	default:
		return nil, fmt.Errorf("unexpected value %v", value)
	}
	return &b, nil
}

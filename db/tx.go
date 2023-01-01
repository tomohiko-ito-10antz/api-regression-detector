package db

import (
	"context"
	"database/sql"
	"reflect"

	"github.com/Jumpaku/api-regression-detector/log"
	"go.uber.org/multierr"
)

type Exec interface {
	Write(ctx context.Context, stmt string, params []any) (err error)
	Read(ctx context.Context, stmt string, params []any) (rows Rows, err error)
}

type exec struct {
	tx *sql.Tx
}

func rollback(ctx context.Context, tx *sql.Tx, err error) error {
	return multierr.Combine(err, tx.Rollback())
}
func commit(ctx context.Context, tx *sql.Tx) (err error) {
	err = tx.Commit()
	if err != nil {
		return rollback(ctx, tx, err)
	}
	return nil
}
func Transaction(ctx context.Context, db *sql.DB, handler func(ctx context.Context, exec Exec) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { err = rollback(ctx, tx, err) }()
	e := exec{tx: tx}
	err = handler(ctx, &e)
	if err != nil {
		return err
	}
	return commit(ctx, tx)
}

func (e *exec) Write(ctx context.Context, stmt string, params []any) (err error) {
	log.Stderr("SQL\n\tstatement: %v\n\tparams   : %v", stmt, params)
	_, err = e.tx.Exec(stmt, params...)
	if err != nil {
		return err
	}
	return nil
}

func (e *exec) Read(ctx context.Context, stmt string, params []any) (rows Rows, err error) {
	log.Stderr("SQL\n\tstatement: %v\n\tparams   : %v", stmt, params)
	itr, err := e.tx.Query(stmt, params...)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = multierr.Combine(err, itr.Close())
	}()
	for itr.Next() {
		columns, err := itr.Columns()
		if err != nil {
			return nil, err
		}
		types, err := itr.ColumnTypes()
		if err != nil {
			return nil, err
		}
		var values []any
		for _, typ := range types {
			switch {
			case isBoolean(typ):
				var v *bool
				values = append(values, &v)
			case isInteger(typ):
				var v *int64
				values = append(values, &v)
			case isFloat(typ):
				var v *float64
				values = append(values, &v)
			default:
				var v *string
				values = append(values, &v)
			}
		}
		err = itr.Scan(values...)
		if err != nil {
			return nil, err
		}
		row := Row{}
		for i, column := range columns {
			row[column] = values[i]
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func isBoolean(typ *sql.ColumnType) bool {
	t := typ.ScanType()
	v := reflect.New(t).Elem().Interface()
	switch v.(type) {
	case bool, *bool, sql.NullBool:
		return true
	default:
		return false
	}
}

func isInteger(typ *sql.ColumnType) bool {
	t := typ.ScanType()
	v := reflect.New(t).Elem().Interface()
	switch v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64,
		*int, *int8, *int16, *int32, *int64, *uint, *uint8, *uint16, *uint32, *uint64,
		sql.NullByte, sql.NullInt16, sql.NullInt32, sql.NullInt64:
		return true
	default:
		return false
	}
}

func isFloat(typ *sql.ColumnType) bool {
	t := typ.ScanType()
	v := reflect.New(t).Elem().Interface()
	switch v.(type) {
	case float32, float64, *float32, *float64, sql.NullFloat64:
		return true
	default:
		return false
	}
}

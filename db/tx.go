package db

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

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
	log.Stderr("SQL\n\tstatement: %v\n\tparams   : %v", stmt, paramsToStrings(params))
	_, err = e.tx.Exec(stmt, params...)
	if err != nil {
		return err
	}
	return nil
}

func (e *exec) Read(ctx context.Context, stmt string, params []any) (rows Rows, err error) {
	log.Stderr("SQL\n\tstatement: %v\n\tparams   : %v", stmt, paramsToStrings(params))
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
		columnCount := len(columns)
		var pointers = make([]any, columnCount)
		var values = make([]any, columnCount)
		for i := 0; i < columnCount; i++ {
			pointers[i] = &values[i]
		}
		err = itr.Scan(pointers...)
		if err != nil {
			return nil, err
		}
		row := Row{}
		for i, column := range columns {
			row[strings.ToLower(column)] = values[i]
		}
		rows = append(rows, row)
	}
	fmt.Printf("%t\n", rows)
	return rows, nil
}

func paramsToStrings(params []any) (strArr []string) {
	strArr = []string{}
	for _, param := range params {
		rv := reflect.ValueOf(param)
		switch {
		case !rv.IsValid(), rv.Kind() == reflect.Pointer && rv.IsNil():
			strArr = append(strArr, "<nil>")
		case rv.Kind() == reflect.Pointer:
			strArr = append(strArr, fmt.Sprintf("%v", rv.Elem().Interface()))
		default:
			strArr = append(strArr, fmt.Sprintf("%v", rv.Interface()))
		}
	}
	return strArr
}

func getType(rt reflect.Type) (ct columnType, err error) {
	switch rt.Kind() {
	case reflect.Bool:
		return ColumnTypeBoolean, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return ColumnTypeInteger, nil
	case reflect.Float32, reflect.Float64:
		return ColumnTypeFloat, nil
	case reflect.String:
		return ColumnTypeString, nil
	case reflect.Pointer:
		return getType(rt.Elem())
	case reflect.Struct:
		rv := reflect.New(rt).Elem()
		if _, isTime := rv.Interface().(time.Time); isTime {
			return ColumnTypeTimestamp, nil
		}
	}
	return ColumnTypeUnknown, fmt.Errorf("unsupported column type %v", rt.String())
}

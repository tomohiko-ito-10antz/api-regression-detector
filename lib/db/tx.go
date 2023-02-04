package db

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"

	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/log"
)

type Tx interface {
	Write(ctx context.Context, stmt string, params []any) (err error)
	Read(ctx context.Context, stmt string, params []any) (rows []Row, err error)
}

type transaction struct {
	tx *sql.Tx
}

func runTransaction(ctx context.Context, db *sql.DB, handler func(ctx context.Context, tx Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(errors.Join(err, errors.DBFailure), "fail to begin transaction")
	}

	defer func() { err = errors.Wrap(rollback(tx, err), "fail runTransaction") }()

	err = handler(ctx, &transaction{tx: tx})
	if err != nil {
		return errors.Wrap(err, "fail to run transaction handler")
	}

	err = commit(tx)
	if err != nil {
		return errors.Wrap(err, "fail runTransaction")
	}

	return nil
}

func (e *transaction) Write(ctx context.Context, stmt string, params []any) error {
	log.Stderr("SQL\n\tstatement: %v\n\tparams   : %v", stmt, paramsToStrings(params))

	_, err := e.tx.Exec(stmt, params...)
	if err != nil {
		return errors.Wrap(errors.Join(err, errors.DBFailure), "fail to write in transaction (stmt=%s,params=%v)", stmt, params)
	}

	return nil
}

func (e *transaction) Read(ctx context.Context, stmt string, params []any) ([]Row, error) {
	log.Stderr("SQL\n\tstatement: %v\n\tparams   : %v", stmt, paramsToStrings(params))

	itr, err := e.tx.Query(stmt, params...)
	if err != nil {
		return nil, errors.Wrap(errors.Join(err, errors.IOFailure), "fail to read in transaction (stmt=%s,params=%v)", stmt, params)
	}

	defer func() {
		err = errors.Wrap(errors.Join(err, itr.Close(), errors.DBFailure), "fail Read")
	}()

	rows := []Row{}

	for itr.Next() {
		columns, err := itr.Columns()
		if err != nil {
			return nil, errors.Wrap(errors.Join(err, errors.DBFailure), "fail to get column names")
		}

		columnCount := len(columns)
		pointers := make([]any, columnCount)
		values := make([]any, columnCount)

		for i := 0; i < columnCount; i++ {
			pointers[i] = &values[i]
		}

		err = itr.Scan(pointers...)
		if err != nil {
			return nil, errors.Wrap(errors.Join(err, errors.DBFailure), "fail to scan column values")
		}

		row := Row{}
		for i, column := range columns {
			row[column] = UnknownTypeColumnValue(values[i])
		}

		rows = append(rows, row)
	}

	return rows, nil
}

func rollback(tx *sql.Tx, err error) error {
	return errors.Wrap(errors.Join(err, tx.Rollback(), errors.DBFailure), "fail to rollback")
}

func commit(tx *sql.Tx) error {
	if err := tx.Commit(); err != nil {
		return errors.Wrap(errors.Join(rollback(tx, err), errors.DBFailure), "fail to commit")
	}

	return nil
}

func paramsToStrings(params []any) []string {
	strArr := []string{}

	for _, param := range params {
		rv := reflect.ValueOf(param)

		switch {
		case !rv.IsValid(), rv.Kind() == reflect.Pointer && rv.IsNil():
			strArr = append(strArr, "<nil>")
		case rv.Kind() == reflect.Pointer:
			strArr = append(strArr, fmt.Sprintf("%v:%T", rv.Elem().Interface(), rv.Elem().Interface()))
		default:
			strArr = append(strArr, fmt.Sprintf("%v:%T", rv.Interface(), rv.Interface()))
		}
	}

	return strArr
}

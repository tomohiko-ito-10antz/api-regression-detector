package db

import (
	"context"
	"database/sql"

	"github.com/Jumpaku/api-regression-detector/lib/errors"
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
		return errors.Wrap(errors.DBFailure.Err(err), "fail to begin transaction")
	}

	defer func() { err = errors.Wrap(errors.DBFailure.Err(rollback(tx, err)), "fail runTransaction") }()

	if err = handler(ctx, &transaction{tx: tx}); err != nil {
		return err
	}

	if err = commit(tx); err != nil {
		return errors.Wrap(err, "fail to commit")
	}

	return nil
}

func (e *transaction) Write(ctx context.Context, stmt string, params []any) error {
	errInfo := errors.Info{"stmt": stmt, "params": params}

	if _, err := e.tx.Exec(stmt, params...); err != nil {
		return errors.Wrap(errors.DBFailure.Err(err), errInfo.AppendTo("fail to write in transaction"))
	}

	return nil
}

func (e *transaction) Read(ctx context.Context, stmt string, params []any) ([]Row, error) {
	errInfo := errors.Info{"stmt": stmt, "params": params}

	itr, err := e.tx.Query(stmt, params...)
	if err != nil {
		return nil, errors.Wrap(errors.DBFailure.Err(err), errInfo.AppendTo("fail to read in transaction"))
	}

	defer func() {
		if errs := errors.Join(err, itr.Close()); err != nil {
			err = errors.Wrap(errors.DBFailure.Err(errs), "fail Read")
		}
	}()

	rows := []Row{}

	for itr.Next() {
		columns, err := itr.Columns()
		if err != nil {
			return nil, errors.Wrap(errors.DBFailure.Err(err), errInfo.AppendTo("fail to get column names"))
		}

		columnCount := len(columns)
		pointers := make([]any, columnCount)
		values := make([]any, columnCount)

		for i := 0; i < columnCount; i++ {
			pointers[i] = &values[i]
		}

		err = itr.Scan(pointers...)
		if err != nil {
			return nil, errors.Wrap(errors.DBFailure.Err(err), errInfo.AppendTo("fail to scan column values"))
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
	return errors.Wrap(errors.Join(err, errors.DBFailure.Err(tx.Rollback())), "fail to rollback")
}

func commit(tx *sql.Tx) error {
	if err := tx.Commit(); err != nil {
		return errors.Wrap(errors.DBFailure.Err(rollback(tx, err)), "fail to commit")
	}

	return nil
}

package db

import (
	"database/sql"
	"fmt"
)

type IExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type IQueryer interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type IAutoTx interface {
	IExecutor
	IQueryer
	Prepare(query string) (*sql.Stmt, error)
	Stmt(stmt *sql.Stmt) *sql.Stmt
}

// A transaction handler function.
// The handler handles any panics in transaction and decides whether to
// commit or rollback automatically according to the error returned.
func Transact(db *sql.DB, txFunc func(IAutoTx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			switch p := p.(type) {
			case error:
				err = p
			default:
				err = fmt.Errorf("%v", p)
			}
		}
		if err != nil {
			// Rollback error ignored.
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	return txFunc(tx)
}

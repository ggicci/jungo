package db

import (
	"database/sql"
	"fmt"
)

type Executor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type Queryer interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type Preparer interface {
	Prepare(query string) (*sql.Stmt, error)
}

type AutoTx interface {
	Executor
	Queryer
	Preparer
	Stmt(stmt *sql.Stmt) *sql.Stmt
}

// A transaction handler function.
// The handler handles any panics in transaction and decides whether to
// commit or rollback automatically according to the error returned.
func Transact(db *sql.DB, txFunc func(AutoTx) error) (err error) {
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

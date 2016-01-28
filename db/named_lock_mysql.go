package db

import (
	"database/sql"
	"errors"
)

type mysqlNamedLock struct {
	db   Queryer
	name string
}

func NewMySQLNamedLock(db Queryer, name string) NamedLock {
	return &mysqlNamedLock{db, name}
}

func (lock *mysqlNamedLock) Name() string { return lock.name }

func (lock *mysqlNamedLock) Acquire(timeout int) error {
	sqlstr := `select get_lock(?, ?);`
	var nullint sql.NullInt64
	if err := lock.db.QueryRow(sqlstr, lock.name, timeout).Scan(&nullint); err != nil {
		return err
	}
	if nullint.Valid {
		if nullint.Int64 == 1 {
			// 1: was obtained successfully.
			return nil
		}
		if nullint.Int64 == 0 {
			// 0: the attempt timed out (for example, because another client has previously locked the name).
			return errors.New("timeout")
		}
	}
	// NULL: an error occurred (such as running out of memory or the thread was killed with mysqladmin kill).
	return errors.New("null value")
}

func (lock *mysqlNamedLock) Release() error {
	sqlstr := `select release_lock(?);`
	var nullint sql.NullInt64
	if err := lock.db.QueryRow(sqlstr, lock.name).Scan(&nullint); err != nil {
		return err
	}
	if nullint.Valid {
		if nullint.Int64 == 1 {
			// 1: lock was released.
			return nil
		}
		if nullint.Int64 == 0 {
			// 0: the lock was not established by this thread (in which case the lock is not released).
			return errors.New("the lock was not established by you")
		}
	}
	// NULL: the named lock did not exist.
	return errors.New("lock does not exist")
}

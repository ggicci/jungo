package db

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
)

type mysqlNamedLock struct {
	tx   *sql.Tx
	name string
}

func NewMySQLNamedLock(db *sql.DB, dbName, tblName, lockName string) (NamedLock, error) {
	name := dbName + "-" + tblName + "-" + lockName
	if len(name) > 64 {
		arr := md5.Sum([]byte(name))
		name = hex.EncodeToString(arr[:])
	}
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	return &mysqlNamedLock{tx, name}, nil
}

func (lock *mysqlNamedLock) Name() string { return lock.name }

func (lock *mysqlNamedLock) Acquire(timeout int) error {
	sqlstr := `select get_lock(?, ?);`
	var nullint sql.NullInt64
	if err := lock.tx.QueryRow(sqlstr, lock.name, timeout).Scan(&nullint); err != nil {
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
	if err := lock.tx.QueryRow(sqlstr, lock.name).Scan(&nullint); err != nil {
		lock.tx.Rollback()
		return err
	}
	if nullint.Valid {
		if nullint.Int64 == 1 {
			// 1: lock was released.
			lock.tx.Commit()
			return nil
		}
		if nullint.Int64 == 0 {
			// 0: the lock was not established by this thread (in which case the lock is not released).
			lock.tx.Rollback()
			return errors.New("the lock was not established by you")
		}
	}
	// NULL: the named lock did not exist.
	lock.tx.Rollback()
	return errors.New("lock does not exist")
}

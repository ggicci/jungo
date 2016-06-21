package db

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
)

type NamedLock interface {
	Acquire(name string, timeout int) error // timeout: seconds
	Release(name string) error
}

type mysqlNamedLock struct {
	tx *sql.Tx
}

func NewMySQLNamedLock(db *sql.DB) (NamedLock, error) {
	// Names are locked on a server-wide basis.
	// Use tx. Since only tx is bound to a single connection. You can't directly use *sql.DB,
	// it's a connection pool.
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	return &mysqlNamedLock{tx}, nil
}

func (lock *mysqlNamedLock) Acquire(name string, timeout int) error {
	// MySQL 5.7.5 and later enforces a maximum length on lock names of 64 characters.
	// Previously, no limit was enforced.
	if len(name) > 64 {
		arr := md5.Sum([]byte(name))
		name = hex.EncodeToString(arr[:])
	}

	sqlstr := `select get_lock(?, ?);`
	var nullint sql.NullInt64
	if err := lock.tx.QueryRow(sqlstr, name, timeout).Scan(&nullint); err != nil {
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

// In MySQL 5.7.5 or later, the second GET_LOCK() acquires a second lock and both RELEASE_LOCK()
// calls return 1 (success). Before MySQL 5.7.5, the second GET_LOCK() releases the first lock ('lock1')
// and the second RELEASE_LOCK() returns NULL (failure) because there is no 'lock1' to release.
func (lock *mysqlNamedLock) Release(name string) error {
	sqlstr := `select release_lock(?);`
	var nullint sql.NullInt64
	if err := lock.tx.QueryRow(sqlstr, name).Scan(&nullint); err != nil {
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

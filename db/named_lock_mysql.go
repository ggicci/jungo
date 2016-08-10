package db

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"
)

type SingleNamedLock interface {
	Acquire(timeout int) error // timeout: seconds
	Release() error
}

func BattleForLock(db *sql.DB, name string) error {
	lock, err := NewMySQLSingleNamedLock(db, name)
	if err != nil {
		return err
	}

	if err = lock.Acquire(0); err != nil {
		return err
	}

	defer func(lock SingleNamedLock, name string) {
		if err := lock.Release(); err != nil {
			println("[MySQLSingleNamedLock] release lock:", err)
		}
	}(lock, name)

	time.Sleep(3 * time.Second)

	return nil
}

type mysqlSingleNamedLock struct {
	tx   *sql.Tx
	name string
}

func NewMySQLSingleNamedLock(db *sql.DB, name string) (SingleNamedLock, error) {
	// Names are locked on a server-wide basis.
	// Use tx. Since only tx is bound to a single connection. You can't directly use *sql.DB,
	// it's a connection pool.
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	// MySQL 5.7.5 and later enforces a maximum length on lock names of 64 characters.
	// Previously, no limit was enforced.
	if len(name) > 64 {
		arr := md5.Sum([]byte(name))
		name = hex.EncodeToString(arr[:])
		println("[MySQLSingleNamedLock], lock named was reduced to \"" + name + "\"")
	}

	return &mysqlSingleNamedLock{tx, name}, nil
}

// name: you'd better use format like {db}-{table}-{lock}.
func (lock *mysqlSingleNamedLock) Acquire(timeout int) error {
	name := lock.name
	sqlstr := `select get_lock(?, ?);`
	var result sql.NullInt64
	if err := lock.tx.QueryRow(sqlstr, name, timeout).Scan(&result); err != nil {
		lock.tx.Rollback()
		return err
	}

	if result.Valid && result.Int64 == 1 {
		// 1: was obtained successfully.
		return nil
	}

	lock.tx.Rollback()

	if result.Valid && result.Int64 == 0 {
		// 0: the attempt timed out (for example, because another client has previously locked the name).
		return errors.New("timeout")
	}
	// NULL: an error occurred (such as running out of memory or the thread was killed with mysqladmin kill).
	return errors.New("null value (unknown error)")
}

// In MySQL 5.7.5 or later, the second GET_LOCK() acquires a second lock and both RELEASE_LOCK()
// calls return 1 (success). Before MySQL 5.7.5, the second GET_LOCK() releases the first lock ('lock1')
// and the second RELEASE_LOCK() returns NULL (failure) because there is no 'lock1' to release.
// You must call this method after calling `Acquire`.
func (lock *mysqlSingleNamedLock) Release() error {
	name := lock.name
	sqlstr := `select release_lock(?);`
	var result sql.NullInt64
	if err := lock.tx.QueryRow(sqlstr, name).Scan(&result); err != nil {
		lock.tx.Rollback()
		return err
	}

	lock.tx.Commit()

	if result.Valid && result.Int64 == 1 {
		// 1: lock was released.
		return nil
	}
	if result.Valid && result.Int64 == 0 {
		// 0: the lock was not established by this thread (in which case the lock is not released).
		return errors.New("the lock was not established by you")
	}
	// NULL: the named lock did not exist.
	return errors.New("lock does not exist")
}

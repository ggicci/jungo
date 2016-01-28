package db

type NamedLock interface {
	Acquire(timeout int) error // timeout: seconds
	Release() error
	Name() string
}

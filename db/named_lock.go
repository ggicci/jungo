package db

type NamedLock interface {
	Acquire(name string, timeout int) error // timeout: seconds
	Release(name string) error
}

package misc

import (
	"errors"
)

type Checker struct {
	err  error
	desc string
}

func NewChecker() *Checker { return &Checker{} }

func (ck *Checker) Check(desc string, fn func() error) {
	if ck.err != nil {
		return
	}
	ck.desc, ck.err = desc, fn()
}

func (ck *Checker) LastError() error { return ck.err }

func (ck *Checker) String() string {
	if ck.err == nil {
		return ""
	}
	if ck.desc == "" {
		return ck.err.Error()
	}
	return ck.desc + ": " + ck.err.Error()
}

// Returns nil if LastError() is nil, or a new error with desc prefixed.
func (ck *Checker) DecoratedError() error {
	if ck.LastError() == nil {
		return nil
	}
	return errors.New(ck.String())
}

func (ck *Checker) Reset() {
	ck.desc, ck.err = "", nil
}

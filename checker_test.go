package jungo

import (
	"errors"
	"testing"
)

func TestChecker(t *testing.T) {
	err1 := errors.New("can't open file")
	err2 := errors.New("permission not allowed")
	err3 := errors.New("unknown error")

	ck := NewChecker()
	ck.Check("init config", func() error { return err1 })
	ck.Check("init log", func() error { return err2 })
	ck.Check("won't touch", func() error { return nil })

	if ck.LastError() == nil {
		t.Fail()
	}
	t.Log(ck)

	ck.Reset()
	ck.Check("init config", func() error { return nil })
	ck.Check("", func() error { return err3 })
	ck.Check("won't touch", func() error { return nil })
	if ck.LastError() == nil {
		t.Fail()
	}
	t.Log(ck)

	ck.Reset()
	ck.Check("step 1", func() error { return nil })
	ck.Check("step 2", func() error { return nil })
	if ck.LastError() != nil {
		t.Fail()
	}
	t.Log(ck)
}

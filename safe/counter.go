package safe

import (
	"strconv"
	"sync/atomic"
)

type Counter struct{ value int64 }

func NewCounter() *Counter                           { return &Counter{} }
func (sc *Counter) ResetTo(v int64)                  { atomic.StoreInt64(&sc.value, v) }
func (sc *Counter) Value() int64                     { return atomic.LoadInt64(&sc.value) }
func (sc *Counter) Add(delta int64) (newValue int64) { return atomic.AddInt64(&sc.value, delta) }
func (sc *Counter) Sub(delta int64) (newValue int64) { return sc.Add(-delta) }
func (sc *Counter) Inc() (newValue int64)            { return sc.Add(1) }
func (sc *Counter) Dec() (newValue int64)            { return sc.Add(-1) }

func (sc *Counter) MarshalJSON() ([]byte, error)    { return sc.MarshalText() }
func (sc *Counter) UnmarshalJSON(text []byte) error { return sc.UnmarshalText(text) }

func (sc *Counter) MarshalText() ([]byte, error) {
	dst := []byte{}
	dst = strconv.AppendInt(dst, sc.Value(), 10)
	return dst, nil
}

func (sc *Counter) UnmarshalText(text []byte) error {
	v, e := strconv.ParseInt(string(text), 10, 64)
	if e != nil {
		return e
	}
	sc.ResetTo(v)
	return nil
}

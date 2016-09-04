package misc

import (
	"time"
)

type Retryer struct {
	MaxRetryTimes int
	RetryInterval time.Duration

	retryProgressFunc func(retry int, err error)
}

func NewRetryer() *Retryer {
	return &Retryer{
		MaxRetryTimes: 3,
		RetryInterval: time.Second,
	}
}

func (r *Retryer) SetRetryProgressFunc(fn func(retry int, err error)) {
	r.retryProgressFunc = fn
}

func (r *Retryer) Do(fn func() error) (err error) {
	for i := 0; i <= r.MaxRetryTimes; i++ {
		if err = fn(); err == nil {
			break
		}
		time.Sleep(r.RetryInterval)
		if i > 0 && r.retryProgressFunc != nil {
			r.retryProgressFunc(i, err)
		}
	}
	return
}

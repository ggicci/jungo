package log

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"log"
)

type RotatePeriod time.Duration

const (
	Hourly RotatePeriod = RotatePeriod(time.Hour)
	Daily  RotatePeriod = RotatePeriod(time.Hour * 24)

	kHourlyFormatter string = "20060102.15"
	kDailyFormatter  string = "20060102"
)

type TimeRotator struct {
	dir      string
	basename string
	period   RotatePeriod

	bufw *bufio.Writer
	fd   *os.File
	mu   sync.Mutex
}

func NewTimeRotator(filename string, period RotatePeriod) (*TimeRotator, error) {
	rotator := &TimeRotator{
		dir:      filepath.Dir(filename),
		basename: filepath.Base(filename),
		period:   period,
	}

	if period != Hourly && period != Daily {
		return nil, fmt.Errorf("no such period")
	}

	err := rotator.startRotate()

	if err != nil {
		return nil, err
	}

	return rotator, nil
}

func (r *TimeRotator) Write(b []byte) (n int, err error) {
	return r.bufw.Write(b)
}

func (r *TimeRotator) Flush() error {
	return r.bufw.Flush()
}

func (r *TimeRotator) startRotate() error {

	if err := r.rotate(); err != nil {
		return err
	}

	var diff time.Duration
	now := time.Now()

	switch r.period {
	case Hourly:
		diff = now.Truncate(time.Hour).Add(time.Hour).Sub(now)
	case Daily:
		diff = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, 1).Sub(now)
	}

	time.AfterFunc(diff, r.protate)
	return nil
}

func (r *TimeRotator) protate() {
	if err := r.rotate(); err != nil {
		log.Printf("failed to rotate log: %v", err)
	}

	time.AfterFunc(time.Duration(r.period), r.protate)
}

func (r *TimeRotator) rotate() error {
	if _, err := os.Stat(r.dir); os.IsNotExist(err) {
		if err = os.MkdirAll(r.dir, 0755); err != nil {
			return err
		}
	}

	fd, err := os.OpenFile(
		r.formFilename(),
		os.O_CREATE|os.O_RDWR|os.O_APPEND,
		0644,
	)
	if err != nil {
		return err
	}

	// Flush and close the old.
	oldfd := r.fd
	oldbufw := r.bufw
	go func() {
		if oldfd == nil {
			return
		}
		time.Sleep(3 * time.Second)
		oldbufw.Flush()
		oldfd.Close()
	}()

	bufw := bufio.NewWriter(fd)
	// Replace with the new bufw and fd.
	r.fd = fd
	r.bufw = bufw
	return nil
}

func (r *TimeRotator) formFilename() string {
	fullname := filepath.Join(r.dir, r.basename)
	switch r.period {
	case Hourly:
		return fmt.Sprintf("%s.%s", fullname, time.Now().Format(kHourlyFormatter))
	case Daily:
		return fmt.Sprintf("%s.%s", fullname, time.Now().Format(kDailyFormatter))
	}
	return fullname
}

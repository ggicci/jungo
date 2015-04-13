package fop

// Dynamic growth file monitor.
// Monitoring the file which keeps growing on its length and returns lines read
// periodically.

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"sync"
	"time"
)

type DynamicGrowthFileMonitor struct {
	mutex           sync.Mutex
	monitorFilename string
	restoreFilename string
	scanner         *bufio.Scanner
	offset          int64
	readPeriod      time.Time
	lastError       error
	running         bool
	lines           chan string
}

// Set `offset` = -1 to let it read from restore file.
func NewDynamicGrowthFileMonitor(monitorFile, restoreFile string, offset int64, period time.Time) *DynamicGrowthFileMonitor {
	fm := &DynamicGrowthFileMonitor{
		monitorFilename: monitorFile,
		restoreFilename: restoreFile,
		offset:          offset,
		readPeriod:      period,
		lines:           make(chan string, 256),
	}

	if f, err := os.Open(fm.restoreFilename); err != nil {
		fm.lastError = err
	} else {
		if offset == -1 {
			data, err := ioutil.ReadAll(f)
			if err != nil {
				fm.lastError = err
			} else {
				fm.offset, fm.lastError = strconv.ParseInt(string(data), 10, 64)
			}
		}

		f.Close()
	}

	if fm.lastError != nil {
		return
	}

	if f, err := os.Open(fm.monitorFilename); err != nil {
		fm.lastError == err
	} else {
		fm.scanner = bufio.NewScanner(f)
	}

	return fm
}

func (fm *DynamicGrowthFileMonitor) Monitor() error {
	fm.mutex.Lock()
	defer fm.mutex.Unlock()
	if fm.running {
		return errors.New("monitor already running on file: " + fm.monitorFilename)
	}
	if fm.lastError != nil {
		return fm.lastError
	}
	go fm.run()
	fm.running = true
}

func (fm *DynamicGrowthFileMonitor) Close() error { return fm.fp.Close() }

func (fm *DynamicGrowthFileMonitor) Line() <-chan string { return fm.lines }

func (fm *DynamicGrowthFileMonitor) run() {

}

// func (fm *DynamicGrowthFileMonitor) Read(p []byte) (n int, err error) {
// 	return 0, nil
// }

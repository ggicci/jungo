package logrus

import (
	"fmt"
	std "log"
	"os"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"
)

type RollPeriod time.Duration

const (
	ROLL_PERIOD_NONE   = RollPeriod(0)
	ROLL_PERIOD_HOURLY = RollPeriod(time.Hour)
	ROLL_PERIOD_DAILY  = RollPeriod(time.Hour * 24)
)

func RollPeriodString(period RollPeriod) string {
	return map[RollPeriod]string{
		ROLL_PERIOD_NONE:   "none",
		ROLL_PERIOD_HOURLY: "hourly",
		ROLL_PERIOD_DAILY:  "daily",
	}[period]
}

func ParseRollPeriod(period string) (RollPeriod, error) {
	switch period {
	case "none":
		return ROLL_PERIOD_NONE, nil
	case "hourly":
		return ROLL_PERIOD_HOURLY, nil
	case "daily":
		return ROLL_PERIOD_DAILY, nil
	}

	return RollPeriod(0), fmt.Errorf("not a valid roll period: %q", period)
}

func startRolling(logger *log.Logger, cfg *config) {
	var diff time.Duration
	now := time.Now()
	switch cfg.logRollPeriod {
	case ROLL_PERIOD_HOURLY:
		diff = now.Truncate(time.Hour).Add(time.Hour).Sub(now)
	case ROLL_PERIOD_DAILY:
		diff = time.Date(now.Year(), now.Month(), now.Day(),
			0, 0, 0, 0, time.Local).AddDate(0, 0, 1).Sub(now)
	default:
		return
	}
	time.AfterFunc(diff, func() { roll(logger, cfg.LogDir, cfg.LogFilename, cfg.logRollPeriod) })
}

func roll(logger *log.Logger, logDir, logFilename string, period RollPeriod) {
	if newFile, err := os.OpenFile(
		getFilename(logDir, logFilename, period),
		os.O_CREATE|os.O_RDWR|os.O_APPEND,
		0644,
	); err != nil {
		std.Printf("failed rolling log file, error: %v", err)
	} else {
		if oldFile, ok := logger.Out.(*os.File); ok {
			go func() {
				time.Sleep(time.Second * 3)
				oldFile.Close()
			}()
		}
		logger.Out = newFile
	}

	time.AfterFunc(time.Duration(period), func() { roll(logger, logDir, logFilename, period) })
}

func getFilename(logDir, logFilename string, period RollPeriod) string {
	logFilename = filepath.Join(logDir, logFilename)
	switch period {
	case ROLL_PERIOD_NONE:
		return logFilename
	case ROLL_PERIOD_HOURLY:
		return logFilename + "." + time.Now().Format("20060102.15")
	case ROLL_PERIOD_DAILY:
		return logFilename + "." + time.Now().Format("20060102")
	}
	return logFilename
}

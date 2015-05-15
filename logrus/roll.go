package logrus

import (
	"fmt"
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

func startRolling(logger *log.Logger, period RollPeriod) {
	println("start rolling")
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

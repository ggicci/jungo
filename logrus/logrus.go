package logrus

import (
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
)

func NewLoggerFromConfigFile(filename string) (*log.Logger, error) {
	cfg, err := loadConfigs(filename)
	if err != nil {
		return nil, err
	}
	return NewLoggerFromConfig(cfg)
}

// Returns a default logrus logger, with the default settings are:
// 	{
// 		"log_dir": "../logs",
// 		"log_filename": "logrus.log",
// 		"log_level": "info",
// 		"log_roll_period": "none"
// 	}
func NewLogger() (*log.Logger, error) {
	return NewLoggerFromConfig(NewConfig())
}

// Instantiate a logrus logger object from a config.
// The default settings are:
// 	{
// 		"log_dir": "../logs",
// 		"log_filename": "logrus.log",
// 		"log_level": "info",
// 		"log_roll_period": "none"
// 	}
func NewLoggerFromConfig(cfg *config) (*log.Logger, error) {
	cfg = normConfigs(cfg)

	logger := log.New()
	formatter := new(log.TextFormatter)
	formatter.DisableColors = true
	formatter.TimestampFormat = time.RFC3339Nano
	logger.Formatter = formatter
	logger.Level = cfg.logLevel

	_, err := os.Stat(cfg.LogDir)
	if os.IsNotExist(err) {
		if err = os.Mkdir(cfg.LogDir, 0755); err != nil {
			return nil, err
		}
	}

	// Common log file.
	if f, err := os.OpenFile(
		getFilename(cfg.LogDir, cfg.LogFilename, cfg.logRollPeriod),
		os.O_CREATE|os.O_RDWR|os.O_APPEND,
		0644,
	); err != nil {
		return nil, err
	} else {
		logger.Out = f
	}

	// Error log file.
	if cfg.LogErrorFilename != "" {
		errorHook, err := NewErrorHook(getErrorFilename(cfg.LogDir, cfg.LogErrorFilename))
		if err != nil {
			return nil, err
		}
		// Hook the error levels.
		logger.Hooks.Add(errorHook)
	}

	// Roll common logs.
	if cfg.logRollPeriod != ROLL_PERIOD_NONE {
		// Start rolling.
		startRolling(logger, cfg)
	}

	return logger, nil
}

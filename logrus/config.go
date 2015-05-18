package logrus

import (
	"encoding/json"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/ggicci/jungo/program"
)

var std = log.StandardLogger()

func NewConfig() *config { return &config{} }

type config struct {
	// The directory in which the log files stay.
	LogDir string `json:"log_dir"`

	// The basic filename of the logs. If not set, defaults to "logrus.log".
	// And according to different roll period, final filenames are:
	//
	// Period    	Filename
	// none		 	logrus.go
	// hourly		logrus.20060102.15
	// daily		logrus.20060102
	//
	// NB: slash("/") in the filename will be replace with an underscore character("_"),
	// and dots(".") and spaces will be trimmed.
	LogFilename string `json:"log_filename"`

	// Like, `LogFilename`.
	// If set, logs with level >= "error" will also be written to the error log file.
	// No rolling.
	LogErrorFilename string `json:"log_error_filename"`

	// Log level, valid levels are: "panic", "fatal", "error", "warn" or "warning", "info", "debug".
	LogLevelString string `json:"log_level"`

	// Log roll period, "none", "hourly", "daily".
	LogRollPeriodString string `json:"log_roll_period"`

	logLevel      log.Level
	logRollPeriod RollPeriod
}

func loadConfigs(filename string) (*config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg := NewConfig()
	if err = json.NewDecoder(f).Decode(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func normConfigs(cfg *config) *config {
	// Check the config items.

	if cfg.LogDir = strings.TrimSpace(cfg.LogDir); cfg.LogDir == "" {
		cfg.LogDir = program.AbsPath("../logs")
		std.Warnf("log_dir not set, defaults to %q", cfg.LogDir)
	} else {
		cfg.LogDir = program.AbsPath(cfg.LogDir)
	}

	cfg.LogFilename = strings.TrimSpace(cfg.LogFilename)
	cfg.LogFilename = strings.Replace(cfg.LogFilename, "/", "_", -1)
	cfg.LogFilename = strings.Trim(cfg.LogFilename, ".")
	if cfg.LogFilename == "" {
		cfg.LogFilename = "logrus.log"
		std.Warnf("log_filename not set, defaults to %q", cfg.LogFilename)
	}

	cfg.LogErrorFilename = strings.TrimSpace(cfg.LogErrorFilename)
	cfg.LogErrorFilename = strings.Replace(cfg.LogErrorFilename, "/", "_", -1)
	cfg.LogErrorFilename = strings.Trim(cfg.LogErrorFilename, ".")
	if cfg.LogErrorFilename == cfg.LogFilename {
		cfg.LogErrorFilename = cfg.LogFilename + ".error"
	}

	if level, err := log.ParseLevel(cfg.LogLevelString); err != nil {
		cfg.LogLevelString = LogLevelString(log.InfoLevel)
		std.Warnf("%s, reset to %q", err, cfg.LogLevelString)
		cfg.logLevel = log.InfoLevel
	} else {
		cfg.logLevel = level
	}

	if period, err := ParseRollPeriod(cfg.LogRollPeriodString); err != nil {
		cfg.LogRollPeriodString = RollPeriodString(ROLL_PERIOD_NONE)
		std.Warnf("%s, reset to %q", err, cfg.LogRollPeriodString)
		cfg.logRollPeriod = ROLL_PERIOD_NONE
	} else {
		cfg.logRollPeriod = period
	}

	return cfg
}

func LogLevelString(level log.Level) string {
	return map[log.Level]string{
		log.PanicLevel: "panic",
		log.FatalLevel: "fatal",
		log.ErrorLevel: "error",
		log.WarnLevel:  "warn",
		log.InfoLevel:  "info",
		log.DebugLevel: "debug",
	}[level]
}

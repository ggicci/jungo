package logrus

import (
	"encoding/json"
	std "log"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/ggicci/jungo/path"
)

func NewConfig() *config { return &config{} }

type config struct {
	LogDir              string `json:"log_dir"`
	LogFilename         string `json:"log_filename"`
	LogLevelString      string `json:"log_level"`
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
		cfg.LogDir = path.AbsPath("../logs")
		std.Printf("log_dir not set, defaults to %q", cfg.LogDir)
	}

	cfg.LogFilename = strings.TrimSpace(cfg.LogFilename)
	if cfg.LogFilename = strings.Trim(cfg.LogFilename, "."); cfg.LogFilename == "" {
		cfg.LogFilename = "logrus.log"
		std.Printf("log_filename not set, defaults to %q", cfg.LogFilename)
	}

	if level, err := log.ParseLevel(cfg.LogLevelString); err != nil {
		cfg.LogLevelString = LogLevelString(log.InfoLevel)
		std.Printf("%s, reset to %q", err, cfg.LogLevelString)
		cfg.logLevel = log.InfoLevel
	} else {
		cfg.logLevel = level
	}

	if period, err := ParseRollPeriod(cfg.LogRollPeriodString); err != nil {
		cfg.LogRollPeriodString = RollPeriodString(ROLL_PERIOD_NONE)
		std.Printf("%s, reset to %q", err, cfg.LogRollPeriodString)
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

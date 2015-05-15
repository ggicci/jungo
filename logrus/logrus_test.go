package logrus

import (
	"testing"
)

func TestNewLogrusLogger(t *testing.T) {
	cfg := NewConfig()
	cfg.LogRollPeriodString = "hourly"
	logger, err := NewLoggerFromConfig(cfg)
	if err != nil {
		t.Fatalf("failed to create a new logrus logger: %v", err)
	}
	logger.WithField("name", "Ggicci").Info("He is funny.")
}

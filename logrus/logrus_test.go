package logrus

import (
	"testing"
	"time"
)

func TestNewLogrusLogger(t *testing.T) {
	cfg := NewConfig()
	cfg.LogRollPeriodString = "daily"
	logger, err := NewLogrusLoggerFromConfig(cfg)
	if err != nil {
		t.Fatalf("failed to create a new logrus logger: %v", err)
	}
	logger.WithField("name", "Ggicci").Info("He is funny.")

	time.Sleep(time.Second)
}

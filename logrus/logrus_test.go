package logrus

import (
	"errors"
	"testing"
)

func TestNewLogrusLogger(t *testing.T) {
	logger, err := NewLoggerFromConfigFile("./config.json")
	if err != nil {
		t.Fatalf("failed to create a new logrus logger: %v", err)
	}
	logger.WithField("game", "cross fire").Info("a fun game")
	logger.WithField("error", errors.New("unable to open file")).Error("initialize a logger")
}

package logrus

import (
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
)

type ErrorHook struct {
	logger *log.Logger
}

func NewErrorHook(filename string) (*ErrorHook, error) {
	logger := log.New()
	if f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644); err != nil {
		return nil, err
	} else {
		formatter := new(log.TextFormatter)
		formatter.DisableColors = true
		formatter.TimestampFormat = time.RFC3339Nano
		logger.Formatter = formatter
		logger.Out = f
		return &ErrorHook{logger}, nil
	}
}

func (eh *ErrorHook) Levels() []log.Level {
	return []log.Level{log.ErrorLevel, log.FatalLevel, log.PanicLevel}
}

func (eh *ErrorHook) Fire(entry *log.Entry) error {
	switch entry.Level {
	case log.ErrorLevel:
		eh.logger.WithFields(entry.Data).Error(entry.Message)
	case log.FatalLevel:
		eh.logger.WithFields(entry.Data).Fatal(entry.Message)
	case log.PanicLevel:
		eh.logger.WithFields(entry.Data).Panic(entry.Message)
	}
	return nil
}

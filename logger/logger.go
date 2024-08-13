package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type (
	Fields = logrus.Fields
)

type Logger struct {
	logger *logrus.Entry
}

var l = logrus.New()

func Init(environment string) {
	var (
		isRelease   = environment == "release"
		level       = logrus.DebugLevel
		forceColors = true
	)

	if isRelease {
		level = logrus.InfoLevel
		forceColors = false
	}

	l.SetLevel(level)
	l.SetFormatter(&logrus.TextFormatter{
		ForceColors:     forceColors,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		DisableQuote:    true,
	})

	// specific level for test
	if environment == "test" {
		l.Level = logrus.FatalLevel
	}

	fmt.Printf("⚡️ [logger]: initialized \n")
}

func NewLogger(data Fields) *Logger {
	return &Logger{
		logger: l.WithFields(data),
	}
}

func (l *Logger) AddData(fields Fields) {
	l.logger = l.logger.WithFields(mergeFields(l.logger.Data, fields))
}

func mergeFields(fields ...Fields) Fields {
	result := make(map[string]interface{})
	for _, m := range fields {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

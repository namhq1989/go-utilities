package logger

import (
	"encoding/json"
	"fmt"
)

func (l *Logger) Debug(message string, fields Fields) {
	l.logger.WithFields(fields).Debug(message)
}

func (l *Logger) Error(message string, err error, fields Fields) {
	if err == nil {
		l.logger.WithFields(fields).Error(message)
	} else {
		l.logger.WithFields(fields).Error(fmt.Sprintf("%s: %s", message, err.Error()))
	}
}

func (l *Logger) Info(message string, fields Fields) {
	l.logger.WithFields(fields).Info(message)
}

func (l *Logger) Text(message string) {
	l.logger.Info(message)
}

func (l *Logger) ErrorText(message string) {
	l.logger.Error(message)
}

func (*Logger) Print(message string, data interface{}) {
	s, _ := json.MarshalIndent(data, "", "  ")
	fmt.Printf("%s: %+v\n", message, string(s))
}

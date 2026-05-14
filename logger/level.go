package logger

import (
	"encoding/json"
	"fmt"
)

func (l *Logger) Debug(message string, fields Fields) {
	l.logger.WithFields(RedactFields(fields)).Debug(message)
}

// Error logs at error level. The raw err is attached via logrus.WithError so
// downstream hooks (e.g. Sentry) can retrieve the original error value with
// its type and chain intact via entry.Data[logrus.ErrorKey].
//
// `message` MUST describe WHAT failed (e.g. "failed to send email via Brevo").
// Pass the original `err` — never err.Error() or fmt.Errorf string wraps that
// destroy the chain, which would break Sentry's type-based issue grouping.
func (l *Logger) Error(message string, err error, fields Fields) {
	entry := l.logger.WithFields(RedactFields(fields))
	if err != nil {
		entry = entry.WithError(err)
	}
	entry.Error(message)
}

// SilentError logs at error level WITHOUT triggering external alert hooks.
// Use in shared/library layers to avoid duplicate alerts — the caller in the
// service layer is expected to translate the returned error into Error().
func (l *Logger) SilentError(message string, err error, fields Fields) {
	merged := RedactFields(fields)
	if merged == nil {
		merged = Fields{}
	}
	merged[MarkerSilent] = true
	entry := l.logger.WithFields(merged)
	if err != nil {
		entry = entry.WithError(err)
	}
	entry.Error(message)
}

func (l *Logger) Info(message string, fields Fields) {
	l.logger.WithFields(RedactFields(fields)).Info(message)
}

func (l *Logger) Text(message string) {
	l.logger.Info(message)
}

func (l *Logger) Warn(message string, fields Fields) {
	l.logger.WithFields(RedactFields(fields)).Warn(message)
}

func (l *Logger) ErrorText(message string) {
	l.logger.Error(message)
}

func (*Logger) Print(message string, data interface{}) {
	s, _ := json.MarshalIndent(data, "", "  ")
	fmt.Printf("%s: %+v\n", message, string(s))
}

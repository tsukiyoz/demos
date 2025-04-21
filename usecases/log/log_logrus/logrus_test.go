package loglogrus

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
)

type contextKey string

const debugModeKey contextKey = "__logger_debug_mode__"

type ForceLogHook struct{}

func (h *ForceLogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func isForceLog(ctx context.Context) bool {
	if ctx == nil {
		return false
	}
	v, ok := ctx.Value(debugModeKey).(bool)
	return ok && v
}

func (h *ForceLogHook) Fire(entry *logrus.Entry) error {
	if isForceLog(entry.Context) {
		entry.Level = logrus.DebugLevel
	}
	return nil
}

func TestLogrus_Usage(t *testing.T) {
	l := logrus.New()
	l.Print("Hello, Logrus!")
	l.WithField("key", "value").Info("This is an info message")

	// Level Usage
	l.SetLevel(logrus.InfoLevel)
	l.Info("This is an info message")
	l.Debug("This debug message will not be shown because the level is set to Info")

	l.SetLevel(logrus.DebugLevel)
	l.Debug("This debug message will be shown because the level is set to Debug")

	// hook usage
	hook := &ForceLogHook{}
	l.AddHook(hook)
	l.SetLevel(logrus.InfoLevel) // Set level to Info, so that the debug messages are not shown
	l.Debug("This debug message will not be shown because the level is set to Info")
	l.WithContext(context.WithValue(context.Background(), debugModeKey, true)).Debug("This debug message will be shown because the context is set to debug mode")
}

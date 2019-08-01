package log

import (
	"sync"

	"github.com/sirupsen/logrus"
)

// Level type
type Level uint32

// Fields type, used to pass to `WithFields`.
type Fields map[string]interface{}

// ILoggerHook todo
type ILoggerHook interface {
	Fire(entry *Entry)
}

// Hook todo...
type Hook struct {
	logger *Logger
	hooks  []ILoggerHook
	hlock  sync.Mutex
}

func newHook(logger *Logger) *Hook {
	hook := &Hook{
		logger: logger,
	}

	// Add "us" as a hook to logrus, and we will delegate Fire to our consumers in Fire()
	hook.logger.logrus.AddHook(hook)
	return hook
}

// AddHook todo...
func (h *Hook) AddHook(logHook ILoggerHook) {
	h.hlock.Lock()
	defer h.hlock.Unlock()

	h.hooks = append(h.hooks, logHook)
}

// Fire todo...
func (h *Hook) Fire(entry *logrus.Entry) error {
	logEntry := convertEntry(entry)

	h.hlock.Lock()
	defer h.hlock.Unlock()

	for _, hook := range h.hooks {
		if h.logger.logrus.IsLevelEnabled(entry.Level) {
			hook.Fire(logEntry)
		}
	}

	return nil
}

// Levels todo...
func (h *Hook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Common logging routines
func convertEntry(entry *logrus.Entry) *Entry {
	e := &Entry{
		Time:    entry.Time,
		Level:   (Level)(entry.Level),
		Caller:  entry.Caller,
		Message: entry.Message,
	}

	e.Fields = convertFieldsFromLogrus(entry.Data)
	return e
}

func convertFieldsFromLogrus(fields logrus.Fields) Fields {

	Fields := make(Fields)
	for k, v := range fields {
		Fields[k] = v
	}
	return Fields
}

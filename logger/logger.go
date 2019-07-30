package logger

import (
	"github.com/sirupsen/logrus"
)

// Common logging routines
func convertEntryToLogEntry(entry *logrus.Entry) *LogEntry {
	logEntry := &LogEntry{
		Time:    entry.Time,
		Level:   (LogLevel)(entry.Level),
		Caller:  entry.Caller,
		Message: entry.Message,
	}

	logEntry.Fields = convertDataToLogData(entry.Data)

	return logEntry
}

func convertDataToLogData(fields logrus.Fields) LogFields {

	logFields := make(LogFields)
	for k, v := range fields {
		logFields[k] = v
	}
	return logFields
}

package log

// ILoggerHook todo
type ILoggerHook interface {
	Fire(entry *LogEntry)
}

package log

// Scope allows you to set types that can be re-used for subsequent log event. Useful for setting username, requestid etc for a Http Web Request.
type Scope struct {
	logger *FieldLogger
	fields Fields
}

func newScope(logger *FieldLogger, fields Fields) *Scope {
	return &Scope{
		logger: logger,
		fields: fields,
	}
}

// Debug writes a debug message with optional types to the underlying standard logger.
// Useful for adding detailed tracing that you don't normally want to appear, but turned on
// when hunting down incorrect behaviour.
// All types values will be automatically quoted (keys will not be).
// Debug adds types {level="debug", time="2006-01-02T15:04:05Z07:00"}
// and prints output in the format "types message"
// Use lower-case keys and values if possible.
func (scope Scope) Debug(message string, fields ...Fields) {
	merged := scope.fields.merge(fields...)
	scope.logger.Debug(message, merged)
}

// Print writes a message with optional types to the underlying standard logger.
// Useful to normal tracing that should be captured during standard operating behaviour.
// All types values will be automatically quoted (keys will not be).
// Debug adds types {time="2006-01-02T15:04:05Z07:00"}
// and prints output in the format "types message"
// Use lower-case keys and values if possible.
func (scope Scope) Print(message string, fields ...Fields) {
	merged := scope.fields.merge(fields...)
	scope.logger.Print(message, merged)
}

// Error writes a error message with optional types to the underlying standard logger.
// Useful to trace errors should be captured always.
// All types values will be automatically quoted (keys will not be).
// Debug adds types {level="error", time="2006-01-02T15:04:05Z07:00"}
// and prints output in the format "types message"
// Use lower-case keys and values if possible.
func (scope Scope) Error(err error, fields ...Fields) {
	merged := scope.fields.merge(fields...)
	scope.logger.Error(err, merged)
}

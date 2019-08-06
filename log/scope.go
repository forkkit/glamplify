package log

// FieldLogger wraps the standard library logger and add structured fields as quoted key value pairs
type Scope struct {
	logger 	*FieldLogger
	fields   Fields
}

func newScope(logger *FieldLogger, fields Fields) *Scope {
	return &Scope{
		logger: logger,
		fields: fields,
	}
}

// Debug writes a debug message with optional fields to the underlying standard logger.
// Useful for adding detailed tracing that you don't normally want to appear, but turned on
// when hunting down incorrect behaviour.
// All field values will be automatically quoted (keys will not be).
// Debug adds fields {level="debug", time="2006-01-02T15:04:05Z07:00"}
// and prints output in the format "fields message"
// Use lower-case keys and values if possible.
func (scope Scope) Debug(message string, fields ...Fields) error {
	merged := scope.fields.merge(fields...)
	return scope.logger.Debug(message, merged)
}

// Print writes a message with optional fields to the underlying standard logger.
// Useful to normal tracing that should be captured during standard operating behaviour.
// All field values will be automatically quoted (keys will not be).
// Debug adds fields {time="2006-01-02T15:04:05Z07:00"}
// and prints output in the format "fields message"
// Use lower-case keys and values if possible.
func (scope Scope) Print(message string, fields ...Fields) error {
	merged := scope.fields.merge(fields...)
	return scope.logger.Print(message, merged)
}

// Error writes a error message with optional fields to the underlying standard logger.
// Useful to trace errors should be captured always.
// All field values will be automatically quoted (keys will not be).
// Debug adds fields {level="error", time="2006-01-02T15:04:05Z07:00"}
// and prints output in the format "fields message"
// Use lower-case keys and values if possible.
func (scope Scope) Error(err error, fields ...Fields) error {
	merged := scope.fields.merge(fields...)
	return scope.logger.Error(err, merged)
}
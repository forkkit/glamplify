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
// Use lower-case keys and values if possible.
func (scope Scope) Debug(message string, fields ...Fields) {
	merged := scope.fields.Merge(fields...)
	scope.logger.Debug(message, merged)
}

// Info writes a message with optional types to the underlying standard logger.
// Useful for normal tracing that should be captured during standard operating behaviour.
// Use lower-case keys and values if possible.
func (scope Scope) Info(message string, fields ...Fields) {
	merged := scope.fields.Merge(fields...)
	scope.logger.Info(message, merged)
}

// Warn writes a message with optional types to the underlying standard logger.
// Useful for unusual but recoverable tracing that should be captured during standard operating behaviour.
// Use lower-case keys and values if possible.
func  (scope Scope) Warn(message string, fields ...Fields) {
	merged := scope.fields.Merge(fields...)
	scope.logger.Warn(message, merged)
}

// Error writes a error message with optional types to the underlying standard logger.
// Useful to trace errors that are usually not recoverable. These should always be logged.
// Use lower-case keys and values if possible.
func (scope Scope) Error(err error, fields ...Fields) {
	merged := scope.fields.Merge(fields...)
	scope.logger.Error(err, merged)
}

// Fatal writes a error message with optional types to the underlying standard logger and then calls panic!
// Panic will terminate the current go routine.
// Useful to trace catastrophic errors that are not recoverable. These should always be logged.
// Use lower-case keys and values if possible.
func (scope Scope) Fatal(err error, fields ...Fields) {
	merged := scope.fields.Merge(fields...)
	scope.logger.Fatal(err, merged)
}
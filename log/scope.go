package log

import "context"

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
func (scope Scope) Debug(ctx context.Context, event string, fields ...Fields) {
	merged := scope.fields.Merge(fields...)
	scope.logger.Debug(ctx, event, merged)
}

// Info writes a message with optional types to the underlying standard logger.
// Useful for normal tracing that should be captured during standard operating behaviour.
// Use lower-case keys and values if possible.
func (scope Scope) Info(ctx context.Context, event string, fields ...Fields) {
	merged := scope.fields.Merge(fields...)
	scope.logger.Info(ctx, event, merged)
}

// Warn writes a message with optional types to the underlying standard logger.
// Useful for unusual but recoverable tracing that should be captured during standard operating behaviour.
// Use lower-case keys and values if possible.
func  (scope Scope) Warn(ctx context.Context, event string, fields ...Fields) {
	merged := scope.fields.Merge(fields...)
	scope.logger.Warn(ctx, event, merged)
}

// Error writes a error message with optional types to the underlying standard logger.
// Useful to trace errors that are usually not recoverable. These should always be logged.
// Use lower-case keys and values if possible.
func (scope Scope) Error(ctx context.Context, err error, fields ...Fields) {
	merged := scope.fields.Merge(fields...)
	scope.logger.Error(ctx, err, merged)
}

// Fatal writes a error message with optional types to the underlying standard logger and then calls panic!
// Panic will terminate the current go routine.
// Useful to trace catastrophic errors that are not recoverable. These should always be logged.
// Use lower-case keys and values if possible.
func (scope Scope) Fatal(ctx context.Context, err error, fields ...Fields) {
	merged := scope.fields.Merge(fields...)
	scope.logger.Fatal(ctx, err, merged)
}
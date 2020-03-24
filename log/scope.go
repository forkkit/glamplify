package log

import (
	"context"
	"github.com/cultureamp/glamplify/constants"
	"strings"
)

// Scope allows you to set types that can be re-used for subsequent log event. Useful for setting username, requestid etc for a Http Web Request.
type Scope struct {
	ctx        context.Context
	logger     *FieldLogger
	fields     Fields
	defValues  *DefaultValues
}

func newScope(ctx context.Context, logger *FieldLogger, fields ...Fields) *Scope {
	merged := Fields{}
	merged = merged.Merge(fields...)
	scope := &Scope{
		ctx:        ctx,
		logger:     logger,
		fields:     merged,
	}
	scope.defValues = newDefaultValues()
	return scope
}

// Debug writes a debug message with optional types to the underlying standard logger.
// Useful for adding detailed tracing that you don't normally want to appear, but turned on
// when hunting down incorrect behaviour.
// Use lower-case keys and values if possible.
func (scope Scope) Debug(event string, fields ...Fields) {
	meta := scope.defValues.getDefaults(scope.ctx, event, constants.DebugSevLogValue)
	merged := scope.fields.Merge(fields...)
	scope.logger.debug(event, meta, merged)
}

// Info writes a message with optional types to the underlying standard logger.
// Useful for normal tracing that should be captured during standard operating behaviour.
// Use lower-case keys and values if possible.
func (scope Scope) Info(event string, fields ...Fields) {
	meta := scope.defValues.getDefaults(scope.ctx, event, constants.InfoSevLogValue)
	merged := scope.fields.Merge(fields...)
	scope.logger.info(event, meta, merged)
}

// Warn writes a message with optional types to the underlying standard logger.
// Useful for unusual but recoverable tracing that should be captured during standard operating behaviour.
// Use lower-case keys and values if possible.
func (scope Scope) Warn(event string, fields ...Fields) {
	meta := scope.defValues.getDefaults(scope.ctx, event, constants.WarnSevLogValue)
	merged := scope.fields.Merge(fields...)
	scope.logger.warn(event, meta, merged)
}

// Error writes a error message with optional types to the underlying standard logger.
// Useful to trace errors that are usually not recoverable. These should always be logged.
// Use lower-case keys and values if possible.
func (scope Scope) Error(err error, fields ...Fields) {
	event := strings.TrimSpace(err.Error())
	meta := scope.defValues.getDefaults(scope.ctx, event, constants.ErrorSevLogValue)
	meta = scope.defValues.getErrorDefaults(err, meta)
	merged := scope.fields.Merge(fields...)
	scope.logger.error(event, meta, merged)
}

// Fatal writes a error message with optional types to the underlying standard logger and then calls panic!
// Panic will terminate the current go routine.
// Useful to trace catastrophic errors that are not recoverable. These should always be logged.
// Use lower-case keys and values if possible.
func (scope Scope) Fatal(err error, fields ...Fields) {
	event := strings.TrimSpace(err.Error())
	meta := scope.defValues.getDefaults(scope.ctx, event, constants.FatalSevLogValue)
	meta = scope.defValues.getErrorDefaults(err, meta)
	merged := scope.fields.Merge(fields...)
	scope.logger.fatal(event, meta, merged)
}

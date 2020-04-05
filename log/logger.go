package log

import (
	"context"
	"github.com/cultureamp/glamplify/helper"
	"github.com/cultureamp/glamplify/jwt"
	"net/http"
	"strings"
)

// Logger
type Logger struct {
	ctx       context.Context
	writer    *FieldWriter
	fields    Fields
	defValues *DefaultValues
}

var (
	internal = newWriter(func(conf *config) {})
)

// NewFromRequest creates a new logger and does all the good things
// like setting the current user, customer, etc from decoding the JWT on the request (if present)
// The error returned indicates a problem with decoding the JWT, but a new *Logger is always returned regardless of error
func NewFromRequest(ctx context.Context, r *http.Request, fields ...Fields) (*Logger, error) {

	token := r.Header.Get("Authorization") // "Authorization: Bearer xxxxx.yyyyy.zzzzz"
	if len(token) > 0 {

		splitToken := strings.Split(token, "Bearer")
		token = splitToken[1]

		jwt, err := jwt.NewDecoder()
		if err != nil {
			return New(ctx, fields...), err
		}

		payload, err := jwt.Decode(token)
		if err != nil {
			return  New(ctx, fields...), err
		}

		ctx = AddCustomer(ctx, payload.Customer)
		ctx = AddUser(ctx, payload.EffectiveUser)
	}

	return New(ctx, fields...), nil
}

// New creates a *Logger with optional fields. Useful for when you want to add a field to all subsequent logging calls eg. request_id, etc.
func New(ctx context.Context, fields ...Fields) *Logger {
	return newLogger(ctx, internal, fields...)
}

func newLogger(ctx context.Context, writer *FieldWriter, fields ...Fields) *Logger {
	merged := Fields{}
	merged = merged.Merge(fields...)
	scope := &Logger{
		ctx:    ctx,
		writer: writer,
		fields: merged,
	}
	scope.defValues = newDefaultValues()
	return scope
}

// Debug writes a debug message with optional types to the underlying standard writer.
// Useful for adding detailed tracing that you don't normally want to appear, but turned on
// when hunting down incorrect behaviour.
// Use lower-case keys and values if possible.
func (logger Logger) Debug(event string, fields ...Fields) {
	event = helper.ToSnakeCase(event)
	meta := logger.defValues.getDefaults(logger.ctx, event, DebugSev)
	merged := logger.fields.Merge(fields...)
	logger.writer.debug(event, meta, merged)
}

// Info writes a message with optional types to the underlying standard writer.
// Useful for normal tracing that should be captured during standard operating behaviour.
// Use lower-case keys and values if possible.
func (logger Logger) Info(event string, fields ...Fields) {
	event = helper.ToSnakeCase(event)
	meta := logger.defValues.getDefaults(logger.ctx, event, InfoSev)
	merged := logger.fields.Merge(fields...)
	logger.writer.info(event, meta, merged)
}

// Warn writes a message with optional types to the underlying standard writer.
// Useful for unusual but recoverable tracing that should be captured during standard operating behaviour.
// Use lower-case keys and values if possible.
func (logger Logger) Warn(event string, fields ...Fields) {
	event = helper.ToSnakeCase(event)
	meta := logger.defValues.getDefaults(logger.ctx, event, WarnSev)
	merged := logger.fields.Merge(fields...)
	logger.writer.warn(event, meta, merged)
}

// Error writes a error message with optional types to the underlying standard writer.
// Useful to trace errors that are usually not recoverable. These should always be logged.
// Use lower-case keys and values if possible.
func (logger Logger) Error(err error, fields ...Fields) {
	event := helper.ToSnakeCase(err.Error())
	meta := logger.defValues.getDefaults(logger.ctx, event, ErrorSev)
	meta = logger.defValues.getErrorDefaults(err, meta)
	merged := logger.fields.Merge(fields...)
	logger.writer.error(event, meta, merged)
}

// Fatal writes a error message with optional types to the underlying standard writer and then calls panic!
// Panic will terminate the current go routine.
// Useful to trace catastrophic errors that are not recoverable. These should always be logged.
// Use lower-case keys and values if possible.
func (logger Logger) Fatal(err error, fields ...Fields) {
	event := helper.ToSnakeCase(err.Error())
	meta := logger.defValues.getDefaults(logger.ctx, event, FatalSev)
	meta = logger.defValues.getErrorDefaults(err, meta)
	merged := logger.fields.Merge(fields...)
	logger.writer.fatal(event, meta, merged)
}

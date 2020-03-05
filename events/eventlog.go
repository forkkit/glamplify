package events

import (
	"context"
	"errors"
	"github.com/cultureamp/glamplify/constants"
	"github.com/cultureamp/glamplify/helper"
	"github.com/cultureamp/glamplify/log"
	"io"
	"os"
	"sync"
)

// Config for setting initial values for EventLog
type Config struct {
	Output      io.Writer
	Product     string
	Application string
}

// EventLog
type EventLog struct {
	mutex       sync.Mutex
	product     string
	application string
	log         *log.FieldLogger
}

// So that you don't even need to create a new logger
var (
	internal = NewEventLog(func(conf *Config) {})
)

// NewEventLog creates a new FieldLogger. The optional configure func lets you set values on the underlying standard logger.
// eg. SetOutput
func NewEventLog(configure ...func(*Config)) *EventLog { // https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis

	eventLog := &EventLog{}
	conf := Config{
		Output: os.Stdout,
	}

	for _, config := range configure {
		config(&conf)
	}

	eventLog.mutex.Lock()
	defer eventLog.mutex.Unlock()

	eventLog.application = conf.Application
	eventLog.product = conf.Product
	eventLog.log = log.New(func(c *log.Config) { c.Output = conf.Output })

	return eventLog
}

// Audit emits a event log entry as per https://cultureamp.atlassian.net/wiki/spaces/TV/pages/959939199/Logging
func Audit(event string, success bool, ctx context.Context, fields log.Fields) {
	internal.Audit(event, success, ctx, fields)
}

// Audit emits a event log entry as per https://cultureamp.atlassian.net/wiki/spaces/TV/pages/959939199/Logging
func (eventlog EventLog) Audit(event string, success bool, ctx context.Context, fields log.Fields) {

	// add event and success fields
	fields[constants.EventLog] = event
	fields[constants.StatusLog] = success

	// add config defaults
	fields = eventlog.addConfigDefaults(fields)

	// add in any missing fields from context and then os env
	fields = eventlog.addIfMissing(ctx, fields)

	if success {
		eventlog.log.Info(constants.EmptyString, fields)
	} else {
		eventlog.log.Error(errors.New(constants.EmptyString), fields)
	}
}

func (eventlog EventLog) addConfigDefaults(fields log.Fields) log.Fields {
	if _, ok := fields[constants.ProductLog]; !ok {
		if eventlog.product != "" {
			fields[constants.ProductLog] = eventlog.product
		}
	}
	if _, ok := fields[constants.AppLog]; !ok {
		if eventlog.application != "" {
			fields[constants.AppLog] = eventlog.application
		}
	}

	return fields
}

func (eventlog EventLog) addIfMissing(ctx context.Context, fields log.Fields) log.Fields {

	fields = eventlog.addFieldIfMissingOrDefault(constants.ProductLog, constants.ProductEnv, constants.ProductCtx, ctx, fields, constants.UnknownString)
	fields = eventlog.addFieldIfMissingOrDefault(constants.AppLog, constants.AppEnv, constants.AccountCtx, ctx, fields, constants.UnknownString)
	fields = eventlog.addFieldIfMissingOrDefault(constants.TraceIdLog, constants.TraceIdEnv, constants.TraceIdCtx, ctx, fields, helper.NewTraceID())
	fields = eventlog.addFieldIfMissingOrDefault(constants.ModuleLog, constants.ModuleEnv, constants.ModuleCtx, ctx, fields, constants.UnknownString)
	fields = eventlog.addFieldIfMissingOrDefault(constants.AccountLog, constants.AccountEnv, constants.AccountCtx, ctx, fields, constants.UnknownString)
	fields = eventlog.addFieldIfMissingOrDefault(constants.UserLog, constants.UserEnv, constants.UserCtx, ctx, fields, constants.UnknownString)

	return fields
}

func (eventlog EventLog) addFieldIfMissingOrDefault(
	fieldName string,
	osVar string,
	ctxKey constants.EventCtxKey,
	ctx context.Context,
	fields log.Fields,
	defValue string) log.Fields {

	// If it contains it already, all good!
	if _, ok := fields[fieldName]; ok {
		return fields
	}

	// first check context
	if prod, ok := ctx.Value(ctxKey).(string); ok {
		fields[fieldName] = prod
		return fields
	}

	// next, check env
	if prod, ok := os.LookupEnv(osVar); ok {
		fields[fieldName] = prod
		return fields
	}

	// how else?

	// still missing, so add default
	fields[fieldName] = defValue
	return fields
}

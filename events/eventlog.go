package events

import (
	"context"
	"github.com/cultureamp/glamplify/types"
	"github.com/cultureamp/glamplify/log"
	"os"
	"sync"
)



// Config for setting initial values for EventLog
type Config struct {
	Product     string
	Application string
}

// EventLog
type EventLog struct {
	mutex       sync.Mutex
	product     string
	application string
}

// So that you don't even need to create a new logger
var (
	internal = New(func(conf *Config) {})
)

// New creates a new FieldLogger. The optional configure func lets you set values on the underlying standard logger.
// eg. SetOutput
func New(configure ...func(*Config)) *EventLog { // https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis

	eventLog := &EventLog{}
	conf := Config{
	}

	for _, config := range configure {
		config(&conf)
	}

	eventLog.mutex.Lock()
	defer eventLog.mutex.Unlock()

	eventLog.application = conf.Application
	eventLog.product = conf.Product

	return eventLog
}

func Audit(event string, success bool, ctx context.Context, fields types.Fields) {
	internal.Audit(event, success, ctx, fields)
}

func (eventlog EventLog) Audit(event string, success bool, ctx context.Context, fields types.Fields) {

	// todo - add event and success fields

	// add in any missing fields from context and then os env
	fields = eventlog.addIfMissing(ctx, fields)

	//log.Print()
}

func (eventlog EventLog) addIfMissing(ctx context.Context, fields types.Fields) types.Fields {

	fields = eventlog.addFieldIfMissing(log.PRODUCT, ProductEnv, ctx, fields)
	fields = eventlog.addFieldIfMissing(log.APP, AppEnv, ctx, fields)
	fields = eventlog.addFieldIfMissing(log.TRACE_ID, TraceIdEnv, traceId, ctx, fields)
	fields = eventlog.addFieldIfMissing(log.MODULE, ModuleEnv, ctx, fields)
	fields = eventlog.addFieldIfMissing(log.ACCOUNT, AccountEnv, ctx, fields)
	fields = eventlog.addFieldIfMissing(log.USER, UserEnv, ctx, fields)

	return fields
}

func (eventlog EventLog) addFieldIfMissing(fieldName string, osVar string, ctxKey eventKey, ctx context.Context, fields types.Fields) types.Fields {

	// If it contains
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

	return fields
}

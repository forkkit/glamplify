package events

import (
	"bytes"
	"context"
	"github.com/cultureamp/glamplify/constants"
	"github.com/cultureamp/glamplify/helper"
	"github.com/cultureamp/glamplify/log"
	"gotest.tools/assert"
	"os"
	"strings"
	"testing"
	"time"
)

func Test_EventLog_AlreadyPresent(t *testing.T) {

	eventLog := NewEventLog(func(config *Config) {})

	ctx := context.TODO()
	fields := log.Fields{"APP": "myapp"}
	fields = eventLog.addFieldIfMissingOrDefault("APP", "APP_VAR", constants.AppCtx, ctx, fields, "default")

	assert.Assert(t, fields["APP"] == "myapp")
}

func Test_EventLog_MissingAddDefault(t *testing.T) {

	eventLog := NewEventLog(func(config *Config) {})

	ctx := context.TODO()
	fields := log.Fields{}
	fields = eventLog.addFieldIfMissingOrDefault("APP", "APP_VAR", constants.AppCtx, ctx, fields, "default")

	assert.Assert(t, fields["APP"] == "default")
}

func Test_EventLog_Missing_InEnv(t *testing.T) {

	eventLog := NewEventLog(func(config *Config) {})

	ctx := context.TODO()
	fields := log.Fields{}
	os.Setenv("APP_VAR", "anotherapp")
	fields = eventLog.addFieldIfMissingOrDefault("APP", "APP_VAR", constants.AppCtx, ctx, fields, "default")

	assert.Assert(t, fields["APP"] == "anotherapp")
}

func Test_EventLog_Missing_InCtx(t *testing.T) {

	eventLog := NewEventLog(func(config *Config) {})

	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.AppCtx, "appctx")
	fields := log.Fields{}
	fields = eventLog.addFieldIfMissingOrDefault("APP", "APP_VAR", constants.AppCtx, ctx, fields, "default")

	assert.Assert(t, fields["APP"] == "appctx")
}

func Test_EventLog_Status_Ok(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	eventLog := NewEventLog(func(config *Config) { config.Output = memBuffer })

	eventLog.Audit("test_event", true, context.TODO(), log.Fields{})

	audit := memBuffer.String()
	assert.Assert(t, strings.Contains(audit, "test_event"))
	assert.Assert(t, strings.Contains(audit, "status\":true"))
}

func Test_EventLog_Status_Fail(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	eventLog := NewEventLog(func(config *Config) { config.Output = memBuffer })

	eventLog.Audit("test_event", false, context.TODO(), log.Fields{})

	audit := memBuffer.String()
	assert.Assert(t, strings.Contains(audit, "test_event"))
	assert.Assert(t, strings.Contains(audit, "status\":false"))
}

func Test_EventLog_Examples(t *testing.T) {

	duration := time.Millisecond * 234
	Audit("report_shared", true, context.TODO(), log.Fields{
		constants.AppLog: "engagement",
		constants.ProductLog: "service",
		constants.TraceIdLog: helper.NewTraceID(),
		constants.AccountLog: "hooli",
		constants.ModuleLog: "report_shared",
		constants.UserLog: "abc-123",
		"report_shared" : log.Fields{
			constants.TimeTakenLog: helper.DurationAsISO8601(duration),
			constants.UserLog: "xyz-456",
			"survey":  "MLPIOASHF98D8",
		},
	})
}

func Test_EventLog_Examples2(t *testing.T) {
	eventLog := NewEventLog(func(config *Config) {
		config.Product =  "engagement"
		config.Application = "service"
	})

	duration := time.Millisecond * 234
	eventLog.Audit("report_shared", true, context.TODO(), log.Fields{
		constants.TraceIdLog: helper.NewTraceID(),
		constants.AccountLog: "hooli",
		constants.ModuleLog: "report_shared",
		constants.UserLog: "abc-123",
		"report_shared" : log.Fields{
			constants.TimeTakenLog: helper.DurationAsISO8601(duration),
			constants.UserLog: "xyz-456",
			"survey":  "MLPIOASHF98D8",
		},
	})
}

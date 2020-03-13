package log

import (
	"context"
	"encoding/hex"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/cultureamp/glamplify/constants"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"
)

// DefaultValues
type DefaultValues struct {
	timeFormat string
}

func NewDefaultValues(timeFormat string) *DefaultValues {
	return &DefaultValues{timeFormat: timeFormat}
}

func (df DefaultValues) GetDefaults(message string, sev string) Fields {
	fields := Fields{
		constants.ArchitectureLogField: df.targetArch(),
		constants.HostLogField:         df.hostName(),
		constants.OsLogField:           df.targetOS(),
		constants.PidLogField:          df.processID(),
		constants.ProcessLogField:      df.processName(),
		constants.SeverityLogField:     sev,
		constants.TimeLogField:         df.timeNow(df.timeFormat),
	}

	// if message is empty (from log.Audit) then don't add it
	if message != constants.EmptyString {
		fields[constants.MessageLogField] = message
	}

	fields = df.getEnvDefaults(fields)

	return fields
}

func (df DefaultValues) GetErrorDefaults(err error, fields Fields) Fields {
	errorMessage := strings.TrimSpace(err.Error())

	stats := &debug.GCStats{}
	buf := debug.Stack()
	debug.ReadGCStats(stats)

	exception := Fields{
		"error":    errorMessage,
		"trace":    string(buf),
		"gc_stats": stats,
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		exception["build_info"] = info
	}

	fields[constants.ExceptionLogField] = exception
	return fields
}

func (df DefaultValues) GetAuditDefaults(ctx context.Context, event string, fields Fields) Fields {
	fields = df.getCtxDefaults(ctx, fields)

	// add sensible defaults for mandatory missing fields!
	fields = df.addMandatoryAuditFieldsIfMissing(ctx, fields)
	fields[constants.EventLogField] = event

	return fields
}

func (df DefaultValues) getEnvDefaults(fields Fields) Fields {

	fields = df.addEnvFieldIfMissing(constants.ProductLogField, constants.ProductEnv, fields)
	fields = df.addEnvFieldIfMissing(constants.AppLogField, constants.AppEnv, fields)
	fields = df.addEnvFieldIfMissing(constants.RegionLogField, constants.RegionEnv, fields)

	return fields
}

func (df DefaultValues) getCtxDefaults(ctx context.Context, fields Fields) Fields {

	fields = df.addCtxFieldIfMissing(ctx, constants.ProductLogField, constants.ProductCtx, fields)
	fields = df.addCtxFieldIfMissing(ctx, constants.AppLogField, constants.AppCtx, fields)
	fields = df.addCtxFieldIfMissing(ctx, constants.TraceIdLogField, constants.TraceIdCtx, fields)
	fields = df.addCtxFieldIfMissing(ctx, constants.ModuleLogField, constants.ModuleCtx, fields)
	fields = df.addCtxFieldIfMissing(ctx, constants.AccountLogField, constants.AccountCtx, fields)
	fields = df.addCtxFieldIfMissing(ctx, constants.UserLogField, constants.UserCtx, fields)
	fields = df.addCtxFieldIfMissing(ctx, constants.RegionLogField, constants.RegionCtx, fields)

	return fields
}

func (df DefaultValues) addMandatoryAuditFieldsIfMissing(ctx context.Context, fields Fields) Fields {
	// Trace_Id
	fields = df.addTraceIdIfMissing(ctx, fields)

	return fields
}

func (df DefaultValues) addTraceIdIfMissing(ctx context.Context, fields Fields) Fields {

	// If it contains it already, all good!
	if _, ok := fields[constants.TraceIdLogField]; ok {
		return fields
	}

	if xray.RequestWasTraced(ctx) {
		fields[constants.TraceIdLogField] = xray.TraceID(ctx)
	} else {
		fields[constants.TraceIdLogField] = df.NewTraceID()
	}

	return fields
}

func (df DefaultValues) addEnvFieldIfMissing(fieldName string, osVar string, fields Fields) Fields {

	// If it contains it already, all good!
	if _, ok := fields[fieldName]; ok {
		return fields
	}

	// next, check env
	if prod, ok := os.LookupEnv(osVar); ok {
		fields[fieldName] = prod
		return fields
	}

	return fields
}

func (df DefaultValues) addCtxFieldIfMissing(ctx context.Context, fieldName string, ctxKey constants.EventCtxKey, fields Fields) Fields {

	// If it contains it already, all good!
	if _, ok := fields[fieldName]; ok {
		return fields
	}

	if prod, ok := ctx.Value(ctxKey).(string); ok {
		fields[fieldName] = prod
		return fields
	}

	return fields
}



func (df DefaultValues) timeNow(format string) string {
	return time.Now().UTC().Format(format)
}

var host string
var hostOnce sync.Once

func (df DefaultValues) hostName() string {

	var err error
	hostOnce.Do(func() {
		host, err = os.Hostname()
		if err != nil {
			host = constants.UnknownString
		}
	})

	return host
}

func (df DefaultValues) processName() string {
	name := os.Args[0]
	if len(name) > 0 {
		name = filepath.Base(name)
	}

	return name
}

var pid int
var pidOnce sync.Once

func (df DefaultValues) processID() int {
	pidOnce.Do(func() {
		pid = os.Getpid()
	})

	return pid
}

func (df DefaultValues) targetArch() string {
	return runtime.GOARCH
}

func (df DefaultValues) targetOS() string {
	return runtime.GOOS
}

var randG = rand.New(rand.NewSource(time.Now().UnixNano()))

func (df DefaultValues) NewTraceID() string {
	epoch := time.Now().Unix()
	hex := df.randHexString(24)

	var sb strings.Builder

	sb.Grow(+40)

	sb.WriteString("1-")
	sb.WriteString(strconv.FormatInt(epoch, 10))
	sb.WriteString("-")
	sb.WriteString(hex)

	return sb.String()
}

func (df DefaultValues) randHexString(n int) string {
	b := make([]byte, (n+1)/2) // can be simplified to n/2 if n is always even

	if _, err := randG.Read(b); err != nil {
		panic(err)
	}

	return hex.EncodeToString(b)[:n]
}



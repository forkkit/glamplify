package log

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/aws/aws-xray-sdk-go/xray"
	"math/rand"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"
)

// DefaultValues
type DefaultValues struct {

}

func DurationAsISO8601(duration time.Duration) string {
	return fmt.Sprintf("P%gS", duration.Seconds())
}

func newDefaultValues() *DefaultValues {
	return &DefaultValues{}
}

func (df DefaultValues) getDefaults(ctx context.Context, event string, sev string) Fields {
	fields := Fields{
		Time:     df.timeNow(RFC3339Milli),
		Event:    event,
		Resource: df.hostName(),
		Os:       df.targetOS(),
		Severity: sev,
	}

	fields = df.getCtxDefaults(ctx, fields)
	fields = df.getEnvDefaults(fields)
	fields = df.addMandatoryFieldsIfMissing(ctx, fields)

	return fields
}

func (df DefaultValues) getErrorDefaults(err error, fields Fields) Fields {
	errorMessage := strings.TrimSpace(err.Error())

	stats := &debug.GCStats{}
	buf := debug.Stack()
	debug.ReadGCStats(stats)

	fields[Exception] = Fields{
		"error":    errorMessage,
		"trace":    string(buf),
		"gc_stats": stats,
	}

	return fields
}

func (df DefaultValues) getEnvDefaults(fields Fields) Fields {

	fields = df.addEnvFieldIfMissing(Product, ProductEnv, fields)
	fields = df.addEnvFieldIfMissing(App, AppEnv, fields)
	fields = df.addEnvFieldIfMissing(AppVer, AppVerEnv, fields)
	fields = df.addEnvFieldIfMissing(Region, RegionEnv, fields)

	return fields
}

func (df DefaultValues) getCtxDefaults(ctx context.Context, fields Fields) Fields {

	fields = df.addCtxFieldIfMissing(ctx, TraceId, TraceIdCtx, fields)
	fields = df.addCtxFieldIfMissing(ctx, Customer, CustomerCtx, fields)
	fields = df.addCtxFieldIfMissing(ctx, User, UserCtx, fields)

	return fields
}

func (df DefaultValues) addMandatoryFieldsIfMissing(ctx context.Context, fields Fields) Fields {
	// Trace_Id
	fields = df.addTraceIdIfMissing(ctx, fields)

	return fields
}

func (df DefaultValues) addTraceIdIfMissing(ctx context.Context, fields Fields) Fields {

	// If it contains it already, all good!
	if _, ok := fields[TraceId]; ok {
		return fields
	}

	if xray.RequestWasTraced(ctx) {
		fields[TraceId] = xray.TraceID(ctx)
	} else {
		fields[TraceId] = df.newTraceID()
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

func (df DefaultValues) addCtxFieldIfMissing(ctx context.Context, fieldName string, ctxKey EventCtxKey, fields Fields) Fields {

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
			host = Unknown
		}
	})

	return host
}

func (df DefaultValues) targetOS() string {
	return runtime.GOOS
}

var randG = rand.New(rand.NewSource(time.Now().UnixNano()))

func (df DefaultValues) newTraceID() string {
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


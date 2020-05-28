package context

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/aws/aws-xray-sdk-go/xray"
	"io"
)


const (
	errorUUID = "00000000-0000-0000-0000-000000000000"
)

func AddRequestFields(ctx context.Context, rsFields RequestScopedFields) context.Context {
	return context.WithValue(ctx, RequestFieldsCtx, rsFields)
}

func GetRequestScopedFields(ctx context.Context) (RequestScopedFields, bool) {
	rsFields, ok := ctx.Value(RequestFieldsCtx).(RequestScopedFields)
	return rsFields, ok
}

func WrapCtx(ctx context.Context) context.Context {

	rsFields, ok := GetRequestScopedFields(ctx)
	if ok {
		return ctx
	}

	traceID := ""
	if xray.RequestWasTraced(ctx) {
		traceID = xray.TraceID(ctx)
	}

	requestID := ""
	correlationID, err := newUUID()
	if err != nil {
		correlationID = errorUUID
	}

	rsFields = RequestScopedFields{
		TraceID: traceID,
		RequestID: requestID,
		CorrelationID: correlationID,
	}

	return AddRequestFields(ctx, rsFields)
}

// newUUID generates a random UUID according to RFC 4122
func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
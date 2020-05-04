package aws

import (
	"context"
	"github.com/aws/aws-xray-sdk-go/xray"
)

func GetTraceId(ctx context.Context) (string, bool) {

	if xray.RequestWasTraced(ctx) {
		return xray.TraceID(ctx), true
	}

	return xray.NewTraceID(), false
}

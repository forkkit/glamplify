package event

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cultureamp/glamplify/log"
)

type lambdaHandler struct {
	impl lambda.Handler

	app Application

	functionName    string
	functionVersion string
	logGroupName    string
	logStreamName   string
	memoryLimitInMB int
}

// Invoke - Invoke API operation for AWS Lambda.
func (handler *lambdaHandler) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	// payload => The JSON that you want to provide to your Lambda function as input.

	// Use handlertrace instead of wrapping the Invoke?
	// - https://godoc.org/github.com/aws/aws-lambda-go/lambda/handlertrace
	// - https://github.com/aws/aws-lambda-go/blob/master/lambda/handler_test.go

	var tin interface{}
	if handler.app.conf.Logging {
		json.Unmarshal(payload, tin)

		handler.app.log("Begin Invoke", log.Fields{
			"function":        handler.functionName,
			"version":         handler.functionVersion,
			"logGroup":        handler.logGroupName,
			"logStream":       handler.logStreamName,
			"memoryLimitInMB": handler.memoryLimitInMB,
			"payload":         tin,
		})
	}

	result, err := handler.impl.Invoke(ctx, payload)

	handler.app.log("End Invoke", log.Fields{
		"function":        handler.functionName,
		"version":         handler.functionVersion,
		"logGroup":        handler.logGroupName,
		"logStream":       handler.logStreamName,
		"memoryLimitInMB": handler.memoryLimitInMB,
		"result":          result,
		"error":           err,
	})

	return result, err
}

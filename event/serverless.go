package event

import (
	"github.com/aws/aws-lambda-go/lambda"
	newrelic "github.com/newrelic/go-agent"
	"github.com/newrelic/go-agent/_integrations/nrlambda"
)

// Start should be used in place of lambda.Start.  Replace:
//
//	lambda.Start(myhandler)
//
// With:
//
//	glServerless.Start(myhandler, app)
//
func Start(handler interface{}, app newrelic.Application) {
	nrlambda.Start(handler, app)
}

// StartHandler should be used in place of lambda.StartHandler.  Replace:
//
//	lambda.StartHandler(myhandler)
//
// With:
//
//	glServerless.StartHandler(myhandler, app)
//
func StartHandler(handler lambda.Handler, app newrelic.Application) {
	nrlambda.StartHandler(handler, app)
}

package aws

import (
	"context"
	"net/http"

	"github.com/aws/aws-xray-sdk-go/awsplugins/ec2"
	"github.com/aws/aws-xray-sdk-go/awsplugins/ecs"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/aws/aws-xray-sdk-go/xraylog"
)

// TracerConfig for setting initial values for Tracer
type TracerConfig struct {
	Environment   string
	AWSService    string
	EnableLogging bool
	Version       string
}

type Tracer struct {
	config TracerConfig
	logger *xrayLogger
}

func NewTracer(ctx context.Context, configure ...func(*TracerConfig)) *Tracer {

	conf := TracerConfig{
		Environment: "development",
	}
	for _, config := range configure {
		config(&conf)
	}

	if conf.Environment == "production" {
		if conf.AWSService == "ECS" {
			ecs.Init()
		} else if conf.AWSService == "EC2" {
			ec2.Init()
		}
	}

	logger := newXrayLogger(ctx)
	if conf.EnableLogging {
		xray.SetLogger(logger)
	}

	if err := xray.Configure(xray.Config{ServiceVersion: conf.Version}); err != nil {
		logger.Log(xraylog.LogLevelError, newPrintArgs(err.Error()))
	}

	return &Tracer{
		config: conf,
		logger: logger,
	}
}

func (tracer Tracer) GetTraceID(ctx context.Context) string {
	if xray.RequestWasTraced(ctx) {
		return xray.TraceID(ctx)
	}

	return ""
}

func (tracer Tracer) RoundTripper(rt http.RoundTripper) http.RoundTripper {
	return xray.RoundTripper(rt)
}

func (tracer Tracer) SegmentHandler(name string, h http.Handler) http.Handler {

	sn := xray.NewFixedSegmentNamer(name)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xray.Handler(sn, h)
	})
}

func (tracer Tracer) DynamicSegmentHandler(fallback string, wildcardHost string, h http.Handler) http.Handler {

	sn := xray.NewDynamicSegmentNamer(fallback, wildcardHost)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xray.Handler(sn, h)
	})
}

// Capture wrapper around xray.Capture as per https://docs.aws.amazon.com/xray/latest/devguide/xray-sdk-go-subsegments.html
func (tracer Tracer) Capture(ctx context.Context, name string, fn func(context.Context) error) (err error) {
	return xray.Capture(ctx, name, fn)
}

// AddMetadata wrapper around xray.AddMetadata as per https://docs.aws.amazon.com/xray/latest/devguide/xray-sdk-go-subsegments.html
func (tracer Tracer) AddMetadata(ctx context.Context,  key string, value interface{}) error {
	return xray.AddMetadata(ctx, key, value)
}
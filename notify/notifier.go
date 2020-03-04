package notify

import (
	"context"
	"github.com/bugsnag/bugsnag-go"
	"github.com/cultureamp/glamplify/types"
	"github.com/cultureamp/glamplify/helper"
	"net/http"
	"os"
	"time"
)

type Config struct {
	Enabled bool
	Logging bool
	License string
	AppName string
	AppVersion string
	ReleaseStage string
	ProjectPackages []string
}

type Notifier struct {
	conf Config
}

const (
	waitFORBugsnag = 2 * time.Second
)

var (
	internal, _ = NewNotifier(helper.GetEnvOrDefault("APP_NAME", "default"), func(conf *Config) {conf.Enabled = true})
)

func NewNotifier(name string, configure ...func(*Config)) (*Notifier, error) {

	if len(name) == 0 {
		name = helper.GetEnvOrDefault("APP_NAME", "default")
	}

	conf := Config{
		Enabled:        	false,
		Logging:			false,
		License:        	os.Getenv("BUGSNAG_LICENSE_KEY"),
		AppName:			name,
		AppVersion: 		helper.GetEnvOrDefault("APP_VERSION", "1.0.0"),
		ReleaseStage:   	helper.GetEnvOrDefault("APP_ENV", "production"),
		ProjectPackages: 	[]string{"github.com/cultureamp"},
	}

	for _, config := range configure {
		config(&conf)
	}

	cfg := bugsnag.Configuration{
		APIKey:          conf.License,
		AppType: 		conf.AppName,
		AppVersion: 	conf.AppVersion,
		ReleaseStage:    conf.ReleaseStage,
		ProjectPackages: conf.ProjectPackages,
		ParamsFilters:[]string{"password", "pwd"}, // todo - add others
	}

	if conf.Logging {
		cfg.Logger = newNotifyLogger()
	}

	bugsnag.Configure(cfg)

	return &Notifier{conf:conf}, nil
}

// Shutdown flushes any remaining data to the SAAS endpoint
func (notify Notifier) Shutdown() {
	time.Sleep(waitFORBugsnag)
}

func (notify *Notifier) WrapHTTPHandler(pattern string, handler func(http.ResponseWriter, *http.Request)) (string, func(http.ResponseWriter, *http.Request)) {
	p, h := notify.wrapHTTPHandler(pattern, http.HandlerFunc(handler))
	return p, func(w http.ResponseWriter, r *http.Request) {
		r = notify.addToHTTPContext(r)
		h.ServeHTTP(w, r)
	}
}

func (notify *Notifier) wrapHTTPHandler(pattern string, handler http.Handler) (string, http.Handler) {
	return pattern, bugsnag.Handler(handler)
}

func Error(err error, fields types.Fields) error {
	return internal.Error(err, fields)
}

func (notify Notifier) Error(err error, fields types.Fields) error {
	if !notify.conf.Enabled { return nil}

	ctx := bugsnag.StartSession(context.Background())
	defer bugsnag.AutoNotify(ctx)

	return notify.ErrorWithContext(err, ctx, fields)
}

func ErrorWithContext(err error, ctx context.Context, fields types.Fields) error {
	return internal.ErrorWithContext(err, ctx, fields)
}

func (notify Notifier) ErrorWithContext(err error, ctx context.Context, fields types.Fields) error {
	if !notify.conf.Enabled { return nil}

	meta := fieldsAsMetaData(fields)
	return bugsnag.Notify(err, ctx, meta)
}

func (notify *Notifier) addToHTTPContext(req *http.Request) *http.Request {
	ctx := notify.addToContext(req.Context())
	return req.WithContext(ctx)
}

func (notify *Notifier) addToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, notifyContextKey, notify)
}

func fieldsAsMetaData(fields types.Fields) bugsnag.MetaData {
	meta := make(bugsnag.MetaData)
	for k, v := range fields {
		meta.Add("app context", k, v)
	}
	return meta
}
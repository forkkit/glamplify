package notify

import (
	"context"
	"github.com/bugsnag/bugsnag-go"
	"github.com/cultureamp/glamplify/field"
	"os"
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

var (
	internal = New(func(conf *Config) {conf.Enabled = true})
)

func New(configure ...func(*Config)) *Notifier {

	conf := Config{
		Enabled:        	false,
		Logging:			false,
		License:        	os.Getenv("BUGSNAG_LICENSE_KEY"),
		AppName: 			os.Getenv("APP_NAME"),
		AppVersion: 		os.Getenv("APP_VERSION"),
		ReleaseStage:   	os.Getenv("APP_ENV"),
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

	return &Notifier{conf:conf}
}

func Error(err error, fields field.Fields) error {
	return internal.Error(err, fields)
}

func (notify Notifier) Error(err error, fields field.Fields) error {
	if !notify.conf.Enabled { return nil}

	ctx := bugsnag.StartSession(context.Background())
	defer bugsnag.AutoNotify(ctx)

	return notify.ErrorWithContext(err, ctx, fields)
}

func ErrorWithContext(err error, ctx context.Context, fields field.Fields) error {
	return internal.ErrorWithContext(err, ctx, fields)
}

func (notify Notifier) ErrorWithContext(err error, ctx context.Context, fields field.Fields) error {
	if !notify.conf.Enabled { return nil}

	meta := fieldsAsMetaData(fields)
	return bugsnag.Notify(err, ctx, meta)
}

func fieldsAsMetaData(fields field.Fields) bugsnag.MetaData {
	meta := make(bugsnag.MetaData)
	for k, v := range fields {
		meta.Add("app context", k, v)
	}
	return meta
}
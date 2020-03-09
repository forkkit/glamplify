# glamplify
Go Amplify Module of useful common tools. The guiding principle is to implement a very light weight wrapper over the standard library (or where not adequate an open source community library), that conforms to our standard practises (12-Factor) and sensible defaults.


## Install

```
go get github.com/cultureamp/glamplify
```

## Usage

### Config
```Go
package main

import (
	"github.com/cultureamp/glamplify/config"
)

func main() {

    // settings will contain configuration data as read in from the config file.
    settings := config.Load()

    // Or if you want to look for a config file from a specific location use
    // settings = config.LoadFrom([]string{"${HOME}/settings"}, "config")

    // Then you can use
    if settings.App.Version > 2.0 {
        // to do
    }
}
```
If no config.yml or config.json can be found, or if it is corrupted, then a config will be created by checking these ENV variables.

- CONFIG_APPNAME (default: "service-name")
- CONFIG_VERSION (default: 1.0)

### Logging

```Go
package main

import (
    "bytes"
    "errors"

    "github.com/cultureamp/glamplify/log"
)

func main() {

    // You can either get a new logger, or just use the public functions which internally use an internal logger
    // eg. log.Debug(), log.Info(), log.Warn(), log.Error() and log.Fatal()

    // Example below shows usage with the package level logger (sensible default), but can 
    // use an instance of a logger by calling mylogger := log.New()

    // Emit debug trace
    // All messages must be static strings (as per Culture Amp Sensibile Default)
    log.Debug("Something happened")

    // Fields can contain any type of variables
    log.Debug("Something happened", log.Fields{
        "aString": "hello",
        "aInt":    123,
        "aFloat":  42.48,
     })

    // Typically Info will be sent onto 3rd party aggregation tools (eg. Splunk)
    log.Info("Executing main")

    // Fields can contain any type of variables
    log.Info("Executing main", log.Fields{
        "program-name": "helloworld.exe",
        "start-up-param":    123,
        "user":  "admin",
    })

    // Errors will always be sent onto 3rd party aggregation tools (eg. Splunk)
    err := errors.New("missing database connection string")
    log.Error(err)

    // Fields can contain any type of variables
    err = errors.New("missing database connection string")
    log.Error(err, log.Fields{
        "program-name": "helloworld.exe",
        "start-up-param":    123,
        "user":  "admin",
     })

    // have a requestID for every log message within that scope) then you can use WithScope()
    scope := log.WithScope(log.Fields { "requestID" : 123 })

    // then just use the scope as you would a normal logger
    scope.Info("Starting web request", log.Fields { "auth": "oauth" })

    // If you want to change the output or time format you can only do this for an
    // instance of the logger you create (not the internal one) by doing this:

    memBuffer := &bytes.Buffer{}
    logger := log.New(func(conf *log.Config) {
        conf.Output = memBuffer                     // can be set to anything that support io.Write
        conf.TimeFormat = "2006-01-02T15:04:05"     // any valid time format
    })

    // The internal logger will always use these default values:
    // output = os.Stdout
    // time format = "2006-01-02T15:04:05.000Z07:00"
    // debugForwardLogTo = "none"
    // printForwardLogTo = "splunk"
    // errorForwardLogTo = "splunk"
}

```
Use `log.Debug` for logging that will only be used when diving deep to uncover bugs. Typically `log.Debug` messages will not automatically be sent to other 3rd party systems (eg. Splunk).

Use `log.Print` for standard log messages that you want to see always. These will never be turned off and will likely be always sent to 3rd party systems for further analysis (eg. Spliunk).

Use `log.Error` when you have encounter a GO error. This will NOT stop the program, it is up to you to call exit() or panic() if this is not recoverable. All error messages will be forwarded to 3rd party systems for monitoring and further analysis.

### Monitor

Make sure you have the environment variable NEW_RELIC_LICENSE_KEY set to the correct 40 character license key. 
Alternatively, you can read it from another environment variable and pass it into the monitor.Config struct.

#### Adding Attributes to a Web Request Transaction
```Go
package main

import (

"github.com/cultureamp/glamplify/log"
"github.com/cultureamp/glamplify/monitor"
"net/http"
)

func main() {

    app, err := monitor.NewApplication("GlamplifyDemo", func(conf *monitor.Config) {
        conf.Enabled = true             // default = "false"
        conf.Logging = true             // default = "false"
        conf.ServerlessMode = false     // default = "false"
     })

    _, handler := app.WrapHTTPHandler("/", rootRequestHandler)
    h := http.HandlerFunc(handler)

    if err = http.ListenAndServe(":8080", h); err != nil {
        panic(err)
    }

    app.Shutdown()
}

func rootRequestHandler(w http.ResponseWriter, r *http.Request) {

    // Do things

    txn, err := monitor.TxnFromRequest(w, r)
    if err != nil {
        txn.AddAttributes(log.Fields{
            "aString": "hello world",
            "aInt":    123,
        })
    }

    // Do more things

    if err != nil {
        txn.AddAttributes(log.Fields{
            "aString2": "goodbye",
            "aInt2":    456,
        })
    }

}
```

#### Custom Events to a Web Request Transaction
```Go
package main

import (

"github.com/cultureamp/glamplify/log"
"github.com/cultureamp/glamplify/monitor"
"net/http"
)

func main() {

    app, err := monitor.NewApplication("GlamplifyDemo", func(conf *monitor.Config) {
		conf.Enabled = true             // default = "false"
		conf.Logging = true             // default = "false"
		conf.ServerlessMode = false     // default = "false"
	})

    _, handler := app.WrapHTTPHandler("/", rootRequestHandler)
    h := http.HandlerFunc(handler)

    if err = http.ListenAndServe(":8080", h); err != nil {
        panic(err)
    }
    
    app.Shutdown()
}

func rootRequestHandler(w http.ResponseWriter, r *http.Request) {

    // Do things
    app, err := monitor.AppFromRequest(w, r)

    err = app.RecordEvent("mycustomEvent", log.Fields{
        "aString": "hello world",
        "aInt":    123,
    })

    // Do more things
}
```

#### Adding Attributes to a Lambda (Serverless)
```Go
package main

import (

"context"
"github.com/cultureamp/glamplify/log"
"github.com/cultureamp/glamplify/monitor"
"net/http"
)

func main() {
    app, err := monitor.NewApplication("GlamplifyDemo", func(conf *monitor.Config) {
        conf.Enabled = true             // default = "false"
        conf.Logging = true             // default = "false"
        conf.ServerlessMode = true      // default = "false"
     })

    monitor.Start(handler, app)
}

func handler(ctx context.Context) {

    // Do things

    txn, err := monitor.TxnFromContext(ctx)
    if err != nil {
        txn.AddAttributes(log.Fields{
            "aString": "hello world",
            "aInt":    123,
        })
    }

    // Do more things

    if err != nil {
        txn.AddAttributes(log.Fields{
            "aString2": "goodbye",
            "aInt2":    456,
        })
    }
}
```

#### Custom Events to a Lambda (Serverless)
```Go
package main

import (

"context"
"github.com/cultureamp/glamplify/log"
"github.com/cultureamp/glamplify/monitor"
"net/http"
)

func main() {
    app, err := monitor.NewApplication("GlamplifyDemo", func(conf *monitor.Config) {
        conf.Enabled = true             // default = "false"
        conf.Logging = true             // default = "false"
        conf.ServerlessMode = true      // default = "false"
    })

    monitor.Start(handler, app)
}

func handler(ctx context.Context) {

    // Do things

    app, err := monitor.AppFromContext(ctx)
    if err != nil {
        err = app.RecordEvent("mycustomEvent", log.Fields{
            "aString": "hello world",
            "aInt":    123,
        })
    }
    // Do more things
}
```

## Notify

Make sure you have the environment variable BUGSNAG_LICENSE_KEY set to the correct license key. 
Alternatively, you can read it from another environment variable and pass it into the notify.Config struct.

```Go
package main

import (

"errors"
"github.com/cultureamp/glamplify/log"
"github.com/cultureamp/glamplify/notify"
"net/http"
)

func main() {

    notifier, err := notify.NewNotifier("GlamplifyDemo", func(conf *notify.Config) {
        conf.Enabled = true             // default = "false"
        conf.Logging = true             // default = "false"
     })

    _, handler := notifier.WrapHTTPHandler("/", rootRequestHandler)
    h := http.HandlerFunc(handler)

    if err = http.ListenAndServe(":8080", h); err != nil {
        panic(err)
    }

    notifier.Shutdown()
}

func rootRequestHandler(w http.ResponseWriter, r *http.Request) {

    // Do things

    // pretend we got an error
    err := errors.New("NPE") // 

    notifier, notifyErr := notify.NotifyFromRequest(w, r)
    if notifyErr != nil {
        log.Error(err)
    }

    notifier.ErrorWithContext(err, r.Context(), log.Fields {
        "key": "value",
    })

}
```

## Event Logging for Audit

So we can generate audit logs, there is a new sensible default for logging [here](https://cultureamp.atlassian.net/wiki/spaces/TV/pages/959939199/Logging)

The eventLog in glamplify implements this by wrapping the logger.

```Go
package main

import (
    "context"
    "github.com/cultureamp/glamplify/constants"
    "github.com/cultureamp/glamplify/events"
    "github.com/cultureamp/glamplify/helper"
    "github.com/cultureamp/glamplify/log"
    "time"
)

func main() {

    // use the default internal event log Audit
    duration := time.Millisecond * 234
    events.Audit("report_shared", true, context.TODO(), log.Fields{
        constants.AppLogField:     "engagement",
        constants.ProductLogField: "service",
        constants.TraceIdLogField: helper.NewTraceID(),
        constants.AccountLogField: "hooli",
        constants.ModuleLogField:  "report_shared",
        constants.UserLogField:    "abc-123",
        "report_shared" : log.Fields{
            constants.TimeTakenLogField: helper.DurationAsISO8601(duration),
            constants.UserLogField:      "xyz-456",
            "survey":                    "MLPIOASHF98D8",
        },
    })

    // Or if you want to create your own instance and set up some defaults
    eventLog := events.NewEventLog(func(config *events.Config) {
		config.Product =  "engagement"
		config.Application = "service"
	})

	duration = time.Millisecond * 567
	eventLog.Audit("report_shared", true, context.TODO(), log.Fields{
		constants.TraceIdLogField: helper.NewTraceID(),
		constants.AccountLogField: "hooli",
		constants.ModuleLogField:  "report_shared",
		constants.UserLogField:    "abc-123",
		"report_shared" : log.Fields{
			constants.TimeTakenLogField: helper.DurationAsISO8601(duration),
			constants.UserLogField:      "xyz-456",
			"survey":                    "MLPIOASHF98D8",
		},
	})

    // you can also add some values to the context and then they will be automatically added to the event
    ctx := context.Background()
    ctx = context.WithValue(ctx, constants.TraceIdCtx, "current_trace_id")
    ctx = context.WithValue(ctx, constants.UserCtx, "current_user_id")
    // see constants for other Ctx constant keys
    eventLog.Audit("report_shared", true, ctx, log.Fields{
		constants.AccountLogField: "hooli",
		constants.ModuleLogField:  "report_shared",
		"report_shared" : log.Fields{
			constants.TimeTakenLogField: helper.DurationAsISO8601(duration),
			constants.UserLogField:      "xyz-456",
			"survey":                    "MLPIOASHF98D8",
		},
	})

    // values can also be read form the ENV, the most useful being 
    // constants.ProductEnv and constants.AppEnv
    // see constants for other Env constant keys
}
```


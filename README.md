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

Logging in GO supports the Culture Amp [sensible default](https://cultureamp.atlassian.net/wiki/spaces/TV/pages/959939199/Logging) 

```Go
package main

import (
    "bytes"
    "context"
    "errors"
    "time"
    "github.com/cultureamp/glamplify/constants"
    "github.com/cultureamp/glamplify/log"
)

func main() {

	// If you aren't passed a context, then you need to create a new one and then you should add
    // all the mandatory values to it so logging can retrieve them automatically
    // Example:
    ctx := context.Background()
   
    // AWS X-ray trace_id normally passed via http headers or by another method
    traceId := "1-58406520-a006649127e371903a2de979"
    ctx = log.AddTraceId(ctx, traceId)  
   
    // If this service deals with a particularly customer, then set that on the context as well
    customer := "FNSNDCJDF343"
    ctx = log.AddCustomer(ctx, customer)

    // And finally if this service deals with a particular user, then set that on the context as well
    user := "JFOSNDJF97S"
    ctx = log.AddUser(ctx, user)

    // You can either get a new logger, or just use the public functions which internally use an internal logger
    // eg. log.Debug(), log.Info(), log.Warn(), log.Error() and log.Fatal()

    // Example below shows usage with the package level logger (sensible default), but can 
    // use an instance of a logger by calling mylogger := log.New()

    // Emit debug trace
    // All messages must be static strings (as per Culture Amp Sensibile Default)
    log.Debug(ctx, "something_happened_event")

    // Fields can contain any type of variables
    // there are a number of constants for common field names (see constants.go for full list)
    // TraceIdLogField             = "trace_id"
    // ResourceLogField            = "resource"
    // CustomerLogField            = "customer"
    // UserLogField                = "user"
    // MessageLogField             = "message"
    // TimeTakenLogField           = "time_taken"
    // MemoryUsedLogField          = "memory_used"
    // MemoryAvailLogField         = "memory_available"
    // ItemsProcessedLogField      = "items_processed"
    // TotalItemsProcessedLogField = "total_items_processed"
    // TotalItemsRequestedLogField = "total_items_requested"
 
    d := time.Millisecond * 123
    log.Debug(ctx, "something_happened", log.Fields{
        constants.MessageLogField: "the thing did what we expected it to do",
        constants.TimeTakenLogField : log.DurationAsISO8601(d), // returns "P0.123S" 
        constants.UserLogField: "MMLKSN443FN",
        "report":  "NVJKSJFJ34NBFN44",
        "aInt":    123,
        "aFloat":  42.48,
        "aString": "more info",
     })

    // Typically Info will be sent onto 3rd party aggregation tools (eg. Splunk)
    log.Info(ctx, "something_happened_event")

    // Fields can contain any type of variables
    d = time.Millisecond * 456
    log.Info(ctx, "something_happened_event", log.Fields{
        "program-name": "helloworld.exe",
        "start-up-param":    123,
        constants.UserLogField:  "admin",
         constants.MessageLogField: "the thing did what we expected it to do",
        constants.TimeTakenLogField: log.DurationAsISO8601(d), // returns "P0.456S" 
    })

    // Errors will always be sent onto 3rd party aggregation tools (eg. Splunk)
    err := errors.New("missing database connection string")
    log.Error(ctx, err)

    // Fields can contain any type of variables
    err = errors.New("missing database connection string")
    log.Error(ctx, err, log.Fields{
        "program-name": "helloworld.exe",
        "start-up-param":    123,
        "user":  "admin",
        "message": "the thing did not do what we expected it to do",
     })

    // have a requestID for every log message within that scope) then you can use WithScope()
    scope := log.WithScope(log.Fields { "requestID" : 123 })

    // then just use the scope as you would a normal logger
    scope.Info(ctx, "something_happened_event", log.Fields { "auth": "oauth" })

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

Use `log.Error` when you have encounter a GO error. This will NOT stop the program, it is assumed that the system has recovered. All error messages will be forwarded to 3rd party systems for monitoring and further analysis.

Use `log.Fatal` when you have encounter a GO error that is not recoverable. This will stop the program by calling panic(). All fatal messages will be forwarded to 3rd party systems for monitoring and further analysis.

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
        log.Error(r.Context(), err)
    }

    notifier.ErrorWithContext(err, r.Context(), log.Fields {
        "key": "value",
    })

}
```


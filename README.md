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
    "github.com/cultureamp/glamplify/log"
    "net/http"
    "time"
)

func main() {

    // Creating loggers is cheap. Create them on every request/run
    // DO NOT CACHE/REUSE THEM
    ctx := context.Background()
    
    // This does all the good things - eg. reads the aws xray Trace_ID, or creates a new trace if missing
    ctx, logger := log.New(ctx)
    
    // or if you want a field to be present on each subsequent logging call do this:
    ctx, logger = log.New(ctx, log.Fields{"request_id": 123})
    
    h := http.HandlerFunc(requestHandler)
    
    if err := http.ListenAndServe(":8080", h); err != nil {
        logger.Error(err)
    }
}

func requestHandler(w http.ResponseWriter, r *http.Request) {

    ctx := r.Context()
    
    /* REQUEST LOGGING */
    
    // This helper does all the good things:
    // Decoded JWT (if present) and set User/Customer on the context
    // Reads the aws xray Trace_ID or creates a new one if missing
    // Can optionally pass in log.Fields{} if you have values you want to
    // scope to every subsequent logging calls..   eg. logger, ctx, err := helper.NewFromRequest(ctx, r, log.Fields{"request_id": 123})
    ctx, logger, err := log.NewFromRequest(ctx, r)
    if err != nil {
        // Error here usually means missing public key or corrupted JWT or such like
        // But a valid logger is ALWAYS returned, so it is safe to use. It just won't have User/Customer logging fields
        logger.Error(err)
    }
    
    // Once you have a logger then you can call logger.Debug/Info/Warn/Error/Fatal
    // NOTE: unlike normal logging, the string you pass in SHOULD BE AN EVENT NAME (past tense)
    
    // Example
    logger.Debug("something_happened_event")
    
    // if you want to pass a detailed message then use
    logger.Debug("something_happed_event", log.Fields {
        log.Message: "something that I was expecting actually did happen!",
    })
    
    // Fields can contain any type of variables, but here are some helpful predefined ones
    // (see constants.go for full list)
    // MessageLogField             = "message"
    // TimeTakenLogField           = "time_taken"
    // MemoryUsedLogField          = "memory_used"
    // MemoryAvailLogField         = "memory_available"
    // ItemsProcessedLogField      = "items_processed"
    // TotalItemsProcessedLogField = "total_items_processed"
    // TotalItemsRequestedLogField = "total_items_requested"
    
    d := time.Millisecond * 123
    logger.Debug("something_happened", log.Fields{
        log.Message: "the thing did what we expected it to do",
        log.TimeTaken : log.DurationAsISO8601(d), // returns "P0.123S" as per sensible default 
        log.User: "MMLKSN443FN",
        "report":  "NVJKSJFJ34NBFN44",
        "aInt":    123,
        "aFloat":  42.48,
        "aString": "more info",
     })
    
    // Typically Info will be sent onto 3rd party aggregation tools (eg. Splunk)
    logger.Info("something_happened_event")
    
    // Fields can contain any type of variables
    d = time.Millisecond * 456
    logger.Info("something_happened_event", log.Fields{
        "program-name": "helloworld.exe",
        "start-up-param":    123,
        log.User:  "admin",
        log.Message: "the thing did what we expected it to do",
        log.TimeTaken: log.DurationAsISO8601(d), // returns "P0.456S" 
    })
    
    // Errors will always be sent onto 3rd party aggregation tools (eg. Splunk)
    err = errors.New("missing database connection string")
    logger.Error(err)
    
    // Fields can contain any type of variables
    err = errors.New("missing database connection string")
    logger.Error(err, log.Fields{
        "program-name": "helloworld.exe",
        "start-up-param":    123,
        log.User:  "admin",
        log.Message: "the thing did not do what we expected it to do",
     })
}

```
Use `log.New()` for logging without a http request. Use `log.NewFromRequest()` when you do have a http request. This initializes a bunch of stuff for you (eg. JWT details are automatically logged for you). NewFromRequest always returns a valid Logger and an optional error (which usually describes problems decoding the JWT etc)

Use `logger.Debug` for logging that will only be used when diving deep to uncover bugs. Typically `scope.Debug` messages will not automatically be sent to other 3rd party systems (eg. Splunk).

Use `logger.Print` for standard log messages that you want to see always. These will never be turned off and will likely be always sent to 3rd party systems for further analysis (eg. Spliunk).

Use `logger.Warnr` when you have encounter a warning that should be looked at by a human, but has been recovered. All warning messages will be forwarded to 3rd party systems for monitoring and further analysis.

Use `logger.Error` when you have encountered a GO error. This will NOT stop the program, it is assumed that the system has recovered. All error messages will be forwarded to 3rd party systems for monitoring and further analysis.

Use `logger.Fatal` when you have encountered a GO error that is not recoverable. This will stop the program by calling panic(). All fatal messages will be forwarded to 3rd party systems for monitoring and further analysis.

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
    
    _, handler := app.WrapHTTPHandler("/", requestHandler)
    h := http.HandlerFunc(handler)
    
    if err = http.ListenAndServe(":8080", h); err != nil {
        panic(err)
    }
    
    app.Shutdown()
}

func requestHandler(w http.ResponseWriter, r *http.Request) {

    ctx := r.Context()
    ctx, logger, err := log.NewFromRequest(ctx, r)
    if err != nil {
        logger.Error(err)
    }
    
    // Do things
    
    txn, err := monitor.TxnFromRequest(w, r)
    if err != nil {
        logger.Error(err)
        txn.AddAttributes(log.Fields{
            "aString": "hello world",
            "aInt":    123,
        })
    }
    
    // Do more things
    
    if err != nil {
        logger.Error(err)
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
    
    _, handler := app.WrapHTTPHandler("/", requestHandler)
    h := http.HandlerFunc(handler)
    
    if err = http.ListenAndServe(":8080", h); err != nil {
        panic(err)
    }
    
    app.Shutdown()
}

func requestHandler(w http.ResponseWriter, r *http.Request) {

     ctx := r.Context()
     ctx, logger, err := log.NewFromRequest(ctx, r)
     if err != nil {
         logger.Error(err)
     }
    
    // Do things
    app, err := monitor.AppFromRequest(w, r)
    if err != nil {
        logger.Error(err)
    }   
    
    err = app.RecordEvent("mycustomEvent", log.Fields{
        "aString": "hello world",
        "aInt":    123,
    })
    if err != nil {
        logger.Error(err)
    }  
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

    app, _ := monitor.NewApplication("GlamplifyDemo", func(conf *monitor.Config) {
        conf.Enabled = true             // default = "false"
        conf.Logging = true             // default = "false"
        conf.ServerlessMode = true      // default = "false"
     })
    
    monitor.Start(handler, app)
}

func handler(ctx context.Context) {
   
    ctx, logger := log.New(ctx)
    
    // Do things
    
    txn, err := monitor.TxnFromContext(ctx)
    if err != nil {
        logger.Error(err)
        txn.AddAttributes(log.Fields{
            "aString": "hello world",
            "aInt":    123,
        })
    }
    
    // Do more things
    
    if err != nil {
        logger.Error(err)
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

    app, _ := monitor.NewApplication("GlamplifyDemo", func(conf *monitor.Config) {
        conf.Enabled = true             // default = "false"
        conf.Logging = true             // default = "false"
        conf.ServerlessMode = true      // default = "false"
    })
    
    monitor.Start(handler, app)
}

func handler(ctx context.Context) {

    ctx, logger := log.New(ctx)
    
    // Do things
    
    app, err := monitor.AppFromContext(ctx)
    if err != nil {
        logger.Error(err)
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
    
    _, handler := notifier.WrapHTTPHandler("/", requestHandler)
    h := http.HandlerFunc(handler)
    
    if err = http.ListenAndServe(":8080", h); err != nil {
        panic(err)
    }
    
    notifier.Shutdown()
}

func requestHandler(w http.ResponseWriter, r *http.Request) {

    ctx := r.Context()
    ctx, logger, err := log.NewFromRequest(ctx, r)
    if err != nil {
       logger.Error(err)
    }
    
    notifier, notifyErr := notify.NotifyFromRequest(w, r)
    if notifyErr != nil {
        logger.Error(err)
    }
    
    // Do things
    
    // pretend we got an error
    err = errors.New("NPE")  
    notifier.ErrorWithContext(ctx, err, log.Fields {
        "key": "value",
    })
}
```


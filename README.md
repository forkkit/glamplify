# glamplify
Go Amplify Module of useful common tools. The guiding principle is to implement a very light weight wrapper over the standard library (or where not adequate an open source community library), that conforms to our standard practises (12-Factor) and sensible defaults.


## Install

```
go get github.com/cultureamp/glamplify
```

## Usage

### Config
```
package main

import (
    "github.com/cultureamp/glamplify/config"
)

func main() {

    // settings will contain configuration data as read in from the config file.
    settings := config.Load()

    // Or if you want to look for a config file from a specific location use
    settings = config.LoadFrom([]string{"${HOME}/settings"}, "config")

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

```
package main

import (
    "errors"

    "github.com/cultureamp/glamplify/log"
)

func main() {

    // Get the logger
    logger := log.New()

    // Emit debug trace
    // All messages must be static strings (as per Culture Amp Sensibile Default)
    logger.Debug("Something happened")

    // Emit debug trace with fields
    // Fields can contain any type of variables
    logger.Debug("Something happened", log.Fields{
		"aString": "hello",
		"aInt":    123,
		"aFloat":  42.48,
	})

    // Emit normal logging (can add optional fields if required)
    logger.Print("Executing main")

    // Emit Error (can add optional fields if required)
    err := errors.New("Main program stopped unexpectedly",
    logger.Error(err)
}

```

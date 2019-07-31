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

    // If you want to look for a config file from a specific location use
    settings = config.LoadFrom([]string{"${HOME}/settings"}, "config")
}
```
If no config.yml or config.json can be found, or if it is corrupted, then a config will be created by checking these ENV variables.

- CONFIG_APPNAME (default: "service-name")
- CONFIG_VERSION (default: 1.0)
- CONFIG_LOGNAME (default: "default")
- CONFIG_LOGLEVEL (default: "warn")

### Logging

```
package main

import (
"github.com/cultureamp/glamplify/log"
)

func main() {

    // Get default logger
    logger := LoggerFactory.Get("default")

    // Emit debug trace
    logger.Debug("Something happened")

    // Emit info trace with formatting
    logger.Infof("Executing %s...", "main")

    // Emit Warning with strutured fields
    logger.WarnWithFields(
        logger.Fields{"cpu": "amd"},
        "Wrong CPU type, expect slow execution times"
    )

    // Emit Warning with strutured fields
    program := "example.exe"
    logger.ErrorfWithFields(
        logger.Fields{"cpu": "amd"},
        "Main program %s stopped unexpectedly",
        program
    )
}

```

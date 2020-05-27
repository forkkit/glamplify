package monitor

import (
	"github.com/cultureamp/glamplify/log"
	"os"
	"sync"
)

// WriterConfig for setting initial values for Logger
type WriterConfig struct {
	// License is your New Relic license key.
	//
	// https://docs.newrelic.com/docs/accounts/install-new-relic/account-setup/license-key
	License string `yaml:"license"`

	// URL
	// US: https://log-api.newrelic.com/log/v1  (default)
	// EU: https://log-api.eu.newrelic.com/log/v1
	Endpoint string `yaml:"endpoint"`
}

// FieldWriter wraps the standard library writer and add structured types as quoted key value pairs
type FieldWriter struct {
	mutex      sync.Mutex
	config     WriterConfig
}


// NewWriter creates a new FieldWriter. The optional configure func lets you set values on the underlying standard writer.
// Useful for CLI apps that want to direct logging to a file or stderr
// eg. SetOutput
func NewWriter(configure ...func(*WriterConfig)) *FieldWriter { // https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
	conf := WriterConfig{
		License:        os.Getenv("NEW_RELIC_LICENSE_KEY"),
		Endpoint: getEnvOrDefault("NEW_RELIC_LOG_ENDPOINT", "https://log-api.newrelic.com/log/v1"),
	}

	for _, config := range configure {
		config(&conf)
	}

	writer := &FieldWriter{
		config: conf,
	}

	return writer
}

func (writer *FieldWriter) WriteFields(system log.Fields, fields ...log.Fields) {
	merged := log.Fields{}
	properties := merged.Merge(fields...)
	if len(properties) > 0 {
		system[log.Properties] = properties
	}
	// TODO - write to end point in NR format
}

func getEnvOrDefault(key string, defaultValue string) string {

	val, found := os.LookupEnv(key)
	if !found {
		return defaultValue
	}
	return val
}
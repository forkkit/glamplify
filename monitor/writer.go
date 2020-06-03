package monitor

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/cultureamp/glamplify/log"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// WriterConfig for setting initial values for Monitor Writer
type WriterConfig struct {
	// License is your New Relic license key.
	//
	// https://docs.newrelic.com/docs/accounts/install-new-relic/account-setup/license-key
	License string

	// URL
	// US: https://log-api.newrelic.com/log/v1  (default)
	// EU: https://log-api.eu.newrelic.com/log/v1
	Endpoint string

	// Timeout
	Timeout  time.Duration
}

// FieldWriter sends logging output to NR as per https://docs.newrelic.com/docs/logs/new-relic-logs/log-api/introduction-log-api
type FieldWriter struct {
	mutex      sync.Mutex
	config     WriterConfig
}

// NewWriter creates a new FieldWriter. The optional configure func lets you set values on the underlying standard writer.
// Useful for CLI apps that want to direct logging to a file or stderr
// eg. SetOutput
func NewWriter(configure ...func(*WriterConfig)) *FieldWriter { // https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
	conf := WriterConfig{
		License:    os.Getenv("NEW_RELIC_LICENSE_KEY"),
		Endpoint:   getEnvOrDefaultString("NEW_RELIC_LOG_ENDPOINT", "https://log-api.newrelic.com/log/v1"),
		Timeout:    time.Second * time.Duration(getEnvOrDefaultInt("NEW_RELIC_TIMEOUT", 5)),
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

	json := system.ToSnakeCase().ToJson()

	go post(writer.config, json)
}

func post(config WriterConfig, jsonStr string) error {
	// https://docs.newrelic.com/docs/logs/new-relic-logs/log-api/introduction-log-api
	jsonBytes := []byte(jsonStr)

	req, err := http.NewRequest("POST", config.Endpoint, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-License-Key", config.License)

	var client = &http.Client{
		Timeout: config.Timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	str := string(body)
	return errors.New(fmt.Sprintf("bad server response: %d. body: %v", resp.StatusCode, str))
}

func getEnvOrDefaultString(key string, defaultValue string) string {

	val, found := os.LookupEnv(key)
	if !found {
		return defaultValue
	}
	return val
}

func getEnvOrDefaultInt(key string, defaultValue int) int {

	val, found := os.LookupEnv(key)
	if !found {
		return defaultValue
	}

	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}

	return i
}

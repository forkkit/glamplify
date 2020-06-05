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

// writerConfig for setting initial values for Monitor Writer
type writerConfig struct {
	// license is your New Relic license key.
	//
	// https://docs.newrelic.com/docs/accounts/install-new-relic/account-setup/license-key
	license string

	// URL
	// US: https://log-api.newrelic.com/log/v1  (default)
	// EU: https://log-api.eu.newrelic.com/log/v1
	endpoint string

	// timeout
	timeout time.Duration
}

// FieldWriter sends logging output to NR as per https://docs.newrelic.com/docs/logs/new-relic-logs/log-api/introduction-log-api
type FieldWriter struct {
	mutex  sync.Mutex
	config writerConfig
}

// newWriter creates a new FieldWriter. The optional configure func lets you set values on the underlying standard writer.
// Useful for CLI apps that want to direct logging to a file or stderr
// eg. SetOutput
func newWriter(configure ...func(*writerConfig)) *FieldWriter { // https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
	conf := writerConfig{
		license:  os.Getenv("NEW_RELIC_LICENSE_KEY"),
		endpoint: getEnvOrDefaultString("NEW_RELIC_LOG_ENDPOINT", "https://log-api.newrelic.com/log/v1"),
		timeout:  time.Second * time.Duration(getEnvOrDefaultInt("NEW_RELIC_TIMEOUT", 5)),
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

func post(config writerConfig, jsonStr string) error {
	// https://docs.newrelic.com/docs/logs/new-relic-logs/log-api/introduction-log-api
	jsonBytes := []byte(jsonStr)

	req, err := http.NewRequest("POST", config.endpoint, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-license-Key", config.license)

	var client = &http.Client{
		Timeout: config.timeout,
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

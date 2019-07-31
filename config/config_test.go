package config_test

import (
	"os"
	"testing"

	"github.com/cultureamp/glamplify/config"
	"gotest.tools/assert"
)

func TestConfig(t *testing.T) {

	config := config.Load()
	assert.Assert(t, config.App.Name == "service-name")
	assert.Assert(t, config.App.Version == 1.0)
	assert.Assert(t, config.App.Loggers[0].Name == "default")
	assert.Assert(t, config.App.Loggers[0].Level == "warn")

}

func TestGet(t *testing.T) {

	settings := config.Load()
	assert.Assert(t, settings.App.Name == "service-name")
	assert.Assert(t, settings.App.Version == 1.0)
	assert.Assert(t, settings.App.Loggers[0].Name == "default")
	assert.Assert(t, settings.App.Loggers[0].Level == "warn")

	os.Setenv("CONFIG_APPNAME", "dummy-micro-service")
	os.Setenv("CONFIG_VERSION", "2.5")
	os.Setenv("CONFIG_LOGNAME", "test")
	os.Setenv("CONFIG_LOGLEVEL", "fatal")

	settings = config.LoadFrom("filename_that_does_not_exist")
	assert.Assert(t, settings.App.Name == "dummy-micro-service")
	assert.Assert(t, settings.App.Version == 2.5)
	assert.Assert(t, settings.App.Loggers[0].Name == "test")
	assert.Assert(t, settings.App.Loggers[0].Level == "fatal")
}

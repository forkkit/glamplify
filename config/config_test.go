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

}

func TestGet(t *testing.T) {

	settings := config.Load()
	assert.Assert(t, settings.App.Name == "service-name")
	assert.Assert(t, settings.App.Version == 1.0)

	os.Setenv("CONFIG_APPNAME", "dummy-micro-service")
	os.Setenv("CONFIG_VERSION", "2.5")

	settings = config.LoadFrom([]string{","}, "filename_that_does_not_exist")
	assert.Assert(t, settings.App.Name == "dummy-micro-service")
	assert.Assert(t, settings.App.Version == 2.5)
}

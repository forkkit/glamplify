package config

import (
	"github.com/spf13/viper"
)

// Configuration corrosponds to the yml config file format
type Configuration struct {
	App ApplicationConfiguration `yaml:"app"`
}

// ApplicationConfiguration contains the 'app:' elements of the yml config file
type ApplicationConfiguration struct {
	Name    string  `yaml:"name"`
	Version float64 `yaml:"version"`
}

// Load the default config.yml file
func Load() *Configuration {
	// Todo - better way to work out where the config.yml file is?
	// Do we have a sensible default for this?
	return LoadFrom(
		[]string{".", "../", "../config", "./config"},
		"config",
	)
}

// LoadFrom the config from these specific paths with this specific filename
func LoadFrom(paths []string, configName string) *Configuration {

	viper.SetDefault("appname", "service-name")
	viper.SetDefault("version", 1.0)

	viper.SetEnvPrefix("CONFIG")
	viper.AutomaticEnv()

	viper.SetConfigName(configName)

	for _, path := range paths {
		viper.AddConfigPath(path)
	}

	config := &Configuration{}
	err := viper.ReadInConfig()
	if err != nil {
		config = createDefaultConfig()
		return config
	}

	err = viper.Unmarshal(config)
	if err != nil {
		config = createDefaultConfig()
		return config
	}
	return config
}

func createDefaultConfig() *Configuration {
	config := &Configuration{}

	config.App = ApplicationConfiguration{
		Name:    viper.GetString("appname"),
		Version: viper.GetFloat64("version"),
	}

	return config
}

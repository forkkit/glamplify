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
	Name    string                `yaml:"name"`
	Version float32               `yaml:"version"`
	Loggers []LoggerConfiguration `yaml:"loggers"`
}

// LoggerConfiguration contains the 'streamloggers:' elements in the config file
type LoggerConfiguration struct {
	Name  string `yaml:"name"`
	Level string `yaml:"level"`
}

// Config contains the current configuration as per config.yml, or if missing
// by the default configuration values (in code)
var Settings *Configuration

func init() {
	Settings = loadConfig()
}

func loadConfig() *Configuration {

	viper.SetConfigName("config")

	// Todo - better way to work out where the config.yml file is?
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../config")

	config := Configuration{}
	err := viper.ReadInConfig()
	if err != nil {
		config = createDefaultConfig()
		return &config
	}

	_ = viper.Unmarshal(&config)
	if err != nil {
		config = createDefaultConfig()
		return &config
	}
	return &config
}

func createDefaultConfig() Configuration {
	config := Configuration{}

	config.App = ApplicationConfiguration{
		Name:    "service-name",
		Version: 1.0,
	}
	config.App.Loggers = make([]LoggerConfiguration, 1)
	config.App.Loggers[0].Name = "default"
	config.App.Loggers[0].Level = "warn"

	return config
}

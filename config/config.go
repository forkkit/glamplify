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
	Name          string                      `yaml:"name"`
	Version       float32                     `yaml:"version"`
	StreamLoggers []StreamLoggerConfiguration `yaml:"streamloggers"`
}

// LoggerConfiguration contains the 'loggers:' elements of of the yml config file
type StreamLoggerConfiguration struct {
	Name          string `yaml:"name"`
	Level         string `yaml:"level"`
	Output        string `yaml:"output"`
	Formatter     string `yaml:"formatter"`
	FullTimestamp bool   `yaml:"fullTimestamp"`
}

// Config contains the current configuration as per config.yml, or if missing
// by the default configuration values (in code)
var Config *Configuration

func init() {
	Config = loadConfig()
}

func loadConfig() *Configuration {
	config := createDefaultConfig()

	viper.SetConfigName("config")

	// Todo - better way to work out where the config.yml file is?
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../config")

	err := viper.ReadInConfig()
	if err != nil {
		return &config
	}
	_ = viper.Unmarshal(&config)
	return &config
}

func createDefaultConfig() Configuration {
	config := Configuration{}

	config.App = ApplicationConfiguration{
		Name:    "service-name",
		Version: 1.0,
	}
	config.App.StreamLoggers = make([]StreamLoggerConfiguration, 1)
	config.App.StreamLoggers[0].Name = "default"
	config.App.StreamLoggers[0].Level = "warn"
	config.App.StreamLoggers[0].Output = "stderr"
	config.App.StreamLoggers[0].Formatter = "text"
	config.App.StreamLoggers[0].FullTimestamp = true

	return config
}

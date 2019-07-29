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
	Name    string              `yaml:"name"`
	Version float32             `yaml:"version"`
	Loggers LoggerConfiguration `yaml:"loggers"`
}

// LoggerConfiguration contains the 'streamloggers:' elements in the config file
type LoggerConfiguration struct {
	StreamLoggers []StreamLoggerConfiguration `yaml:"streamloggers"`
	SlackLoggers  []SlackLoggerConfiguration  `yaml:"slackloggers"`
	SplunkLoggers []SplunkLoggerConfiguration `yaml:"splunkloggers"`
}

// StreamLoggerConfiguration contains the 'streamloggers:' elements in the config file
type StreamLoggerConfiguration struct {
	Name          string `yaml:"name"`
	Level         string `yaml:"level"`
	Output        string `yaml:"output"`
	Formatter     string `yaml:"formatter"`
	FullTimestamp bool   `yaml:"fullTimestamp"`
}

// SlackLoggerConfiguration contains the 'slackLoggers:' elements in the config file
type SlackLoggerConfiguration struct {
	Name          string `yaml:"name"`
	Level         string `yaml:"level"`
	Formatter     string `yaml:"formatter"`
	FullTimestamp bool   `yaml:"fullTimestamp"`
	URL           string `yaml:"url"`
	Channel       string `yaml:"channel"`
	Emoji         bool   `yaml:"emoji"`
}

// SplunkLoggerConfiguration contains the 'slackLoggers:' elements in the config file
type SplunkLoggerConfiguration struct {
	Name          string `yaml:"name"`
	Level         string `yaml:"level"`
	Formatter     string `yaml:"formatter"`
	FullTimestamp bool   `yaml:"fullTimestamp"`
	URL           string `yaml:"url"`
	Token         string `yaml:"token"`
	Source        string `yaml:"source"`
	SourceType    string `yaml:"sourceType"`
	Index         string `yaml:"index"`
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
	config.App.Loggers = LoggerConfiguration{}

	config.App.Loggers.StreamLoggers = make([]StreamLoggerConfiguration, 1)
	config.App.Loggers.StreamLoggers[0].Name = "default"
	config.App.Loggers.StreamLoggers[0].Level = "warn"
	config.App.Loggers.StreamLoggers[0].Output = "stderr"
	config.App.Loggers.StreamLoggers[0].Formatter = "text"
	config.App.Loggers.StreamLoggers[0].FullTimestamp = true

	return config
}

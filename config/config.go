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
	Rules   []RuleConfiguration `yaml:"rules"`
	Targets TargetConfiguration `yaml:"targets"`
}

// RuleConfiguration todo
type RuleConfiguration struct {
	Name    string                    `yaml:"name"`
	Level   string                    `yaml:"level"`
	WriteTo []RuleTargetConfiguration `yaml:"writeto"`
}

// RuleTargetConfiguration todo
type RuleTargetConfiguration struct {
	Target string `yaml:"target"`
}

// TargetConfiguration todo
type TargetConfiguration struct {
	Stream []StreamTargetConfiguration `yaml:"stream"`
	Slack  []SlackTargetConfiguration  `yaml:"slack"`
	Splunk []SplunkTargetConfiguration `yaml:"splunk"`
}

// StreamTargetConfiguration contains the 'streaml:' elements in the config file
type StreamTargetConfiguration struct {
	Name          string `yaml:"name"`
	Output        string `yaml:"output"`
	Formatter     string `yaml:"formatter"`
	FullTimestamp bool   `yaml:"fullTimestamp"`
}

// SlackTargetConfiguration contains the 'slack' elements in the config file
type SlackTargetConfiguration struct {
	Name          string `yaml:"name"`
	Formatter     string `yaml:"formatter"`
	FullTimestamp bool   `yaml:"fullTimestamp"`
	URL           string `yaml:"url"`
	Channel       string `yaml:"channel"`
	Emoji         bool   `yaml:"emoji"`
}

// SplunkTargetConfiguration contains the 'slack:' elements in the config file
type SplunkTargetConfiguration struct {
	Name          string `yaml:"name"`
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
	config.App.Loggers = LoggerConfiguration{}

	config.App.Loggers.Rules = make([]RuleConfiguration, 1)
	config.App.Loggers.Rules[0].Name = "default"
	config.App.Loggers.Rules[0].Level = "warn"
	config.App.Loggers.Rules[0].WriteTo = make([]RuleTargetConfiguration, 1)
	config.App.Loggers.Rules[0].WriteTo[0].Target = "stderr-text"

	config.App.Loggers.Targets.Stream = make([]StreamTargetConfiguration, 1)
	config.App.Loggers.Targets.Stream[0].Name = "default"
	config.App.Loggers.Targets.Stream[0].Output = "stderr"
	config.App.Loggers.Targets.Stream[0].Formatter = "text"
	config.App.Loggers.Targets.Stream[0].FullTimestamp = true

	return config
}

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
	Version float64               `yaml:"version"`
	Loggers []LoggerConfiguration `yaml:"loggers"`
}

// LoggerConfiguration contains the 'streamloggers:' elements in the config file
type LoggerConfiguration struct {
	Name  string `yaml:"name"`
	Level string `yaml:"level"`
}

// Load todo
func Load() *Configuration {
	return LoadFrom("config")
}

// LoadFrom todo...
func LoadFrom(configName string) *Configuration {

	viper.SetDefault("appname", "service-name")
	viper.SetDefault("version", 1.0)
	viper.SetDefault("logname", "default")
	viper.SetDefault("loglevel", "warn")

	viper.SetEnvPrefix("CONFIG")
	viper.AutomaticEnv()

	viper.SetConfigName(configName)

	// Todo - better way to work out where the config.yml file is?
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("./config")

	config := &Configuration{}
	err := viper.ReadInConfig()
	if err != nil {
		config = createDefaultConfig()
		return config
	}

	_ = viper.Unmarshal(config)
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
	config.App.Loggers = make([]LoggerConfiguration, 1)
	config.App.Loggers[0].Name = viper.GetString("logname")
	config.App.Loggers[0].Level = viper.GetString("loglevel")

	return config
}

package logger

import (
	conf "github.com/cultureamp/gamplify/config"
)

// LogFactory contains all the registered loggers
type LogFactory struct {
	loggers    map[string]ILogger
	nullLogger ILogger
}

// LoggerFactory to retrieve registered loggers
var LoggerFactory *LogFactory

// Get a registered logger by name
func (factory *LogFactory) Get(loggerName string) ILogger {

	logger, ok := factory.loggers[loggerName]
	if !ok {
		return factory.nullLogger
	}

	return logger
}

func init() {

	LoggerFactory = &LogFactory{}
	LoggerFactory.loggers = make(map[string]ILogger)

	// Create the default, NullLogger
	LoggerFactory.nullLogger = newNullLogger()

	// convert targets to map's
	streamMap, slackMap, splunkMap := convertTargetsToMaps(conf.Config.App.Loggers.Targets.Stream, conf.Config.App.Loggers.Targets.Slack, conf.Config.App.Loggers.Targets.Splunk)

	// Loop through all the Rules in the config and create specific loggers and add them to the LoggerFactory
	for _, rule := range conf.Config.App.Loggers.Rules {

		for _, writeTo := range rule.WriteTo {
			// For each writeTo, find the target that matches (ignore non-matches) and create logger

			var ok bool

			ok = createStreamLogger(streamMap, rule, writeTo)

			if !ok {
				ok = createSlackLogger(slackMap, rule, writeTo)
			}

			if !ok {
				ok = createSplunkLogger(splunkMap, rule, writeTo)
			}
		}
	}
}

func convertTargetsToMaps(streamTargets []conf.StreamTargetConfiguration, slackTargets []conf.SlackTargetConfiguration, splunkTargets []conf.SplunkTargetConfiguration) (map[string]conf.StreamTargetConfiguration, map[string]conf.SlackTargetConfiguration, map[string]conf.SplunkTargetConfiguration) {
	streamMap := make(map[string]conf.StreamTargetConfiguration)
	for _, stream := range conf.Config.App.Loggers.Targets.Stream {
		streamMap[stream.Name] = stream
	}

	slackMap := make(map[string]conf.SlackTargetConfiguration)
	for _, slack := range conf.Config.App.Loggers.Targets.Slack {
		slackMap[slack.Name] = slack
	}

	splunkMap := make(map[string]conf.SplunkTargetConfiguration)
	for _, splunk := range conf.Config.App.Loggers.Targets.Splunk {
		splunkMap[splunk.Name] = splunk
	}

	return streamMap, slackMap, splunkMap
}

func createStreamLogger(streamMap map[string]conf.StreamTargetConfiguration, rule conf.RuleConfiguration, writeTo conf.RuleTargetConfiguration) bool {

	stream, ok := streamMap[writeTo.Target]
	if ok {
		logger := newStreamLogger(
			rule.Name,
			stream.Formatter,
			stream.FullTimestamp,
			stream.Output,
			rule.Level,
		)
		LoggerFactory.loggers[rule.Name] = logger
	}

	return ok
}

func createSlackLogger(slackMap map[string]conf.SlackTargetConfiguration, rule conf.RuleConfiguration, writeTo conf.RuleTargetConfiguration) bool {

	slack, ok := slackMap[writeTo.Target]
	if ok {
		logger := newSlackLogger(
			rule.Name,
			slack.Formatter,
			slack.FullTimestamp,
			slack.URL,
			slack.Channel,
			slack.Emoji,
			rule.Level,
		)
		LoggerFactory.loggers[rule.Name] = logger
	}

	return ok
}

func createSplunkLogger(splunkMap map[string]conf.SplunkTargetConfiguration, rule conf.RuleConfiguration, writeTo conf.RuleTargetConfiguration) bool {

	_, ok := splunkMap[writeTo.Target]
	if ok {
		/*
			logger := newSplunkLogger(
				rule.Name,
				....
			)
			LoggerFactory.loggers[rule.Name] = logger
		*/
	}

	return ok
}

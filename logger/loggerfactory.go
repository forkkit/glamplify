package logger

import (
	. "github.com/cultureamp/gamplify/config"
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
	streamMap, slackMap, splunkMap := convertTargetsToMaps(Config.App.Loggers.Targets.Stream, Config.App.Loggers.Targets.Slack, Config.App.Loggers.Targets.Splunk)

	// Loop through all the Rules in the config and create specific loggers and add them to the LoggerFactory
	for _, rule := range Config.App.Loggers.Rules {

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

func convertTargetsToMaps(streamTargets []StreamTargetConfiguration, slackTargets []SlackTargetConfiguration, splunkTargets []SplunkTargetConfiguration) (map[string]StreamTargetConfiguration, map[string]SlackTargetConfiguration, map[string]SplunkTargetConfiguration) {
	streamMap := make(map[string]StreamTargetConfiguration)
	for _, stream := range Config.App.Loggers.Targets.Stream {
		streamMap[stream.Name] = stream
	}

	slackMap := make(map[string]SlackTargetConfiguration)
	for _, slack := range Config.App.Loggers.Targets.Slack {
		slackMap[slack.Name] = slack
	}

	splunkMap := make(map[string]SplunkTargetConfiguration)
	for _, splunk := range Config.App.Loggers.Targets.Splunk {
		splunkMap[splunk.Name] = splunk
	}

	return streamMap, slackMap, splunkMap
}

func createStreamLogger(streamMap map[string]StreamTargetConfiguration, rule RuleConfiguration, writeTo RuleTargetConfiguration) bool {

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

func createSlackLogger(slackMap map[string]SlackTargetConfiguration, rule RuleConfiguration, writeTo RuleTargetConfiguration) bool {

	_, ok := slackMap[writeTo.Target]
	if ok {
		/*
			logger := newSlackLogger(
				rule.Name,
				....
			)
			LoggerFactory.loggers[rule.Name] = logger
		*/
	}

	return ok
}

func createSplunkLogger(splunkMap map[string]SplunkTargetConfiguration, rule RuleConfiguration, writeTo RuleTargetConfiguration) bool {

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

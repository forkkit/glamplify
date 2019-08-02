package log

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Fields type, used to pass to `WithFields`.
type Fields map[string]interface{}

// FieldLogger todo...
type FieldLogger struct {
	stdLogger *log.Logger
}

// New todo...
// https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
func New(configure ...func(*log.Logger)) *FieldLogger {

	logger := &FieldLogger{}
	logger.stdLogger = log.New(os.Stdout, "", 0)

	for _, config := range configure {
		config(logger.stdLogger)
	}

	return logger
}

// Debug todo...
func (logger FieldLogger) Debug(message string, fields ...Fields) {
	meta := Fields{
		"level": "debug",
		"time":  time.Now().Format(time.RFC3339),
	}

	str := combine(meta, message, fields...)
	logger.stdLogger.Print(str)
}

// Print todo...
func (logger FieldLogger) Print(message string, fields ...Fields) {
	meta := Fields{
		"time": time.Now().Format(time.RFC3339),
	}

	str := combine(meta, message, fields...)
	logger.stdLogger.Print(str)
}

// Error todo...
func (logger FieldLogger) Error(err error, fields ...Fields) {
	meta := Fields{
		"level": "error",
		"time":  time.Now().Format(time.RFC3339),
	}

	str := combine(meta, err.Error(), fields...)
	logger.stdLogger.Print(str)
}

func combine(meta Fields, message string, fields ...Fields) string {

	var str strings.Builder

	_, pre := serialize(meta)
	str.WriteString(pre)

	for _, f := range fields {
		count, post := serialize(f)
		if count > 0 {
			str.WriteString(" ")
			str.WriteString(post)
		}
	}

	str.WriteString(" ")
	str.WriteString(message)

	return str.String()
}

func serialize(fields Fields) (int, string) {
	var pairs []string
	for k, v := range fields {
		vs := fmt.Sprintf("%v", v)
		pairs = append(pairs, k+"="+strconv.Quote(vs))
	}
	sort.Strings(pairs)
	return len(pairs), strings.Join(pairs, " ")
}

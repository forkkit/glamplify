package log

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	// List of standard keys used for logging
	ARCHITECTURE = "arch"
	ERROR        = "error"
	HOST         = "host"
	MESSAGE      = "msg"
	OS           = "os"
	PID          = "pid"
	PROCESS      = "process"
	SEVERITY     = "severity"
	TIME         = "time"
	FORWARD      = "forward-log"

	// Severity Values
	DEBUG_SEV = "DEBUG"
	INFO_SEV  = "INFO"
	ERROR_SEV = "ERROR"
)

// Fields type, used to pass to Debug, Print and Error.
type Fields map[string]interface{}

var first = []string{TIME, SEVERITY, OS, ARCHITECTURE, HOST, PID, PROCESS}
var last = []string{MESSAGE, ERROR}

func (fields Fields) merge(other ...Fields) Fields {
	merged := Fields{}

	for k, v := range fields {
		merged[k] = v
	}

	for _, f := range other {
		for k, v := range f {
			merged[k] = v
		}
	}

	return merged
}

func (fields Fields) serialize() string {
	var pairs []string

	// Do 'first' feids
	pairs = fields.accumulate(pairs, first)

	// everything else in the middle - sorted
	pairs = fields.sortMiddle(pairs)

	// finish with 'last' fields
	pairs = fields.accumulate(pairs, last)

	return strings.Join(pairs, " ")
}

func (fields Fields) sortMiddle(pairs []string) []string {
	var middle []string
	for k, v := range fields {
		if !stringInSlice(k, last) {
			middle = appendTo(middle, k, v)
		}
	}
	if len(middle) > 0 {
		sort.Strings(middle)
		pairs = append(pairs, middle...)
	}

	return pairs
}

func (fields Fields) accumulate(pairs []string, from []string) []string {
	for _, k := range from {
		v, ok := fields[k]
		if ok {
			pairs = appendTo(pairs, k, v)
			delete(fields, k)
		}
	}
	return pairs
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func appendTo(pairs []string, key string, val interface{}) []string {
	vs, ok := val.(string)
	if !ok {
		// only Sptrinf non-strings
		vs = fmt.Sprintf("%v", val)
	}

	return append(pairs, quoteIfRequired(key)+"="+quoteIfRequired(vs))
}

func quoteIfRequired(input string) string {
	if strings.Contains(input, " ") {
		// strconv.Quote is slow(ish) and does a lot of extra work we don't need
		// input = strconv.Quote(input)

		var sb strings.Builder

		sb.Grow(len(input) + 2)
		sb.WriteString("\"")
		sb.WriteString(input)
		sb.WriteString("\"")

		input = sb.String()
	}
	return input
}

func timeNow(format string) string {
	return time.Now().UTC().Format(format)
}

var host string
var hostOnce sync.Once

func hostName() string {

	var err error
	hostOnce.Do(func() {
		host, err = os.Hostname()
		if err != nil {
			host = "<unknown>"
		}
	})

	return host
}

func processName() string {
	name := os.Args[0]
	if len(name) > 0 {
		name = filepath.Base(name)
	}

	return name
}

var pid int
var pidOnce sync.Once

func processID() int {
	pidOnce.Do(func() {
		pid = os.Getpid()
	})

	return pid
}

func targetArch() string {
	return runtime.GOARCH
}

func targetOS() string {
	return runtime.GOOS
}

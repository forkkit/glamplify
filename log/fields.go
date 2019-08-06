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

// Fields type, used to pass to Debug, Print and Error.
type Fields map[string]interface{}

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
	for k, v := range fields {
		vs, ok := v.(string)
		if !ok {
			// only Sptrinf non-strings
			vs = fmt.Sprintf("%v", v)
		}

		pairs = append(pairs, quoteIfRequired(k)+"="+quoteIfRequired(vs))
	}

	sort.Strings(pairs)
	return strings.Join(pairs, " ")
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
	return time.Now().Format(format)
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

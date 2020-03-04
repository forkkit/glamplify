package log

import (
	"encoding/hex"
	json2 "encoding/json"
	"github.com/cultureamp/glamplify/types"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)


type Fields types.Fields

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

	bytes, err := json2.Marshal(fields)
	if err != nil {
		// REVISIT - panic?
	}

	return string(bytes)
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

var randG = rand.New(rand.NewSource(time.Now().UnixNano()))
func traceID() string {
	epoch := time.Now().Unix()
	hex := randHexString(24)

	var sb strings.Builder

	sb.Grow( +40)

	sb.WriteString("1-")
	sb.WriteString(strconv.FormatInt(epoch, 10))
	sb.WriteString("-")
	sb.WriteString(hex)

	return sb.String()
}

func randHexString(n int) string {
	b := make([]byte, (n+1)/2) // can be simplified to n/2 if n is always even

	if _, err := randG.Read(b); err != nil {
		panic(err)
	}

	return hex.EncodeToString(b)[:n]
}
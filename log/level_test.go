package log

import (
	"gotest.tools/assert"
	"os"
	"testing"
)

func Test_Sev_Log(t *testing.T) {

	sev := newSystemLogLevel()

	ok := sev.shouldLog(DebugSev)
	assert.Assert(t, ok, ok)
	ok = sev.shouldLog(InfoSev)
	assert.Assert(t, ok, ok)
	ok = sev.shouldLog(WarnSev)
	assert.Assert(t, ok, ok)
	ok = sev.shouldLog(ErrorSev)
	assert.Assert(t, ok, ok)
	ok = sev.shouldLog(FatalSev)
	assert.Assert(t, ok, ok)
}

func Test_Sev_Log_Env(t *testing.T) {

	os.Setenv("LOG_LEVEL", WarnSev)
	defer os.Unsetenv("LOG_LEVEL")

	sev := newSystemLogLevel()

	ok := sev.shouldLog(DebugSev)
	assert.Assert(t, !ok, ok)
	ok = sev.shouldLog(InfoSev)
	assert.Assert(t, !ok, ok)
	ok = sev.shouldLog(WarnSev)
	assert.Assert(t, ok, ok)
	ok = sev.shouldLog(ErrorSev)
	assert.Assert(t, ok, ok)
	ok = sev.shouldLog(FatalSev)
	assert.Assert(t, ok, ok)
}

func Test_Sev_Log_Unknown(t *testing.T) {

	sev := newSystemLogLevel()

	ok := sev.shouldLog("unknown")
	assert.Assert(t, !ok, ok)
}
package constants

type EventCtxKey int

const (
	UnknownString = "unknown"
	EmptyString   = ""

	// RFC3339Milli is the standard RFC3339 format with added milliseconds
	RFC3339Milli  = "2006-01-02T15:04:05.000Z07:00"

	// JSON LOG KEYS
	// List of standard keys used for logging as per https://cultureamp.atlassian.net/wiki/spaces/TV/pages/959939199/Logging
	TimeLog         = "time"
	ProductLog      = "product"
	AppLog          = "app"
	EventLog        = "event"
	SeverityLog     = "severity"
	TraceIdLog      = "trace_id"
	ModuleLog       = "module"
	StatusLog       = "status"
	ArchitectureLog = "arch"
	HostLog         = "host"
	OsLog           = "os"
	PidLog          = "pid"
	ProcessLog      = "process"

	AccountLog = "account"
	UserLog    = "user"

	TimeTakenLog           = "time_taken"
	MemoryUsedLog          = "memory_used"
	MemoryAvailLog         = "memory_available"
	ItemsProcessedLog      = "items_processed"
	TotalItemsProcessedLog = "total_items_processed"
	TotalItemsRequestedLog  = "total_items_requested"

	MessageLog = "message"

	// Severity Values
	DebugSevLog = "DEBUG"
	InfoSevLog  = "INFO"
	WarnSevLog  = "WARN"
	ErrorSevLog = "ERROR"
	FatalSevLog = "FATAL"

	// ENVIRONMENT VARIABLES
	// List of  Environment Variables keys
	ProductEnv = "PRODUCT"
	AppEnv     = "APPL"
	TraceIdEnv = "TRACE_ID"
	ModuleEnv  = "MODULE"
	AccountEnv = "ACCOUNT"
	UserEnv    = "USER"
	/*
		TIME_TAKEN	= "time_taken"
		MEMORY_USED	= "memory_used"
		MEMORY_AVAIL	= "memory_available"
		ITEMS_PROCESSED	= "items_processed"
		TOTAL_ITEMS_PROCESSED = "total_items_processed"
		TOTAL_ITEMS_REQESTED = "total_items_requested"

	*/

	// CONTEXT KEYS
	ProductCtx EventCtxKey = iota
	AppCtx     EventCtxKey = iota
	TraceIdCtx EventCtxKey = iota
	ModuleCtx  EventCtxKey = iota
	AccountCtx EventCtxKey = iota
	UserCtx    EventCtxKey = iota
)

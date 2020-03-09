package constants

type EventCtxKey int

const (
	UnknownString = "unknown"
	EmptyString   = ""

	// RFC3339Milli is the standard RFC3339 format with added milliseconds
	RFC3339Milli  = "2006-01-02T15:04:05.000Z07:00"

	// JSON LOG KEYS
	// List of standard keys used for logging as per https://cultureamp.atlassian.net/wiki/spaces/TV/pages/959939199/Logging
	TimeLogField         = "time"
	ProductLogField      = "product"
	AppLogField          = "app"
	EventLogField        = "event"
	SeverityLogField     = "severity"
	TraceIdLogField      = "trace_id"
	ModuleLogField       = "module"
	StatusLogField       = "status"
	ArchitectureLogField = "arch"
	HostLogField         = "host"
	OsLogField           = "os"
	PidLogField          = "pid"
	ProcessLogField      = "process"

	AccountLogField = "account"
	UserLogField    = "user"

	TimeTakenLogField           = "time_taken"
	MemoryUsedLogField          = "memory_used"
	MemoryAvailLogField         = "memory_available"
	ItemsProcessedLogField      = "items_processed"
	TotalItemsProcessedLogField = "total_items_processed"
	TotalItemsRequestedLogField = "total_items_requested"

	MessageLogField = "message"

	// Severity Values
	DebugSevLogValue = "DEBUG"
	InfoSevLogValue  = "INFO"
	WarnSevLogValue  = "WARN"
	ErrorSevLogValue = "ERROR"
	FatalSevLogValue = "FATAL"

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

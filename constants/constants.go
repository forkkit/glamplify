package constants

type EventCtxKey int

const (
	UnknownString = "unknown"
	EmptyString   = ""

	// RFC3339Milli is the standard RFC3339 format with added milliseconds
	RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"

	// JSON LOG KEYS
	// List of standard keys used for logging as per https://cultureamp.atlassian.net/wiki/spaces/TV/pages/959939199/Logging
	TraceIdLogField             = "trace_id"
	TimeLogField                = "time"
	EventLogField               = "event"
	ProductLogField             = "product"
	AppLogField                 = "app"
	AppVerLogField              = "app_version"
	SeverityLogField            = "severity"
	RegionLogField              = "region"
	ResourceLogField            = "resource"
	OsLogField                  = "os"
	CustomerLogField            = "customer"
	UserLogField                = "user"
	ExceptionLogField           = "exception"
	MessageLogField				= "message"
	TimeTakenLogField           = "time_taken"
	MemoryUsedLogField          = "memory_used"
	MemoryAvailLogField         = "memory_available"
	ItemsProcessedLogField      = "items_processed"
	TotalItemsProcessedLogField = "total_items_processed"
	TotalItemsRequestedLogField = "total_items_requested"

	// Severity Values
	DebugSevLogValue = "DEBUG"
	InfoSevLogValue  = "INFO"
	WarnSevLogValue  = "WARN"
	ErrorSevLogValue = "ERROR"
	FatalSevLogValue = "FATAL"
	AuditSevLogValue = "AUDIT"

	// ENVIRONMENT VARIABLES
	// List of  Environment Variables keys
	ProductEnv = "PRODUCT"
	AppEnv     = "APP"
	AppVerEnv  = "APP_VERSION"
	RegionEnv  = "REGION"

	// CONTEXT KEYS
	TraceIdCtx  EventCtxKey = iota
	CustomerCtx EventCtxKey = iota
	UserCtx     EventCtxKey = iota
)

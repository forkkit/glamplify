package types

// Fields type, used to pass to Debug, Print and Error.
type Fields map[string]interface{}

type eventKey int

const (
	// List of standard keys used for logging as per https://cultureamp.atlassian.net/wiki/spaces/TV/pages/959939199/Logging
	TIME_LOG     = "time"
	PRODUCT_LOG      = "product"
	APP_LOG          = "app"
	EVENT_LOG        = "event"
	SEVERITY_LOG     = "severity"
	TRACE_ID_LOG     = "trace_id"
	MODULE_LOG       = "module"
	STATUS_LOG       = "status"
	ARCHITECTURE_LOG = "arch"
	HOST_LOG         = "host"
	OS_LOG           = "os"
	PID_LOG          = "pid"
	PROCESS_LOG      = "process"

	ACCOUNT_LOG	 = "account"
	USER_LOG		 = "user"

	TIME_TAKEN	= "time_taken"
	MEMORY_USED	= "memory_used"
	MEMORY_AVAIL	= "memory_available"
	ITEMS_PROCESSED	= "items_processed"
	TOTAL_ITEMS_PROCESSED = "total_items_processed"
	TOTAL_ITEMS_REQESTED = "total_items_requested"

	MESSAGE_LOG      = "message"

	// Severity Values
	DEBUG_SEV_LOG = "DEBUG"
	INFO_SEV_LOG  = "INFO"
	ERROR_SEV_LOG = "ERROR"


	// List of standard keys used for logging as per https://cultureamp.atlassian.net/wiki/spaces/TV/pages/959939199/Logging
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

	ProductCtx eventKey = iota
	AppCtx     eventKey = iota
	TraceIdCtx eventKey = iota
	ModuleCtx  eventKey = iota
	AccountCtx eventKey = iota
	UserCtx    eventKey = iota
)
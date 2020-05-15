package context

type EventCtxKey int

const (
	// CONTEXT KEYS
	TraceIDCtx   EventCtxKey = iota
	RequestIDCtx EventCtxKey = iota
	CustomerCtx  EventCtxKey = iota
	UserCtx      EventCtxKey = iota
)

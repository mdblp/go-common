package context

import (
	"context"
	"net/http"
	"regexp"

	"github.com/google/uuid"
)

var TRACE_SESSION_HEADER = "x-tidepool-trace-session"

type traceSessionKeyType int

const traceSessionKey traceSessionKeyType = iota + 1

// Return the trace session id from the context
func GetTraceSessionId(ctx context.Context) (string, bool) {
	traceSessionId, ok := ctx.Value(traceSessionKey).(string)
	return traceSessionId, ok
}

// SetTraceSessionIdInRequest set the trace session id in the request context
func SetTraceSessionIdInRequest(r *http.Request) *http.Request {
	re := regexp.MustCompile(`[\w]{8}-[\w]{4}-[\w]{4}-[\w]{4}-[\w]{12}`)
	var traceSessionId string
	if traceSessionId = re.FindString(r.Header.Get(TRACE_SESSION_HEADER)); traceSessionId == "" {
		// first occurrence
		traceSessionId = uuid.New().String()
	}
	return r.WithContext(context.WithValue(r.Context(), traceSessionKey, traceSessionId))
}

// set the context id of the context
func SetTraceSessionId(ctx context.Context, traceSessionId string) context.Context {
	return context.WithValue(ctx, traceSessionKey, traceSessionId)
}

// Deprecated: SetTraceSessionIdCtx exists for historical compatibility, please use
// SetTraceSessionIdInRequest instead
func SetTraceSessionIdCtx(r *http.Request) *http.Request {
	return SetTraceSessionIdInRequest(r)
}

// Deprecated: GetTraceSessionIdCtx exists for historical compatibility, please use
// GetTraceSessionId instead
func GetTraceSessionIdCtx(ctx context.Context) (string, bool) {
	return GetTraceSessionId(ctx)
}

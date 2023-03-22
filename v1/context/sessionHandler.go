package context

import (
	"context"
	"net/http"
	"regexp"

	"github.com/google/uuid"
)

var TRACE_SESSION_HEADER = "x-tidepool-trace-session"

type traceSessionKeyType int

const TraceSessionKey traceSessionKeyType = iota + 1

// Return the trace session id from the context
func GetTraceSessionIdCtx(ctx context.Context) (string, bool) {
	traceSessionId, ok := ctx.Value(TraceSessionKey).(string)
	return traceSessionId, ok
}

// Set the trace session id in the request context
func SetTraceSessionIdCtx(r *http.Request) *http.Request {
	re := regexp.MustCompile(`[\w]{8}-[\w]{4}-[\w]{4}-[\w]{4}-[\w]{12}`)
	var traceSessionId string
	if traceSessionId = re.FindString(r.Header.Get(TRACE_SESSION_HEADER)); traceSessionId == "" {
		// first occurrence
		traceSessionId = uuid.New().String()
	}
	return r.WithContext(context.WithValue(r.Context(), TraceSessionKey, traceSessionId))
}

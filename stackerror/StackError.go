// Package stackerror: custom error which carry the stack trace in the error message
package stackerror

import (
	"fmt"
	"runtime"
	"strings"
)

type ClientErrorWriter interface {
	Kind() string
	Message() string
}

func newStackError(message string) error {
	pc := make([]uintptr, 15)
	n := runtime.Callers(4, pc)
	frames := runtime.CallersFrames(pc[:n])
	stackTrace := ""
	for {
		frame, more := frames.Next()
		if !strings.Contains(frame.File, "gin-gonic") &&
			!strings.Contains(frame.File, "gin-contrib") &&
			!strings.Contains(frame.File, "middleware") &&
			!strings.Contains(frame.File, "go-common") &&
			!strings.Contains(frame.File, "go-router") {
			stackTrace += fmt.Sprintln("[", frame.File, frame.Line, frame.Function, "]")
		}
		if !more {
			break
		}
	}
	return fmt.Errorf("%s \n %s", message, stackTrace)
}

type StackError struct {
	error
	kind    string
	message string
	details map[string]interface{}
}

func New(kind string, msg string) StackError {
	return StackError{
		message: msg,
		error:   newStackError(msg),
		details: map[string]interface{}{},
		kind:    kind,
	}
}

func Newf(kind string, message string, args ...interface{}) StackError {
	formatErr := fmt.Sprintf(message, args...)
	return New(kind, formatErr)
}

func NewWithDetails(kind string, msg string, details map[string]interface{}) StackError {
	detailsStr := "details : "
	for key, value := range details {
		detailsStr += fmt.Sprintf("[key=%s,value=%v]", key, value)
	}
	detailsErr := New(kind, msg)
	detailsErr.details = details
	return detailsErr
}

func (ce StackError) Unwrap() error {
	return ce.error
}

func (ce StackError) Message() string {
	return ce.message
}

func (ce StackError) Details() map[string]interface{} {
	return ce.details
}

func (ce StackError) Kind() string {
	return ce.kind
}

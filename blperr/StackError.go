// Package blperr: custom error which carry the stack trace in the error message
package blperr

import (
	"fmt"
	"runtime"
	"strings"
)

type ClientErrorWriter interface {
	Kind() string
	Message() string
}

func newStackError(message string, kind string, details map[string]interface{}) error {
	/*Building stack trace*/
	pc := make([]uintptr, 15)
	n := runtime.Callers(1, pc)
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
	/*Building details map*/
	detailsStr := ""
	if len(details) > 0 {
		for key, value := range details {
			detailsStr += fmt.Sprintf("[key=%s,value=%v]", key, value)
		}
	}
	return fmt.Errorf("kind=[%s] message=[%s] details=[%s] stackTrace=[%s]", kind, message, detailsStr, stackTrace)
}

type StackError struct {
	error
	kind    string
	message string
	details map[string]interface{}
}

func New(kind string, msg string) StackError {
	details := map[string]interface{}{}
	return StackError{
		message: msg,
		error:   newStackError(msg, kind, details),
		details: details,
		kind:    kind,
	}
}

func Newf(kind string, message string, args ...interface{}) StackError {
	formatErr := fmt.Sprintf(message, args...)
	return New(kind, formatErr)
}

func NewWithDetails(kind string, msg string, details map[string]interface{}) StackError {
	return StackError{
		message: msg,
		error:   newStackError(msg, kind, details),
		details: details,
		kind:    kind,
	}
}

func (se StackError) Unwrap() error {
	return se.error
}

func (se StackError) Message() string {
	return se.message
}

func (se StackError) Details() map[string]interface{} {
	return se.details
}

func (se StackError) Kind() string {
	return se.kind
}

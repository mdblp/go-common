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

type PrivateError struct {
	error
	message string
	details map[string]interface{}
}

type PublicError struct {
	PrivateError
	kind string
}

func NewPrivate(msg string) PrivateError {
	return PrivateError{
		message: msg,
		error:   newStackError(msg),
		details: map[string]interface{}{},
	}
}

func NewPrivatef(message string, args ...interface{}) PrivateError {
	formatErr := fmt.Sprintf(message, args...)
	return NewPrivate(formatErr)
}

func NewPrivateWithDetails(msg string, details map[string]interface{}) PrivateError {
	detailsStr := "details : "
	for key, value := range details {
		detailsStr += fmt.Sprintf("[key=%s,value=%v]", key, value)
	}
	detailsErr := NewPrivate(msg)
	detailsErr.details = details
	return detailsErr
}

func (ce PrivateError) Unwrap() error {
	return ce.error
}

func (ce PrivateError) Message() string {
	return ce.message
}

func (ce PrivateError) Details() map[string]interface{} {
	return ce.details
}

func New(kind string, msg string) PublicError {
	return PublicError{
		PrivateError: NewPrivate(msg),
		kind:         kind,
	}
}

func Newf(kind string, message string, args ...interface{}) PublicError {
	formatErr := fmt.Sprintf(message, args...)
	return New(kind, formatErr)
}

func NewWithDetails(kind string, msg string, details map[string]interface{}) PublicError {
	return PublicError{
		PrivateError: NewPrivateWithDetails(msg, details),
		kind:         kind,
	}
}

func (ce PublicError) Kind() string {
	return ce.kind
}

// Package stackerror: custom error which carry the stack trace in the error message
package stackerror

import (
	"fmt"
	"runtime"
	"strings"
)

type ClientError interface {
	Type() string
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

type PublicError struct {
	error
	kind    string
	message string
	details map[string]interface{}
}

func New(kind string, msg string) PublicError {
	return PublicError{
		kind:    kind,
		message: msg,
		error:   newStackError(msg),
		details: map[string]interface{}{},
	}
}

func Newf(kind string, message string, args ...interface{}) PublicError {
	formatErr := fmt.Sprintf(message, args)
	return PublicError{
		kind:    kind,
		message: formatErr,
		error:   newStackError(formatErr),
		details: map[string]interface{}{},
	}
}

func NewWithDetails(kind string, msg string, details map[string]interface{}) PublicError {
	detailsStr := "details : "
	for key, value := range details {
		detailsStr += fmt.Sprintf("[key=%s,value=%v]", key, value)
	}
	return PublicError{
		kind:    kind,
		message: msg,
		error:   newStackError(detailsStr),
		details: map[string]interface{}{},
	}
}

func (ce PublicError) Type() string {
	return ce.kind
}

func (ce PublicError) Message() string {
	return ce.message
}

func (ce PublicError) Unwrap() error {
	return ce.error
}

type PrivateError struct {
	error
	message string
	details map[string]interface{}
}

func NewPrivate(kind string, msg string) PrivateError {
	return PrivateError{
		message: msg,
		error:   newStackError(msg),
		details: map[string]interface{}{},
	}
}

func NewPrivatef(message string, args ...interface{}) PrivateError {
	formatErr := fmt.Sprintf(message, args)
	return PrivateError{
		message: formatErr,
		error:   newStackError(formatErr),
		details: map[string]interface{}{},
	}
}

func NewPrivateWithDetails(msg string, details map[string]interface{}) PrivateError {
	detailsStr := "details : "
	for key, value := range details {
		detailsStr += fmt.Sprintf("[key=%s,value=%v]", key, value)
	}
	return PrivateError{
		message: msg,
		error:   newStackError(detailsStr),
		details: map[string]interface{}{},
	}
}

func (ce PrivateError) Unwrap() error {
	return ce.error
}

// Package stackerror: custom error which carry the stack trace in the error message
package stackerror

import (
	"fmt"
	"runtime"
)

type ClientError interface {
	Type() string
	Message() string
}

func NewLineError(message string) error {
	pc := make([]uintptr, 15)
	n := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return fmt.Errorf("[%s:%d %s] %s", frame.File, frame.Line, frame.Function, message)
}

func WrapLineError(err error) error {
	pc := make([]uintptr, 15)
	n := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return fmt.Errorf("[%s:%d %s] %w", frame.File, frame.Line, frame.Function, err)
}

type StackError struct {
	error
	kind    string
	message string
	details map[string]interface{}
}

func New(kind string, msg string) StackError {
	return StackError{
		kind:    kind,
		message: msg,
		error:   NewLineError(msg),
		details: map[string]interface{}{},
	}
}

func Newf(kind string, message string, args ...interface{}) StackError {
	formatErr := fmt.Sprintf(message, args)
	return StackError{
		kind:    kind,
		message: formatErr,
		error:   NewLineError(formatErr),
		details: map[string]interface{}{},
	}
}

func NewWithDetails(kind string, msg string, details map[string]interface{}) StackError {
	detailsStr := "details : "
	for key, value := range details {
		detailsStr += fmt.Sprintf("[key=%s,value=%v]", key, value)
	}
	return StackError{
		kind:    kind,
		message: msg,
		error:   NewLineError(detailsStr),
		details: map[string]interface{}{},
	}
}

func (ce *StackError) Type() string {
	return ce.kind
}

func (ce *StackError) Message() string {
	return ce.message
}

func (ce *StackError) Unwrap() error {
	return ce.error
}

func (ce *StackError) AddDetail(key string, value interface{}) {
	ce.details[key] = value
}

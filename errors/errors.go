// Package errors: custom errors which carry the stack trace in the error message
package errors

import (
	"fmt"
	"runtime"
)

// New creates a new error composed of an error message and the stack trace
func New(msg string) error {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return fmt.Errorf("[%s:%d %s] : %s", frame.File, frame.Line, frame.Function, msg)
}

// Newf is like New() but it uses the Printf formatting
func Newf(message string, args ...interface{}) error {
	return New(fmt.Sprintf(message, args...))
}

// Wrap returns an error based on an existing error and add stack trace details
func Wrap(errorToWrap error) error {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return fmt.Errorf("[%s:%d %s] : %s", frame.File, frame.Line, frame.Function, errorToWrap.Error())
}

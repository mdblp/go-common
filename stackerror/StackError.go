// Package stackerror: custom error which carry the stack trace in the error message
package stackerror

import (
	"errors"
	"fmt"
	"runtime"
)

type ClientError interface {
	Type() string
	Message() string
}

// Error defines an error with details about the source (function, line number...) and other details
type StackError struct {
	message        string                 // error message
	errType        string                 // error type
	nextError      *StackError            // next error in the chain
	details        map[string]interface{} // optional details
	sourceFilename string                 // name of the file from where the error was fired
	sourceFunction string                 // name of the function from where the error was fired
	lineNumber     int                    // line number where the error was fired
}

// New creates a new error composed of an error message and the stack trace
func newStackError(errType string, msg string, nextError *StackError) StackError {
	pc := make([]uintptr, 15)
	n := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return StackError{
		message:        msg,
		errType:        errType,
		nextError:      nextError,
		sourceFilename: frame.File,
		sourceFunction: frame.Function,
		lineNumber:     frame.Line,
		details:        make(map[string]interface{}),
	}
}

func New(errType string, msg string) StackError {
	return newStackError(errType, msg, nil)
}

// Newf is like New() but it uses the Printf formatting
func Newf(errType string, message string, args ...interface{}) StackError {
	return New(errType, fmt.Sprintf(message, args...))
}

func NewWithDetails(errType string, message string, details map[string]interface{}) StackError {
	err := New(errType, message)
	err.details = details
	return err
}

// Wrap returns an error based on an existing error and add stack trace details
func Wrap(errorToWrap StackError) error {
	return newStackError("errorWrap", errorToWrap.Error(), &errorToWrap)
}

func (err StackError) Unwrap() error {
	return err.nextError
}

func (err StackError) As(target interface{}) bool {
	return errors.As(err.nextError, target)
}

func (err StackError) Error() string {
	detailsString := ""
	for key, value := range err.details {
		detailsString += fmt.Sprintf(" [%s=%v] ", key, value)
	}
	return fmt.Sprintf("[%s:%d %s] : %s => %s", err.sourceFilename, err.lineNumber, err.sourceFunction, err.message, detailsString)
}

func (err StackError) AddDetail(key string, value interface{}) StackError {
	err.details[key] = value
	return err
}

func (err StackError) Type() string {
	return err.errType
}

func (err StackError) Message() string {
	return err.message
}

// Package stackerror: custom error which carry the stack trace in the error message
package stackerror

import (
	"fmt"
	"runtime"
)

// Error defines an error with details about the source (function, line number...) and other details
type Error struct {
	message        string                 // error message
	details        map[string]interface{} // optional details
	sourceFilename string                 // name of the file from where the error was fired
	sourceFunction string                 // name of the function from where the error was fired
	lineNumber     int                    // line number where the error was fired
}

// New creates a new error composed of an error message and the stack trace
func New(msg string) Error {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return Error{
		message:        msg,
		sourceFilename: frame.File,
		sourceFunction: frame.Function,
		lineNumber:     frame.Line,
	}
}

// Newf is like New() but it uses the Printf formatting
func Newf(message string, args ...interface{}) Error {
	return New(fmt.Sprintf(message, args...))
}

func NewWithDetails(message string, details map[string]interface{}) Error {
	err := New(message)
	err.details = details
	return err
}

func (err Error) Error() string {
	return fmt.Sprintf("[%s:%d %s] : %s", err.sourceFilename, err.lineNumber, err.sourceFunction, err.message)
}

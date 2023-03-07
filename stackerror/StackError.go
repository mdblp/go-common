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
	n := runtime.Callers(4, pc)
	frames := runtime.CallersFrames(pc[:n])
	framesStr := ""
	for {
		frame, more := frames.Next()
		framesStr += " " + fmt.Sprintf("[%s:%d %s] \n", frame.File, frame.Line, frame.Function)
		if !more {
			break
		}
	}
	return fmt.Errorf("%s %s", framesStr, message)
}

func WrapLineError(err error) error {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return fmt.Errorf("[%s:%d %s] %w", frame.File, frame.Line, frame.Function, err)
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
		error:   NewLineError(msg),
		details: map[string]interface{}{},
	}
}

func Newf(kind string, message string, args ...interface{}) PublicError {
	formatErr := fmt.Sprintf(message, args)
	return PublicError{
		kind:    kind,
		message: formatErr,
		error:   NewLineError(formatErr),
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
		error:   NewLineError(detailsStr),
		details: map[string]interface{}{},
	}
}

func (ce *PublicError) Type() string {
	return ce.kind
}

func (ce *PublicError) Message() string {
	return ce.message
}

func (ce *PublicError) Unwrap() error {
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
		error:   NewLineError(msg),
		details: map[string]interface{}{},
	}
}

func NewPrivatef(message string, args ...interface{}) PrivateError {
	formatErr := fmt.Sprintf(message, args)
	return PrivateError{
		message: formatErr,
		error:   NewLineError(formatErr),
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
		error:   NewLineError(detailsStr),
		details: map[string]interface{}{},
	}
}

func (ce *PrivateError) Unwrap() error {
	return ce.error
}

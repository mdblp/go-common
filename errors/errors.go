package errors

import (
	"errors"
	"fmt"
	"runtime"
)

type WrappedError struct {
	error
	m string
}

func (err *WrappedError) Error() string {
	return err.m + ": " + err.error.Error()
}

func New(message string) error {
	return errors.New(message)
}

func Newf(message string, args ...interface{}) error {
	return errors.New(fmt.Sprintf(message, args...))
}

func Wrap(e error, m string) error {
	return &WrappedError{error: e, m: m}
}

func Wrapf(e error, m string, args ...interface{}) error {
	return Wrap(e, fmt.Sprintf(m, args...))
}

// Create a new error which contains the stack trace
func NewDetailedError(msg string) error {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return fmt.Errorf("[%s:%d %s] : %s", frame.File, frame.Line, frame.Function, msg)
}

func WrapDetailedError(errorToWrap error) error {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return fmt.Errorf("[%s:%d %s] : %s", frame.File, frame.Line, frame.Function, errorToWrap.Error())
}

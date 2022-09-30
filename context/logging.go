package context

import (
	"context"
	"os"
	log "github.com/sirupsen/logrus"
)

type loggerKeyType int
const LoggerKey loggerKeyType = iota

var logger *log.Logger

func init() {
	logger = log.New()
	logger.Out = os.Stdout
	logger.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	logger.SetReportCaller(true)
	envLogLevel := os.Getenv("LOG_LEVEL")
	logLevel, err := log.ParseLevel(envLogLevel)
	if err != nil {
		logLevel = log.WarnLevel
	}
	logger.SetLevel(logLevel)
}
func NewLog() *log.Entry {
	return log.NewEntry(logger)
}


// WithLogger returns a context with the provided logger. Use in
// combination with logger.WithField(s) for great effect.
func WithLogger(ctx context.Context, logger *log.Entry) context.Context {
	//l := logger.WithContext(ctx)
	return context.WithValue(ctx, LoggerKey, logger)
}

// GetLogger retrieves the current logger from the context. If no logger is
// available, the default logger is returned.
func GetLogger(ctx context.Context) *log.Entry {
	if ctx != nil {
		newLogger := ctx.Value(LoggerKey)

		if newLogger == nil {
			log.Debug("Logger is missing in the context, use a the root one") // panics
			newLogger = log.NewEntry(logger)
		}

		return newLogger.(*log.Entry)
	} else {
		return log.NewEntry(logger)
	}
}
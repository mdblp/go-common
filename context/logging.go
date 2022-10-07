package context

import (
	"context"
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"os"
	"strconv"
	"strings"
)

type loggerKeyType int

const LoggerKey loggerKeyType = iota

var logger *log.Logger

func init() {
	logger = log.New()
	logger.Out = os.Stdout
	envLogFormatter := os.Getenv("LOG_FORMATTER")
	logFormatter := parseFormatter(envLogFormatter)
	logger.SetFormatter(logFormatter)
	reportCaller := false
	envLogReportCaller, err := strconv.ParseBool(os.Getenv("LOG_REPORT_CALLER"))
	if err == nil && envLogReportCaller {
		reportCaller = true
	}
	logger.SetReportCaller(reportCaller)
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

// parseFormatter takes a string that describe as the format and returns the Logrus log formatter.
func parseFormatter(formatter string) log.Formatter {
	switch strings.ToLower(formatter) {
	case "json":
		return &DBLJSONFormatter{}
	case "text":
		return &log.TextFormatter{
			DisableColors: false,
			FullTimestamp: true,
		}
	case "simple":
		return &easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat:       "%time% [%lvl%]:  %msg%",
		}
	default:
		return &log.TextFormatter{
			DisableColors: false,
			FullTimestamp: true,
		}
	}
}

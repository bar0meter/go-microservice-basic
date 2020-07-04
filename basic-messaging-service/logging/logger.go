package logging

import (
	"io"
	logging "log"
	"os"

	"github.com/sirupsen/logrus"
)

// LogWrapper wraps logrus logger. We can use this to wrap other logger also.
type LogWrapper struct {
	logrusLogger *logrus.Logger
}

// NewLogger returns a new LogWrapper instance
// We are wrapping logrus logger here. But we can use this for any other loggers also.
// We are also defining Info, Warn, Error methods for Wrapper.
// Even though we change logger, we can still use log.Info, log.Warn,
// log.Error methods without worrying about the logger
func NewLogger() *LogWrapper {

	// Here for development purpose we are overwriting log file
	// While deploying to production, make sure to use append here
	f, err := os.OpenFile("logs/application.log", os.O_RDWR|os.O_CREATE, 0666)
	//f, err := os.OpenFile("logs/application.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		logging.Fatalf("error opening file: %v", err)
	}

	log := &logrus.Logger{
		Out:          os.Stderr,
		Formatter:    new(logrus.TextFormatter),
		Hooks:        make(logrus.LevelHooks),
		Level:        logrus.DebugLevel,
		ReportCaller: false,
	}

	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)

	return &LogWrapper{
		logrusLogger: log,
	}
}

// Info ...
func (lw *LogWrapper) Info(format string, v ...interface{}) {
	lw.logrusLogger.Infof(format, v...)
}

// Warn ...
func (lw *LogWrapper) Warn(format string, v ...interface{}) {
	lw.logrusLogger.Warnf(format, v...)
}

// Error ...
func (lw *LogWrapper) Error(format string, v ...interface{}) {
	lw.logrusLogger.Errorf(format, v...)
}

var (
	// ConfigError ...
	ConfigError = "%v type=config.error"

	// HTTPError ...
	HTTPError = "%v type=http.error"

	// HTTPWarn ...
	HTTPWarn = "%v type=http.warn"

	// HTTPInfo ...
	HTTPInfo = "%v type=http.info"
)

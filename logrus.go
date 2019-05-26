package logrus

import (
	"fmt"
	"github.com/go-coder/logr"
	slog "github.com/sirupsen/logrus"
)

type ruslog struct {
	name string
	logger *slog.Logger
}

var _ logr.InfoLogger = (*ruslog)(nil)
var _ logr.Logger = (*ruslog)(nil)

// NewLogger creates a new logr.Logger using logrus
func NewLogger(name string, l *slog.Logger)  logr.Logger {
	if l == nil {
		panic("non-nil l *logrus.Logger must be provided as 2nd parameter")
	}
	return ruslog {
		name: name,
		logger: l,
	}
}

// Enabled tests whether this InfoLogger is enabled.  For example,
// commandline flags might be used to set the logging verbosity and disable
// some info logs.
func (l ruslog) Enabled() bool {
	return l.logger.IsLevelEnabled(l.logger.GetLevel())
}

// Info logs the given message and key/value pairs, 
// string keys and arbitrary values are required.
func (l ruslog) Info(msg string, keysAndValues ...interface{}) {
	fields := getFields(keysAndValues...)
	fields["name"] = l.name
	fields["msg"] = msg
	l.logger.WithFields(fields).Info(msg)
}

// Error logs an error, with the given message and key/value pairs as context.
func (l ruslog) Error(err error, msg string, keysAndValues ...interface{}) {
	fields := getFields(keysAndValues...)
	fields["name"] = l.name
	fields["msg"] = msg
	fields["err"] = err
	l.logger.WithFields(fields).Error(msg)
}

// V returns an InfoLogger value for a specific verbosity level.  
// A higher verbosity level means a log message is less important. 
// It's illegal to pass a log level less than zero.
func (l ruslog) V(level int) logr.InfoLogger {
	if level < 0 || level >= len(slog.AllLevels) {
		panic(fmt.Sprintf("invalid log level [%d] is set", level))
	}
	newLogger := ruslog { 
		name: l.name, 
		logger: l.logger,
	}
	newLogger.logger.SetLevel(slog.AllLevels[level])
	return newLogger
}

// WithFields adds some key-value pairs of context to a logger.
func (l ruslog) WithFields(keysAndValues ...interface{}) logr.Logger {
	newLogger := NewLogger(l.name, l.logger)
	newLogger.WithFields(getFields(keysAndValues...))
	return newLogger
}

// WithName adds a new element to the logger's name.
func (l ruslog) WithName(name string) logr.Logger {
	newLogger := ruslog { 
		name: l.name, 
		logger: l.logger, 
	}
	if len(l.name) > 0 {
		newLogger.name += "." + name
	} else {
		newLogger.name = name
	}
	return newLogger
}

// getFields returns logrus.Fields (aka. map[string]interface{}) used for structured log
func getFields(keysAndValues ...interface{}) slog.Fields {
	if len(keysAndValues) % 2 != 0 {
		panic("keysAndValues is not valid")
	}
	fields := slog.Fields {}
	for i := 0; i < len(keysAndValues); i+=2 {
		key, ok := keysAndValues[i].(string)
		if !ok {
			panic(fmt.Sprintf("key [%v] is not a string type", keysAndValues[i]))
		}
		fields[key] = keysAndValues[i+1]
	}
	return fields
}
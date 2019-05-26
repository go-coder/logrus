package logrus

import (
	"fmt"
	"sort"

	"github.com/go-coder/logr"
	slog "github.com/sirupsen/logrus"
)

const (
	Prefix = "prefix"
	Time = slog.FieldKeyTime // "time"
	Level = slog.FieldKeyLevel // "level"
	Message = slog.FieldKeyMsg //"msg"
	Error = slog.FieldKeyLogrusError // "logrus_error"
)

var (
	Header = []string {
		Prefix, Time, Level, Message, Error,
	}
	formatter = &slog.TextFormatter{
		DisableSorting: false,
		SortingFunc: func(keys []string) {
			sort.Slice(keys, func(i, j int) bool {
				for _, key:= range Header {
					if keys[i] == key {
						return true
					} else if keys[j] == key {
						return false
					} 
				} // Header 顺序
				return keys[i] < keys[j]
			})
		},
	}
)

type ruslog struct {
	name string
	entry *slog.Entry
}

var _ logr.InfoLogger = (*ruslog)(nil)
var _ logr.Logger = (*ruslog)(nil)

// NewLogger creates a new logr.Logger using logrus
func NewLogger(name string, l *slog.Logger)  logr.Logger {
	if l == nil {
		panic("non-nil l *logrus.Logger must be provided as 2nd parameter")
	}
	l.SetFormatter(formatter)
	return ruslog {
		name: name,
		entry: slog.NewEntry(l),
	}
}

// Enabled tests whether this InfoLogger is enabled.  For example,
// commandline flags might be used to set the logging verbosity and disable
// some info logs.
func (l ruslog) Enabled() bool {
	return l.entry.Logger.IsLevelEnabled(l.entry.Logger.GetLevel())
}

// Info logs the given message and key/value pairs, 
// string keys and arbitrary values are required.
func (l ruslog) Info(msg string, keysAndValues ...interface{}) {
	fields := getFields(keysAndValues...)
	fields[Prefix] = l.name
	l.entry.WithFields(fields).Info(msg)
}

// Error logs an error, with the given message and key/value pairs as context.
func (l ruslog) Error(err error, msg string, keysAndValues ...interface{}) {
	fields := getFields(keysAndValues...)
	fields[Prefix] = l.name
	l.entry.WithFields(fields).WithError(err).Error(msg)
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
		entry: l.entry,
	}
	newLogger.entry.Logger.SetLevel(slog.AllLevels[level])
	return newLogger
}

// WithFields adds some key-value pairs of context to a logger.
func (l ruslog) WithFields(keysAndValues ...interface{}) logr.Logger {
	newLogger := ruslog { 
		name: l.name, 
		entry: l.entry.WithFields(getFields(keysAndValues...)), 
	}
	return newLogger
}

// WithName adds a new element to the logger's name.
func (l ruslog) WithName(name string) logr.Logger {
	newLogger := ruslog { 
		name: l.name, 
		entry: l.entry, 
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
		panic(fmt.Sprintf("keysAndValues is not valid: %v", keysAndValues))
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
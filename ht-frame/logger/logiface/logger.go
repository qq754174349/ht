package logiface

import (
	"io"
)

type Log struct {
	Log Config
}

type Config struct {
	Level       string
	OutputPaths string `json:"output_paths" yaml:"output_paths" mapstructure:"output_paths"`
}

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
	WithTraceID(traceId string) Logger
	WithFields(fields map[string]interface{}) Logger
	RedirectStdLog()
	Writer() io.Writer
}

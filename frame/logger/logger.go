package logger

import (
	"io"

	"github.com/qq754174349/ht/ht-frame/autoconfigure"
	"github.com/qq754174349/ht/ht-frame/logger/internal"
	"github.com/qq754174349/ht/ht-frame/logger/logiface"
)

var config *logiface.Log
var log logiface.Logger

type AutoConfig struct{}

const (
	defaultLevel       = "info"
	defaultOutputPaths = "logs/"
)

func init() {
	autoconfigure.Register(&AutoConfig{})
	config = &logiface.Log{}
	autoconfigure.ConfigRead(config)
	logConfig := config.Log
	InitLogger(logConfig)
}

func (AutoConfig) Close() error {
	return nil
}

func InitLogger(config logiface.Config) {
	if config.Level == "" {
		config.Level = defaultLevel
	}
	if config.OutputPaths == "" {
		config.OutputPaths = defaultOutputPaths
	}

	log = internal.NewZapLog(&config)
	// 替换标准日志
	log.RedirectStdLog()
}

func Debug(args ...interface{})                   { log.Debug(args...) }
func Info(args ...interface{})                    { log.Info(args...) }
func Warn(args ...interface{})                    { log.Warn(args...) }
func Error(args ...interface{})                   { log.Error(args...) }
func Fatal(args ...interface{})                   { log.Fatal(args...) }
func Debugf(template string, args ...interface{}) { log.Debugf(template, args...) }
func Infof(template string, args ...interface{})  { log.Infof(template, args...) }
func Warnf(template string, args ...interface{})  { log.Warnf(template, args...) }
func Errorf(template string, args ...interface{}) { log.Errorf(template, args...) }
func Fatalf(template string, args ...interface{}) { log.Fatalf(template, args...) }

func WithTraceID(traceId string) logiface.Logger {
	return log.WithTraceID(traceId)
}

func WithFields(fields map[string]interface{}) logiface.Logger {
	return log.WithFields(fields)
}

func Writer() io.Writer {
	return log.Writer()
}

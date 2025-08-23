package internal

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/qq754174349/ht/ht-frame/logger/logiface"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const SKIP = 2

type zapLog struct {
	logger *zap.SugaredLogger
}

func NewZapLog(logCfg *logiface.Config) logiface.Logger {
	level, err := zap.ParseAtomicLevel(logCfg.Level)
	if err != nil {
		panic(err)
	}
	encoder := zapcore.NewConsoleEncoder(*newEncoderConfig())

	// 控制台输出日志
	stdWriter := zapcore.AddSync(os.Stdout)
	stdoutEncoder := zapcore.NewCore(encoder, stdWriter, level)

	// 文件输出日志
	fileWriter := zapcore.AddSync(newFileWriteSyncer(logCfg.OutputPaths, logCfg))
	fileEncoder := zapcore.NewCore(encoder, fileWriter, level)

	// 日志采样
	sampler := zapcore.NewSamplerWithOptions(
		zapcore.NewTee(stdoutEncoder, fileEncoder),
		time.Second, // 采样间隔
		10,          // 每秒最多记录 10 条日志
		10,          // 每秒最多记录 10 条日志
	)

	logger := zap.New(sampler, zap.WithCaller(true))
	// 全局使用
	zap.ReplaceGlobals(logger)
	sugar := logger.Sugar()
	sugar.Info("init log complete")
	sugar.Info("use zap log")
	sugar.Infof("current log level %s", logCfg.Level)

	return &zapLog{logger: sugar.WithOptions(zap.AddCallerSkip(SKIP))}
}

func newEncoderConfig() *zapcore.EncoderConfig {
	return &zapcore.EncoderConfig{
		MessageKey:    "msg",
		LevelKey:      "level",
		TimeKey:       "time",
		CallerKey:     "caller",
		StacktraceKey: "stacktrace",
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime:    zapcore.TimeEncoderOfLayout(time.DateTime),
		EncodeCaller:  zapcore.FullCallerEncoder,
	}
}

func newFileWriteSyncer(outPutPath string, cfg *logiface.Config) *lumberjack.Logger {
	logFile := fmt.Sprintf("%s-%s.log", outPutPath, time.Now().Format("2006-01-02-15"))
	return &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10,   // 考虑从配置中读取
		MaxBackups: 3,    // 考虑从配置中读取
		MaxAge:     28,   // 考虑从配置中读取
		Compress:   true, // 考虑从配置中读取
	}
}

func (log zapLog) Debug(args ...interface{})                   { log.logger.Debug(args...) }
func (log zapLog) Info(args ...interface{})                    { log.logger.Info(args...) }
func (log zapLog) Warn(args ...interface{})                    { log.logger.Warn(args...) }
func (log zapLog) Error(args ...interface{})                   { log.logger.Error(args...) }
func (log zapLog) Fatal(args ...interface{})                   { log.logger.Fatal(args...) }
func (log zapLog) Debugf(template string, args ...interface{}) { log.logger.Debugf(template, args...) }
func (log zapLog) Infof(template string, args ...interface{})  { log.logger.Infof(template, args...) }
func (log zapLog) Warnf(template string, args ...interface{})  { log.logger.Warnf(template, args...) }
func (log zapLog) Errorf(template string, args ...interface{}) { log.logger.Errorf(template, args...) }
func (log zapLog) Fatalf(template string, args ...interface{}) { log.logger.Fatalf(template, args...) }

func (log zapLog) WithTraceID(traceId string) logiface.Logger {
	return &zapLog{logger: log.logger.With("traceId", traceId).WithOptions(zap.AddCallerSkip(-1))}
}

func (log zapLog) WithFields(fields map[string]interface{}) logiface.Logger {
	zapFields := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		zapFields = append(zapFields, k, v)
	}
	return &zapLog{logger: log.logger.With(zapFields...).WithOptions(zap.AddCallerSkip(-1))}
}

func (log zapLog) RedirectStdLog() {
	zap.RedirectStdLog(log.logger.Desugar().WithOptions(zap.AddCallerSkip(-SKIP)))
}

func (log zapLog) Writer() io.Writer {
	ginLog := log.logger.Desugar().WithOptions(zap.AddCallerSkip(-SKIP - 1))
	return zap.NewStdLog(ginLog).Writer()
}

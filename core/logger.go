package core

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
	"time"
)

type Logger struct {
	logger *zap.Logger
}

func NewLogger() (*Logger, error) {
	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	zapConfig.EncoderConfig.TimeKey = "time"
	zapConfig.EncoderConfig.StacktraceKey = "stacktrace"
	zapConfig.EncoderConfig.EncodeLevel = ColorLevelEncoder
	zapConfig.EncoderConfig.EncodeCaller = PackagePathEncoder
	zapConfig.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)

	l, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{
		logger: l.WithOptions(
			zap.AddCaller(),
			zap.AddStacktrace(zapcore.ErrorLevel),
		),
	}, nil
}

func (o *Logger) Get() *zap.Logger {
	return o.logger
}

func (o *Logger) Close() {
	err := o.logger.Sync()
	if err != nil {
		println(err.Error())
	}
}

func PackagePathEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	if caller.Defined {

		cwd, err := os.Getwd()
		if err != nil {
			enc.AppendString(fmt.Sprintf("error: %v", err))
			return
		}

		fullPath := caller.FullPath()

		trimmedPath := strings.TrimPrefix(fullPath, cwd+"/")
		enc.AppendString(fmt.Sprintf("%s:%d", trimmedPath, caller.Line))
	} else {
		enc.AppendString("undefined")
	}
}

func ColorLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	const (
		ColorReset  = "\033[0m"
		ColorRed    = "\033[31m"
		ColorGreen  = "\033[32m"
		ColorYellow = "\033[33m"
		ColorBlue   = "\033[34m"
		ColorPurple = "\033[35m"
		ColorWhite  = "\033[37m"
	)

	var color string
	switch level {
	case zapcore.DebugLevel:
		color = ColorBlue
	case zapcore.InfoLevel:
		color = ColorGreen
	case zapcore.WarnLevel:
		color = ColorYellow
	case zapcore.ErrorLevel:
		color = ColorRed
	case zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		color = ColorPurple
	default:
		color = ColorWhite
	}

	//goland:noinspection GoDfaNilDereference
	enc.AppendString(fmt.Sprintf("%s%s%s", color, level.CapitalString(), ColorReset))
}

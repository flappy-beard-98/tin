package core

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	zapConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)

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

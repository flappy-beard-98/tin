package adapter

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

type Logger interface {
	Infof(template string, args ...any)
	Errorf(template string, args ...any)
	Fatalf(template string, args ...any)
	Close()
}

type zapLoggerAdapter struct {
	logger *zap.Logger
}

func NewLogger() (Logger, error) {
	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	zapConfig.EncoderConfig.TimeKey = "time"

	l, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}

	return &zapLoggerAdapter{
		logger: l.WithOptions(zap.AddCallerSkip(1)),
	}, nil
}

func (o *zapLoggerAdapter) Close() {
	_ = o.logger.Sync()
}

func (o *zapLoggerAdapter) Infof(template string, args ...any) {
	o.logger.Sugar().Infof(template, args...)
}
func (o *zapLoggerAdapter) Errorf(template string, args ...any) {
	o.logger.Sugar().Errorf(template, args...)
}
func (o *zapLoggerAdapter) Fatalf(template string, args ...any) {
	o.logger.Sugar().Fatalf(template, args...)
}

package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var l *zap.Logger

type Logger interface {
	Error(string)
	Warning(string)
	Info(string)
	LogRest(string, string, string)
}

type LogService struct {
	*zap.Logger
}

func NewLogger() Logger {
	var cfg zap.Config = zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	cfg.EncoderConfig.CallerKey = ""
	cfg.EncoderConfig.StacktraceKey = ""
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stdout"}
	l, _ = cfg.Build()
	return &LogService{}
}

//Logs errors message
func (log *LogService) Error(msg string) {
	l.Error(msg)
}

func (log *LogService) Warning(msg string) {
	l.Warn(msg)
}

func (log *LogService) Info(msg string) {
	l.Info(msg)
}

// ******************************************************************************************************************************************
//Logs REST API event
func (log *LogService) LogRest(msg, route, method string) {
	l.Info("",
		zap.String("route", route),
		zap.String("method", method),
		zap.String("body", msg),
	)
}

package logger

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

type Logger struct {
    *zap.Logger
}

func NewLogger() *Logger {
    config := zap.NewProductionConfig()
    config.EncoderConfig.TimeKey = "timestamp"
    config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

    logger, _ := config.Build()
    return &Logger{logger}
}

func (l *Logger) Info(msg string, fields ...interface{}) {
    if len(fields) > 0 {
        if fieldMap, ok := fields[0].(map[string]interface{}); ok {
            zapFields := make([]zap.Field, 0, len(fieldMap))
            for k, v := range fieldMap {
                zapFields = append(zapFields, zap.Any(k, v))
            }
            l.Logger.Info(msg, zapFields...)
            return
        }
    }
    l.Logger.Info(msg)
}

func (l *Logger) Error(msg string, err error, fields ...interface{}) {
    zapFields := []zap.Field{zap.Error(err)}
    if len(fields) > 0 {
        if fieldMap, ok := fields[0].(map[string]interface{}); ok {
            for k, v := range fieldMap {
                zapFields = append(zapFields, zap.Any(k, v))
            }
        }
    }
    l.Logger.Error(msg, zapFields...)
}

func (l *Logger) Fatal(msg string, err error) {
    l.Logger.Fatal(msg, zap.Error(err))
}

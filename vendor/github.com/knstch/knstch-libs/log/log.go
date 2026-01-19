package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Level string

type Logger struct {
	lg *zap.Logger
}

type Message struct {
	key   string
	value interface{}
}

func AddMessage(key string, value interface{}) Message {
	return Message{key, value}
}

const (
	DebugLevel  Level = "debug"
	InfoLevel   Level = "info"
	WarnLevel   Level = "warn"
	DPanicLevel Level = "dpanic"
	PanicLevel  Level = "panic"
	FalalLevel  Level = "fatal"
)

func NewLogger(serviceName string, level Level) *Logger {
	levelsToZap := map[Level]zapcore.Level{
		DebugLevel:  zapcore.DebugLevel,
		InfoLevel:   zapcore.InfoLevel,
		WarnLevel:   zapcore.WarnLevel,
		DPanicLevel: zapcore.DPanicLevel,
		PanicLevel:  zapcore.PanicLevel,
		FalalLevel:  zapcore.FatalLevel,
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(&lumberjack.Logger{
			Filename:   `./log/` + serviceName + `_logfile.log`,
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     28,
		}), levelsToZap[level]),
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(&lumberjack.Logger{
			Filename:   `./log/` + serviceName + `_error.log`,
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     28,
		}), zap.ErrorLevel),
	)

	return &Logger{
		lg: zap.New(core),
	}
}

func getFields(fields ...Message) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))

	for _, v := range fields {
		zapFields = append(zapFields, zap.Any(v.key, v.value))
	}

	return zapFields
}

func (l *Logger) Error(msg string, err error, fields ...Message) {
	allFields := append([]zap.Field{zap.Error(err)}, getFields(fields...)...)
	l.lg.Error(msg, allFields...)
}

func (l *Logger) Info(msg string, fields ...Message) {
	l.lg.Info(msg, getFields(fields...)...)
}

func (l *Logger) Debug(msg string, fields ...Message) {
	l.lg.Debug(msg, getFields(fields...)...)
}

func (l *Logger) With(fields ...Message) *Logger {
	return &Logger{
		lg: l.lg.With(getFields(fields...)...),
	}
}

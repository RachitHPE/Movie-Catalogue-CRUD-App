package log

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	log *zap.Logger
}

// ZapOption enables extending the default logger.
type ZapOption func(l *zapLogger)

// NewZapLog returns zap logger which implements Logger interface
func NewZapLog(level string, opts ...ZapOption) Logger {
	cfg := zap.NewProductionConfig()

	switch level {
	case "debug":
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	default:
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	cfgE := zap.NewProductionEncoderConfig()
	cfgE.EncodeTime = zapcore.ISO8601TimeEncoder
	cfgE.EncodeCaller = zapcore.FullCallerEncoder
	cfgE.TimeKey = time
	cfgE.CallerKey = caller
	cfg.EncoderConfig = cfgE
	cfg.DisableStacktrace = true

	log, err := cfg.Build(zap.AddCallerSkip(skipLevel))
	if err != nil {
		panic(err)
	}

	zapLog := &zapLogger{log: log}

	// process log options
	for _, o := range opts {
		if o != nil {
			o(zapLog)
		}
	}

	return zapLog
}

// Info use zap log to log info level log
func (zapLog *zapLogger) Info(ctx context.Context, args ...interface{}) {
	defer zapLog.log.Sync()
	zapLog.addAdditionalField(ctx, logType).Sugar().Info(args...)
}

// Infof use zap log to log info level log
func (zapLog *zapLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	defer zapLog.log.Sync()
	zapLog.addAdditionalField(ctx, logType).Sugar().Infof(format, args...)
}

// Warn use zap log to log warning level log
func (zapLog *zapLogger) Warn(ctx context.Context, args ...interface{}) {
	defer zapLog.log.Sync()
	zapLog.addAdditionalField(ctx, logType).Sugar().Warn(args...)
}

// Warnf use zap log to log warning level log
func (zapLog *zapLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	defer zapLog.log.Sync()
	zapLog.addAdditionalField(ctx, logType).Sugar().Warnf(format, args...)
}

// Error use zap log to log error level log
func (zapLog *zapLogger) Error(ctx context.Context, args ...interface{}) {
	defer zapLog.log.Sync()
	logger := zapLog.addAdditionalField(ctx, logType)
	zapLog.addStackTrace(logger, args...).Sugar().Error(args...)
}

// Errorf use zap log to log error level log
func (zapLog *zapLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	defer zapLog.log.Sync()
	logger := zapLog.addAdditionalField(ctx, logType)
	zapLog.addStackTrace(logger, args...).Sugar().Errorf(format, args...)
}

// Debug use zap log to log debug level log
func (zapLog *zapLogger) Debug(ctx context.Context, args ...interface{}) {
	defer zapLog.log.Sync()
	zapLog.addAdditionalField(ctx, logType).Sugar().Debug(args...)
}

// Debugf use zap log to log debug level log
func (zapLog *zapLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	defer zapLog.log.Sync()
	zapLog.addAdditionalField(ctx, logType).Sugar().Debugf(format, args...)
}

// Audit use zap log to log audit log
func (zapLog *zapLogger) Audit(ctx context.Context, args ...interface{}) {
	defer zapLog.log.Sync()
	zapLog.addAdditionalField(ctx, auditType).Sugar().Info(args...)
}

// Auditf use zap log to log info level log
func (zapLog *zapLogger) Auditf(ctx context.Context, format string, args ...interface{}) {
	defer zapLog.log.Sync()
	zapLog.addAdditionalField(ctx, auditType).Sugar().Infof(format, args...)
}

func getRequestIDFromContext(ctx context.Context) string {
	rID, ok := ctx.Value("X-Request-ID").(string)
	if !ok {
		rID = ""
	}

	return rID
}

func (zapLog *zapLogger) addAdditionalField(ctx context.Context, lType string) *zap.Logger {
	rID := getRequestIDFromContext(ctx)

	reqField := zap.String(reqID, rID)
	typeField := zap.String(typeKey, lType)

	return zapLog.log.With(reqField).With(typeField)
}

func (zapLog *zapLogger) createZapFields(fields map[string]interface{}) []interface{} {
	zapFields := make([]interface{}, 0)
	for key, value := range fields {
		zapFields = append(zapFields, key)
		zapFields = append(zapFields, value)
	}

	return zapFields
}

func (zapLog *zapLogger) addStackTrace(logger *zap.Logger, args ...interface{}) *zap.Logger {
	if len(args) == 0 {
		return logger
	}
	if st := createStackTraceMap(args[0]); st != nil {
		trace := zap.Any(stackTrace, st)
		return logger.With(trace)
	}

	return logger
}

func getEnvOrDefault(envVar string, defaultValue string) string {
	if v := os.Getenv(envVar); v != "" {
		return v
	}
	return defaultValue
}
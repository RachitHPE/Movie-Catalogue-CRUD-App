package log

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

const (
	debug = "debug"
	info  = "info"

	caller     = "context"
	reqID      = "requestID"
	tenantID   = "tenantID"
	userID     = "userID"
	time       = "@timestamp"
	stackTrace = "stacktrace"
	typeKey    = "type"
	dedKey     = "DED"

	auditType = "audit"
	logType   = "log"
)

var (
	stackTraceDepth = 3
	skipLevel       = 3
	level           = getEnvOrDefault("LOG_LEVEL", info)
)

type Config struct {
	SkipLevel int
	DED       string
}

type causer interface {
	Cause() error
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

var logger = NewZapLog(info)

// ConfigureLogger configures logger base on env variables
func ConfigureLogger() {
	logger = NewZapLog(level)
}

// Info logs info level
func Info(ctx context.Context, args ...interface{}) {
	logger.Info(ctx, args...)
}

// Infof logs info level
func Infof(ctx context.Context, format string, args ...interface{}) {
	logger.Infof(ctx, format, args...)
}

// Warn logs warning level
func Warn(ctx context.Context, args ...interface{}) {
	logger.Warn(ctx, args...)
}

// Warnf logs warning level
func Warnf(ctx context.Context, format string, args ...interface{}) {
	logger.Warnf(ctx, format, args...)
}

// Debug logs debug level
func Debug(ctx context.Context, args ...interface{}) {
	logger.Debug(ctx, args...)
}

// Debugf logs debug level
func Debugf(ctx context.Context, format string, args ...interface{}) {
	logger.Debugf(ctx, format, args...)
}

// Error logs error level
func Error(ctx context.Context, args ...interface{}) {
	logger.Error(ctx, args...)
}

// Errorf logs error level
func Errorf(ctx context.Context, format string, args ...interface{}) {
	logger.Errorf(ctx, format, args...)
}

// Audit logs info level
func Audit(ctx context.Context, args ...interface{}) {
	logger.Audit(ctx, args...)
}

// Auditf logs info level
func Auditf(ctx context.Context, format string, args ...interface{}) {
	logger.Auditf(ctx, format, args...)
}

func createStackTraceMap(err interface{}) []map[string]string {
	trace := getTopStack(err)
	if trace == nil {
		return nil
	}
	depth := stackTraceDepth
	if len(trace) <= stackTraceDepth {
		depth = len(trace) - 1
	}

	stackTrace := make([]map[string]string, depth)
	for i := 1; i <= depth; i++ {
		valued := fmt.Sprintf("%+v", trace[i])
		valued = strings.Trim(valued, "\n")
		valued = strings.Replace(valued, "\t", "", -1)
		stack := strings.Split(valued, "\n")
		stackTrace[i-1] = map[string]string{"file": stack[1], "func": stack[0]}
	}

	return stackTrace
}

func getTopStack(err interface{}) errors.StackTrace {
	var topStackInfo stackTracer
	for {
		stackErr, ok := err.(stackTracer)
		if ok {
			topStackInfo = stackErr
		}

		cause, ok := err.(causer)
		if !ok {
			break
		}

		err = cause.Cause()
	}

	if topStackInfo != nil {
		return topStackInfo.StackTrace()
	}

	return nil
}

// SetSkipLevel adjust skip level as needed.
func SetSkipLevel(level int) {
	skipLevel = level
}
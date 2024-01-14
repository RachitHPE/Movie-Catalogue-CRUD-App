package log

import (
	"context"
)

// Logger defines interface for logging
type Logger interface {
	Info(ctx context.Context, args ...interface{})
	Infof(ctx context.Context, format string, args ...interface{})

	Warn(ctx context.Context, args ...interface{})
	Warnf(ctx context.Context, format string, args ...interface{})

	Error(ctx context.Context, args ...interface{})
	Errorf(ctx context.Context, format string, args ...interface{})

	Debug(ctx context.Context, args ...interface{})
	Debugf(ctx context.Context, format string, args ...interface{})

	Audit(ctx context.Context, args ...interface{})
	Auditf(ctx context.Context, format string, args ...interface{})
}
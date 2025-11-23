package logger

import (
	"context"

	"go.uber.org/zap"
)

type ctxKey struct{}

func ToContext(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, l)
}

func FromContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return zap.L()
	}

	if l, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		return l
	}

	return zap.L()
}

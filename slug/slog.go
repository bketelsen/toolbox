package slug

import (
	"context"
	"log/slog"
)

type ctxKey string

const ctxSlogKey ctxKey = "slog"

func CtxWithSlog(ctx context.Context, slog *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxSlogKey, slog)
}

func CtxSlog(ctx context.Context) (logger *slog.Logger, ok bool) {
	logger, ok = ctx.Value(ctxSlogKey).(*slog.Logger)
	return logger, ok
}

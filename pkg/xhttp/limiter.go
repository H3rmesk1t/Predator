package xhttp

import "context"

const ctxQPSLimiterKey = "QPSLimiter"

type QPSLimiter interface {
	Wait(context.Context, string) error
}

type WaitFunc func(context.Context, string) error

func (f WaitFunc) Wait(ctx context.Context, key string) error {
	return f(ctx, key)
}

func DumbWait(_ context.Context, _ string) error {
	return nil
}

func ExtractQPSLimiter(ctx context.Context) QPSLimiter {
	if value := ctx.Value(ctxQPSLimiterKey); value != nil {
		if limiter, ok := value.(QPSLimiter); ok {
			return limiter
		}
	}
	return WaitFunc(DumbWait)
}

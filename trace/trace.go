package trace

import (
	"context"
	"github.com/alioth-center/infrastructure/utils/generate"
	"github.com/alioth-center/infrastructure/utils/values"
	"time"
)

func Background() context.Context {
	return background
}

func Trace(ctx context.Context) context.Context {
	if _, ok := GetTrace(ctx); ok {
		return ctx
	}

	return context.WithValue(ctx, BasicType, &Basic{
		TraceID:  generate.TraceID(),
		TracedAt: time.Now(),
		Instance: instance,
		Service:  service,
	})
}

func GetTrace(ctx context.Context) (info *Basic, exist bool) {
	return getFromContext[*Basic](ctx, BasicType)
}

func GetRequestInfo(ctx context.Context) (info *Request, exist bool) {
	return getFromContext[*Request](ctx, RequestType)
}

func getFromContext[T any](ctx context.Context, t Type) (info T, exist bool) {
	ptr := ctx.Value(t)
	if ptr == nil {
		return values.Nil[T](), false
	}

	content, ok := ptr.(T)
	if !ok {
		return values.Nil[T](), false
	}

	return content, true
}

func SetBackground(i, s string) {
	instance, service = i, s
}

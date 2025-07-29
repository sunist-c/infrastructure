package aslog

import (
	"context"
	"fmt"
	"github.com/alioth-center/infrastructure/trace"
	"strings"
	"time"
)

type field struct {
	content *LogContent
}

func NewField(ctx context.Context) LogField {
	content := LogContent{}
	if ctx == nil {
		panic("nil context")
	}

	// attach basic trace info
	if pb := ctx.Value(trace.BasicType); pb == nil {
		panic("non-traced context")
	} else if cb, ok := pb.(*trace.Basic); !ok || cb == nil {
		panic("non-traced context")
	} else {
		content.Trace = cb.TraceID
		content.Instance = cb.Instance
		content.Timestamp = time.Now().Format(time.RFC3339)
		content.TimeCost = time.Since(cb.TracedAt).String()
	}

	// initialize log content structure
	content.LogPoint = trace.Caller(1)
	content.Labels = map[string]string{}

	// attach request info if exist
	if pr := ctx.Value(trace.RequestType); pr == nil {
		content.RequestInfo = nil
	} else if cr, ok := pr.(*trace.Request); !ok || cr == nil {
		content.RequestInfo = nil
	} else {
		content.RequestInfo = &RequestContent{
			Method:   cr.Method,
			Path:     cr.Path,
			Host:     cr.Host,
			ClientIP: cr.ClientIP,
		}
	}

	return &field{content: &content}
}

func (f *field) Messagef(format string, args ...any) LogField {
	if len(args) == 0 {
		f.content.Message = format
		return f
	}

	f.content.Message = fmt.Sprintf(format, args...)
	return f
}

func (f *field) Data(data any) LogField {
	f.content.Content = data
	return f
}

func (f *field) Labels(key string, values ...string) LogField {
	switch len(values) {
	case 0:
		return f
	case 1:
		f.content.Labels[key] = values[0]
	default:
		f.content.Labels[key] = strings.Join(values, ",")
	}

	return f
}

func (f *field) Service(service string) LogField {
	f.content.Service = service
	return f
}

func (f *field) Format() *LogContent {
	return f.content
}

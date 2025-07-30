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
	basic, existBasic := trace.GetTrace(ctx)
	if !existBasic {
		panic("non-traced context")
	}

	// initialize log content structure
	content.Trace = basic.TraceID
	content.Instance = basic.Instance
	content.Service = basic.Service
	content.Timestamp = time.Now().Format(time.RFC3339)
	content.TimeCost = time.Since(basic.TracedAt).String()
	content.LogPoint = trace.Caller(0)
	content.Labels = map[string]string{}

	// attach request info if exist
	if request, existRequest := trace.GetRequestInfo(ctx); existRequest {
		content.RequestInfo = &RequestContent{
			Method:   request.Method,
			Path:     request.Path,
			Host:     request.Host,
			ClientIP: request.ClientIP,
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

func (f *field) Format() *LogContent {
	return f.content
}

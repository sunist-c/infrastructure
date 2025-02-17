package logger

import (
	"context"
	"fmt"
	"github.com/alioth-center/infrastructure/trace"
)

func NewFields(ctx ...context.Context) Fields {
	f := &field{
		tid:     "",
		message: "",
		fields:  make(map[string]any),
		data:    nil,
	}
	if len(ctx) > 0 {
		return f.Context(ctx[0])
	}

	return f
}

type Fields interface {
	export(entry *entry)
	Context(ctx context.Context) Fields
	Message(message string, args ...any) Fields
	Field(key string, value any) Fields
	Data(data any) Fields
}

type field struct {
	tid     string
	message string
	fields  map[string]any
	data    any
}

func (f *field) Context(ctx context.Context) Fields {
	f.tid = trace.GetTid(ctx)
	return f
}

func (f *field) Message(message string, args ...any) Fields {
	if len(args) == 0 {
		f.message = message
		return f
	}

	f.message = fmt.Sprintf(message, args...)
	return f
}

func (f *field) Field(key string, value any) Fields {
	f.fields[key] = value
	return f
}

func (f *field) Data(data any) Fields {
	f.data = data
	return f
}

func (f *field) export(entry *entry) {
	entry.TraceID = f.tid
	entry.Message = f.message
	entry.Data = f.data
	entry.Extra = f.fields
}

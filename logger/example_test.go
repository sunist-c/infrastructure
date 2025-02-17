package logger_test

import (
	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/trace"
)

func ExampleNewLogger() {
	ctx := trace.NewContext()
	payload := map[string]any{"foo": "bar"}
	myLogger := logger.NewLogger(logger.LevelTrace)
	myLogger.Info(logger.NewFields(ctx).Data(payload).Message("hello"))

	// Output:
	// {"file":"example_test.go","level":"info","service":"example_test","trace_id":"trace_id","message":"hello","call_time":"call_time","data":{"foo":"bar"}}
}

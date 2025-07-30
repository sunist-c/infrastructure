package logger_test

import (
	"github.com/alioth-center/infrastructure/grace"
	"github.com/alioth-center/infrastructure/logger"
	aslog "github.com/alioth-center/infrastructure/logger/service"
	"github.com/alioth-center/infrastructure/trace"
	"sync"
	"testing"
	"time"
)

var (
	so sync.Once
	eo sync.Once
)

//
//func TestLogger(t *testing.T) {
//	rotator := logger.NewNoRotator("./log")
//	writer := logger.NewFileWriter(rotator)
//	grace.ServeGraceful()
//
//	slog := aslog.NewServiceLogger(logger.LevelDebug, writer)
//
//	trace.SetBackground("test-instance", "unittest")
//	ctx := trace.Trace(context.Background())
//	slog.Info(aslog.NewField(ctx).Messagef("hello! %s", "aslog"))
//
//	grace.CloseGraceful()
//}

var (
	slog aslog.ServiceLogger
	once sync.Once
)

func init() {
	trace.SetBackground("test-instance", "unittest")
	rotator := logger.NewNoRotator("/dev/null")
	writer := logger.NewFileWriter(rotator)
	grace.ServeGraceful()
	slog = aslog.NewServiceLogger(logger.LevelInfo, writer)
}

func BenchmarkLogger(b *testing.B) {
	// $ go test -bench=. -benchtime=20s
	// goos: linux
	// goarch: amd64
	// pkg: github.com/alioth-center/infrastructure/logger
	// cpu: 12th Gen Intel(R) Core(TM) i3-12100F
	// BenchmarkLogger-8       15118152              1606 ns/op
	// PASS
	// ok      github.com/alioth-center/infrastructure/logger  25.654s
	ctx := trace.Background()
	for i := 0; i < b.N; i++ {
		slog.Info(aslog.NewField(ctx).Messagef("hello! %s", "aslog"))
	}

	once.Do(func() {
		time.Sleep(time.Second)
	})
}

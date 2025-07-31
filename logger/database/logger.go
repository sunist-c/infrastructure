package adlog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/trace"
	glogger "gorm.io/gorm/logger"
	"time"
)

type databaseLogger struct {
	logLevel  logger.Level
	writer    logger.LogWriter
	enableSQL bool
}

func (dl *databaseLogger) LogMode(level glogger.LogLevel) glogger.Interface {
	switch level {
	case glogger.Silent:
		dl.logLevel = logger.LevelFatal
	case glogger.Error:
		dl.logLevel = logger.LevelError
	case glogger.Warn:
		dl.logLevel = logger.LevelWarn
	case glogger.Info:
		dl.logLevel = logger.LevelInfo
	}

	return dl
}

func (dl *databaseLogger) Info(ctx context.Context, s string, i ...interface{}) {
	dl.logMessage(ctx, logger.LevelInfo, s, i...)
}

func (dl *databaseLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	dl.logMessage(ctx, logger.LevelWarn, s, i...)
}

func (dl *databaseLogger) Error(ctx context.Context, s string, i ...interface{}) {
	dl.logMessage(ctx, logger.LevelError, s, i...)
}

func (dl *databaseLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, affected := fc()
	dl.logSQL(ctx, begin, sql, affected, err)
}

func (dl *databaseLogger) logMessage(ctx context.Context, level logger.Level, format string, args ...any) {
	if level < dl.logLevel {
		return
	}

	if ctx == nil {
		panic("nil context")
	}

	content := LogContent{}
	if basic, success := trace.GetTrace(ctx); success {
		content.Trace, content.Service, content.Instance = basic.TraceID, basic.Service, basic.Instance
	}

	content.Level, content.Message = level.String(), fmt.Sprintf(format, args...)
	buffer := &bytes.Buffer{}
	_ = json.NewEncoder(buffer).Encode(&content)
	dl.writer.WriteRaw(buffer)
}

func (dl *databaseLogger) logSQL(ctx context.Context, started time.Time, sql string, affected int64, err error) {
	if !dl.enableSQL {
		return
	}

	if ctx == nil {
		panic("nil context")
	}

	content := LogContent{}
	if basic, success := trace.GetTrace(ctx); success {
		content.Trace, content.Service, content.Instance = basic.TraceID, basic.Service, basic.Instance
	}

	content.Level, content.Message, content.SQL, content.ExecTime, content.RowsAffected =
		logger.LevelInfo.String(), "sql executed", sql, time.Since(started).String(), affected
	if err != nil {
		content.ExecError = err.Error()
	}
}

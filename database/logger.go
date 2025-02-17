package database

import (
	"context"
	"time"

	"github.com/alioth-center/infrastructure/logger"
	glog "gorm.io/gorm/logger"
)

// DBLogger used alioth-center infrastructure logger to log sql.
// Do not call LogMode when database is in use because it is not thread safe.
type DBLogger struct {
	log logger.Logger
}

func NewDBLogger(log logger.Logger) *DBLogger {
	return &DBLogger{log: log}
}

func (dl *DBLogger) LogMode(level glog.LogLevel) glog.Interface {
	switch level {
	case glog.Info:
		dl.log = logger.NewLogger(logger.LevelInfo)
	case glog.Silent:
		dl.log = logger.NewLogger(logger.LevelPanic)
	case glog.Error:
		dl.log = logger.NewLogger(logger.LevelError)
	case glog.Warn:
		dl.log = logger.NewLogger(logger.LevelWarn)
	}

	return dl
}

func (dl *DBLogger) Info(ctx context.Context, s string, i ...interface{}) {
	dl.log.Info(logger.NewFields(ctx).Message(s, i...))
}

func (dl *DBLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	dl.log.Warn(logger.NewFields(ctx).Message(s, i...))
}

func (dl *DBLogger) Error(ctx context.Context, s string, i ...interface{}) {
	dl.log.Error(logger.NewFields(ctx).Message(s, i...))
}

func (dl *DBLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, rows := fc()
	if err != nil {
		logMessage := map[string]any{"sql": sql, "error": err.Error(), "rows": rows, "duration": time.Since(begin)}
		dl.log.Error(logger.NewFields(ctx).Message("tracing sql with error").Data(logMessage))
		return
	}

	logMessage := map[string]any{"sql": sql, "rows": rows, "duration": time.Since(begin)}
	dl.log.Trace(logger.NewFields(ctx).Message("tracing sql").Data(logMessage))
}

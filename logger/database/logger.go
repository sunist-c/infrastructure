package database

import (
	"context"
	"github.com/alioth-center/infrastructure/logger"
	glogger "gorm.io/gorm/logger"
	"time"
)

type databaseLogger struct {
	logLevel logger.Level
	writer   logger.LogWriter
}

func (S *databaseLogger) LogMode(level glogger.LogLevel) glogger.Interface {
	switch level {
	case glogger.Silent:
		S.logLevel = logger.LevelFatal
	case glogger.Error:
		S.logLevel = logger.LevelError
	case glogger.Warn:
		S.logLevel = logger.LevelWarn
	case glogger.Info:
		S.logLevel = logger.LevelInfo
	}

	return S
}

func (S *databaseLogger) Info(ctx context.Context, s string, i ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (S *databaseLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (S *databaseLogger) Error(ctx context.Context, s string, i ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (S *databaseLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	//TODO implement me
	panic("implement me")
}

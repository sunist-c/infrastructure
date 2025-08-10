package aslog

import (
	"bytes"
	"encoding/json"

	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/trace"
)

type serviceLogger struct {
	logLevel logger.Level
	writer   logger.LogWriter
}

func NewServiceLogger(logLevel logger.Level, writer logger.LogWriter) ServiceLogger {
	slog := &serviceLogger{
		logLevel: logLevel,
		writer:   writer,
	}
	slog.Debug(NewField(trace.Background()).Messagef("service logger initialized"))

	return slog
}

func (s *serviceLogger) Debug(field LogField) {
	s.log(logger.LevelDebug, field)
}

func (s *serviceLogger) Info(field LogField) {
	s.log(logger.LevelInfo, field)
}

func (s *serviceLogger) Warn(field LogField) {
	s.log(logger.LevelWarn, field)
}

func (s *serviceLogger) Error(field LogField) {
	s.log(logger.LevelError, field)
}

func (s *serviceLogger) Fatal(field LogField) {
	s.log(logger.LevelFatal, field)
}

func (s *serviceLogger) log(level logger.Level, field LogField) {
	if level < s.logLevel {
		return
	}

	content := field.Format()
	content.Level = level.String()

	buffer := &bytes.Buffer{}
	_ = json.NewEncoder(buffer).Encode(content)

	s.writer.WriteRaw(buffer)
}

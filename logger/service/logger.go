package aslog

import (
	"bytes"
	"encoding/json"
	"github.com/alioth-center/infrastructure/logger"
)

type serviceLogger struct {
	writer logger.LogWriter
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
	content := field.Format()
	content.Level = level.String()

	buffer := &bytes.Buffer{}
	_ = json.NewEncoder(buffer).Encode(content)

	s.writer.WriteRaw(buffer)
}

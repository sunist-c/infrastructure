package logger

import (
	"bytes"
	"github.com/alioth-center/infrastructure/grace"
)

type LogWriter interface {
	WriteRaw(log *bytes.Buffer)
}

type GracefulLogWriter interface {
	grace.Graceful
	LogWriter
}

type LogRotator interface {
	CalculateRotationKey() (key string)
}

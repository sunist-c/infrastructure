package logger

import "bytes"

type LogWriter interface {
	WriteRaw(log *bytes.Buffer)
}

type LogRotator interface {
	CalculateRotationKey() (key string)
}

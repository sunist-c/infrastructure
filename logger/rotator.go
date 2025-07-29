package logger

import (
	"fmt"
	"path/filepath"
	"sync/atomic"
	"time"
)

type noRotator struct {
	basePath string
}

func NewNoRotator(basePath string) LogRotator {
	return &noRotator{
		basePath: basePath,
	}
}

func (r *noRotator) CalculateRotationKey() (key string) {
	return r.basePath
}

type timeRotator struct {
	basePath   string
	fileName   string
	timeFormat string
}

func NewTimeRotator(basePath, fileName, timeFormat string) LogRotator {
	return &timeRotator{
		basePath:   basePath,
		fileName:   fileName,
		timeFormat: timeFormat,
	}
}

func (r *timeRotator) CalculateRotationKey() (key string) {
	return filepath.Join(r.basePath, fmt.Sprintf(r.fileName, time.Now().Format(r.timeFormat)))
}

type fileRotator struct {
	basePath  string
	fileName  string
	limitLine int64
	counter   atomic.Int64
}

func NewFileRotator(basePath, fileName string, limitLine int64) LogRotator {
	if limitLine <= 0 {
		limitLine = 1000
	}

	return &fileRotator{
		basePath:  basePath,
		fileName:  fileName,
		limitLine: limitLine,
		counter:   atomic.Int64{},
	}
}

func (r *fileRotator) CalculateRotationKey() (key string) {
	return filepath.Join(r.basePath, fmt.Sprintf(r.fileName, r.counter.Add(1)%r.limitLine))
}

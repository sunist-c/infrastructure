package logger

import (
	"fmt"
	"path/filepath"
	"sync/atomic"
	"time"
)

type TimeRotator struct {
	basePath   string
	fileName   string
	timeFormat string
}

func NewTimeRotator(basePath, fileName, timeFormat string) LogRotator {
	return &TimeRotator{
		basePath:   basePath,
		fileName:   fileName,
		timeFormat: timeFormat,
	}
}

func (r *TimeRotator) CalculateRotationKey() (key string) {
	return filepath.Join(r.basePath, fmt.Sprintf(r.fileName, time.Now().Format(r.timeFormat)))
}

type FileRotator struct {
	basePath  string
	fileName  string
	limitLine int64
	counter   atomic.Int64
}

func NewFileRotator(basePath, fileName string, limitLine int64) LogRotator {
	if limitLine <= 0 {
		limitLine = 1000
	}

	return &FileRotator{
		basePath:  basePath,
		fileName:  fileName,
		limitLine: limitLine,
		counter:   atomic.Int64{},
	}
}

func (r *FileRotator) CalculateRotationKey() (key string) {
	return filepath.Join(r.basePath, fmt.Sprintf(r.fileName, r.counter.Add(1)%r.limitLine))
}

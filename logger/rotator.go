package logger

import (
	"fmt"
	"github.com/alioth-center/infrastructure/env"
	"github.com/alioth-center/infrastructure/utils/values"
	"time"
)

var (
	RotatorType     string
	TimeRotatorOpts string
	LineRotatorOpts int64
	SizeRotatorOpts int64
	LogPath         string

	NoRotator   Rotator[string]    = noRotator{}
	TimeRotator Rotator[time.Time] = timeRotator{}
	LineRotator Rotator[int64]     = lineRotator{}
	SizeRotator Rotator[int64]     = sizeRotator{}
)

const (
	RotatorTypeNone = "NONE"
	RotatorTypeTime = "TIME"
	RotatorTypeLine = "LINE"
	RotatorTypeSize = "SIZE"
)

func init() {
	LogPath = env.GetEnv(env.AliothFrameworkLogPathKey)
	RotatorType = env.GetEnv(env.AliothFrameworkLogRotateTypeKey)

	TimeRotatorOpts = env.GetEnv(env.AliothFrameworkLogRotateTimeKey)
	LineRotatorOpts = values.StringToInt(env.GetEnv(env.AliothFrameworkLogRotateLineKey), int64(1000))
	SizeRotatorOpts = env.ParseEnv(env.AliothFrameworkLogRotateSizeKey, env.SizeParser)
}

type Rotator[T any] interface {
	Rotate(source T) (filename string)
}

type noRotator struct{}

func (r noRotator) Rotate(source string) string {
	return source
}

type timeRotator struct{}

func (r timeRotator) Rotate(source time.Time) string {
	return source.Format(TimeRotatorOpts)
}

type sizeRotator struct{}

func (r sizeRotator) Rotate(source int64) string {
	return fmt.Sprintf("%d", (source%SizeRotatorOpts)+1)
}

type lineRotator struct{}

func (r lineRotator) Rotate(source int64) string {
	return fmt.Sprintf("%d", (source%LineRotatorOpts)+1)
}

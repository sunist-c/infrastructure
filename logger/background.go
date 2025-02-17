package logger

import (
	"github.com/alioth-center/infrastructure/env"
	"github.com/alioth-center/infrastructure/trace"
	"github.com/alioth-center/infrastructure/utils/timezone"
	"strings"
	"time"
)

type background struct {
	file     string
	level    Level
	service  string
	callTime time.Time
}

// simplifyFilename Simplify the filename with cut the prefix of the filename and the version of the package.
//
// Use the environment variable "AF_SIMPLIFY_PATH" and "AF_PKG_PATH" to control the simplify behavior.
//
// example:
//
//	base   = /go/pkg/mod/
//	file   = /go/pkg/mod/github.com/alioth-center/infrastructure@v1.2.20/thirdparty/openai/client.go:185
//	result =             github.com/alioth-center/infrastructure        /thirdparty/openai/client.go:185
func (b *background) simplifyFilename(file string) string {
	if enableSimplify := env.GetEnv(env.AliothFrameworkSimplifyPathKey); strings.ToUpper(enableSimplify) != "TRUE" {
		return file
	}

	base := env.GetEnv(env.AliothFrameworkPackagePathKey)
	result := strings.TrimPrefix(file, base)

	if atIndex := strings.Index(result, "@"); atIndex != -1 {
		if slashIndex := strings.Index(result[atIndex:], "/"); slashIndex != -1 {
			result = result[:atIndex] + result[atIndex+slashIndex:]
		}
	}

	return result
}

func (b *background) export(entry *entry) {
	entry.File = b.simplifyFilename(b.file)
	entry.Level = string(b.level)
	entry.Service = b.service
	entry.CallTime = b.callTime.Format(timeFormat)
}

func newBackground(level Level) background {
	return background{
		level:    level,
		file:     trace.Caller(2),
		service:  env.GetEnv(env.AliothFrameworkServiceNameKey),
		callTime: timezone.NowInLocalTime(),
	}
}

package env

import (
	"fmt"
	"github.com/alioth-center/infrastructure/utils/values"
	"strings"
	"time"
)

func GetEnv(key string) string {
	return envMap[key]
}

func ParseEnv[T any](key string, parser func(value string) T) T {
	return parser(GetEnv(key))
}

func SizeParser(value string) (result int64) {
	sizeNum, sizeUnit := 0, ""
	scanNum, scanErr := fmt.Sscanf(value, "%d%s", &sizeNum, &sizeUnit)
	if scanErr == nil && scanNum == 2 {
		switch strings.ToUpper(sizeUnit) {
		case "B":
			result = int64(sizeNum)
		case "K", "KB", "KIB":
			result = int64(sizeNum) * 1024
		case "M", "MB", "MIB":
			result = int64(sizeNum) * 1024 * 1024
		case "G", "GB", "GIB":
			result = int64(sizeNum) * 1024 * 1024 * 1024
		}
	}

	return result
}

func TimeParser(value string) (result time.Duration) {
	timeNum, timeUnit := 0, ""
	scanNum, scanErr := fmt.Sscanf(value, "%d%s", &timeNum, &timeUnit)
	if scanErr == nil && scanNum == 2 {
		switch strings.ToUpper(timeUnit) {
		case "MS":
			result = time.Duration(timeNum) * time.Millisecond
		case "S":
			result = time.Duration(timeNum) * time.Second
		case "M":
			result = time.Duration(timeNum) * time.Minute
		case "H":
			result = time.Duration(timeNum) * time.Hour
		}
	}

	return result
}

func IntParser(value string) (result int64) {
	return values.StringToInt(value, result)
}

func BoolParser(value string) (result bool) {
	return values.StringToBool(value, result)
}

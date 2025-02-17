package logger

type Level string

const (
	LevelTrace  Level = "trace"
	LevelDebug  Level = "debug"
	LevelInfo   Level = "info"
	LevelNotice Level = "notice"
	LevelWarn   Level = "warn"
	LevelError  Level = "error"
	LevelPanic  Level = "panic"
)

func (l Level) shouldLog(level Level) bool {
	// logger.Level.shouldLog(input.Level)
	return LevelValueMap[l] <= LevelValueMap[level]
}

var (
	timeFormat    = "2006.01.02-15:04:05.000Z07:00"
	LevelValueMap = map[Level]int{
		LevelTrace:  0,
		LevelDebug:  1,
		LevelInfo:   2,
		LevelNotice: 3,
		LevelWarn:   4,
		LevelError:  5,
		LevelPanic:  6,
	}
)

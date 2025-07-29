package logger

type Level string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
	LevelFatal Level = "fatal"
	LevelPanic Level = "panic"
)

var lvlMap = map[Level]int{LevelDebug: 0, LevelInfo: 1, LevelWarn: 2, LevelError: 3, LevelFatal: 4, LevelPanic: 5}

func (l Level) String() string {
	return string(l)
}

func (l Level) ShouldLog(level Level) bool {
	// logger.Level.shouldLog(input.Level)
	return lvlMap[l] <= lvlMap[level]
}

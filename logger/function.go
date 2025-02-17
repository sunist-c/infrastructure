package logger

var (
	defaultLogger = &logger{level: LevelInfo, writer: DefaultFileWriter}

	Default Logger = defaultLogger
)

func SetLevel(level Level) {
	defaultLogger.level = level
}

func Trace(fields Fields) {
	defaultLogger.Trace(fields)
}

func Debug(fields Fields) {
	defaultLogger.Debug(fields)
}

func Info(fields Fields) {
	defaultLogger.Info(fields)
}

func Notice(fields Fields) {
	defaultLogger.Notice(fields)
}

func Warn(fields Fields) {
	defaultLogger.Warn(fields)
}

func Error(fields Fields) {
	defaultLogger.Error(fields)
}

func Panic(fields Fields) {
	defaultLogger.Panic(fields)
}

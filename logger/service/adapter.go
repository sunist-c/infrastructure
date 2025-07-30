package aslog

type LogField interface {
	Messagef(format string, args ...any) LogField
	Data(data any) LogField
	Labels(key string, values ...string) LogField
	Format() *LogContent
}

type ServiceLogger interface {
	Debug(field LogField)
	Info(field LogField)
	Warn(field LogField)
	Error(field LogField)
	Fatal(field LogField)
}

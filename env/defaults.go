package env

var envMap = map[string]string{
	AliothFrameworkRuntimeModeKey:      "DEBUG",
	AliothFrameworkBannerPrintKey:      "TRUE",
	AliothFrameworkServiceTimezoneKey:  "Asia/Shanghai",
	AliothFrameworkLogRotateTypeKey:    "TIME",
	AliothFrameworkLogRotateTimeKey:    "20060102_15",
	AliothFrameworkLogPathKey:          "./logs",
	AliothFrameworkServiceNameKey:      "new-alioth-app",
	AliothFrameworkExitWaitDurationKey: "5s",
	AliothFrameworkTimeFormatKey:       "2006.01.02-15:04:05.000Z07:00",
}

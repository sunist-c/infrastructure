package env

const (
	AliothFrameworkRuntimeModeKey      = "AF_RUN_MODE"        // AliothFrameworkRuntimeModeKey set to "DEBUG" to enable debug mode, print debug information
	AliothFrameworkPackagePathKey      = "AF_PKG_PATH"        // AliothFrameworkPackagePathKey set to the package path to simplify the filename
	AliothFrameworkSimplifyPathKey     = "AF_SIMPLIFY_PATH"   // AliothFrameworkSimplifyPathKey set to "TRUE" to simplify the filename
	AliothFrameworkBannerPrintKey      = "AF_BANNER_PRINT"    // AliothFrameworkBannerPrintKey set to "TRUE" to print the banner
	AliothFrameworkServiceNameKey      = "AF_SERVICE_NAME"    // AliothFrameworkServiceNameKey set to the service name for logging and tracing
	AliothFrameworkServiceTimezoneKey  = "AF_SERVICE_TZ"      // AliothFrameworkServiceTimezoneKey set to the application timezone
	AliothFrameworkLogPathKey          = "AF_LOG_PATH"        // AliothFrameworkLogPathKey set to the log path
	AliothFrameworkLogRotateTypeKey    = "AF_LOG_ROTATE_TYPE" // AliothFrameworkLogRotateTypeKey set to the log rotate type
	AliothFrameworkLogRotateTimeKey    = "AF_LOG_ROTATE_TIME" // AliothFrameworkLogRotateTimeKey set to the log rotate time
	AliothFrameworkLogRotateSizeKey    = "AF_LOG_ROTATE_SIZE" // AliothFrameworkLogRotateSizeKey set to the log rotate size
	AliothFrameworkLogRotateLineKey    = "AF_LOG_ROTATE_LINE" // AliothFrameworkLogRotateLineKey set to the log rotate line
	AliothFrameworkExitWaitDurationKey = "AF_EXIT_DURATION"   // AliothFrameworkExitWaitDurationKey set to the exit wait duration
	AliothFrameworkTimeFormatKey       = "AF_TIME_FORMAT"     // AliothFrameworkTimeFormatKey set to the time format
)

var Keys = [13]string{
	AliothFrameworkRuntimeModeKey,
	AliothFrameworkPackagePathKey,
	AliothFrameworkSimplifyPathKey,
	AliothFrameworkBannerPrintKey,
	AliothFrameworkServiceNameKey,
	AliothFrameworkServiceTimezoneKey,
	AliothFrameworkLogPathKey,
	AliothFrameworkLogRotateTypeKey,
	AliothFrameworkLogRotateTimeKey,
	AliothFrameworkLogRotateSizeKey,
	AliothFrameworkLogRotateLineKey,
	AliothFrameworkExitWaitDurationKey,
	AliothFrameworkTimeFormatKey,
}

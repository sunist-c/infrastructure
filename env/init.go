package env

import (
	"os"
)

func init() {
	loadLocalEnvironment()
	checkLogPath()
	printBanner()
}

func loadLocalEnvironment() {
	for _, key := range Keys {
		if envValue := os.Getenv(key); envValue != "" {
			envMap[key] = envValue
		}
	}
}

func checkLogPath() {
	if envMap[AliothFrameworkLogPathKey] != "" {
		// check path exist, if not, create it
		if _, err := os.Stat(envMap[AliothFrameworkLogPathKey]); os.IsNotExist(err) {
			createErr := os.MkdirAll(envMap[AliothFrameworkLogPathKey], os.ModePerm)
			if createErr != nil {
				panic(createErr)
			}
		}

		// check path is directory, if not, panic
		if fileInfo, err := os.Stat(envMap[AliothFrameworkLogPathKey]); err != nil || !fileInfo.IsDir() {
			panic("Log path is not a directory")
		}
	}
}

package exit_test

import (
	"github.com/alioth-center/infrastructure/exit"
	"github.com/alioth-center/infrastructure/utils/console"
)

type ifrace interface {
	Save()
	Flush()
	GracefulClose()
	Shutdown()
}

var files, cache, database, api ifrace

func ExampleRegisterExitAction() {
	exitFunc := func(message console.Messages) {
		files.Save()
		cache.Flush()
		database.GracefulClose()
		api.Shutdown()

		message.Message("All resources released")
	}

	exit.RegisterExitAction("Graceful Shutdown", exitFunc)

	// Output:
	// [Alioth Framework Exit Actions]: 2025.02.17-13:36:42.634+08:00
	// Exit Signal: interrupt
	// 	[Exit Action Result]: 2025.02.17-13:36:42.634+08:00
	//	Action: Graceful Shutdown
	//	Status: Success
	//	Function: /path/to/your/project/main.go:42
	//	Message: All resources released
}

package exit

import (
	"github.com/alioth-center/infrastructure/utils/concurrency"
	"github.com/alioth-center/infrastructure/utils/console"
)

var Actions = concurrency.NewSlice[Action]()

// Action is a struct that defines the exit action to be executed when the program is about to exit.
type Action struct {
	Name string
	Func func(message console.Messages)
}

// RegisterExitAction registers an exit action to be executed
// when the program is about to exit.
//
// Parameters:
//
//		actionName(string): the name of the action to be registered.
//
//		callback (func(message console.Messages)): the function to be executed when the program is about to exit.
//	                                             You can use console.Messages to print messages to the console
//	                                             when the action is executed.
func RegisterExitAction(actionName string, callback func(message console.Messages)) {
	Actions.Append(Action{Name: actionName, Func: callback})
}

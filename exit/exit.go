package exit

import (
	"github.com/alioth-center/infrastructure/env"
	"github.com/alioth-center/infrastructure/trace"
	"github.com/alioth-center/infrastructure/utils/concurrency"
	"github.com/alioth-center/infrastructure/utils/console"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	blockedChannel = make(chan struct{}, 1)
	signalChannel  = make(chan os.Signal, 1)
)

func init() {
	signal.Notify(signalChannel, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go awaitExitActions()
}

// awaitExitActions waits for a termination signal and executes all registered exit functions.
func awaitExitActions() {
	exitSignal := <-signalChannel
	task := &sync.WaitGroup{}
	task.Add(Actions.Length())
	output := console.NewBlock().
		Title("Alioth Framework Exit Actions").
		Message(console.NewMessage().Message("Exit Signal: %s", exitSignal.String()))

	go func() {
		waitDuration := env.ParseEnv(env.AliothFrameworkExitWaitDurationKey, env.TimeParser)
		<-time.After(waitDuration)
		output = output.SubBlock(console.NewBlock().Title("Exit Action Timeout").Message(
			console.NewMessage().Message("Timeout: %s", waitDuration.String()),
		))

		console.Print(output)
		os.Exit(1)
	}()

	for _, action := range Actions.Items() {
		go func(action Action) {
			defer task.Done()

			actionMessage := console.NewMessage().Message("Message: ")
			defer func() {
				if err := concurrency.RecoverErr(recover()); err != nil {
					output = output.SubBlock(console.NewBlock().Title("Exit Action Result").Message(
						console.NewMessage().Message("Action: %s", action.Name),
						console.NewMessage().Message("Status: ").Red("Failed"),
						console.NewMessage().Message("Function: %s", trace.FunctionLocation(action.Func)),
						console.NewMessage().Message("Error: ").Red(err.Error()),
						actionMessage,
					))
				}
			}()

			action.Func(actionMessage)
			output = output.SubBlock(console.NewBlock().Title("Exit Action Result").Message(
				console.NewMessage().Message("Action: %s", action.Name),
				console.NewMessage().Message("Status: ").Green("Success"),
				console.NewMessage().Message("Function: %s", trace.FunctionLocation(action.Func)),
				actionMessage,
			))
		}(action)
	}

	task.Wait()
	console.Print(output)
	blockedChannel <- struct{}{}
}

// BlockedUntilTerminate blocks the current goroutine until a termination signal is received
// and all exit functions have completed. This function should be called to ensure the program
// does not exit immediately and waits for a proper shutdown sequence.
func BlockedUntilTerminate() {
	sync.OnceFunc(func() {
		<-blockedChannel
		os.Exit(0)
	})()
}

// Exit sends a termination signal to the signal channel, initiating the shutdown process.
func Exit() {
	signalChannel <- syscall.SIGTERM
}

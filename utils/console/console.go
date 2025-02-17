package console

import (
	"sync"
)

type Control string

const (
	Reset        Control = "\033[0m"  // Reset reset all format
	FormatBold   Control = "\033[1m"  // FormatBold start bold format
	FormatDim    Control = "\033[2m"  // FormatDim start dim format
	FormatUnder  Control = "\033[4m"  // FormatUnder start underline format
	FormatBlink  Control = "\033[5m"  // FormatBlink start blink format
	FormatRev    Control = "\033[7m"  // FormatRev start reverse format
	ColorRed     Control = "\033[31m" // ColorRed start red color
	ColorGreen   Control = "\033[32m" // ColorGreen start green color
	ColorYellow  Control = "\033[33m" // ColorYellow start yellow color
	ColorBlue    Control = "\033[34m" // ColorBlue start blue color
	ColorMagenta Control = "\033[35m" // ColorMagenta start magenta color
	ColorCyan    Control = "\033[36m" // ColorCyan start cyan color
	ColorWhite   Control = "\033[37m" // ColorWhite start white color
)

var (
	syncMutex = sync.Mutex{}
)

func Print(block Block) {
	syncMutex.Lock()
	defer syncMutex.Unlock()
	block.print(0)
}

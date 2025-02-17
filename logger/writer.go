package logger

import (
	"fmt"
	"github.com/alioth-center/infrastructure/utils/console"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/alioth-center/infrastructure/utils/concurrency"
	"github.com/alioth-center/infrastructure/utils/timezone"
)

var DefaultFileWriter = NewFileWriter()

type FileWriter interface {
	WriteLog(content []byte)
}

type fileWriter struct {
	rotatorArgs atomic.Int64
	files       concurrency.Map[string, *os.File]
	alerts      concurrency.Map[string, *sync.Once]
	exitMtx     sync.RWMutex
}

// WriteLog writes the log content to the file.
// This function will write the log content to the file based on the rotator type.
// Any error occurred during the write operation will be written to the stdout.
//
// Parameters:
//
//	content ([]byte): The log content to be written to the file.
func (fw *fileWriter) WriteLog(content []byte) {
	file := fw.getLogFile(int64(len(content)))
	_, writeErr := file.Write(content)
	if writeErr != nil {
		_, _ = os.Stdout.Write(content)
	}
}

// getLogFile returns the file to write the log content.
// This function will return the non-nil file based on the rotator type.
func (fw *fileWriter) getLogFile(delta int64) *os.File {
	writeFile := ""
	switch strings.ToUpper(RotatorType) {
	case RotatorTypeTime:
		writeFile = TimeRotator.Rotate(timezone.NowInLocalTime())
	case RotatorTypeLine:
		writeFile = LineRotator.Rotate(fw.rotatorArgs.Add(1))
	case RotatorTypeSize:
		writeFile = SizeRotator.Rotate(fw.rotatorArgs.Add(delta))
	case RotatorTypeNone:
		writeFile = NoRotator.Rotate("logs")
	}

	if writeFile == "" {
		return os.Stdout
	}
	writeFile = fmt.Sprintf("%s.jsonl", filepath.Join(LogPath, writeFile))
	if file, ok := fw.files.Get(writeFile); ok {
		return file
	}

	file, err := os.OpenFile(writeFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o755)
	if err != nil {
		onceAlert := sync.Once{}
		onceAlert.Do(func() {
			console.Print(console.NewBlock().
				Title("Alioth Framework File Logger Open File Error").
				Message(
					console.NewMessage().Index(1).Message("Failed to open file: %s", writeFile),
					console.NewMessage().Index(2).Message("Error: ").Red(err.Error()),
				),
			)
		})
		fw.alerts.Set(writeFile, &onceAlert)
		return os.Stdout
	}
	fw.files.Set(writeFile, file)
	return file
}

func NewFileWriter() FileWriter {
	writer := &fileWriter{
		rotatorArgs: atomic.Int64{},
		files:       concurrency.NewMap[string, *os.File](),
		alerts:      concurrency.NewMap[string, *sync.Once](),
		exitMtx:     sync.RWMutex{},
	}

	return writer
}

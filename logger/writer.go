package logger

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	"github.com/alioth-center/infrastructure/grace"
)

type fileWriter struct {
	rotator     LogRotator
	keyCache    string
	writerCache *os.File
	cacheMutex  sync.RWMutex
	logBuffer   chan *bytes.Buffer
	closer      chan struct{}
	closed      atomic.Bool
	wait        sync.WaitGroup
}

func NewFileWriter(rotator LogRotator) LogWriter {
	writer := &fileWriter{
		rotator:    rotator,
		cacheMutex: sync.RWMutex{},
		logBuffer:  make(chan *bytes.Buffer, 4096),
		closer:     make(chan struct{}, 2),
		closed:     atomic.Bool{},
		wait:       sync.WaitGroup{},
	}
	writer.wait.Add(1)
	writer.closed.Store(false)
	grace.RegisterGraceful(writer)

	return writer
}

func NewGracefulFileWriter(rotator LogRotator) GracefulLogWriter {
	writer := &fileWriter{
		rotator:    rotator,
		cacheMutex: sync.RWMutex{},
		logBuffer:  make(chan *bytes.Buffer, 4096),
		closer:     make(chan struct{}, 2),
		closed:     atomic.Bool{},
		wait:       sync.WaitGroup{},
	}
	writer.wait.Add(1)
	writer.closed.Store(false)

	return writer
}

func (fw *fileWriter) WriteRaw(log *bytes.Buffer) {
	if !fw.closed.Load() {
		fw.logBuffer <- log

		return
	}

	_, _ = fmt.Fprintln(os.Stderr, log.String())
}

func (fw *fileWriter) writeLog(log *bytes.Buffer) (err error) {
	key := fw.rotator.CalculateRotationKey()
	fw.cacheMutex.RLock()
	if key == fw.keyCache {
		if fw.writerCache != nil {
			_, err = fw.writerCache.Write(log.Bytes())
			fw.cacheMutex.RUnlock()
			return err
		}

		fw.cacheMutex.RUnlock()
		return errors.New("nil writer cache")
	}
	fw.cacheMutex.RUnlock()

	fw.cacheMutex.Lock()
	if fw.writerCache != nil {
		_ = fw.writerCache.Close()
	}

	fw.writerCache, err = os.OpenFile(key, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o755)
	if err != nil {
		fw.cacheMutex.Unlock()
		return errors.New("open log file failed")
	}

	fw.keyCache = key
	_, err = fw.writerCache.Write(log.Bytes())
	if err != nil {
		fw.cacheMutex.Unlock()
		return err
	}

	fw.cacheMutex.Unlock()
	return nil
}

func (fw *fileWriter) GracefulClose() {
	fw.closed.Store(true)
	fw.closer <- struct{}{}
	fw.wait.Wait()

	fw.cacheMutex.Lock()
	defer fw.cacheMutex.Unlock()

	close(fw.logBuffer)
	for message := range fw.logBuffer {
		if fw.writerCache != nil {
			_, _ = fw.writerCache.Write(message.Bytes())
		}
	}

	if fw.writerCache != nil {
		_ = fw.writerCache.Close()
		fw.writerCache = nil
	}
}

func (fw *fileWriter) ListenAndServe() {
	for {
		select {
		case <-fw.closer:
			fw.wait.Done()
			return
		case log := <-fw.logBuffer:
			if writeErr := fw.writeLog(log); writeErr != nil {
				_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf(`{"log_error": %s, "raw_log": %s}`, writeErr.Error(), log.String()))
			}
		}
	}
}

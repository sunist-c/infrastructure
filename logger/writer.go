package logger

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"sync"
)

type fileWriter struct {
	rotator     LogRotator
	keyCache    string
	writerCache *os.File
	cacheMutex  sync.RWMutex
	logBuffer   chan *bytes.Buffer
	closer      chan struct{}
}

func (fw *fileWriter) WriteRaw(log *bytes.Buffer) {
	fw.logBuffer <- log
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
	fw.closer <- struct{}{}

	close(fw.logBuffer)
	for message := range fw.logBuffer {
		if fw.writerCache != nil {
			_, _ = fw.writerCache.Write(message.Bytes())
		}
	}

	_ = fw.writerCache.Close()
}

func (fw *fileWriter) ListenAndServe() {
	for {
		select {
		case log := <-fw.logBuffer:
			if writeErr := fw.writeLog(log); writeErr != nil {
				_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf(`{"log_error": %s, "raw_log": %s}`, writeErr.Error(), log.String()))
			}
		case <-fw.closer:
			return
		}
	}
}

package local

import (
	"bytes"
	"context"
	"github.com/alioth-center/infrastructure/storage"
	"io"
	"os"
	"path/filepath"
)

type localStorage struct {
	basePath string
}

func NewLocalStorage(basePath string) storage.Storage {
	return &localStorage{basePath: basePath}
}

func (s *localStorage) LoadFile(_ context.Context, key string) (content *bytes.Buffer, err error) {
	fileContent, readErr := os.ReadFile(filepath.Join(s.basePath, key))
	if readErr != nil {
		return nil, readErr
	}

	return bytes.NewBuffer(fileContent), nil
}

func (s *localStorage) StoreFile(_ context.Context, key string, content *bytes.Buffer) (err error) {
	file, openErr := os.OpenFile(filepath.Join(s.basePath, key), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if openErr != nil {
		return openErr
	}

	defer func() { _ = file.Close() }()
	if _, copyErr := io.Copy(file, content); copyErr != nil {
		return copyErr
	}

	return nil
}

func (s *localStorage) DeleteFile(_ context.Context, key string) (err error) {
	return os.Remove(filepath.Join(s.basePath, key))
}

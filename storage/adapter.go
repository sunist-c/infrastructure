package storage

import (
	"bytes"
	"context"
)

type Storage interface {
	LoadFile(ctx context.Context, key string) (content *bytes.Buffer, err error)
	StoreFile(ctx context.Context, key string, content *bytes.Buffer) (err error)
	DeleteFile(ctx context.Context, key string) (err error)
}

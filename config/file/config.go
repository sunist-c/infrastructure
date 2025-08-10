package afloader

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/alioth-center/infrastructure/config"
	"gopkg.in/yaml.v3"
)

type fileLoader struct{}

func NewFileLoader() config.Loader {
	return &fileLoader{}
}

func (l *fileLoader) Load(source, namespace string, receiver any) error {
	path := filepath.Join(source, namespace)
	content, readErr := os.ReadFile(path)
	if readErr != nil {
		return readErr
	}

	switch filepath.Ext(path) {
	case ".yaml", ".yml":
		return yaml.Unmarshal(content, receiver)
	case ".json":
		return json.Unmarshal(content, receiver)
	default:
		return errors.New("invalid file type")
	}
}

func (l *fileLoader) MustLoad(source, namespace string, receiver any) {
	if err := l.Load(source, namespace, receiver); err != nil {
		panic(err)
	}
}

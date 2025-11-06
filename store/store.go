package store

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

// FileStore is a thread-safe JSON-based storage handler.
type FileStore[T any] struct {
	FilePath string
	mu       sync.RWMutex
}

// Load loads data from a JSON file into memory.
func (fs *FileStore[T]) Load() ([]T, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	file, err := os.Open(fs.FilePath)
	if errors.Is(err, os.ErrNotExist) {
		return []T{}, nil // empty if not found
	} else if err != nil {
		return nil, err
	}
	defer file.Close()

	var data []T
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return []T{}, nil // treat empty file as empty slice
	}
	return data, nil
}

// Save writes data back to JSON file.
func (fs *FileStore[T]) Save(data []T) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	file, err := os.Create(fs.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

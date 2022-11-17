// Package storage
package storage

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Compile time check to verify implements the Storage interface.
var _ Storage = (*fileSystem)(nil)

// fileSystem implements Storage and provides the ability
// write files to the file system.
type fileSystem struct{}

// NewFileSystem creates a Storage compatible storage for the
// filesystem.
func NewFileSystem() Storage {
	return &fileSystem{}
}

// Put a new object on the file system or overwrites an existing
// contentType is ignored for this storage implementation.
func (s *fileSystem) Put(_ context.Context, dirPath, fileName string, contents []byte, _ bool, _ string) error {
	pth := filepath.Join(dirPath, fileName)
	if err := ioutil.WriteFile(pth, contents, 0644); err != nil {
		return fmt.Errorf("failed to create object: %w", err)
	}
	return nil
}

// Delete deletes an object from the filesystem. It returns nil if the
// object was deleted or if the object doesn't  exists.
func (s *fileSystem) Delete(_ context.Context, dirPath, fileName string) error {
	fp := filepath.Join(dirPath, fileName)
	if err := os.Remove(fp); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete object: %w", err)
	}
	return nil
}

// Get returns the contents for the given object. If the object does not
// exist, it returns error not found.
func (s *fileSystem) Get(_ context.Context, dirPath, fileName string) ([]byte, error) {
	fp := filepath.Join(dirPath, fileName)

	b, err := ioutil.ReadFile(fp)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return b, nil
}

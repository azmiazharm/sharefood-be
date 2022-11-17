// Package storage
package storage

import (
	"context"
	"fmt"
)

var ErrNotFound = fmt.Errorf("storage object not found")

// Storage defines the minimum interface for a blob storage system.
type Storage interface {
	// Put creates or overwrites an object in the storage system.
	// If contentType is blank, the default for the chosen storage implementation is used.
	Put(ctx context.Context, parent, name string, contents []byte, cacheAble bool, contentType string) error

	// Delete deletes an object or does nothing if the object doesn't exist.
	Delete(ctx context.Context, parent, bame string) error

	// Get fetches the object's contents.
	Get(ctx context.Context, parent, name string) ([]byte, error)
}

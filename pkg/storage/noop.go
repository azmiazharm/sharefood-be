// Package storage
package storage

import "context"

// Compile time check to verify implements the Storage interface.
var _ Storage = (*noop)(nil)

// noop is a Storage  dummy for testing purpose
type noop struct{}

func NewNoop() Storage {
	return &noop{}
}

// Put do nothing
func (s *noop) Put(_ context.Context, _, _ string, _ []byte, _ bool, _ string) error {
	return nil
}

// Delete do nothing
func (s *noop) Delete(_ context.Context, _, _ string) error {
	return nil
}

// Get do nothing
func (s *noop) Get(_ context.Context, _, _ string) ([]byte, error) {
	return nil, nil
}

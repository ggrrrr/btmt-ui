package state

import (
	"context"
)

type MockStore struct{}

// Fetch implements StateStore.
func (m *MockStore) Fetch(ctx context.Context, key string) (EntityState, error) {
	return EntityState{}, nil
}

// History implements StateStore.
func (m *MockStore) History(ctx context.Context, key string) ([]EntityState, error) {
	return []EntityState{}, nil
}

// Push implements StateStore.
func (m *MockStore) Push(ctx context.Context, entity NewEntity) (uint64, error) {
	return 0, nil
}

var _ (StateStore) = (*MockStore)(nil)

func NewMockStore() *MockStore {
	return &MockStore{}
}

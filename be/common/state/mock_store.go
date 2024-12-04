package state

import (
	"context"
)

type MockStore struct {
	ReadError  error
	WriteError error
	Data       EntityState
}

// Fetch implements StateStore.
func (m *MockStore) Fetch(ctx context.Context, key string) (EntityState, error) {
	if m.ReadError != nil {
		return EntityState{}, m.ReadError
	}
	return m.Data, nil
}

// History implements StateStore.
func (m *MockStore) History(ctx context.Context, key string) ([]EntityState, error) {
	if m.ReadError != nil {
		return []EntityState{}, m.ReadError
	}
	return []EntityState{}, nil
}

// Push implements StateStore.
func (m *MockStore) Push(ctx context.Context, entity NewEntity) (uint64, error) {
	if m.WriteError != nil {
		return 0, m.WriteError
	}
	return 0, nil
}

var _ (StateStore) = (*MockStore)(nil)

func NewMockStore() *MockStore {
	return &MockStore{}
}

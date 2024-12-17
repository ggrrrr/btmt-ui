package state

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	mock.Mock
}

var _ (StateStore) = (*MockStore)(nil)

// Fetch implements StateStore.
func (m *MockStore) Fetch(ctx context.Context, key string) (EntityState, error) {
	args := m.Called(key)
	return args.Get(0).(EntityState), args.Error(1)
}

// History implements StateStore.
func (m *MockStore) History(ctx context.Context, key string) ([]EntityState, error) {
	args := m.Called(key)
	return args.Get(0).([]EntityState), args.Error(1)
}

// Push implements StateStore.
func (m *MockStore) Push(ctx context.Context, entity NewEntity) (uint64, error) {
	args := m.Called(entity)
	return args.Get(0).(uint64), args.Error(1)
}

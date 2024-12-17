package blob

import (
	"context"
	"io"

	"github.com/stretchr/testify/mock"
)

type MockBlobStore struct {
	mock.Mock
}

var _ (Store) = (*MockBlobStore)(nil)

// Fetch implements Store.
func (m *MockBlobStore) Fetch(ctx context.Context, tenant string, blobId BlobId) (BlobReader, error) {
	args := m.Called(tenant, blobId)
	return args.Get(0).(BlobReader), args.Error(1)
}

// Head implements Store.
func (m *MockBlobStore) Head(ctx context.Context, tenant string, blobId BlobId) (BlobMD, error) {
	args := m.Called(tenant, blobId)
	return args.Get(0).(BlobMD), args.Error(1)
}

// List implements Store.
func (m *MockBlobStore) List(ctx context.Context, tenant string, blobId BlobId) ([]ListResult, error) {
	args := m.Called(tenant, blobId)
	return args.Get(0).([]ListResult), args.Error(1)
}

// Push implements Store.
func (m *MockBlobStore) Push(ctx context.Context, tenant string, blobId BlobId, blobInfo BlobMD, reader io.ReadSeeker) (BlobId, error) {
	args := m.Called(tenant, blobId)
	return args.Get(0).(BlobId), args.Error(1)
}

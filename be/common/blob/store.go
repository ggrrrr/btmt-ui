package blob

import (
	"context"
	"io"
)

type (
	Store interface {
		Fetcher
		Pusher
	}

	Fetcher interface {
		Fetch(ctx context.Context, tenant string, blobId BlobId) (BlobReader, error)
		Head(ctx context.Context, tenant string, blobId BlobId) (BlobMD, error)
		List(ctx context.Context, tenant string, blobId BlobId) ([]ListResult, error)
	}

	Pusher interface {
		// must apply NameRegExp for folder and Id
		Push(ctx context.Context, tenant string, blobId BlobId, blobInfo BlobMD, reader io.ReadSeeker) (BlobId, error)
	}
)

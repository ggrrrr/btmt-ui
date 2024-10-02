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
		// must apply NameRegExp for folder and Id
		Fetch(ctx context.Context, tenant string, blobId string) (*FetchResult, error)
		Head(ctx context.Context, tenant string, blobId string) (*BlobInfo, error)
	}

	Pusher interface {
		// must apply NameRegExp for folder and Id
		Push(ctx context.Context, tenant string, idString string, blobInfo *BlobInfo, reader io.ReadSeeker) (*BlobId, error)
	}
)

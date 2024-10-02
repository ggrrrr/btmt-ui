package blob

import (
	"context"
	"io"
	"time"
)

type (
	BlobMetadata struct {
		// Template, Attachment,
		Type string
		// text/html, text/plan, image/png
		ContentType string
		// Name of the file when downloading or rendering template
		Name string
		// TODO: for future ACL rules
		Owner     string
		CreatedAt time.Time
	}

	FetchResult struct {
		Id         BlockId
		Metadata   BlobMetadata
		ReadCloser io.ReadCloser
	}

	Store interface {
		Fetcher
		Pusher
	}

	Fetcher interface {
		// must apply NameRegExp for folder and Id
		Fetch(ctx context.Context, tenant string, idString string) (*FetchResult, error)
	}

	Pusher interface {
		// must apply NameRegExp for folder and Id
		Push(ctx context.Context, tenant string, idString string, metadata *BlobMetadata, reader io.ReadSeeker) (*BlockId, error)
	}
)

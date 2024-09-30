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
		Name        string
		Version     string
		Id          string
		Owner       string
		CreatedAt   time.Time
	}

	FetchResult struct {
		Metadata BlobMetadata
		Writer   io.WriterTo
	}

	Store interface {
		Fetcher
		Uploader
	}

	Fetcher interface {
		Fetch(ctx context.Context, blobId string, version string) (*FetchResult, error)
	}

	Uploader interface {
		Upload(ctx context.Context, metadata *BlobMetadata, reader io.ReadSeeker) error
	}
)

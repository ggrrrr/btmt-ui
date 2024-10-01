package blob

import (
	"context"
	"io"
	"regexp"
	"time"
)

type (
	BlobMetadata struct {
		// Uniq ID in each folder
		Id string
		// Template, Attachment,
		Type string
		// text/html, text/plan, image/png
		ContentType string
		// Name of the file when downloading or rendering template
		Name    string
		Version string
		// TODO: for future ACL rules
		Owner     string
		CreatedAt time.Time
	}

	FetchResult struct {
		Metadata   BlobMetadata
		ReadCloser io.ReadCloser
	}

	Store interface {
		Fetcher
		Uploader
	}

	Fetcher interface {
		Fetch(ctx context.Context, folder string, blobId string, version string) (*FetchResult, error)
	}

	Uploader interface {
		Upload(ctx context.Context, folder string, metadata *BlobMetadata, reader io.ReadSeeker) error
	}
)

var NameRegExp = regexp.MustCompile(`^[a-zA-Z][0-9a-zA-Z\-]`)

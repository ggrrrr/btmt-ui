package awss3

import (
	"context"
	"fmt"

	"github.com/ggrrrr/btmt-ui/be/common/blob"
)

var _ (blob.Fetcher) = (*Client)(nil)

func (u *Client) Fetch(ctx context.Context, blobId string, version string) (*blob.FetchResult, error) {
	var err error

	versions, err := u.list(ctx, blobId)
	if err != nil {
		return nil, fmt.Errorf("unable to list folder[%s], %w", blobId, err)
	}
	if len(versions) == 0 {
		return nil, blob.NewNotFoundError(blobId, nil)
	}
	lastVersion := versions[len(versions)-1]
	object, err := u.get(ctx, lastVersion)
	if err != nil {
		return nil, err
	}
	object.Metadata.Id = blobId
	// object.Metadata.Version =
	return object, nil
}

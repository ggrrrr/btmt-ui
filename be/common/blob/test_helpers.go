package blob

import (
	"bytes"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBlob(t *testing.T, got *FetchResult, want *BlobMetadata, wantData string, duration int) {
	if duration > 0 {
		gotInLocal := got.Metadata.CreatedAt.In(want.CreatedAt.Location())
		if !assert.WithinDuration(t, want.CreatedAt, gotInLocal, 1+time.Second) {
			fmt.Printf("\t\t %+v \n", want.CreatedAt)
			fmt.Printf("\t\t %+v \n", got.Metadata.CreatedAt)
			fmt.Printf("\t\t %#v \n", want.CreatedAt)
			fmt.Printf("\t\t %#v \n", gotInLocal)
			fmt.Printf("\t\t %#v \n", got.Metadata.CreatedAt)
		}
	}
	got.Metadata.CreatedAt = time.Time{}
	want.CreatedAt = time.Time{}
	assert.Equal(t, want, &got.Metadata)
	gotData := new(bytes.Buffer)
	_, err := io.Copy(gotData, got.ReadCloser)
	assert.NoError(t, err)
	defer got.ReadCloser.Close()

	assert.Equal(t, wantData, gotData.String())
}

package blob

import (
	"bytes"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBlobInfo(t *testing.T, got *BlobInfo, want *BlobInfo, duration int) {
	if duration > 0 {
		gotInLocal := got.CreatedAt.In(want.CreatedAt.Location())
		if !assert.WithinDuration(t, want.CreatedAt, gotInLocal, 1+time.Second) {
			fmt.Printf("\t\t %+v \n", want.CreatedAt)
			fmt.Printf("\t\t %+v \n", got.CreatedAt)
			fmt.Printf("\t\t %#v \n", want.CreatedAt)
			fmt.Printf("\t\t %#v \n", gotInLocal)
			fmt.Printf("\t\t %#v \n", got.CreatedAt)
			fmt.Printf("\t\t %#v \n", got.ContentLength)
		}
	}

	got.CreatedAt = time.Time{}
	want.CreatedAt = time.Time{}
	assert.Equal(t, want, got)
}

func TestBlobInfoData(t *testing.T, got *FetchResult, want *BlobInfo, wantData string, duration int) {
	gotData := new(bytes.Buffer)
	_, err := io.Copy(gotData, got.ReadCloser)
	assert.NoError(t, err)
	defer got.ReadCloser.Close()

	assert.Equal(t, wantData, gotData.String())
}

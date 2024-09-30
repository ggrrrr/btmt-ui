package blob

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBlob(t *testing.T, got *FetchResult, want *BlobMetadata, wantData string, duration int) {
	if duration > 0 {
		gotInLocal := got.Metadata.CreatedAt.In(want.CreatedAt.Location())
		fmt.Printf("\t\t %+v \n", want.CreatedAt)
		fmt.Printf("\t\t %+v \n", got.Metadata.CreatedAt)
		fmt.Printf("\t\t %#v \n", want.CreatedAt)
		fmt.Printf("\t\t %#v \n", gotInLocal)
		fmt.Printf("\t\t %#v \n", got.Metadata.CreatedAt)
		assert.WithinDuration(t, want.CreatedAt, gotInLocal, 1+time.Second)
	}
	got.Metadata.CreatedAt = time.Time{}
	want.CreatedAt = time.Time{}
	assert.Equal(t, want, &got.Metadata)
	gotData := new(bytes.Buffer)
	_, err := got.Writer.WriteTo(gotData)
	assert.NoError(t, err)

	assert.Equal(t, wantData, gotData.String())
}

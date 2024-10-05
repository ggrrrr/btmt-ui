package blob

import (
	"bytes"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBlobInfo(subTest string, t *testing.T, actual BlobInfo, expected BlobInfo, duration int) {
	if duration > 0 {
		gotInLocal := actual.CreatedAt.In(expected.CreatedAt.Location())
		if !assert.WithinDurationf(t, expected.CreatedAt, gotInLocal, 1+time.Second, subTest) {
			fmt.Printf("\t\t%10s %+v \n", "expected", expected.CreatedAt)
			fmt.Printf("\t\t%10s %+v \n", "actual", actual.CreatedAt)
			fmt.Printf("\t\t%10s %#v \n", "expected", expected.CreatedAt)
			fmt.Printf("\t\t%10s %#v \n", "actual", gotInLocal)
			fmt.Printf("\t\t%10s %#v \n", "actual", actual.CreatedAt)
			fmt.Printf("\t\t%10s %#v \n", "actual", actual.ContentLength)
		}
	}

	actual.CreatedAt = time.Time{}
	expected.CreatedAt = time.Time{}
	assert.Equalf(t, expected, actual, subTest)
}

func TestBlobInfoData(t *testing.T, got FetchResult, want BlobInfo, wantData string, duration int) {
	gotData := new(bytes.Buffer)
	_, err := io.Copy(gotData, got.ReadCloser)
	assert.NoError(t, err)
	defer got.ReadCloser.Close()

	assert.Equal(t, wantData, gotData.String())
}

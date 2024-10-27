package image

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/ggrrrr/btmt-ui/be/common/blob"
	"github.com/ggrrrr/btmt-ui/be/help"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHead(t *testing.T) {

	pwd := help.RepoDir()

	want := blob.MDImageInfo{
		Width:  256,
		Height: 256,
	}

	actual, err := HeadImage(fmt.Sprintf("%s/glass-mug-variant.png", pwd))
	require.NoError(t, err)
	assert.Equal(t, want, actual)

	_, err = HeadImage(fmt.Sprintf("%s/notfound", pwd))
	require.Error(t, err)
	fmt.Printf("%v \n", err)

	_, err = HeadImage(fmt.Sprintf("%s/test.txt", pwd))
	require.Error(t, err)
	fmt.Printf("%v \n", err)

}

func TestResize(t *testing.T) {

	pwd := help.RepoDir()

	sourceFile := fmt.Sprintf("%s/big_image.jpg", pwd)
	fmt.Printf("source: %s \n", sourceFile)

	image1, err := os.Open(sourceFile)
	require.NoError(t, err)

	stat, err := image1.Stat()
	require.NoError(t, err)

	fmt.Printf("source: %s %+v \n", sourceFile, stat)

	resizedImage, err := ResizeImage(context.TODO(), 256/2, blob.BlobReader{
		Blob:       blob.Blob{},
		ReadCloser: image1,
	})
	require.NoError(t, err)

	fmt.Printf("resizedImage: %+v \n", resizedImage)
	// defer os.Remove(resizedImage.TmpFile)

}

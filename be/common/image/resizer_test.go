package image

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ggrrrr/btmt-ui/be/common/blob"
	"github.com/ggrrrr/btmt-ui/be/help"
)

func TestHead(t *testing.T) {

	pwd := help.RepoDir()
	fmt.Println("RepoDir: ", pwd)

	want := blob.MDImageInfo{
		Width:  256,
		Height: 256,
	}

	actual, err := HeadImage(fmt.Sprintf("%s/glass-mug-variant.png", pwd))
	require.NoError(t, err)
	assert.Equal(t, want, actual)

	_, err = HeadImage(fmt.Sprintf("%s/notfound", pwd))
	require.Error(t, err)
	expErr := &fs.PathError{}
	assert.ErrorAs(t, err, &expErr)
	fmt.Printf("HeadImage: %#v \n", err)

	_, err = HeadImage(fmt.Sprintf("%s/test.txt", pwd))
	require.Error(t, err)
	fmt.Printf("HeadImage: %#v \n", err)

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

	fmt.Printf("tmpfile: %v \n", resizedImage.ReadCloser.tmpFile.Name())
	fmt.Printf("resizedImage: %+v \n", resizedImage)
	defer func() {
		err = resizedImage.ReadCloser.Close()
		fmt.Printf("ReadCloser.Close: %+v \n", err)
	}()

}

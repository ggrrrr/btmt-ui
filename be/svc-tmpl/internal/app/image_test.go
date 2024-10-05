package app

import (
	"fmt"
	"testing"

	"github.com/ggrrrr/btmt-ui/be/common/blob"
	"github.com/ggrrrr/btmt-ui/be/help"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImage(t *testing.T) {

	pwd := help.RepoDir()

	want := blob.ImageInfo{
		Width:  256,
		Height: 256,
	}

	actual, err := headImage(fmt.Sprintf("%s/glass-mug-variant.png", pwd))
	require.NoError(t, err)
	assert.Equal(t, want, actual)

	_, err = headImage(fmt.Sprintf("%s/notfound", pwd))
	require.Error(t, err)
	fmt.Printf("%v \n", err)

	_, err = headImage(fmt.Sprintf("%s/test.txt", pwd))
	require.Error(t, err)
	fmt.Printf("%v \n", err)

}

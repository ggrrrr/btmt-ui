package blob

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseBlobId(t *testing.T) {
	tests := []struct {
		name    string
		fromStr string
		id      BlockId
		err     error
	}{
		{
			name:    "ok",
			fromStr: "Mydir-1/My-id@1",
			id: BlockId{
				folder:  "Mydir-1",
				id:      "My-id",
				version: "1",
			},
			err: nil,
		},
		{
			name:    "ok no ver",
			fromStr: "Mydir-1/My-id",
			id: BlockId{
				folder:  "Mydir-1",
				id:      "My-id",
				version: "",
			},
			err: nil,
		},
		{
			name:    "ok empty ver",
			fromStr: "Mydir-1/My-id@",
			id: BlockId{
				folder:  "Mydir-1",
				id:      "My-id",
				version: "",
			},
			err: nil,
		},
		{
			name:    "err empty str",
			fromStr: "",
			// id:      nil,
			err: &BlobIdInputEmptyError{},
		},
		{
			name:    "from",
			fromStr: "from",
			// id:      ,
			err: &BlobIdInputError{},
		},
		{
			name:    "123/123@123",
			fromStr: "123/123@123",
			// id:      nil,
			err: &BlobIdInputError{},
		},
		{
			name:    "123asdasd/asd@123",
			fromStr: "123asdasd/asd@123",
			// id:      nil,
			err: &BlobIdInputError{},
		},
		{
			name:    "asdasd  ",
			fromStr: "asdasd  ",
			// id:      nil,
			err: &BlobIdInputError{},
		},
		{
			name:    "asdasd/asdasd  ",
			fromStr: "asdasd/asdasd  ",
			// id:      nil,
			err: &BlobIdInputError{},
		},
		{
			name:    "asd/asd asd",
			fromStr: "asd/asd asd",
			// id:      nil,
			err: &BlobIdInputError{},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			toId, err := ParseBlobId(tc.fromStr)
			if tc.err != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.id, toId)
				// assert.Equal(t, tc.fromStr, toId.String())
			}
		})
	}

	testId1, _ := ParseBlobId("folder1/id-1@ver1")
	assert.Equal(t, "folder1/id-1@ver1", testId1.String())

	testId2, _ := ParseBlobId("folder1/id-1")
	assert.Equal(t, "folder1/id-1", testId2.String())
}

func TestBlockIdRefExp(t *testing.T) {
	tests := []struct {
		fromStr string
		result  []string
	}{
		{
			fromStr: "",
			result:  nil,
		},
		{
			fromStr: "asd",
			result:  nil,
		},
		{
			fromStr: "asd/id",
			result:  []string{"asd/id", "asd", "id", ""},
		},
		{
			fromStr: "folder/id@ver1",
			result:  []string{"folder/id@ver1", "folder", "id", "ver1"},
		},
		{
			fromStr: "folder/id@22",
			result:  []string{"folder/id@22", "folder", "id", "22"},
		},
		{
			fromStr: "1folder/id@2",
			result:  nil,
		},
		{
			fromStr: "folder/1id@2",
			result:  nil,
		},
	}
	for _, tc := range tests {
		result := BlockIdRefExp.FindStringSubmatch(tc.fromStr)
		assert.Equalf(t, tc.result, result, "fromStr: %s\n\t\t%#v\n\t\t%#v", tc.fromStr, tc.result, result)
	}
}

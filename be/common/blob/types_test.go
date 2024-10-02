package blob

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBasicTypes(t *testing.T) {
	tests := []struct {
		name string
		f    func(t *testing.T)
	}{
		{
			name: "blob id",
			f: func(t *testing.T) {
				id := BlobId{
					folder:  "f",
					id:      "id",
					version: "ver",
				}
				assert.Equal(t, "f", id.Folder())
				assert.Equal(t, "id", id.Id())
				assert.Equal(t, "ver", id.Version())
				assert.Equal(t, "f/id@ver", id.String())
				assert.Equal(t, "f/id", id.Key())
				id.SetVersion("")
				assert.Equal(t, "f/id", id.String())

				id2 := New("folder", "id", "ver")
				assert.Equal(t, "folder/id@ver", id2.String())

			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, tc.f)
	}
}

func TestParseBlobId(t *testing.T) {
	tests := []struct {
		name    string
		fromStr string
		id      BlobId
		err     error
		skip    bool
	}{
		{
			name:    "ok",
			fromStr: "Mydir-1/My-id@1",
			id: BlobId{
				folder:  "Mydir-1",
				id:      "My-id",
				version: "1",
			},
			err: nil,
		},
		{
			name:    "ok no ver",
			fromStr: "Mydir-1/My-id",
			id: BlobId{
				folder:  "Mydir-1",
				id:      "My-id",
				version: "",
			},
			err: nil,
		},
		{
			name:    "ok empty ver",
			fromStr: "Mydir-1/My-id@",
			id: BlobId{
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
			err:  &BlobIdInputError{},
			skip: true,
		},
		{
			name:    "asd/asd asd",
			fromStr: "asd/asd asd",
			// id:      nil,
			err:  &BlobIdInputError{},
			skip: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			toId, err := ParseBlobId(tc.fromStr)
			if tc.skip {
				// fmt.Printf("err[%s]: toId:%v \n", err, toId.String())
				t.Skipf("skip from:%s  err[%s]: toId:%v \n", tc.fromStr, err, toId.String())
			}
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
	require.Equal(t, "folder1/id-1@ver1", testId1.String())

	testId2, _ := ParseBlobId("folder1/id-1")
	require.Equal(t, "folder1/id-1", testId2.String())
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

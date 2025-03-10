package blob

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ggrrrr/btmt-ui/be/common/app"
)

func TestFileNameFilter(t *testing.T) {
	newName := FileNameFilter("file cdfile ASD (12) .png")
	require.Equal(t, "filecdfileASD12.png", newName)
}

func TestBasicTypes(t *testing.T) {
	id := BlobId{
		path:    "f",
		id:      "id",
		version: "ver",
	}
	assert.Equal(t, "f", id.Path())
	assert.Equal(t, "id", id.Id())
	assert.Equal(t, "ver", id.Version())
	assert.Equal(t, "f/id:ver", id.String())

	assert.Equal(t, "f/id", id.PathId())
	assert.Equal(t, "id:ver", id.IdVersion())

	tests := []struct {
		name string
		f    func(t *testing.T)
	}{
		{
			name: "blob id",
			f: func(t *testing.T) {
				id2, err := NewBlobId("folder", "id", "ver")
				require.NoError(t, err)
				assert.Equal(t, "folder/id:ver", id2.String())

			},
		},
		{
			name: "err folder",
			f: func(t *testing.T) {
				_, err := NewBlobId("folder ", "id", "ver")
				require.Error(t, err)
				// assert.Equal(t, "folder/id:ver", id2.String())
				_, err = NewBlobId(" folder", "id", "ver")
				require.Error(t, err)
				// assert.Equal(t, "folder/id:ver", id2.String())
				_, err = NewBlobId("fol:der", "id", "ver")
				require.Error(t, err)
				_, err = NewBlobId("fol der", "id", "ver")
				require.Error(t, err)

			},
		},
		{
			name: "err name",
			f: func(t *testing.T) {
				_, err := NewBlobId("folder", "asd/asd", "ver")
				require.Error(t, err)
				_, err = NewBlobId("folder", "asd asd", "ver")
				require.Error(t, err)
				_, err = NewBlobId("folder", "asd_asd", "ver")
				require.Error(t, err)
			},
		},
		{
			name: "err folder",
			f: func(t *testing.T) {
				_, err := NewBlobId("folde", "id", "ver:")
				require.Error(t, err)
				_, err = NewBlobId("folde", "id", "ver asd")
				require.Error(t, err)
				_, err = NewBlobId("folde", "id", "ver/asd")
				require.Error(t, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, tc.f)
	}
}
func TestBlobDir(t *testing.T) {
	var err error
	var blobIdEmpty BlobId
	var blobId BlobId

	blobId, err = ParseBlobDir("")
	require.Error(t, err)
	require.Equal(t, blobId, blobIdEmpty)

	blobId, err = ParseBlobDir("fodler1")
	require.NoError(t, err)
	require.Equal(t, blobId, BlobId{path: "fodler1"})

	blobId, err = ParseBlobDir("fodler1/folder2")
	require.NoError(t, err)
	require.Equal(t, blobId, BlobId{path: "fodler1/folder2"})

}

func TestBlobIdSetIdVersion(t *testing.T) {

	tests := []struct {
		name       string
		idVersion  string
		fromBlobID BlobId
		toBlobId   BlobId
		err        error
	}{
		{
			name:       "ok id version",
			idVersion:  "filename:version",
			fromBlobID: BlobId{path: "folder-1"},
			toBlobId:   BlobId{path: "folder-1", id: "filename", version: "version"},
			err:        nil,
		},
		{
			name:       "ok id version",
			idVersion:  "filename:version",
			fromBlobID: BlobId{path: "folder-1"},
			toBlobId:   BlobId{path: "folder-1", id: "filename", version: "version"},
			err:        nil,
		},
		{
			name:       "ok id no version ",
			idVersion:  "filename",
			fromBlobID: BlobId{path: "folder-1"},
			toBlobId:   BlobId{path: "folder-1", id: "filename", version: ""},
			err:        nil,
		},
		{
			name:       "ok id: no version ",
			idVersion:  "filename:",
			fromBlobID: BlobId{path: "folder-1"},
			toBlobId:   BlobId{path: "folder-1", id: "filename", version: ""},
			err:        nil,
		},
		{
			name:       "err id ",
			idVersion:  "fi  lename:ad",
			fromBlobID: BlobId{path: "folder-1"},
			err:        &app.AppError{},
		},
		{
			name:       "err id ",
			idVersion:  "fi/lename:ad",
			fromBlobID: BlobId{},
			err:        &app.AppError{},
		},
		{
			name:       "err version ",
			idVersion:  "fi  lename://ad",
			fromBlobID: BlobId{path: "folder-1"},
			toBlobId:   BlobId{},
			err:        &app.AppError{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			newBlobId, err := tc.fromBlobID.SetIdVersionFromString(tc.idVersion)
			if tc.err == nil {
				require.NoError(t, err)
				assert.Equal(t, tc.toBlobId, newBlobId)
			} else {
				require.Error(t, err)

			}

		})
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
			fromStr: "Mydir-1/Mydir-2/My-id:1",
			id: BlobId{
				path:    "Mydir-1/Mydir-2",
				id:      "My-id",
				version: "1",
			},
			err: nil,
		},
		{
			name:    "ok no ver",
			fromStr: "Mydir-1/My-id",
			id: BlobId{
				path:    "Mydir-1",
				id:      "My-id",
				version: "",
			},
			err: nil,
		},
		{
			name:    "ok empty ver",
			fromStr: "Mydir-1/My-id:",
			id: BlobId{
				path:    "Mydir-1",
				id:      "My-id",
				version: "",
			},
			err: nil,
		},
		{
			name:    "err empty str",
			fromStr: "",
			// id:      nil,
			// err: &BlobIdInputEmptyError{},
			err: &app.AppError{},
		},
		{
			name:    "from",
			fromStr: "from",
			id:      BlobId{path: "from", id: "", version: ""},
			// err: &BlobIdInputError{},
			err: nil,
		},
		{
			name:    "123/123@123",
			fromStr: "123/123@123",
			// id:      nil,
			// err: &BlobIdInputError{},
			err: &app.AppError{},
		},
		{
			name:    "123asdasd/asd@123",
			fromStr: "123asdasd/asd@123",
			// id:      nil,
			// err: &BlobIdInputError{},
			err: &app.AppError{},
		},
		{
			name:    "asdasd  ",
			fromStr: "asdasd  ",
			// id:      nil,
			// err: &BlobIdInputError{},
			err: &app.AppError{},
		},
		{
			name:    "asdasd/asdasd  ",
			fromStr: "asdasd/asdasd  ",
			// id:      nil,
			// err:  &BlobIdInputError{},
			err:  &app.AppError{},
			skip: false,
		},
		{
			name:    "asd/asd asd",
			fromStr: "asd/asd asd",
			// id:      nil,
			// err:  &BlobIdInputError{},
			err: &app.AppError{},
			// skip: true,
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
				require.Errorf(t, err, "result: %+v", toId)

			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.id, toId)
				// assert.Equal(t, tc.fromStr, toId.String())
			}
		})
	}

	testId1, _ := ParseBlobId("folder1/id-1:ver1")
	require.Equal(t, "folder1/id-1:ver1", testId1.String())

	testId2, _ := ParseBlobId("folder1/id-1")
	require.Equal(t, "folder1/id-1", testId2.String())
}

func TestNameRegExp(t *testing.T) {
	tests := []struct {
		fromStr  string
		expected bool
	}{
		{
			fromStr:  "asdasda/sdas/d123-123.png:123",
			expected: false,
		},
		{
			fromStr:  "asdasdasdasd123123-/asdea123asd123:123123",
			expected: false,
		},
		{
			fromStr:  "asdas123123ASD/asASDd-asd123:aASDsd123123",
			expected: false,
		},
		{
			fromStr:  " asd",
			expected: true,
		},
		{
			fromStr:  "asd ",
			expected: true,
		},
		//
		{
			fromStr:  "a sd/a sd asd",
			expected: true,
		},
		{
			fromStr:  "as#d",
			expected: true,
		},
		{
			fromStr:  `as%dasd`,
			expected: true,
		},
		{
			fromStr:  "as`d",
			expected: true,
		},
		{
			fromStr:  "as[d",
			expected: true,
		},
		{
			fromStr:  "as]d",
			expected: true,
		},
		{
			fromStr:  `as\d`,
			expected: true,
		},
		{
			fromStr:  `as=d`,
			expected: true,
		},
		{
			fromStr:  `folder1/folder2 /name.png:1`,
			expected: true,
		},
		//
		// {
		// fromStr:  "a sdasdasdasd123123-/asdea123asd123:123123",
		// expected: false,
		// },
	}

	for _, tc := range tests {
		result := NotAllowedCharRegEx.MatchString(tc.fromStr)
		// result := NameRegExp.FindAllStringSubmatch(tc.fromStr, 10)
		if !assert.Equalf(t, tc.expected, result, "input string:[%s] dost not Match %+v", tc.fromStr, tc.expected) {
			submatch := NotAllowedCharRegEx.FindAllStringSubmatch(tc.fromStr, 10)
			fmt.Printf("\n\t\tsubmatch: %+v \n", submatch)
		}
	}
}

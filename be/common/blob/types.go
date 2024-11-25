package blob

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/ggrrrr/btmt-ui/be/common/app"
)

type (
	BlobType string

	// Id of each blob object consist of path, name and version
	// version can be empty
	BlobId struct {
		// somefolder1/someothjerfolder2
		path string
		// somefilename.png
		id string
		// ver1 or 1
		version string
	}

	BlobMD struct {
		// Template, Attachment,
		Type BlobType
		// text/html, text/plan, image/png
		ContentType string
		ImageInfo   MDImageInfo
		// Name of the file when downloading or rendering template
		Name string
		// TODO: for future ACL rules
		Owner         string
		CreatedAt     time.Time
		ContentLength int64
	}

	Blob struct {
		Id BlobId
		MD BlobMD
	}

	MDImageInfo struct {
		Width  int64
		Height int64
	}

	BlobReader struct {
		Blob       Blob
		ReadCloser io.ReadCloser
	}

	ListResult struct {
		Blob
		Versions []Blob
	}
)

var (
	BlobTypeImage      BlobType = "image"
	BlobTypeTemplate   BlobType = "template"
	BlobTypeAttachment BlobType = "Attachment"
	BlobTypeGPX        BlobType = "gpx"
)

// https://yourbasic.org/golang/regexp-cheat-sheet/
var NotAllowedCharRegEx = regexp.MustCompile(`[^0-9a-zA-Z\:\-\/\.]`)
var PathRegExp = regexp.MustCompile(`^[0-9a-zA-Z\-\/\.]*$`)
var NameRegExp = regexp.MustCompile(`^[0-9a-zA-Z\-\.]*$`)
var NameFilterRegExp = regexp.MustCompile(`[^0-9a-zA-Z\-\.]`)

func FileNameFilter(fromName string) string {
	return NameFilterRegExp.ReplaceAllString(fromName, "")
}

func (id BlobId) Id() string {
	return id.id
}

// return folder/name
func (id BlobId) PathId() string {
	return fmt.Sprintf("%s/%s", id.path, id.id)
}

// return folder/name
func (id BlobId) IdVersion() string {
	return fmt.Sprintf("%s:%s", id.id, id.version)
}

func (id BlobId) Path() string {
	return id.path
}

func (id BlobId) Version() string {
	return id.version
}

func (blobId BlobId) String() string {
	if blobId.version == "" {
		return fmt.Sprintf("%s/%s", blobId.path, blobId.id)
	}
	return fmt.Sprintf("%s/%s:%s", blobId.path, blobId.id, blobId.version)
}

// Set ID and Version of a blob
// if Version is empty sets only ID (filename)
func (blobId BlobId) SetIdVersionFromString(idVersion string) (BlobId, error) {
	if idVersion == "" {
		return BlobId{}, app.BadRequestError("id is empty", nil)
	}
	split := strings.Split(idVersion, ":")
	id := ""
	version := ""

	ok := NameRegExp.MatchString(split[0])
	if !ok {
		return BlobId{}, app.BadRequestError("id incorrect string", nil)
	}
	id = split[0]

	if len(split) > 1 && len(split[1]) > 0 {
		ok = NameRegExp.MatchString(split[1])
		if !ok {
			return BlobId{}, app.BadRequestError("version incorrect string", nil)
		}
		version = split[1]
	}

	newBlobId := BlobId{
		path:    blobId.path,
		id:      id,
		version: version,
	}

	return newBlobId, nil

}

func NewBlobId(path, id, ver string) (BlobId, error) {
	// TODO validate strings
	result := PathRegExp.MatchString(path)
	if !result {
		return BlobId{}, fmt.Errorf("path incorrect string")
	}
	result = NameRegExp.MatchString(id)
	if !result {
		return BlobId{}, fmt.Errorf("name incorrect string")
	}
	result = NameRegExp.MatchString(ver)
	if !result {
		return BlobId{}, fmt.Errorf("version incorrect string")
	}

	return BlobId{
		path:    path,
		id:      id,
		version: ver,
	}, nil

}

func ParseBlobDir(fromString string) (BlobId, error) {
	if fromString == "" {
		return BlobId{}, app.BadRequestError("id is empty", nil)
	}

	result := NotAllowedCharRegEx.MatchString(fromString)
	if result {
		return BlobId{}, app.BadRequestError("string with not allowed characters", nil)
	}

	path := ""

	folders := strings.Split(fromString, "/")
	fLen := len(folders)
	if fLen == 0 {
		return BlobId{}, app.BadRequestError("empty path", nil)
	}

	path = strings.Join(folders[:fLen], "/")
	if fLen == 1 {
		path = folders[fLen-1]
	}

	return BlobId{
		path: path,
	}, nil
}

func ParseBlobId(fromString string) (BlobId, error) {
	if fromString == "" {
		return BlobId{}, app.BadRequestError("id is empty", nil)
	}

	result := NotAllowedCharRegEx.MatchString(fromString)
	if result {
		return BlobId{}, app.BadRequestError("string with not allowed characters", nil)
	}

	path := ""
	id := ""
	version := ""

	folders := strings.Split(fromString, "/")
	fLen := len(folders)
	if fLen == 0 {
		return BlobId{}, app.BadRequestError("empty path", nil)
	}

	path = strings.Join(folders[:fLen-1], "/")
	if fLen == 1 {
		path = folders[fLen-1]
	}

	if fLen > 1 {
		fileVer := strings.Split(folders[fLen-1], ":")
		if len(fileVer) > 2 {
			return BlobId{}, app.BadRequestError("wrong file version", nil)
		}
		if len(fileVer) == 0 {
			return BlobId{}, app.BadRequestError("no file name", nil)
		}
		id = fileVer[0]
		if len(fileVer) > 1 {
			version = fileVer[1]
		}
	}

	return BlobId{
		path:    path,
		id:      id,
		version: version,
	}, nil
}

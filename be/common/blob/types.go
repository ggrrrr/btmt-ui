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
	// Id of each blob object consist of path, name and version
	// version can be empty
	BlobId struct {
		// somefolder1/someothjerfolder2
		path string
		// somefilename.png
		name string
		// ver1 or 1
		version string
	}

	ImageInfo struct {
		Width  int64
		Height int64
	}

	BlobInfo struct {
		// Template, Attachment,
		Type string
		// text/html, text/plan, image/png
		ContentType string
		ImageInfo   ImageInfo
		// Name of the file when downloading or rendering template
		Name string
		// TODO: for future ACL rules
		Owner         string
		CreatedAt     time.Time
		ContentLength int64
	}

	FetchResult struct {
		Id         BlobId
		Info       BlobInfo
		ReadCloser io.ReadCloser
	}

	HeadResult struct {
		Id       BlobId
		Metadata BlobInfo
		// ReadCloser io.ReadCloser
	}
)

// https://yourbasic.org/golang/regexp-cheat-sheet/
var NotAllowedCharRegEx = regexp.MustCompile(`[^0-9a-zA-Z\:\-\/\.]`)
var PathRegExp = regexp.MustCompile(`^[0-9a-zA-Z\-\/\.]*$`)
var NameRegExp = regexp.MustCompile(`^[0-9a-zA-Z\-\.]*$`)
var NameFilterRegExp = regexp.MustCompile(`[^0-9a-zA-Z\-\.]`)

func FileNameFilter(fromName string) string {
	return NameFilterRegExp.ReplaceAllString(fromName, "")
}

func (id *BlobId) Name() string {
	return id.name
}

// return folder/id
func (id *BlobId) Key() string {
	return fmt.Sprintf("%s/%s", id.path, id.name)
}

func (id *BlobId) Path() string {
	return id.path
}

func (id *BlobId) Version() string {
	return id.version
}

func (id *BlobId) SetVersion(ver string) {
	id.version = ver
}

func (id *BlobId) String() string {
	if id.version == "" {
		return fmt.Sprintf("%s/%s", id.path, id.name)
	}
	return fmt.Sprintf("%s/%s:%s", id.path, id.name, id.version)
}

func NewBlobId(path, name, ver string) (BlobId, error) {
	// TODO validate strings
	result := PathRegExp.MatchString(path)
	if !result {
		return BlobId{}, fmt.Errorf("path incorrect string")
	}
	result = NameRegExp.MatchString(name)
	if !result {
		return BlobId{}, fmt.Errorf("name incorrect string")
	}
	result = NameRegExp.MatchString(ver)
	if !result {
		return BlobId{}, fmt.Errorf("version incorrect string")
	}

	return BlobId{
		path:    path,
		name:    name,
		version: ver,
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
	name := ""
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
		name = fileVer[0]
		if len(fileVer) > 1 {
			version = fileVer[1]
		}
	}

	return BlobId{
		path:    path,
		name:    name,
		version: version,
	}, nil
}

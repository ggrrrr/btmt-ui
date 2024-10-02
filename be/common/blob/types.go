package blob

import (
	"fmt"
	"io"
	"regexp"
	"time"
)

// https://yourbasic.org/golang/regexp-cheat-sheet/
var NameRegExp = regexp.MustCompile(`^[a-zA-Z][0-9a-zA-Z\-]`)

var BlockIdRefExp = regexp.MustCompile(`(^[a-zA-Z][a-zA-Z0-9\-]*)/([a-zA-Z][a-zA-Z0-9\-]*)[\@]?([a-zA-Z0-9\-]*)*`)

type (
	// Id of each blob object with or without version
	BlobId struct {
		folder  string
		id      string
		version string
	}

	BlobInfo struct {
		// Template, Attachment,
		Type string
		// text/html, text/plan, image/png
		ContentType string
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

func (id *BlobId) Id() string {
	return id.id
}

// return folder/id
func (id *BlobId) Key() string {
	return fmt.Sprintf("%s/%s", id.folder, id.id)
}

func (id *BlobId) Folder() string {
	return id.folder
}

func (id *BlobId) Version() string {
	return id.version
}

func (id *BlobId) SetVersion(ver string) {
	id.version = ver
}

func (id *BlobId) String() string {
	if id.version == "" {
		return fmt.Sprintf("%s/%s", id.folder, id.id)
	}
	return fmt.Sprintf("%s/%s@%s", id.folder, id.id, id.version)
}

func NewBlobId(folder, id, ver string) BlobId {
	return BlobId{
		folder:  folder,
		id:      id,
		version: ver,
	}
}

func ParseBlobId(fromId string) (BlobId, error) {
	if fromId == "" {
		return BlobId{}, &BlobIdInputEmptyError{}
	}

	result := BlockIdRefExp.FindStringSubmatch(fromId)
	if result == nil {
		return BlobId{}, &BlobIdInputError{from: fromId}
	}

	if len(result) == 4 {
		return BlobId{
			folder:  result[1],
			id:      result[2],
			version: result[3],
		}, nil
	}

	return BlobId{}, &BlobIdInputError{from: fromId}
}

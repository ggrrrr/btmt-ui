package blob

import (
	"fmt"
	"regexp"
)

// https://yourbasic.org/golang/regexp-cheat-sheet/
var NameRegExp = regexp.MustCompile(`^[a-zA-Z][0-9a-zA-Z\-]`)

var BlockIdRefExp = regexp.MustCompile(`(^[a-zA-Z][a-zA-Z0-9\-]*)/([a-zA-Z][a-zA-Z0-9\-]*)[\@]?([a-zA-Z0-9\-]*)*`)

type BlockId struct {
	folder  string
	id      string
	version string
}

func (id *BlockId) Id() string {
	return id.id
}

// return folder/id
func (id *BlockId) Key() string {
	return fmt.Sprintf("%s/%s", id.folder, id.id)
}

func (id *BlockId) Folder() string {
	return id.folder
}

func (id *BlockId) Version() string {
	return id.version
}

func (id *BlockId) SetVersion(ver string) {
	id.version = ver
}

func (id *BlockId) String() string {
	if id.version == "" {
		return fmt.Sprintf("%s/%s", id.folder, id.id)
	}
	return fmt.Sprintf("%s/%s@%s", id.folder, id.id, id.version)
}

func New(folder, id, ver string) BlockId {
	return BlockId{
		folder:  folder,
		id:      id,
		version: ver,
	}
}

func ParseBlobId(fromId string) (BlockId, error) {
	if fromId == "" {
		return BlockId{}, &BlobIdInputEmptyError{}
	}

	result := BlockIdRefExp.FindStringSubmatch(fromId)
	if result == nil {
		return BlockId{}, &BlobIdInputError{from: fromId}
	}

	if len(result) == 4 {
		return BlockId{
			folder:  result[1],
			id:      result[2],
			version: result[3],
		}, nil
	}

	return BlockId{}, &BlobIdInputError{from: fromId}
}

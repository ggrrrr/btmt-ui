package awss3

import (
	"fmt"
	"time"

	"github.com/ggrrrr/btmt-ui/be/common/blob"
)

type awsId struct {
	folder string
	id     string
	ver    string
}

func timeInLocal(from *time.Time) time.Time {
	if from == nil {
		return time.Now()
	}
	to := from.Local()
	return to
}

func awsIdFromString(fromStr string) (awsId, error) {
	blobId, err := blob.ParseBlobId(fromStr)
	if err != nil {
		return awsId{}, err
	}
	return awsId{
		folder: blobId.Folder(),
		id:     blobId.Id(),
		ver:    blobId.Version(),
	}, nil

}

// folder/id
func (i awsId) idFolder() string {
	return fmt.Sprintf("%s/%s", i.folder, i.id)
}

// folder/id/ver
func (i awsId) keyVer() string {
	if i.ver == "" {
		return fmt.Sprintf("%s/%s", i.folder, i.id)
	}
	return fmt.Sprintf("%s/%s/%s", i.folder, i.id, i.ver)
}

// folder/id/ver
func (i awsId) String() string {
	return i.keyVer()
}

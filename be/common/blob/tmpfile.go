package blob

import (
	"context"
	"os"

	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
	// "github.com/ggrrrr/btmt-ui/be/common/logger"
)

type TempFile struct {
	FileName      string
	TempFileName  string
	ContentType   string
	ContentLength int64
}

func (f TempFile) Delete(ctx context.Context) {
	err := os.Remove(f.TempFileName)
	if err != nil {
		log.Log().WarnCtx(ctx, err, "Delete",
			log.WithString("temp.file", f.TempFileName),
		)

		return
	}

	if log.Log().IsTrace() {
		log.Log().DebugCtx(ctx, "Deleted", log.WithString("temp.file", f.TempFileName))
	}
}

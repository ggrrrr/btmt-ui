package blob

import (
	"context"
	"os"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
	// "github.com/ggrrrr/btmt-ui/be/common/logger"
)

type TempFile struct {
	FileName     string
	TempFileName string
	ContentType  string
}

func (f TempFile) Delete(ctx context.Context) {
	err := os.Remove(f.TempFileName)
	if err != nil {
		logger.WarnCtx(ctx).
			Err(err).
			Str("temp.file", f.TempFileName).
			Msg("TempFile")
		return
	}

	if logger.IsTrace() {
		logger.DebugCtx(ctx).
			Str("temp.file", f.TempFileName).
			Msg("Deleted")
	}
}